package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rag-agent/config"
	"rag-agent/internal/domain/aisearch"
	"rag-agent/internal/domain/seckill"
	httpserver "rag-agent/internal/server/http"
	"rag-agent/internal/server/http/handler"
	"rag-agent/pkg/llm"

	"rag-agent/internal/infrastructure/rag"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("./config.yaml"); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}
	cfg := config.GetConfig()

	ctx := context.Background()

	// 初始化RAG引擎
	ragEngine, err := rag.NewRAGEngine(ctx, cfg.RAG.IndexName, cfg.RAG.Prefix)
	if err != nil {
		log.Fatalf("初始化RAG引擎失败: %v", err)
	}

	// 初始化LLM客户端
	llmClient, err := llm.NewLLMClient(ctx)
	if err != nil {
		log.Fatalf("初始化LLM客户端失败: %v", err)
	}

	// 构建AI搜索 Graph - 整合了LLM和RAG能力
	graph, err := aisearch.BuildGraph(ragEngine, llmClient)
	if err != nil {
		log.Fatalf("构建AI搜索Graph失败: %v", err)
	}

	// 初始化服务
	// AI搜索服务 - 三大主要功能之一（整合了LLM和RAG）
	aisearchService := aisearch.NewService(graph, ragEngine, llmClient)

	// 秒杀服务 - 三大主要功能之二
	seckillService := seckill.NewService(nil, nil, nil, &cfg.Seckill) 

	// 初始化处理器
	aisearchHandler := handler.NewAISearchHandler(aisearchService)
	seckillHandler := handler.NewSeckillHandler(seckillService)

	// 设置路由
	router := httpserver.NewRouter(seckillHandler, aisearchHandler)
	engine := router.Setup()

	// 启动HTTP服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	log.Printf("服务器启动成功，监听地址: %s", addr)
	log.Println("三大主要功能:")
	log.Println("  1. AI搜索 (整合了LLM和RAG能力) - /api/v1/aisearch")
	log.Println("  2. 秒杀系统 - /api/v1/seckill")
	log.Println("  3. [待实现的第三个功能]")

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭:", err)
	}

	log.Println("服务器已关闭")
}
