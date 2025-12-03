package handler

import (
	"net/http"
	"rag-agent/internal/domain/aisearch"

	"github.com/gin-gonic/gin"
)

// AISearchHandler AI搜索处理器 - 整合了LLM和RAG能力
type AISearchHandler struct {
	service *aisearch.Service
}

// NewAISearchHandler 创建AI搜索处理器
func NewAISearchHandler(service *aisearch.Service) *AISearchHandler {
	return &AISearchHandler{
		service: service,
	}
}

// Search AI智能搜索接口
func (h *AISearchHandler) Search(c *gin.Context) {
	var req aisearch.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Search(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddDocument 添加文档到RAG索引
func (h *AISearchHandler) AddDocument(c *gin.Context) {
	var req aisearch.AddDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddDocument(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, aisearch.AddDocumentResponse{
		Success: true,
		Message: "文档添加成功",
	})
}

// SearchStream 流式搜索接口
func (h *AISearchHandler) SearchStream(c *gin.Context) {
	var req aisearch.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置SSE响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// TODO: 实现流式响应
	c.JSON(http.StatusOK, gin.H{"message": "流式搜索待实现"})
}
