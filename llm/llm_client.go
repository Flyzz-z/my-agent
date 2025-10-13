package llm

import (
	"context"
	"fmt"
	"rag-agent/config"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/eino-contrib/ollama/api"
)

// LLMClient 与大语言模型交互的客户端

type LLMClient struct {
	// 使用Eino的组件来与大语言模型交互
	chatModel model.ToolCallingChatModel
}

// NewLLMClient 创建一个新的LLMClient实例
func NewLLMClient(ctx context.Context) (*LLMClient, error) {
	// 初始化Eino的ChatModel组件
	// 这里使用Eino提供的MockChatModel作为示例
	// 在实际应用中，应该替换为真实的大语言模型客户端

	cfg := config.GetConfig()
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.LLM.BaseURL,
		Timeout: cfg.LLM.Timeout,
		Model: cfg.LLM.Model,

		    // 模型参数
    Options: &api.Options{
       Runner: api.Runner{
          NumCtx:    4096, // 上下文窗口大小
          NumGPU:    1,    // GPU 数量
          NumThread: 4,    // CPU 线程数
       },
       Temperature:   0.7,        // 温度
       TopP:          0.9,        // Top-P 采样
       TopK:          40,         // Top-K 采样
       Seed:          42,         // 随机种子
       NumPredict:   200,        // 最大生成长度
       Stop:          []string{}, // 停止词
       RepeatPenalty: 1.1,        // 重复惩罚
    },
	})


	if err != nil {
		return nil, fmt.Errorf("创建ollama模型失败: %v", err)
	}

	return &LLMClient{
		chatModel: chatModel,
	}, nil
}

