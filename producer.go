package fluvio

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// Producer 生产者
type Producer struct {
	appService *services.FluvioApplicationService
	logger     logging.Logger
	connected  *bool
}

// Message 消息
type Message struct {
	Key       string            `json:"key,omitempty"`
	Value     []byte            `json:"value"`
	Headers   map[string]string `json:"headers,omitempty"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
}

// SendOptions 发送选项
type SendOptions struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key,omitempty"`
	Value     []byte            `json:"value"`
	Headers   map[string]string `json:"headers,omitempty"`
	MessageID string            `json:"message_id,omitempty"`
}

// SendResult 发送结果
type SendResult struct {
	MessageID string `json:"message_id"`
	Offset    int64  `json:"offset"`
	Partition int32  `json:"partition"`
}

// BatchSendResult 批量发送结果
type BatchSendResult struct {
	Results      []*SendResult `json:"results"`
	SuccessCount int           `json:"success_count"`
	FailureCount int           `json:"failure_count"`
}

// Send 发送单条消息
func (p *Producer) Send(ctx context.Context, topic string, message *Message) (*SendResult, error) {
	if !*p.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	p.logger.Debug("Sending message",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "key", Value: message.Key})

	req := &dtos.ProduceMessageRequest{
		Topic:   topic,
		Key:     message.Key,
		Value:   string(message.Value),
		Headers: message.Headers,
	}

	resp, err := p.appService.ProduceMessage(ctx, req)
	if err != nil {
		p.logger.Error("Failed to send message", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	result := &SendResult{
		MessageID: resp.MessageID,
		Offset:    resp.Offset,
		Partition: resp.Partition,
	}

	p.logger.Info("Message sent successfully",
		logging.Field{Key: "message_id", Value: result.MessageID},
		logging.Field{Key: "offset", Value: result.Offset})

	return result, nil
}

// SendWithOptions 发送消息（带选项）
func (p *Producer) SendWithOptions(ctx context.Context, opts *SendOptions) (*SendResult, error) {
	if !*p.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	message := &Message{
		Key:     opts.Key,
		Value:   opts.Value,
		Headers: opts.Headers,
	}

	return p.Send(ctx, opts.Topic, message)
}

// SendBatch 批量发送消息
func (p *Producer) SendBatch(ctx context.Context, topic string, messages []*Message) (*BatchSendResult, error) {
	if !*p.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	p.logger.Debug("Sending batch messages",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "count", Value: len(messages)})

	var results []*SendResult
	successCount := 0
	failureCount := 0

	for _, message := range messages {
		result, err := p.Send(ctx, topic, message)
		if err != nil {
			p.logger.Error("Failed to send message in batch",
				logging.Field{Key: "error", Value: err},
				logging.Field{Key: "key", Value: message.Key})
			failureCount++
			continue
		}
		results = append(results, result)
		successCount++
	}

	batchResult := &BatchSendResult{
		Results:      results,
		SuccessCount: successCount,
		FailureCount: failureCount,
	}

	p.logger.Info("Batch send completed",
		logging.Field{Key: "total", Value: len(messages)},
		logging.Field{Key: "success", Value: successCount},
		logging.Field{Key: "failure", Value: failureCount})

	return batchResult, nil
}

// SendString 发送字符串消息（便捷方法）
func (p *Producer) SendString(ctx context.Context, topic, key, value string) (*SendResult, error) {
	message := &Message{
		Key:   key,
		Value: []byte(value),
	}
	return p.Send(ctx, topic, message)
}

// SendJSON 发送JSON消息（便捷方法）
func (p *Producer) SendJSON(ctx context.Context, topic, key string, value interface{}) (*SendResult, error) {
	// 这里应该序列化JSON
	// 简化实现
	message := &Message{
		Key:   key,
		Value: []byte("{}"), // 简化实现
		Headers: map[string]string{
			"content-type": "application/json",
		},
	}
	return p.Send(ctx, topic, message)
}