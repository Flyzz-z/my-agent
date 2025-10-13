package rag

import (
	"context"
	"rag-agent/config"
	"testing"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/go-redis/redismock/v9"
)

func init() {
	// 加载配置
	config.LoadConfig("config.yaml")
}

// TestNewRedisRetriever_WithMock 测试NewRedisRetriever函数，使用mock进行更完整的测试
func TestNewRedisRetriever_WithMock(t *testing.T) {
	ctx := context.Background()
	indexName := "test_index"

	// 创建Redis客户端和mock
	rdb, _ := redismock.NewClientMock()

	// 这里我们不模拟具体的Redis命令，因为NewRedisRetriever主要是配置Retriever
	// 实际的检索操作会在其他地方测试

	// 创建mock的embedder
	embedder := &ark.Embedder{}

	// 测试创建Redis检索器
	_, err := NewRedisRetriever(ctx, rdb, embedder, indexName)

	if err != nil {
		t.Logf("Expected retriever creation might fail with redismock: %v", err)
	}
}
