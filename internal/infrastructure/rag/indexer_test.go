package rag

import (
	"context"
	"rag-agent/config"
	"testing"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化配置
	config.LoadConfig(config.DefaultConfigPath)
}

func TestRealInitVectorIndex(t *testing.T) {
	ctx := context.Background()
	
	cfg := config.GetConfig()

	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Protocol: cfg.Redis.Protocol,
		UnstableResp3: cfg.Redis.UnstableResp3,
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


