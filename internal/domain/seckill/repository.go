package seckill

import "context"

// Repository 秒杀数据仓库接口
type Repository interface {
	// GetCoupon 获取优惠券信息
	GetCoupon(ctx context.Context, couponID int64) (*Coupon, error)

	// DecrStock 减少库存
	DecrStock(ctx context.Context, couponID int64) error

	// CreateOrder 创建订单
	CreateOrder(ctx context.Context, order *Order) error

	// GetOrder 获取订单
	GetOrder(ctx context.Context, orderID int64) (*Order, error)

	// UpdateOrderStatus 更新订单状态
	UpdateOrderStatus(ctx context.Context, orderID int64, status int) error
}

// CacheRepository 缓存仓库接口
type CacheRepository interface {
	// GetStock 获取缓存中的库存
	GetStock(ctx context.Context, couponID int64) (int64, error)

	// DecrStock 减少缓存中的库存
	DecrStock(ctx context.Context, couponID int64) (int64, error)

	// SetStock 设置缓存中的库存
	SetStock(ctx context.Context, couponID int64, stock int64) error

	// TryLock 尝试获取分布式锁
	TryLock(ctx context.Context, key string) (bool, error)

	// Unlock 释放分布式锁
	Unlock(ctx context.Context, key string) error
}
