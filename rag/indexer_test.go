package rag

import (
	"context"
	"testing"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// TestNewRedisIndexer 测试NewRedisIndexer函数
func TestNewRedisIndexer(t *testing.T) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // 实际测试环境可能需要修改
	})
	prefix := "test_prefix"
	
	// 由于embedder依赖外部服务，这里我们主要测试参数验证和错误处理
	// 在实际可连接的环境中，应该使用真实的embedder进行测试
	// 这里使用nil来测试错误处理
	indexer, err := NewRedisIndexer(ctx, client, nil, prefix)
	
	// 由于我们传入了nil的embedder，应该返回错误
	assert.Error(t, err, "Should return error when embedder is nil")
	assert.Nil(t, indexer, "Indexer should be nil when error occurs")
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