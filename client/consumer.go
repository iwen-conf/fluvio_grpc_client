package client

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
	"github.com/iwen-conf/fluvio_grpc_client/types"
)

// Consumer 消息消费者（向后兼容）
type Consumer struct {
	client *Client
}

// NewConsumer 创建消息消费者
func NewConsumer(client *Client) *Consumer {
	return &Consumer{
		client: client,
	}
}

// Consume 消费消息（简化实现）
func (c *Consumer) Consume(ctx context.Context, opts types.ConsumeOptions) ([]*types.Message, error) {
	if c.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Group == "" {
		opts.Group = "default"
	}

	if opts.MaxMessages <= 0 {
		opts.MaxMessages = 1
	}

	// 简化实现：返回模拟消息
	messages := []*types.Message{
		{
			Topic:     opts.Topic,
			Key:       "example-key",
			Value:     "example-message",
			Headers:   make(map[string]string),
			Offset:    0,
			Partition: 0,
			MessageID: "msg-001",
			Timestamp: time.Now(),
		},
	}

	return messages, nil
}

// ConsumeStream 流式消费消息（简化实现）
func (c *Consumer) ConsumeStream(ctx context.Context, opts types.StreamConsumeOptions) (<-chan *StreamMessage, error) {
	if c.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	// 简化实现：返回空的channel
	messageChan := make(chan *StreamMessage)
	close(messageChan)

	return messageChan, nil
}

// StreamMessage 流式消息
type StreamMessage struct {
	Message *types.Message
	Error   error
}

// CommitOffset 提交偏移量（简化实现）
func (c *Consumer) CommitOffset(ctx context.Context, opts types.CommitOffsetOptions) error {
	if c.client.isClosed() {
		return errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Group == "" {
		opts.Group = "default"
	}

	// 简化实现：总是成功
	return nil
}

// ConsumeWithRetry 带重试的消息消费（简化实现）
func (c *Consumer) ConsumeWithRetry(ctx context.Context, opts types.ConsumeOptions) ([]*types.Message, error) {
	return c.Consume(ctx, opts)
}

// CommitOffsetWithRetry 带重试的偏移量提交（简化实现）
func (c *Consumer) CommitOffsetWithRetry(ctx context.Context, opts types.CommitOffsetOptions) error {
	return c.CommitOffset(ctx, opts)
}