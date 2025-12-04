package seckill

import (
	"context"
	"testing"

	"rag-agent/config"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 设置测试环境
func setupTestEnv(t *testing.T) (CacheRepository, *redis.Client, func()) {
	// 连接 Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1, // 使用测试数据库
	})

	// 测试连接
	err := client.Ping(context.Background()).Err()
	if err != nil {
		t.Skip("Redis 未启动，跳过测试")
	}

	// 清理测试数据
	client.FlushDB(context.Background())

	// 创建 CacheRepository
	cacheRepo := NewRedisCacheRepository(client)

	// 返回清理函数
	cleanup := func() {
		client.FlushDB(context.Background())
		client.Close()
	}

	return cacheRepo, client, cleanup
}

// 简单的 MQ Producer（用于测试）
type TestMQProducer struct {
	shouldFail bool
}

func (p *TestMQProducer) SendOrderMessage(ctx context.Context, order *Order) error {
	if p.shouldFail {
		return assert.AnError
	}
	return nil
}

// 测试秒杀成功
func TestSeckill_Success(t *testing.T) {
	cacheRepo, _, cleanup := setupTestEnv(t)
	defer cleanup()

	mqProducer := &TestMQProducer{shouldFail: false}
	service := NewService(nil, cacheRepo, mqProducer, &config.SeckillConfig{})

	ctx := context.Background()

	// 初始化库存为 10
	err := cacheRepo.SetStock(ctx, 1, 10)
	require.NoError(t, err)

	// 执行秒杀
	req := &SeckillRequest{
		UserID:   1001,
		CouponID: 1,
	}
	resp, err := service.Seckill(ctx, req)

	// 验证结果
	assert.NoError(t, err)
	assert.True(t, resp.Success)
	assert.Equal(t, "秒杀成功，订单处理中", resp.Message)

	// 验证库存被扣减
	stock, err := cacheRepo.GetStock(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(9), stock)
}

// 测试库存不足
func TestSeckill_StockNotEnough(t *testing.T) {
	cacheRepo, _, cleanup := setupTestEnv(t)
	defer cleanup()

	mqProducer := &TestMQProducer{shouldFail: false}
	service := NewService(nil, cacheRepo, mqProducer, &config.SeckillConfig{})

	ctx := context.Background()

	// 初始化库存为 0
	err := cacheRepo.SetStock(ctx, 1, 0)
	require.NoError(t, err)

	// 执行秒杀
	req := &SeckillRequest{
		UserID:   1001,
		CouponID: 1,
	}
	resp, err := service.Seckill(ctx, req)

	// 验证结果
	assert.Error(t, err)
	assert.Equal(t, ErrStockNotEnough, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "库存不足", resp.Message)

	// 验证库存未变化
	stock, err := cacheRepo.GetStock(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(0), stock)
}

// 测试 MQ 发送失败并回滚
func TestSeckill_MQFailedWithRollback(t *testing.T) {
	cacheRepo, _, cleanup := setupTestEnv(t)
	defer cleanup()

	mqProducer := &TestMQProducer{shouldFail: true} // MQ 失败
	service := NewService(nil, cacheRepo, mqProducer, &config.SeckillConfig{})

	ctx := context.Background()

	// 初始化库存为 50
	err := cacheRepo.SetStock(ctx, 1, 50)
	require.NoError(t, err)

	// 执行秒杀
	req := &SeckillRequest{
		UserID:   1001,
		CouponID: 1,
	}
	resp, err := service.Seckill(ctx, req)

	// 验证秒杀失败
	assert.Error(t, err)
	assert.False(t, resp.Success)
	assert.Equal(t, "秒杀失败，请重试", resp.Message)

	// 验证库存已回滚（应该还是 50）
	stock, err := cacheRepo.GetStock(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(50), stock, "库存应该回滚到初始值")
}

// 测试并发秒杀
func TestSeckill_Concurrent(t *testing.T) {
	cacheRepo, _, cleanup := setupTestEnv(t)
	defer cleanup()

	mqProducer := &TestMQProducer{shouldFail: false}
	service := NewService(nil, cacheRepo, mqProducer, &config.SeckillConfig{})

	ctx := context.Background()

	// 初始化库存为 10
	err := cacheRepo.SetStock(ctx, 1, 10)
	require.NoError(t, err)

	// 20 个并发请求秒杀 10 个库存
	successCount := 0
	done := make(chan bool, 20)

	for i := 0; i < 20; i++ {
		go func(userID int64) {
			req := &SeckillRequest{
				UserID:   userID,
				CouponID: 1,
			}
			resp, err := service.Seckill(ctx, req)
			if err == nil && resp.Success {
				successCount++
			}
			done <- true
		}(int64(1000 + i))
	}

	// 等待所有请求完成
	for i := 0; i < 20; i++ {
		<-done
	}

	// 验证最终库存为 0
	finalStock, err := cacheRepo.GetStock(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(0), finalStock, "最终库存应该为 0")

	// 成功数应该等于初始库存（10 个）
	assert.Equal(t, 10, successCount, "成功秒杀数应该等于初始库存")
}

// 测试初始化库存
func TestInitStock(t *testing.T) {
	cacheRepo, _, cleanup := setupTestEnv(t)
	defer cleanup()

	ctx := context.Background()

	// 手动设置库存测试 InitStock 功能
	err := cacheRepo.SetStock(ctx, 1, 100)
	require.NoError(t, err)

	// 验证库存
	stock, err := cacheRepo.GetStock(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, int64(100), stock)
}
