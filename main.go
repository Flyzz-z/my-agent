package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"rag-agent/config"
	"rag-agent/rag"
)

func main() {
	// 创建上下文
	// ctx := context.Background()

	// 加载配置
	err := config.LoadConfig(config.DefaultConfigPath)
	if err != nil {
		log.Printf("加载配置文件失败，使用默认配置: %v", err)
		config.LoadDefaultConfig()
	}
	cfg := config.GetConfig()

	ctx := context.Background()
	// 初始化RAG引擎
	ragEngine, err := rag.NewRAGEngine(ctx, cfg.RAG.IndexName, cfg.RAG.Prefix)
	if err != nil {
		log.Fatalf("初始化RAG引擎失败: %v", err)
	}

	// 添加文件到索引
	err = ragEngine.AddFile(ctx, "/home/flyzz/agent/doc/kafka.md")
	if err != nil {
		log.Fatalf("添加文件到索引失败: %v", err)
	}

	stream, err := ragEngine.QueryStream(ctx, "Kafka消息队列如何保证消息不丢失")
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}
	defer stream.Close()

	// 处理流式响应
	for {
		msg, err := stream.Recv()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Printf("接收流式响应失败: %v", err)
			break
		}

		fmt.Print(msg.ReasoningContent)
	}
}