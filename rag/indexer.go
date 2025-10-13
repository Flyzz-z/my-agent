package rag

import (
	"context"
	"github.com/cloudwego/eino-ext/components/embedding/ark"
	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisIndexer(ctx context.Context, client *redis.Client, embedder *ark.Embedder, prefix string) (*redisInd.Indexer, error) {
	indexer, err := redisInd.NewIndexer(ctx, &redisInd.IndexerConfig{
		Client:    client,
		Embedding: embedder,
		KeyPrefix: prefix,
	})
	if err != nil {
		log.Printf("NewRedisIndexer err: %v", err)
		return nil, err
	}
	return indexer, nil
}

/*
InitVectorIndex 初始化向量索引
*/
func InitVectorIndex(ctx context.Context, client *redis.Client, indexName, prefix string, dimension int64) error {
	if _, err := client.Do(ctx, "FT.INFO", indexName).Result(); err == nil {
		return nil
	}

	indexArgs := []interface{}{
		"FT.CREATE", indexName,
		"ON", "HASH",
		"PREFIX", "1", prefix,
		"SCHEMA",
		"vector_content", "VECTOR", "FLAT", "6",
		"TYPE", "FLOAT32",
		"DIM", dimension,
		"DISTANCE_METRIC", "COSINE",
		"content", "TEXT",
	}
	
	if _, err := client.Do(ctx, indexArgs...).Result(); err != nil {
		log.Printf("InitVectorIndex err: %v", err)
		return err
	}
	return nil
}
