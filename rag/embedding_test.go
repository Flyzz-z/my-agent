package rag

import (
	"context"
	"log"
	"rag-agent/config"
	"testing"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 加载配置
	err := config.LoadConfig(config.DefaultConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

// TestEmbedder 测试NewArkEmbedder函数
func TestEmbedder(t *testing.T) {
	ctx := context.Background()
	
	// 创建embedder
	embedder, err := NewArkEmbedder(ctx)
	
	// 由于实际环境中可能无法连接到Ark API，这里只测试错误处理逻辑
	// 在实际可连接的环境中，应该测试embedder是否成功创建
	if err != nil {
		t.Fatalf("Expected embedder creation might fail in test environment: %v", err)
	} else {
		// 如果成功创建，验证embedder不为nil
		assert.NotNil(t, embedder, "Embedder should not be nil")
	}

	// 测试EmbedStrings方法
	_, err = embedder.EmbedStrings(ctx, []string{"测试文本"})
	if err != nil {
		t.Fatalf("Expected embedding might fail in test environment: %v", err)
	}
}