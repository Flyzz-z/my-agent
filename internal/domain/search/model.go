package search

import "time"

// Document 表示Elasticsearch中的文档
type Document struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	Content   string                 `json:"content"`
	Category  string                 `json:"category,omitempty"`
	Tags      []string               `json:"tags,omitempty"`
	Author    string                 `json:"author,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query    string   `json:"query" binding:"required"`    // 搜索关键词
	Index    string   `json:"index,omitempty"`             // 索引名称，不指定则搜索所有索引
	Category string   `json:"category,omitempty"`          // 分类过滤
	Tags     []string `json:"tags,omitempty"`              // 标签过滤
	From     int      `json:"from,omitempty"`              // 分页起始位置，默认0
	Size     int      `json:"size,omitempty"`              // 返回结果数量，默认10
	SortBy   string   `json:"sort_by,omitempty"`           // 排序字段，如"created_at", "score"
	SortDesc bool     `json:"sort_desc,omitempty"`         // 是否降序，默认true
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Total    int64            `json:"total"`     // 总结果数
	Took     int64            `json:"took"`      // 搜索耗时(ms)
	MaxScore float64          `json:"max_score"` // 最高得分
	Hits     []SearchHit      `json:"hits"`      // 搜索结果
	Query    string           `json:"query"`     // 原始查询
}

// SearchHit 单个搜索结果
type SearchHit struct {
	Score    float64   `json:"score"`     // 相关度得分
	Document Document  `json:"document"`  // 文档内容
	Highlight map[string][]string `json:"highlight,omitempty"` // 高亮片段
}

// IndexDocumentRequest 索引文档请求
type IndexDocumentRequest struct {
	Index    string                 `json:"index" binding:"required"`    // 索引名称
	ID       string                 `json:"id,omitempty"`                // 文档ID，不指定则自动生成
	Document Document               `json:"document" binding:"required"` // 文档内容
}

// IndexDocumentResponse 索引文档响应
type IndexDocumentResponse struct {
	ID      string `json:"id"`      // 文档ID
	Index   string `json:"index"`   // 索引名称
	Result  string `json:"result"`  // 操作结果: created/updated
	Version int64  `json:"version"` // 文档版本
}

// DeleteDocumentRequest 删除文档请求
type DeleteDocumentRequest struct {
	Index string `json:"index" binding:"required"` // 索引名称
	ID    string `json:"id" binding:"required"`    // 文档ID
}

// DeleteDocumentResponse 删除文档响应
type DeleteDocumentResponse struct {
	ID     string `json:"id"`     // 文档ID
	Index  string `json:"index"`  // 索引名称
	Result string `json:"result"` // 操作结果: deleted
}

// BulkIndexRequest 批量索引请求
type BulkIndexRequest struct {
	Index     string     `json:"index" binding:"required"`     // 索引名称
	Documents []Document `json:"documents" binding:"required"` // 文档列表
}

// BulkIndexResponse 批量索引响应
type BulkIndexResponse struct {
	Took     int64  `json:"took"`      // 耗时(ms)
	Errors   bool   `json:"errors"`    // 是否有错误
	Total    int    `json:"total"`     // 总文档数
	Succeed  int    `json:"succeed"`   // 成功数
	Failed   int    `json:"failed"`    // 失败数
}

// GetDocumentRequest 获取文档请求
type GetDocumentRequest struct {
	Index string `json:"index" binding:"required"` // 索引名称
	ID    string `json:"id" binding:"required"`    // 文档ID
}

// GetDocumentResponse 获取文档响应
type GetDocumentResponse struct {
	Found    bool     `json:"found"`    // 是否找到
	Document Document `json:"document"` // 文档内容
}
