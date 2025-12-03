package handler

import (
	"net/http"
	"rag-agent/internal/domain/seckill"

	"github.com/gin-gonic/gin"
)

// SeckillHandler 秒杀处理器
type SeckillHandler struct {
	service *seckill.Service
}

// NewSeckillHandler 创建秒杀处理器
func NewSeckillHandler(service *seckill.Service) *SeckillHandler {
	return &SeckillHandler{
		service: service,
	}
}

// Seckill 秒杀接口
func (h *SeckillHandler) Seckill(c *gin.Context) {
	var req seckill.SeckillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Seckill(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetCoupon 获取优惠券信息
func (h *SeckillHandler) GetCoupon(c *gin.Context) {
	var couponID int64
	if err := c.ShouldBindUri(&struct {
		ID int64 `uri:"id" binding:"required"`
	}{ID: couponID}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	coupon, err := h.service.GetCoupon(c.Request.Context(), couponID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

// InitStock 初始化库存
func (h *SeckillHandler) InitStock(c *gin.Context) {
	var couponID int64
	if err := c.ShouldBindJSON(&struct {
		CouponID int64 `json:"coupon_id" binding:"required"`
	}{CouponID: couponID}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.InitStock(c.Request.Context(), couponID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "库存初始化成功"})
}
