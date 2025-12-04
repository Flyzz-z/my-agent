package mq

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// Producer RocketMQ 生产者（基础设施层，通用）
type Producer struct {
	producer rocketmq.Producer
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	NameServerAddr []string // NameServer 地址列表
	GroupName      string   // 生产者组名
	RetryTimes     int      // 重试次数
}

// NewProducer 创建 RocketMQ 生产者
func NewProducer(cfg ProducerConfig) (*Producer, error) {
	// 创建生产者
	p, err := rocketmq.NewProducer(
		producer.WithNameServer(cfg.NameServerAddr),
		producer.WithGroupName(cfg.GroupName),
		producer.WithRetry(cfg.RetryTimes),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 RocketMQ Producer 失败: %w", err)
	}

	// 启动生产者
	err = p.Start()
	if err != nil {
		return nil, fmt.Errorf("启动 RocketMQ Producer 失败: %w", err)
	}

	return &Producer{
		producer: p,
	}, nil
}

// SendMessage 发送消息（通用方法）
func (p *Producer) SendMessage(ctx context.Context, topic, tag string, body []byte, keys ...string) error {
	// 构建消息
	msg := &primitive.Message{
		Topic: topic,
		Body:  body,
	}

	// 设置消息 Tag
	if tag != "" {
		msg.WithTag(tag)
	}

	// 设置消息 Key
	if len(keys) > 0 {
		msg.WithKeys(keys)
	}

	// 发送消息
	result, err := p.producer.SendSync(ctx, msg)
	if err != nil {
		return fmt.Errorf("发送消息失败: %w", err)
	}

	// 检查发送状态
	if result.Status != primitive.SendOK {
		return fmt.Errorf("消息发送状态异常: %v", result.Status)
	}

	return nil
}

// Shutdown 关闭生产者
func (p *Producer) Shutdown() error {
	return p.producer.Shutdown()
}

// Consumer RocketMQ 消费者（基础设施层，通用）
type Consumer struct {
	consumer rocketmq.PushConsumer
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	NameServerAddr []string // NameServer 地址列表
	GroupName      string   // 消费者组名
	Topic          string   // 主题
	Tag            string   // 消息标签过滤
}

// MessageHandler 消息处理函数类型
type MessageHandler func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error)

// NewConsumer 创建 RocketMQ 消费者
func NewConsumer(cfg ConsumerConfig, handler MessageHandler) (*Consumer, error) {
	// 创建消费者
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(cfg.NameServerAddr),
		consumer.WithGroupName(cfg.GroupName),
		consumer.WithConsumerModel(consumer.Clustering), // 集群模式
	)
	if err != nil {
		return nil, fmt.Errorf("创建 RocketMQ Consumer 失败: %w", err)
	}

	// 订阅主题
	selector := consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: cfg.Tag,
	}

	err = c.Subscribe(cfg.Topic, selector, handler)
	if err != nil {
		return nil, fmt.Errorf("订阅主题失败: %w", err)
	}

	return &Consumer{
		consumer: c,
	}, nil
}

// Start 启动消费者
func (c *Consumer) Start() error {
	return c.consumer.Start()
}

// Shutdown 关闭消费者
func (c *Consumer) Shutdown() error {
	return c.consumer.Shutdown()
}
