package search

import "context"

// Repository 定义搜索领域的数据访问接口
type Repository interface {
	// Search 执行搜索
	Search(ctx context.Context, req *SearchRequest) (*SearchResponse, error)

	// IndexDocument 索引单个文档
	IndexDocument(ctx context.Context, req *IndexDocumentRequest) (*IndexDocumentResponse, error)

	// BulkIndex 批量索引文档
	BulkIndex(ctx context.Context, req *BulkIndexRequest) (*BulkIndexResponse, error)

	// GetDocument 获取文档
	GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentResponse, error)

	// DeleteDocument 删除文档
	DeleteDocument(ctx context.Context, req *DeleteDocumentRequest) (*DeleteDocumentResponse, error)

	// CreateIndex 创建索引
	CreateIndex(ctx context.Context, index string, mapping map[string]interface{}) error

	// DeleteIndex 删除索引
	DeleteIndex(ctx context.Context, index string) error

	// IndexExists 检查索引是否存在
	IndexExists(ctx context.Context, index string) (bool, error)
}
