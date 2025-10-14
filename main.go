package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"rag-agent/config"
	"rag-agent/llm"
	"rag-agent/rag"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var systemPrompt = `
# 角色: 你是一个专业的问答助手
# 任务: 根据用户的问题,有需要时使用提供的文档内容，生成一个准确的回答
- 提供帮助时：
  • 表达清晰简洁
  • 相关时提供实际示例
  • 有帮助时引用文档
  • 适用时提出改进建议或下一步操作

这里是文档内容：
---- 文档开始 -----
	{documents}
---- 文档结束 ----
`

func main() {

	// 加载配置
	err := config.LoadConfig(config.DefaultConfigPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
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

	// 构建graph
	graph := compose.NewGraph[string,*schema.Message]()

	// 添加Retriever节点
	graph.AddRetrieverNode("retriever", ragEngine.Retriever)

	graph.AddLambdaNode("format_docs", compose.InvokableLambda(func(ctx context.Context, docs []*schema.Document) (map[string]any, error) {
		return map[string]any{
			"documents": docs,
			"content":   ctx.Value("query").(string),
		}, nil
	}))

	// 添加ChatTemplate节点
	graph.AddChatTemplateNode("chat_template", prompt.FromMessages(schema.FString,
		schema.SystemMessage(
			systemPrompt,
		),
		schema.UserMessage(
			"{content}",
		),
	))

	llmClient, err := llm.NewLLMClient(ctx)
	if err != nil {
		log.Fatalf("初始化LLM客户端失败: %v", err)
	}
	// 添加ChatModel节点
	graph.AddChatModelNode("chat_model", llmClient.ChatModel)


	//todo 添加工具节点，当前不知道添加什么

	graph.AddEdge(compose.START, "retriever")
	graph.AddEdge("retriever", "format_docs")
	graph.AddEdge("format_docs", "chat_template")
	graph.AddEdge("chat_template", "chat_model")
	graph.AddEdge("chat_model", compose.END)

	runner, err := graph.Compile(ctx)
	if err != nil {
		log.Fatalf("编译graph失败: %v", err)
	}

	// 运行graph
	ctx = context.WithValue(ctx, "query", "Kafka如何阻止重复消费")
	s, err := runner.Stream(ctx,
		ctx.Value("query").(string),
	)
	if err != nil {
		log.Fatalf("运行graph失败: %v", err)
	}
	for {
		msg, err := s.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalf("接收消息失败: %v", err)
		}
		fmt.Printf("%s", msg.Content)
	}
}
