package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"rag-agent/config"
	"rag-agent/internal/domain/search"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Client Elasticsearch客户端实现
type Client struct {
	client *elasticsearch.Client
	config *config.ElasticsearchConfig
}

// NewClient 创建新的Elasticsearch客户端
func NewClient(ctx context.Context, cfg *config.ElasticsearchConfig) (*Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("创建ES客户端失败: %w", err)
	}

	// 测试连接
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("连接ES失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES返回错误: %s", res.String())
	}

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// Search 执行搜索
func (c *Client) Search(ctx context.Context, req *search.SearchRequest) (*search.SearchResponse, error) {
	// 设置默认值
	if req.Size == 0 {
		req.Size = 10
	}
	if req.SortBy == "" {
		req.SortBy = "_score"
	}

	// 构建查询DSL
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  req.Query,
							"fields": []string{"title^2", "content", "category", "tags"},
							"type":   "best_fields",
						},
					},
				},
			},
		},
		"from": req.From,
		"size": req.Size,
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title":   map[string]interface{}{},
				"content": map[string]interface{}{},
			},
			"pre_tags":  []string{"<em>"},
			"post_tags": []string{"</em>"},
		},
	}

	// 添加过滤条件
	filters := []map[string]interface{}{}
	if req.Category != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"category": req.Category,
			},
		})
	}
	if len(req.Tags) > 0 {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"tags": req.Tags,
			},
		})
	}
	if len(filters) > 0 {
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = filters
	}

	// 添加排序
	sortOrder := "desc"
	if !req.SortDesc {
		sortOrder = "asc"
	}
	query["sort"] = []map[string]interface{}{
		{
			req.SortBy: map[string]interface{}{
				"order": sortOrder,
			},
		},
	}

	// 序列化查询
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("序列化查询失败: %w", err)
	}

	// 确定索引
	indices := []string{}
	if req.Index != "" {
		indices = append(indices, req.Index)
	} else if c.config.Index != "" {
		indices = append(indices, c.config.Index)
	}

	// 执行搜索
	start := time.Now()
	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(indices...),
		c.client.Search.WithBody(strings.NewReader(buf.String())),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("执行搜索失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES返回错误: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	took := time.Since(start).Milliseconds()

	// 构建响应
	response := &search.SearchResponse{
		Took:  took,
		Query: req.Query,
		Hits:  []search.SearchHit{},
	}

	// 解析结果
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				response.Total = int64(value)
			}
		}
		if maxScore, ok := hits["max_score"].(float64); ok {
			response.MaxScore = maxScore
		}
		if hitsArray, ok := hits["hits"].([]interface{}); ok {
			for _, hit := range hitsArray {
				hitMap := hit.(map[string]interface{})
				searchHit := search.SearchHit{}

				if score, ok := hitMap["_score"].(float64); ok {
					searchHit.Score = score
				}

				if source, ok := hitMap["_source"].(map[string]interface{}); ok {
					doc := c.parseDocument(source)
					if id, ok := hitMap["_id"].(string); ok {
						doc.ID = id
					}
					searchHit.Document = doc
				}

				if highlight, ok := hitMap["highlight"].(map[string]interface{}); ok {
					searchHit.Highlight = make(map[string][]string)
					for field, fragments := range highlight {
						if frags, ok := fragments.([]interface{}); ok {
							strFrags := make([]string, len(frags))
							for i, frag := range frags {
								strFrags[i] = frag.(string)
							}
							searchHit.Highlight[field] = strFrags
						}
					}
				}

				response.Hits = append(response.Hits, searchHit)
			}
		}
	}

	return response, nil
}

