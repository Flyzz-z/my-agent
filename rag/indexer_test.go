package rag

import (
	"context"
	"rag-agent/config"
	"testing"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化配置
	config.LoadDefaultConfig()
}

// TestInitVectorIndex 测试InitVectorIndex函数
func TestInitVectorIndex(t *testing.T) {
	ctx := context.Background()
	
	// 创建Redis客户端和mock
	rdb, mock := redismock.NewClientMock()
	
	// 模拟FT.INFO命令返回错误（表示索引不存在）
	mock.ExpectDo("FT.INFO", "test_index").SetErr(redis.Nil)
	
	// 模拟FT.CREATE命令成功
	mock.ExpectDo("FT.CREATE", "test_index", "ON", "HASH", "TYPE", "FLOAT32", "DIM", int64(768), "DISTANCE_METRIC", "COSINE", "PREFIX", "1", "test_prefix", "SCHEMA", "vector_content", "VECTOR", "FLAT", "6", "TYPE", "FLOAT32", "VEC_TYPE", "FLOAT32", "DIM", int64(768), "content", "TEXT").SetVal("OK")
	
	// 测试索引初始化
	err := InitVectorIndex(ctx, rdb, "test_index", "test_prefix", 768)
	
	// 验证没有错误
	assert.NoError(t, err, "Should not return error when initializing index")
	
	// 验证mock期望是否都被满足
	assert.NoError(t, mock.ExpectationsWereMet(), "All mock expectations should be met")
}

// TestInitVectorIndex_ExistingIndex 测试索引已存在的情况
func TestInitVectorIndex_ExistingIndex(t *testing.T) {
	ctx := context.Background()
	
	// 创建Redis客户端和mock
	rdb, mock := redismock.NewClientMock()
	
	// 模拟FT.INFO命令成功（表示索引已存在）
	mock.ExpectDo("FT.INFO", "test_index").SetVal("Index exists")
	
	// 测试索引初始化
	err := InitVectorIndex(ctx, rdb, "test_index", "test_prefix", 768)
	
	// 验证没有错误
	assert.NoError(t, err, "Should not return error when index already exists")
	
	// 验证mock期望是否都被满足
	assert.NoError(t, mock.ExpectationsWereMet(), "All mock expectations should be met")
}

func TestRealInitVectorIndex(t *testing.T) {
	ctx := context.Background()
	
	cfg := config.GetConfig()

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	
	// 测试索引初始化
	err := InitVectorIndex(ctx, rdb, "test_index", "test_prefix", 768)
	
	// 验证没有错误
	assert.NoError(t, err, "Should not return error when initializing index")
	
	// 验证索引已存在
	info, err := rdb.Do(ctx, "FT.INFO", "test_index").Result()
	assert.NoError(t, err, "Should not return error when checking index existence")
	assert.NotEmpty(t, info, "Index should exist after initialization")
}


