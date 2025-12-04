package seckill

import (
	"context"
	"encoding/json"
	"fmt"

	"rag-agent/internal/infrastructure/mq"
)

// OrderMQProducer 订单消息生产者（domain 层适配器）
type OrderMQProducer struct {
	producer *mq.Producer
	topic    string
}

// NewOrderMQProducer 创建订单消息生产者
func NewOrderMQProducer(producer *mq.Producer, topic string) MQProducer {
	return &OrderMQProducer{
		producer: producer,
		topic:    topic,
	}
}

// SendOrderMessage 发送订单消息
func (p *OrderMQProducer) SendOrderMessage(ctx context.Context, order *Order) error {
	// 序列化订单
	orderData, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("序列化订单失败: %w", err)
	}

	// 发送消息
	key := fmt.Sprintf("order_%d", order.ID)
	err = p.producer.SendMessage(ctx, p.topic, "seckill_order", orderData, key)
	if err != nil {
		return fmt.Errorf("发送订单消息失败: %w", err)
	}

	return nil
}
