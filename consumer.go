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
	Group      string        `json:"group,omitempty"`
	Offset     int64         `json:"offset,omitempty"`
	BufferSize int           `json:"buffer_size,omitempty"`
	Timeout    time.Duration `json:"timeout,omitempty"`
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
	// 调用真实的流式消费gRPC API
	// 注意：这里使用partition 0作为默认值，实际应用中可能需要支持多分区
	appMessageChan, err := c.appService.StreamConsume(ctx, topic, 0, opts.Offset)
	if err != nil {
		c.logger.Error("Failed to start stream consumption", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	// 启动goroutine转换消息格式
	go func() {
		defer close(messageChan)

		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Stream consumption cancelled")
				return
			case entityMsg, ok := <-appMessageChan:
				if !ok {
					c.logger.Debug("Stream consumption ended")
					return
				}

				// 转换为Consumer API的消息格式
				consumedMsg := &ConsumedMessage{
					Message: &Message{
						Key:     entityMsg.Key,
						Value:   entityMsg.Value,
						Headers: entityMsg.Headers,
					},
					Topic:     entityMsg.Topic,
					Partition: entityMsg.Partition,
					Offset:    entityMsg.Offset,
					Timestamp: entityMsg.Timestamp,
				}

				select {
				case messageChan <- consumedMsg:
					// 消息发送成功
				case <-ctx.Done():
					c.logger.Debug("Stream consumption cancelled")
					return
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

	// 调用真实的提交偏移量方法
	// 注意：这里使用partition 0作为默认值，实际应用中可能需要支持多分区
	err := c.appService.CommitOffset(ctx, topic, 0, group, offset)
	if err != nil {
		c.logger.Error("Failed to commit offset",
			logging.Field{Key: "topic", Value: topic},
			logging.Field{Key: "group", Value: group},
			logging.Field{Key: "offset", Value: offset},
			logging.Field{Key: "error", Value: err})
		return err
	}

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
