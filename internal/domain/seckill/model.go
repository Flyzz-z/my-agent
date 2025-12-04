package seckill

import "time"

// Coupon 优惠券模型
type Coupon struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	TotalStock  int64     `json:"total_stock" db:"total_stock"`
	RemainStock int64     `json:"remain_stock" db:"remain_stock"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Status      int       `json:"status" db:"status"` // 0-未开始, 1-进行中, 2-已结束
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Order 订单模型
type Order struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	CouponID  int64     `json:"coupon_id" db:"coupon_id"`
	Status    int       `json:"status" db:"status"` // 0-待支付, 1-已支付, 2-已取消
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SeckillRequest 秒杀请求
type SeckillRequest struct {
	UserID   int64 `json:"user_id" binding:"required"`
	CouponID int64 `json:"coupon_id" binding:"required"`
}

// SeckillResponse 秒杀响应
type SeckillResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	OrderID int64  `json:"order_id,omitempty"`
}

// CompensationTask 补偿任务模型（用于MQ发送失败后的补偿处理）
type CompensationTask struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	CouponID   int64     `json:"coupon_id" db:"coupon_id"`
	OrderID    int64     `json:"order_id" db:"order_id"`
	Status     int       `json:"status" db:"status"`         // 0-待处理, 1-处理中, 2-已完成, -1-失败
	RetryCount int       `json:"retry_count" db:"retry_count"` // 重试次数
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
