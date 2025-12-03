package rag

import (
	"context"
	"log"
	"rag-agent/config"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)


func NewArkEmbedder(ctx context.Context) (*ark.Embedder, error) {
	cfg := config.GetConfig()
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		Model: cfg.Embedding.Model,
		APIKey: cfg.Embedding.APIKey,
	})

	if err != nil {
		log.Printf("init ArkEmbedder failed: %v", err)
		return nil, err
	}
	return embedder, nil
}
