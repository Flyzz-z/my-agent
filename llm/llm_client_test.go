package llm

import (
	"context"
	"rag-agent/config"
	"testing"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.LoadConfig(config.DefaultConfigPath)
}

/*
TestLLMClient_NewLLMClient 测试LLM客户端的创建
*/
func TestLLMClient(t *testing.T) {
	// 创建上下文
	ctx := context.Background()

	// 创建LLM客户端
	llmClient, err := NewLLMClient(ctx)
	assert.NoError(t, err, "Should not return error when creating LLM client")

	promptTempalte := prompt.FromMessages(schema.FString, []schema.MessagesTemplate{
		schema.SystemMessage("你是一个{role}，你的任务是回答用户的问题"),
		schema.UserMessage("问题: {content}"),
	}...)

	messages, err := promptTempalte.Format(ctx, map[string]any{
		"role":    "程序员鼓励师",
		"content": "我的代码一直报错，感觉好沮丧，该怎么办?",
	})
	assert.NoError(t, err, "Should not return error when formatting messages")

	// 测试Generate方法
	resp, err := llmClient.ChatModel.Generate(ctx, messages)
	assert.NoError(t, err, "Should not return error when generating response")
	assert.NotEmpty(t, resp, "Response should not be empty")
	if resp != nil {
		t.Logf("Response: %s", resp)
	}
}
