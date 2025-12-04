package seckill

import (
	"context"
	"encoding/json"
	"log"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

// OrderConsumer 订单消费者服务（业务逻辑层）
type OrderConsumer struct {
	repo Repository
}

// NewOrderConsumer 创建订单消费者
func NewOrderConsumer(repo Repository) *OrderConsumer {
	return &OrderConsumer{
		repo: repo,
	}
}

// HandleMessage 处理订单消息（业务逻辑）
func (c *OrderConsumer) HandleMessage(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		// 解析订单消息
		var order Order
		if err := json.Unmarshal(msg.Body, &order); err != nil {
			log.Printf("解析订单消息失败: %v, msgID: %s", err, msg.MsgId)
			// 解析失败，直接返回成功，避免重复消费
			continue
		}

		log.Printf("收到订单消息: userID=%d, couponID=%d, orderID=%d", order.UserID, order.CouponID, order.ID)

		// 处理订单：扣减 MySQL 库存 + 创建订单记录
		if err := c.processOrder(ctx, &order); err != nil {
			log.Printf("处理订单失败: %v, 将重试", err)
			// 返回失败，RocketMQ 会自动重试
			return consumer.ConsumeRetryLater, err
		}

		log.Printf("订单处理成功: orderID=%d", order.ID)
	}

	return consumer.ConsumeSuccess, nil
}

// processOrder 处理订单：扣减 MySQL 库存 + 创建订单记录
func (c *OrderConsumer) processOrder(ctx context.Context, order *Order) error {
	// 1. 扣减 MySQL 库存
	if err := c.repo.DecrStock(ctx, order.CouponID); err != nil {
		log.Printf("扣减 MySQL 库存失败: %v", err)
		return err
	}

	// 2. 创建订单记录
	if err := c.repo.CreateOrder(ctx, order); err != nil {
		log.Printf("创建订单失败: %v", err)
		// TODO: 这里需要回滚 MySQL 库存，或者通过补偿机制处理
		return err
	}

	return nil
}
