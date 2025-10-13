package rag

import (
	"context"
	"fmt"
	"log"

	"rag-agent/config"
	"rag-agent/llm"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	redisInd "github.com/cloudwego/eino-ext/components/indexer/redis"
	redisRet "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/document"
	"github.com/redis/go-redis/v9"
)

/*
	RAGEngine 管理RAG相关的功能
*/
type RAGEngine struct {
	IndexName string
	Prefix    string
	Dimension int64

	redis *redis.Client
	FileLoader *file.FileLoader
	Splitter   document.Transformer
	Retriever  *redisRet.Retriever
	Indexer    *redisInd.Indexer
	LLM  *llm.LLMClient
}

func NewRAGEngine(ctx context.Context, indexName, prefix string) (*RAGEngine, error) {

	// 初始化redis
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
	})
	// 初始化fileLoader
	fileLoader,err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true,
		Parser: nil,
	})
	if err != nil {
		log.Printf("new engine failed, file loader failed: %v", err)
		return nil, fmt.Errorf("file loader failed: %v", err)
	}

	// 初始化embedder
	embedder, err := NewArkEmbedder(ctx)
	if err != nil {
		log.Printf("new engine failed, embedder failed: %v", err)
		return nil, fmt.Errorf("embedder failed: %v", err)
	}

	// 初始化splitter
	splitter, err := NewMarkdownSplitter(ctx)
	if err != nil {
		log.Printf("new engine failed, splitter failed: %v", err)
		return nil, fmt.Errorf("splitter failed: %v", err)
	}

	// 初始化retriever
	retriever, err := NewRedisRetriever(ctx, redisCli, embedder, indexName)
	if err != nil {
		log.Printf("new engine failed, retriever failed: %v", err)
		return nil, fmt.Errorf("retriever failed: %v", err)
	}

	// 初始化indexer
	indexer, err := NewRedisIndexer(ctx, redisCli, embedder, prefix)
	if err != nil {
		log.Printf("new engine failed, indexer failed: %v", err)
		return nil, fmt.Errorf("indexer failed: %v", err)
	}

	// 初始化llm
	llm, err := llm.NewLLMClient(ctx)
	if err != nil {
		log.Printf("new engine failed, llm failed: %v", err)
		return nil, fmt.Errorf("llm failed: %v", err)
	}

	cfg := config.GetConfig()

	return &RAGEngine{
		IndexName: cfg.RAG.IndexName,
		Prefix:    cfg.RAG.Prefix,
		Dimension: cfg.RAG.Dimension,

		redis:      redisCli,
		FileLoader: fileLoader,
		Splitter:   splitter,
		Retriever:  retriever,
		Indexer:    indexer,
		LLM:        llm,
	}, nil
}