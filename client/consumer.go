package client

import (
	"context"
	"io"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	"github.com/iwen-conf/fluvio_grpc_client/types"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"

	"google.golang.org/grpc"
)

// Consumer 消息消费者
type Consumer struct {
	client *Client
}

// NewConsumer 创建消息消费者
func NewConsumer(client *Client) *Consumer {
	return &Consumer{
		client: client,
	}
}

// Consume 消费消息
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

	req := &pb.ConsumeRequest{
		Topic:       opts.Topic,
		Group:       opts.Group,
		Offset:      opts.Offset,
		Partition:   opts.Partition,
		MaxMessages: opts.MaxMessages,
	}

	c.client.logger.Debug("消费消息", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset},
		logger.Field{Key: "max_messages", Value: opts.MaxMessages})

	var messages []*types.Message
	err := c.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.Consume(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "消费消息失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换消息
		for _, msg := range resp.GetMessages() {
			message := &types.Message{
				Topic:     opts.Topic,
				Key:       msg.GetKey(),
				Value:     msg.GetMessage(),
				Headers:   msg.GetHeaders(),
				Offset:    msg.GetOffset(),
				Partition: msg.GetPartition(),
			}

			if msg.GetTimestamp() != nil {
				message.Timestamp = msg.GetTimestamp().AsTime()
			}

			messages = append(messages, message)
		}

		return nil
	})

	if err != nil {
		c.client.logger.Error("消费消息失败", 
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "group", Value: opts.Group},
			logger.Field{Key: "error", Value: err})
		return nil, err
	}

	c.client.logger.Info("消息消费成功", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "count", Value: len(messages)})

	return messages, nil
}

// ConsumeStream 流式消费消息
func (c *Consumer) ConsumeStream(ctx context.Context, opts types.StreamConsumeOptions) (<-chan *StreamMessage, error) {
	if c.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Group == "" {
		opts.Group = "default"
	}

	req := &pb.StreamConsumeRequest{
		Topic:     opts.Topic,
		Group:     opts.Group,
		Offset:    opts.Offset,
		Partition: opts.Partition,
	}

	c.client.logger.Debug("开始流式消费", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset})

	messageChan := make(chan *StreamMessage, 100) // 缓冲通道

	go func() {
		defer close(messageChan)

		err := c.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
			client := pb.NewFluvioServiceClient(conn)
			
			stream, err := client.StreamConsume(ctx, req)
			if err != nil {
				return errors.Wrap(errors.ErrInternal, "创建流式消费失败", err)
			}

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				resp, err := stream.Recv()
				if err == io.EOF {
					c.client.logger.Info("流式消费结束", 
						logger.Field{Key: "topic", Value: opts.Topic})
					return nil
				}
				if err != nil {
					return errors.Wrap(errors.ErrInternal, "接收流式消息失败", err)
				}

				if resp.GetError() != "" {
					streamMsg := &StreamMessage{
						Message: nil,
						Error:   errors.New(errors.ErrInternal, resp.GetError()),
					}
					select {
					case messageChan <- streamMsg:
					case <-ctx.Done():
						return ctx.Err()
					}
					continue
				}

				// 转换消息
				message := &types.Message{
					Topic:     opts.Topic,
					Key:       resp.GetKey(),
					Value:     resp.GetMessage(),
					Headers:   resp.GetHeaders(),
					Offset:    resp.GetOffset(),
					Partition: resp.GetPartition(),
				}

				if resp.GetTimestamp() != nil {
					message.Timestamp = resp.GetTimestamp().AsTime()
				}

				streamMsg := &StreamMessage{
					Message: message,
					Error:   nil,
				}

				select {
				case messageChan <- streamMsg:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})

		if err != nil {
			c.client.logger.Error("流式消费失败", 
				logger.Field{Key: "topic", Value: opts.Topic},
				logger.Field{Key: "error", Value: err})
			
			// 发送错误消息
			streamMsg := &StreamMessage{
				Message: nil,
				Error:   err,
			}
			select {
			case messageChan <- streamMsg:
			default:
			}
		}
	}()

	return messageChan, nil
}

// StreamMessage 流式消息
type StreamMessage struct {
	Message *types.Message
	Error   error
}

// CommitOffset 提交偏移量
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

	req := &pb.CommitOffsetRequest{
		Topic:     opts.Topic,
		Group:     opts.Group,
		Offset:    opts.Offset,
		Partition: opts.Partition,
	}

	c.client.logger.Debug("提交偏移量", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset})

	err := c.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.CommitOffset(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "提交偏移量失败", err)
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		c.client.logger.Error("提交偏移量失败", 
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "group", Value: opts.Group},
			logger.Field{Key: "offset", Value: opts.Offset},
			logger.Field{Key: "error", Value: err})
		return err
	}

	c.client.logger.Info("偏移量提交成功", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset})

	return nil
}