// IndexDocument 索引单个文档
func (c *Client) IndexDocument(ctx context.Context, req *search.IndexDocumentRequest) (*search.IndexDocumentResponse, error) {
	// 序列化文档
	doc := req.Document
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = time.Now()
	}
	doc.UpdatedAt = time.Now()

	body, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("序列化文档失败: %w", err)
	}

	// 执行索引
	var res *esapi.Response
	if req.ID != "" {
		res, err = c.client.Index(
			req.Index,
			strings.NewReader(string(body)),
			c.client.Index.WithContext(ctx),
			c.client.Index.WithDocumentID(req.ID),
		)
	} else {
		res, err = c.client.Index(
			req.Index,
			strings.NewReader(string(body)),
			c.client.Index.WithContext(ctx),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("索引文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES返回错误: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	response := &search.IndexDocumentResponse{
		Index: req.Index,
	}

	if id, ok := result["_id"].(string); ok {
		response.ID = id
	}
	if resultStr, ok := result["result"].(string); ok {
		response.Result = resultStr
	}
	if version, ok := result["_version"].(float64); ok {
		response.Version = int64(version)
	}

	return response, nil
}

// BulkIndex 批量索引文档
func (c *Client) BulkIndex(ctx context.Context, req *search.BulkIndexRequest) (*search.BulkIndexResponse, error) {
	start := time.Now()

	var buf strings.Builder
	for _, doc := range req.Documents {
		if doc.CreatedAt.IsZero() {
			doc.CreatedAt = time.Now()
		}
		doc.UpdatedAt = time.Now()

		// 写入action
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": req.Index,
			},
		}
		if doc.ID != "" {
			action["index"].(map[string]interface{})["_id"] = doc.ID
		}
		if err := json.NewEncoder(&buf).Encode(action); err != nil {
			return nil, fmt.Errorf("序列化action失败: %w", err)
		}

		// 写入文档
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return nil, fmt.Errorf("序列化文档失败: %w", err)
		}
	}

	res, err := c.client.Bulk(
		strings.NewReader(buf.String()),
		c.client.Bulk.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("批量索引失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES返回错误: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	response := &search.BulkIndexResponse{
		Took:    time.Since(start).Milliseconds(),
		Total:   len(req.Documents),
		Succeed: 0,
		Failed:  0,
	}

	if errors, ok := result["errors"].(bool); ok {
		response.Errors = errors
	}

	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			itemMap := item.(map[string]interface{})
			if indexResult, ok := itemMap["index"].(map[string]interface{}); ok {
				if status, ok := indexResult["status"].(float64); ok {
					if status >= 200 && status < 300 {
						response.Succeed++
					} else {
						response.Failed++
					}
				}
			}
		}
	}

	return response, nil
}

// GetDocument 获取文档
func (c *Client) GetDocument(ctx context.Context, req *search.GetDocumentRequest) (*search.GetDocumentResponse, error) {
	res, err := c.client.Get(
		req.Index,
		req.ID,
		c.client.Get.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("获取文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return &search.GetDocumentResponse{
			Found: false,
		}, nil
	}

	if res.IsError() {
		return nil, fmt.Errorf("ES返回错误: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	response := &search.GetDocumentResponse{
		Found: false,
	}

	if found, ok := result["found"].(bool); ok && found {
		response.Found = true
		if source, ok := result["_source"].(map[string]interface{}); ok {
			doc := c.parseDocument(source)
			if id, ok := result["_id"].(string); ok {
				doc.ID = id
			}
			response.Document = doc
		}
	}

	return response, nil
}

// DeleteDocument 删除文档
func (c *Client) DeleteDocument(ctx context.Context, req *search.DeleteDocumentRequest) (*search.DeleteDocumentResponse, error) {
	res, err := c.client.Delete(
		req.Index,
		req.ID,
		c.client.Delete.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("删除文档失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES返回错误: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	response := &search.DeleteDocumentResponse{
		Index: req.Index,
		ID:    req.ID,
	}

	if resultStr, ok := result["result"].(string); ok {
		response.Result = resultStr
	}

	return response, nil
}

// CreateIndex 创建索引
func (c *Client) CreateIndex(ctx context.Context, index string, mapping map[string]interface{}) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("序列化mapping失败: %w", err)
	}

	res, err := c.client.Indices.Create(
		index,
		c.client.Indices.Create.WithContext(ctx),
		c.client.Indices.Create.WithBody(strings.NewReader(string(body))),
	)
	if err != nil {
		return fmt.Errorf("创建索引失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ES返回错误: %s", res.String())
	}

	return nil
}

// DeleteIndex 删除索引
func (c *Client) DeleteIndex(ctx context.Context, index string) error {
	res, err := c.client.Indices.Delete(
		[]string{index},
		c.client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("删除索引失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ES返回错误: %s", res.String())
	}

	return nil
}

// IndexExists 检查索引是否存在
func (c *Client) IndexExists(ctx context.Context, index string) (bool, error) {
	res, err := c.client.Indices.Exists(
		[]string{index},
		c.client.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, fmt.Errorf("检查索引失败: %w", err)
	}
	defer res.Body.Close()

	return res.StatusCode == 200, nil
}

// parseDocument 解析ES文档为Domain模型
func (c *Client) parseDocument(source map[string]interface{}) search.Document {
	doc := search.Document{}

	if id, ok := source["id"].(string); ok {
		doc.ID = id
	}
	if title, ok := source["title"].(string); ok {
		doc.Title = title
	}
	if content, ok := source["content"].(string); ok {
		doc.Content = content
	}
	if category, ok := source["category"].(string); ok {
		doc.Category = category
	}
	if tags, ok := source["tags"].([]interface{}); ok {
		doc.Tags = make([]string, len(tags))
		for i, tag := range tags {
			doc.Tags[i] = tag.(string)
		}
	}
	if author, ok := source["author"].(string); ok {
		doc.Author = author
	}
	if createdAt, ok := source["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, createdAt); err == nil {
			doc.CreatedAt = t
		}
	}
	if updatedAt, ok := source["updated_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339, updatedAt); err == nil {
			doc.UpdatedAt = t
		}
	}
	if metadata, ok := source["metadata"].(map[string]interface{}); ok {
		doc.Metadata = metadata
	}

	return doc
}
