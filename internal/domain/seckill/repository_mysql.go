package seckill

import (
	"context"
	"database/sql"
	"fmt"
)

// MySQLRepository MySQL 秒杀仓库实现（秒杀专用）
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository 创建秒杀仓库
func NewMySQLRepository(db *sql.DB) Repository {
	return &MySQLRepository{
		db: db,
	}
}

// GetCoupon 获取优惠券信息
func (r *MySQLRepository) GetCoupon(ctx context.Context, couponID int64) (*Coupon, error) {
	query := `
		SELECT id, name, description, total_stock, remain_stock,
		       start_time, end_time, status, created_at, updated_at
		FROM coupons
		WHERE id = ?
	`

	var coupon Coupon
	err := r.db.QueryRowContext(ctx, query, couponID).Scan(
		&coupon.ID,
		&coupon.Name,
		&coupon.Description,
		&coupon.TotalStock,
		&coupon.RemainStock,
		&coupon.StartTime,
		&coupon.EndTime,
		&coupon.Status,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrCouponNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询优惠券失败: %w", err)
	}

	return &coupon, nil
}

// DecrStock 减少库存（使用乐观锁防止超卖）
func (r *MySQLRepository) DecrStock(ctx context.Context, couponID int64) error {
	query := `
		UPDATE coupons
		SET remain_stock = remain_stock - 1,
		    updated_at = NOW()
		WHERE id = ? AND remain_stock > 0
	`

	result, err := r.db.ExecContext(ctx, query, couponID)
	if err != nil {
		return fmt.Errorf("扣减库存失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		return ErrStockNotEnough
	}

	return nil
}

// CreateOrder 创建订单
func (r *MySQLRepository) CreateOrder(ctx context.Context, order *Order) error {
	query := `
		INSERT INTO orders (user_id, coupon_id, status, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query, order.UserID, order.CouponID, order.Status)
	if err != nil {
		return fmt.Errorf("创建订单失败: %w", err)
	}

	// 获取插入的订单ID
	orderID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取订单ID失败: %w", err)
	}

	order.ID = orderID
	return nil
}

// GetOrder 获取订单
func (r *MySQLRepository) GetOrder(ctx context.Context, orderID int64) (*Order, error) {
	query := `
		SELECT id, user_id, coupon_id, status, created_at, updated_at
		FROM orders
		WHERE id = ?
	`

	var order Order
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.CouponID,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("订单不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询订单失败: %w", err)
	}

	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (r *MySQLRepository) UpdateOrderStatus(ctx context.Context, orderID int64, status int) error {
	query := `
		UPDATE orders
		SET status = ?,
		    updated_at = NOW()
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取影响行数失败: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("订单不存在")
	}

	return nil
}

// SaveCompensationTask 保存补偿任务（暂不实现）
func (r *MySQLRepository) SaveCompensationTask(ctx context.Context, task *CompensationTask) error {
	// TODO: 后续实现
	return nil
}

// GetPendingCompensationTasks 获取待处理的补偿任务（暂不实现）
func (r *MySQLRepository) GetPendingCompensationTasks(ctx context.Context, limit int) ([]*CompensationTask, error) {
	// TODO: 后续实现
	return nil, nil
}

// UpdateCompensationTaskStatus 更新补偿任务状态（暂不实现）
func (r *MySQLRepository) UpdateCompensationTaskStatus(ctx context.Context, taskID int64, status int, retryCount int) error {
	// TODO: 后续实现
	return nil
}
