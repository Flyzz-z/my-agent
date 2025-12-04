package seckill

import (
	"context"
	"errors"
	"log"
	"rag-agent/config"
)

var (
	ErrStockNotEnough = errors.New("库存不足")
	ErrCouponNotFound = errors.New("优惠券不存在")
	ErrCouponExpired  = errors.New("优惠券已过期")
	ErrLockFailed     = errors.New("获取锁失败")
)

// Service 秒杀服务
type Service struct {
	repo       Repository
	cache      CacheRepository
	mqProducer MQProducer
	cfg        *config.SeckillConfig
}

// MQProducer 消息队列生产者接口
type MQProducer interface {
	SendOrderMessage(ctx context.Context, order *Order) error
}

// NewService 创建秒杀服务
func NewService(repo Repository, cache CacheRepository, mq MQProducer, cfg *config.SeckillConfig) *Service {
	return &Service{
		repo:       repo,
		cache:      cache,
		mqProducer: mq,
		cfg:        cfg,
	}
}

// Seckill 秒杀接口
func (s *Service) Seckill(ctx context.Context, req *SeckillRequest) (*SeckillResponse, error) {
	// 1. 使用 Lua 脚本原子性扣减库存
	stock, err := s.cache.DecrStock(ctx, req.CouponID)
	if err != nil {
		return &SeckillResponse{
			Success: false,
			Message: "系统错误",
		}, err
	}

	// 2. 检查库存是否充足（Lua 脚本返回 -1 表示库存不足）
	if stock < 0 {
		return &SeckillResponse{
			Success: false,
			Message: "库存不足",
		}, ErrStockNotEnough
	}

	// 3. 创建订单
	order := &Order{
		UserID:   req.UserID,
		CouponID: req.CouponID,
		Status:   0, // 待支付
	}

	// 4. 发送到MQ
	err = s.mqProducer.SendOrderMessage(ctx, order)
	if err != nil {
		log.Printf("发送MQ失败: %v, 回滚Redis库存", err)

		// 原子性回滚 Redis 库存（+1）
		if incrErr := s.cache.IncrStock(ctx, req.CouponID, 1); incrErr != nil {
			log.Printf("回滚库存失败: %v", incrErr)
		}

		// 返回失败
		return &SeckillResponse{
			Success: false,
			Message: "秒杀失败，请重试",
		}, err
	}

	// 5. 发送成功
	return &SeckillResponse{
		Success: true,
		Message: "秒杀成功，订单处理中",
		OrderID: order.ID,
	}, nil
}

// GetCoupon 获取优惠券信息
func (s *Service) GetCoupon(ctx context.Context, couponID int64) (*Coupon, error) {
	return s.repo.GetCoupon(ctx, couponID)
}

// InitStock 初始化库存到Redis
func (s *Service) InitStock(ctx context.Context, couponID int64) error {
	coupon, err := s.repo.GetCoupon(ctx, couponID)
	if err != nil {
		return err
	}

	return s.cache.SetStock(ctx, couponID, coupon.RemainStock)
}
