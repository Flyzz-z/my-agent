package seckill

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// RedisCacheRepository Redis 缓存仓库实现（秒杀专用）
type RedisCacheRepository struct {
	client *redis.Client
	prefix string
}

// NewRedisCacheRepository 创建 Redis 缓存仓库
func NewRedisCacheRepository(client *redis.Client) CacheRepository {
	return &RedisCacheRepository{
		client: client,
		prefix: "seckill:stock:",
	}
}

// getStockKey 获取库存的 Redis key
func (r *RedisCacheRepository) getStockKey(couponID int64) string {
	return fmt.Sprintf("%s%d", r.prefix, couponID)
}

// GetStock 获取缓存中的库存
func (r *RedisCacheRepository) GetStock(ctx context.Context, couponID int64) (int64, error) {
	key := r.getStockKey(couponID)
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("库存不存在")
	}
	if err != nil {
		return 0, fmt.Errorf("获取库存失败: %w", err)
	}

	stock, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("解析库存失败: %w", err)
	}

	return stock, nil
}

// DecrStock 使用 Lua 脚本原子性扣减库存
// 返回扣减后的库存，如果库存不足返回 -1
func (r *RedisCacheRepository) DecrStock(ctx context.Context, couponID int64) (int64, error) {
	key := r.getStockKey(couponID)

	// Lua 脚本：检查库存并原子性扣减
	// 如果库存 > 0，则扣减并返回扣减后的库存
	// 如果库存 <= 0，返回 -1 表示库存不足
	luaScript := `
		local stock = redis.call('GET', KEYS[1])
		if not stock then
			return -1
		end
		stock = tonumber(stock)
		if stock <= 0 then
			return -1
		end
		redis.call('DECR', KEYS[1])
		return stock - 1
	`

	result, err := r.client.Eval(ctx, luaScript, []string{key}).Result()
	if err != nil {
		return 0, fmt.Errorf("执行 Lua 脚本失败: %w", err)
	}

	stock, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("Lua 脚本返回类型错误")
	}

	return stock, nil
}

// SetStock 设置缓存中的库存
func (r *RedisCacheRepository) SetStock(ctx context.Context, couponID int64, stock int64) error {
	key := r.getStockKey(couponID)
	err := r.client.Set(ctx, key, stock, 0).Err()
	if err != nil {
		return fmt.Errorf("设置库存失败: %w", err)
	}
	return nil
}

// IncrStock 原子性增加库存（使用 Redis INCRBY 命令）
func (r *RedisCacheRepository) IncrStock(ctx context.Context, couponID int64, delta int64) error {
	key := r.getStockKey(couponID)
	// INCRBY 是原子操作，线程安全
	err := r.client.IncrBy(ctx, key, delta).Err()
	if err != nil {
		return fmt.Errorf("增加库存失败: %w", err)
	}
	return nil
}
