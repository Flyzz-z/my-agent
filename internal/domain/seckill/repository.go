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

	// SaveCompensationTask 保存补偿任务
	SaveCompensationTask(ctx context.Context, task *CompensationTask) error

	// GetPendingCompensationTasks 获取待处理的补偿任务
	GetPendingCompensationTasks(ctx context.Context, limit int) ([]*CompensationTask, error)

	// UpdateCompensationTaskStatus 更新补偿任务状态
	UpdateCompensationTaskStatus(ctx context.Context, taskID int64, status int, retryCount int) error
}

// CacheRepository 缓存仓库接口
type CacheRepository interface {
	// GetStock 获取缓存中的库存
	GetStock(ctx context.Context, couponID int64) (int64, error)

	// DecrStock 使用 Lua 脚本原子性扣减库存
	// 返回扣减后的库存，如果库存不足返回 -1
	DecrStock(ctx context.Context, couponID int64) (int64, error)

	// SetStock 设置缓存中的库存
	SetStock(ctx context.Context, couponID int64, stock int64) error

	// IncrStock 原子性增加库存
	IncrStock(ctx context.Context, couponID int64, delta int64) error
}
