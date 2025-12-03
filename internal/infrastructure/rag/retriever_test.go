package rag

import (
	"context"
	"rag-agent/config"
	"testing"

	"github.com/redis/go-redis/v9"
)

func init() {
	// 加载配置
	config.LoadConfig(config.DefaultConfigPath)
}


func TestRedisRetriever(t *testing.T) {
	ctx := context.Background()
	cfg := config.GetConfig()
	indexName := cfg.RAG.IndexName

	// 创建Redis客户端和mock
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
		Protocol: cfg.Redis.Protocol,
		UnstableResp3: cfg.Redis.UnstableResp3,
	})

	embedder, err := NewArkEmbedder(ctx)
	if err != nil {
		t.Fatalf("ArkEmbedder creation failed: %v", err)
	}

	// 测试创建Redis检索器
	retriever, err := NewRedisRetriever(ctx, rdb, embedder, indexName)

	if err != nil {
		t.Fatalf("Expected retriever creation might fail : %v", err)
	}

	docs, err := retriever.Retrieve(ctx, "kafka的消息丢失问题")
	if err != nil {
		t.Fatalf("Expected retriever retrieval might fail : %v", err)
	}

	if len(docs) > 0 {
		t.Logf("Expected retriever retrieval might return documents: %v", docs)
	}
}
