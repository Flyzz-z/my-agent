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
	
	if err != nil {
		t.Fatalf("Embedder creation failed: %v", err)
	} else {
		// 如果成功创建，验证embedder不为nil
		assert.NotNil(t, embedder, "Embedder should not be nil")
	}

	// 测试EmbedStrings方法
	_, err = embedder.EmbedStrings(ctx, []string{"测试文本"})
	if err != nil {
		t.Fatalf("EmbedStrings failed: %v", err)
	}
}