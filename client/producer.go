package client

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	"github.com/iwen-conf/fluvio_grpc_client/types"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Producer 消息生产者
type Producer struct {
	client *Client
}

// NewProducer 创建消息生产者
func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client,
	}
}

// Produce 发送单条消息
func (p *Producer) Produce(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	if p.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if message == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "消息内容不能为空")
	}

	// 构建请求
	req := &pb.ProduceRequest{
		Topic:   opts.Topic,
		Message: message,
		Key:     opts.Key,
		Headers: opts.Headers,
	}

	// 设置时间戳
	if opts.Timestamp != nil {
		req.Timestamp = timestamppb.New(*opts.Timestamp)
	} else {
		req.Timestamp = timestamppb.New(time.Now())
	}

	p.client.logger.Debug("发送消息", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "key", Value: opts.Key},
		logger.Field{Key: "message_length", Value: len(message)})

	var result *types.ProduceResult
	err := p.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.Produce(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "发送消息失败", err)
		}

		result = &types.ProduceResult{
			MessageID: resp.GetMessageId(),
			Offset:    resp.GetOffset(),
			Success:   resp.GetSuccess(),
			Error:     resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		p.client.logger.Error("发送消息失败", 
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "error", Value: err})
		return nil, err
	}

	p.client.logger.Info("消息发送成功", 
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "message_id", Value: result.MessageID},
		logger.Field{Key: "offset", Value: result.Offset})

	return result, nil
}

// ProduceBatch 批量发送消息
func (p *Producer) ProduceBatch(ctx context.Context, messages []types.Message) (*types.BatchProduceResult, error) {
	if p.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if len(messages) == 0 {
		return nil, errors.New(errors.ErrInvalidArgument, "消息列表不能为空")
	}

	// 构建批量请求
	var requests []*pb.ProduceRequest
	for _, msg := range messages {
		if msg.Topic == "" {
			return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
		}

		req := &pb.ProduceRequest{
			Topic:   msg.Topic,
			Message: msg.Value,
			Key:     msg.Key,
			Headers: msg.Headers,
		}

		if !msg.Timestamp.IsZero() {
			req.Timestamp = timestamppb.New(msg.Timestamp)
		} else {
			req.Timestamp = timestamppb.New(time.Now())
		}

		requests = append(requests, req)
	}

	batchReq := &pb.BatchProduceRequest{
		Requests: requests,
	}

	p.client.logger.Debug("批量发送消息", 
		logger.Field{Key: "count", Value: len(messages)})

	var result *types.BatchProduceResult
	err := p.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.BatchProduce(ctx, batchReq)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "批量发送消息失败", err)
		}

		// 转换结果
		var results []*types.ProduceResult
		var errorMsgs []string

		for _, r := range resp.GetResults() {
			results = append(results, &types.ProduceResult{
				MessageID: r.GetMessageId(),
				Offset:    r.GetOffset(),
				Success:   r.GetSuccess(),
				Error:     r.GetError(),
			})

			if !r.GetSuccess() && r.GetError() != "" {
				errorMsgs = append(errorMsgs, r.GetError())
			}
		}

		result = &types.BatchProduceResult{
			Results:    results,
			TotalCount: len(results),
			Success:    resp.GetSuccess(),
			Errors:     errorMsgs,
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, "批量发送部分失败")
		}

		return nil
	})

	if err != nil {
		p.client.logger.Error("批量发送消息失败", 
			logger.Field{Key: "count", Value: len(messages)},
			logger.Field{Key: "error", Value: err})
		return nil, err
	}

	successCount := 0
	for _, r := range result.Results {
		if r.Success {
			successCount++
		}
	}

	p.client.logger.Info("批量消息发送完成", 
		logger.Field{Key: "total", Value: len(messages)},
		logger.Field{Key: "success", Value: successCount},
		logger.Field{Key: "failed", Value: len(messages) - successCount})

	return result, nil
}

// ProduceAsync 异步发送消息
func (p *Producer) ProduceAsync(ctx context.Context, message string, opts types.ProduceOptions) <-chan *AsyncProduceResult {
	resultChan := make(chan *AsyncProduceResult, 1)

	go func() {
		defer close(resultChan)

		result, err := p.Produce(ctx, message, opts)
		resultChan <- &AsyncProduceResult{
			Result: result,
			Error:  err,
		}
	}()

	return resultChan
}

// AsyncProduceResult 异步生产结果
type AsyncProduceResult struct {
	Result *types.ProduceResult
	Error  error
}

// ProduceWithRetry 带重试的消息发送
func (p *Producer) ProduceWithRetry(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	return p.client.withRetry(ctx, func(retryCtx context.Context) error {
		result, err := p.Produce(retryCtx, message, opts)
		if err != nil {
			return err
		}
		// 这里需要一个方式来返回result，暂时简化处理
		return nil
	})
	// TODO: 完善重试逻辑的结果返回
	return p.Produce(ctx, message, opts)
}
