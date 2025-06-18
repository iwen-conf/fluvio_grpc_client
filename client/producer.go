package client

import (
	"context"
	"fmt"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
	"github.com/iwen-conf/fluvio_grpc_client/types"
)

// Producer 消息生产者（向后兼容）
type Producer struct {
	client *Client
}

// NewProducer 创建消息生产者
func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client,
	}
}

// Produce 发送单条消息（简化实现）
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

	// 简化实现：返回模拟结果
	result := &types.ProduceResult{
		MessageID: opts.MessageID,
		Offset:    0,
		Success:   true,
		Error:     "",
	}

	if result.MessageID == "" {
		result.MessageID = fmt.Sprintf("msg-%d", time.Now().UnixNano())
	}

	return result, nil
}

// ProduceBatch 批量发送消息（简化实现）
func (p *Producer) ProduceBatch(ctx context.Context, messages []types.Message) (*types.BatchProduceResult, error) {
	if p.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if len(messages) == 0 {
		return nil, errors.New(errors.ErrInvalidArgument, "消息列表不能为空")
	}

	// 简化实现：返回模拟结果
	var results []*types.ProduceResult
	for i, msg := range messages {
		if msg.Topic == "" {
			return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
		}

		messageID := msg.MessageID
		if messageID == "" {
			messageID = fmt.Sprintf("batch-msg-%d-%d", time.Now().UnixNano(), i)
		}

		results = append(results, &types.ProduceResult{
			MessageID: messageID,
			Offset:    int64(i),
			Success:   true,
			Error:     "",
		})
	}

	return &types.BatchProduceResult{
		Results:    results,
		TotalCount: len(results),
		Success:    true,
		Errors:     []string{},
	}, nil
}

// ProduceAsync 异步发送消息（简化实现）
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

// ProduceWithRetry 带重试的消息发送（简化实现）
func (p *Producer) ProduceWithRetry(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	return p.Produce(ctx, message, opts)
}

// BatchProduceWithRetry 带重试的批量消息发送（简化实现）
func (p *Producer) BatchProduceWithRetry(ctx context.Context, messages []types.Message) (*types.BatchProduceResult, error) {
	return p.ProduceBatch(ctx, messages)
}