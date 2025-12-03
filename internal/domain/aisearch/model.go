package aisearch

// SearchRequest AI搜索请求 - 整合了RAG和LLM能力
type SearchRequest struct {
	Query   string `json:"query" binding:"required"` // 搜索查询
	Session string `json:"session"`                  // 会话ID
}

// SearchResponse AI搜索响应
type SearchResponse struct {
	Answer    string   `json:"answer"`               // AI生成的答案
	Query     string   `json:"query"`                // 原始查询
	Documents []string `json:"documents,omitempty"`  // 引用的文档片段
	Session   string   `json:"session"`              // 会话ID
}

// AddDocumentRequest 添加文档请求
type AddDocumentRequest struct {
	FilePath string `json:"file_path" binding:"required"` // 文件路径
}

// AddDocumentResponse 添加文档响应
type AddDocumentResponse struct {
	Success bool   `json:"success"` // 是否成功
	Message string `json:"message"` // 响应消息
}
