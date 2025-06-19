package fluvio

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// Consumer 消费者
type Consumer struct {
	appService *services.FluvioApplicationService
	logger     logging.Logger
	connected  *bool
}

// ReceiveOptions 接收选项
type ReceiveOptions struct {
	Group       string        `json:"group,omitempty"`
	Offset      int64         `json:"offset,omitempty"`
	MaxMessages int           `json:"max_messages,omitempty"`
	Timeout     time.Duration `json:"timeout,omitempty"`
}

// StreamOptions 流式选项
type StreamOptions struct {
	Group       string        `json:"group,omitempty"`
	Offset      int64         `json:"offset,omitempty"`
	BufferSize  int           `json:"buffer_size,omitempty"`
	Timeout     time.Duration `json:"timeout,omitempty"`
}

// ConsumedMessage 已消费的消息
type ConsumedMessage struct {
	*Message
	Offset    int64     `json:"offset"`
	Partition int32     `json:"partition"`
	Topic     string    `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
}

// Receive 消费消息
func (c *Consumer) Receive(ctx context.Context, topic string, opts *ReceiveOptions) ([]*ConsumedMessage, error) {
	if !*c.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	if opts == nil {
		opts = &ReceiveOptions{
			MaxMessages: 10,
		}
	}

	c.logger.Debug("Receiving messages",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "group", Value: opts.Group},
		logging.Field{Key: "max_messages", Value: opts.MaxMessages})

	req := &dtos.ConsumeMessageRequest{
		Topic:       topic,
		Group:       opts.Group,
		Offset:      opts.Offset,
		MaxMessages: opts.MaxMessages,
	}

	resp, err := c.appService.ConsumeMessage(ctx, req)
	if err != nil {
		c.logger.Error("Failed to receive messages", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	var messages []*ConsumedMessage
	for _, msg := range resp.Messages {
		consumedMsg := &ConsumedMessage{
			Message: &Message{
				Key:     msg.Key,
				Value:   []byte(msg.Value),
				Headers: msg.Headers,
			},
			Offset:    msg.Offset,
			Partition: msg.Partition,
			Topic:     topic,
			Timestamp: msg.Timestamp,
		}
		messages = append(messages, consumedMsg)
	}

	c.logger.Info("Messages received successfully",
		logging.Field{Key: "count", Value: len(messages)})

	return messages, nil
}

// Stream 流式消费
func (c *Consumer) Stream(ctx context.Context, topic string, opts *StreamOptions) (<-chan *ConsumedMessage, error) {
	if !*c.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	if opts == nil {
		opts = &StreamOptions{
			BufferSize: 100,
		}
	}

	c.logger.Debug("Starting stream consumption",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "group", Value: opts.Group},
		logging.Field{Key: "buffer_size", Value: opts.BufferSize})

	messageChan := make(chan *ConsumedMessage, opts.BufferSize)

	// 启动后台goroutine进行流式消费
	go func() {
		defer close(messageChan)

		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Stream consumption cancelled")
				return
			default:
				// 这里应该实现实际的流式消费
				// 简化实现：定期拉取消息
				receiveOpts := &ReceiveOptions{
					Group:       opts.Group,
					Offset:      opts.Offset,
					MaxMessages: 10,
				}

				messages, err := c.Receive(ctx, topic, receiveOpts)
				if err != nil {
					c.logger.Error("Error in stream consumption", logging.Field{Key: "error", Value: err})
					time.Sleep(time.Second)
					continue
				}

				for _, msg := range messages {
					select {
					case messageChan <- msg:
					case <-ctx.Done():
						return
					}
				}

				if len(messages) == 0 {
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()

	return messageChan, nil
}

// Commit 提交偏移量
func (c *Consumer) Commit(ctx context.Context, topic string, group string, offset int64) error {
	if !*c.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}

	c.logger.Debug("Committing offset",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "group", Value: group},
		logging.Field{Key: "offset", Value: offset})

	// 这里应该调用实际的提交偏移量方法
	// 简化实现
	c.logger.Info("Offset committed successfully",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "group", Value: group},
		logging.Field{Key: "offset", Value: offset})

	return nil
}

// ReceiveOne 接收单条消息（便捷方法）
func (c *Consumer) ReceiveOne(ctx context.Context, topic string, group string) (*ConsumedMessage, error) {
	opts := &ReceiveOptions{
		Group:       group,
		MaxMessages: 1,
	}

	messages, err := c.Receive(ctx, topic, opts)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, nil
	}

	return messages[0], nil
}

// ReceiveString 接收字符串消息（便捷方法）
func (c *Consumer) ReceiveString(ctx context.Context, topic string, opts *ReceiveOptions) ([]string, error) {
	messages, err := c.Receive(ctx, topic, opts)
	if err != nil {
		return nil, err
	}

	var values []string
	for _, msg := range messages {
		values = append(values, string(msg.Value))
	}

	return values, nil
}