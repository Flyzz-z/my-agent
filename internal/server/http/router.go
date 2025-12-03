package http

import (
	"rag-agent/internal/server/http/handler"
	"rag-agent/internal/server/http/middleware"

	"github.com/gin-gonic/gin"
)

// Router 路由配置
type Router struct {
	seckillHandler  *handler.SeckillHandler
	aisearchHandler *handler.AISearchHandler
}

// NewRouter 创建路由
func NewRouter(
	seckillHandler *handler.SeckillHandler,
	aisearchHandler *handler.AISearchHandler,
) *Router {
	return &Router{
		seckillHandler:  seckillHandler,
		aisearchHandler: aisearchHandler,
	}
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	router := gin.New()

	// 使用中间件
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// API版本组
	v1 := router.Group("/api/v1")
	{
		// 秒杀相关路由
		seckill := v1.Group("/seckill")
		{
			seckill.POST("/", r.seckillHandler.Seckill)
			seckill.GET("/coupon/:id", r.seckillHandler.GetCoupon)
			seckill.POST("/init-stock", r.seckillHandler.InitStock)
		}

		// AI搜索相关路由 - 整合了LLM和RAG能力
		aisearch := v1.Group("/aisearch")
		{
			aisearch.POST("/search", r.aisearchHandler.Search)
			aisearch.POST("/search-stream", r.aisearchHandler.SearchStream)
			aisearch.POST("/document", r.aisearchHandler.AddDocument)
		}
	}

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
