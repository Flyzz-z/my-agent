package aisearch

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

var systemPrompt = `
# 角色: 你是一个专业的AI搜索助手
# 任务: 根据用户的问题,结合RAG检索的文档内容和LLM能力，生成一个准确的回答
- 提供帮助时：
  • 表达清晰简洁
  • 相关时提供实际示例
  • 有帮助时引用文档
  • 适用时提出改进建议或下一步操作

这里是检索到的文档内容：
---- 文档开始 -----
	{documents}
---- 文档结束 ----
`

// Service AI搜索服务 - 整合了LLM和RAG能力
type Service struct {
	graph     GraphRunner
	ragEngine RAGEngine
	llmClient ChatModel
}

// GraphRunner graph运行器接口
type GraphRunner interface {
	Stream(ctx context.Context, input string) (StreamReader, error)
}

// StreamReader 流式读取器接口
type StreamReader interface {
	Recv() (*schema.Message, error)
}

// NewService 创建AI搜索服务
func NewService(graph GraphRunner, ragEngine RAGEngine, llmClient ChatModel) *Service {
	return &Service{
		graph:     graph,
		ragEngine: ragEngine,
		llmClient: llmClient,
	}
}

// Search AI智能搜索 - 处理搜索请求并返回AI增强的结果
func (s *Service) Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
	// 运行graph进行AI搜索
	ctx = context.WithValue(ctx, "query", req.Query)
	stream, err := s.graph.Stream(ctx, req.Query)
	if err != nil {
		return nil, fmt.Errorf("运行AI搜索失败: %w", err)
	}

	// 收集流式响应
	var answer string
	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("接收消息失败: %w", err)
		}
		answer += msg.Content
	}

	return &SearchResponse{
		Answer:  answer,
		Query:   req.Query,
		Session: req.Session,
	}, nil
}

// AddDocument 添加文档到RAG索引
func (s *Service) AddDocument(ctx context.Context, req *AddDocumentRequest) error {
	return s.ragEngine.AddFile(ctx, req.FilePath)
}

// BuildGraph 构建eino graph
func BuildGraph(ragEngine RAGEngine, chatModel ChatModel) (compose.GraphStreamRunnable[string, *schema.Message], error) {
	ctx := context.Background()
	graph := compose.NewGraph[string, *schema.Message]()

	// 添加Retriever节点
	graph.AddRetrieverNode("retriever", ragEngine.GetRetriever())

	// 添加格式化文档节点
	graph.AddLambdaNode("format_docs", compose.InvokableLambda(func(ctx context.Context, docs []*schema.Document) (map[string]any, error) {
		return map[string]any{
			"documents": docs,
			"content":   ctx.Value("query").(string),
		}, nil
	}))

	// 添加ChatTemplate节点
	graph.AddChatTemplateNode("chat_template", prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemPrompt),
		schema.UserMessage("{content}"),
	))

	// 添加ChatModel节点
	graph.AddChatModelNode("chat_model", chatModel.GetModel())

	// 连接节点
	graph.AddEdge(compose.START, "retriever")
	graph.AddEdge("retriever", "format_docs")
	graph.AddEdge("format_docs", "chat_template")
	graph.AddEdge("chat_template", "chat_model")
	graph.AddEdge("chat_model", compose.END)

	// 编译graph
	return graph.Compile(ctx)
}

// RAGEngine RAG引擎接口
type RAGEngine interface {
	GetRetriever() interface{}
	AddFile(ctx context.Context, filePath string) error
}

// ChatModel 聊天模型接口
type ChatModel interface {
	GetModel() interface{}
}
