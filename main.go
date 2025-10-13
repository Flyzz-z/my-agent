package main

import (
	"fmt"
	"log"
	"rag-agent/config"
	"rag-agent/llm"
	"rag-agent/rag"
)

func main() {
	// 创建上下文
	// ctx := context.Background()

	// 加载配置
	err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Printf("加载配置文件失败，使用默认配置: %v", err)
		config.LoadDefaultConfig()
	}

}

// Agent 结构体定义
type Agent struct {
	ragManager *rag.RAGEngine
	llm        *llm.LLMClient
}


// buildPrompt 构建提示词
func (a *Agent) buildPrompt(question string, docs []string) string {
	prompt := fmt.Sprintf("根据以下信息回答问题:\n\n问题: %s\n\n相关信息:\n", question)
	for i, doc := range docs {
		prompt += fmt.Sprintf("%d. %s\n", i+1, doc)
	}
	prompt += "\n请基于上述信息，准确回答问题，不要添加额外信息。\n"
	return prompt
}