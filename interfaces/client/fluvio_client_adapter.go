package client

import (
	"context"
	"time"
	
	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/interfaces/api"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// FluvioClientAdapter Fluvio客户端适配器
// 实现向后兼容的公共API，内部使用新的架构
type FluvioClientAdapter struct {
	appService *services.FluvioApplicationService
	connected  bool
}

// NewFluvioClientAdapter 创建客户端适配器
func NewFluvioClientAdapter(appService *services.FluvioApplicationService) api.FluvioAPI {
	return &FluvioClientAdapter{
		appService: appService,
		connected:  false,
	}
}

// Connect 连接到服务器
func (c *FluvioClientAdapter) Connect() error {
	// 这里应该实际连接到gRPC服务器
	// 简化实现
	c.connected = true
	return nil
}

// Close 关闭连接
func (c *FluvioClientAdapter) Close() error {
	c.connected = false
	return nil
}

// HealthCheck 健康检查
func (c *FluvioClientAdapter) HealthCheck(ctx context.Context) error {
	if !c.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}
	// 简化实现
	return nil
}

// Ping 测试连接
func (c *FluvioClientAdapter) Ping(ctx context.Context) (time.Duration, error) {
	if !c.connected {
		return 0, errors.New(errors.ErrConnection, "client not connected")
	}
	
	start := time.Now()
	// 这里应该实际发送ping请求
	// 简化实现
	return time.Since(start), nil
}

// Producer 获取生产者API
func (c *FluvioClientAdapter) Producer() api.ProducerAPI {
	return &ProducerAdapter{
		appService: c.appService,
		connected:  &c.connected,
	}
}

// Consumer 获取消费者API
func (c *FluvioClientAdapter) Consumer() api.ConsumerAPI {
	return &ConsumerAdapter{
		appService: c.appService,
		connected:  &c.connected,
	}
}

// Topic 获取主题API
func (c *FluvioClientAdapter) Topic() api.TopicAPI {
	return &TopicAdapter{
		appService: c.appService,
		connected:  &c.connected,
	}
}

// Admin 获取管理API
func (c *FluvioClientAdapter) Admin() api.AdminAPI {
	return &AdminAdapter{
		appService: c.appService,
		connected:  &c.connected,
	}
}

// ProducerAdapter 生产者适配器
type ProducerAdapter struct {
	appService *services.FluvioApplicationService
	connected  *bool
}

// Produce 生产消息
func (p *ProducerAdapter) Produce(ctx context.Context, value string, opts api.ProduceOptions) (*api.ProduceResult, error) {
	if !*p.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	req := &dtos.ProduceMessageRequest{
		Topic:     opts.Topic,
		Key:       opts.Key,
		Value:     value,
		MessageID: opts.MessageID,
		Headers:   opts.Headers,
	}
	
	resp, err := p.appService.ProduceMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &api.ProduceResult{
		MessageID: resp.MessageID,
		Topic:     resp.Topic,
		Partition: resp.Partition,
		Offset:    resp.Offset,
	}, nil
}

// ProduceBatch 批量生产消息
func (p *ProducerAdapter) ProduceBatch(ctx context.Context, messages []api.Message) (*api.ProduceBatchResult, error) {
	if !*p.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	// 转换消息
	dtoMessages := make([]*dtos.ProduceMessageRequest, len(messages))
	for i, msg := range messages {
		dtoMessages[i] = &dtos.ProduceMessageRequest{
			Topic:     msg.Topic,
			Key:       msg.Key,
			Value:     msg.Value,
			MessageID: msg.MessageID,
			Headers:   msg.Headers,
		}
	}
	
	req := &dtos.ProduceBatchRequest{
		Messages: dtoMessages,
	}
	
	resp, err := p.appService.ProduceBatch(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 转换结果
	results := make([]*api.ProduceResult, len(resp.Results))
	for i, result := range resp.Results {
		results[i] = &api.ProduceResult{
			MessageID: result.MessageID,
			Topic:     result.Topic,
			Partition: result.Partition,
			Offset:    result.Offset,
		}
	}
	
	return &api.ProduceBatchResult{
		Results:       results,
		TotalMessages: resp.TotalMessages,
		SuccessCount:  resp.SuccessCount,
		FailureCount:  resp.FailureCount,
	}, nil
}

// ConsumerAdapter 消费者适配器
type ConsumerAdapter struct {
	appService *services.FluvioApplicationService
	connected  *bool
}

// Consume 消费消息
func (c *ConsumerAdapter) Consume(ctx context.Context, opts api.ConsumeOptions) ([]*api.Message, error) {
	if !*c.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	req := &dtos.ConsumeMessageRequest{
		Topic:       opts.Topic,
		Group:       opts.Group,
		Partition:   opts.Partition,
		Offset:      opts.Offset,
		MaxMessages: opts.MaxMessages,
	}
	
	resp, err := c.appService.ConsumeMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 转换消息
	messages := make([]*api.Message, len(resp.Messages))
	for i, msg := range resp.Messages {
		messages[i] = &api.Message{
			Topic:     msg.Topic,
			Key:       msg.Key,
			Value:     msg.Value,
			MessageID: msg.MessageID,
			Headers:   msg.Headers,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Timestamp: msg.Timestamp,
		}
	}
	
	return messages, nil
}

// ConsumeFiltered 过滤消费
func (c *ConsumerAdapter) ConsumeFiltered(ctx context.Context, opts api.FilteredConsumeOptions) (*api.FilteredConsumeResult, error) {
	if !*c.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	// 转换过滤条件
	filters := make([]*dtos.FilterCondition, len(opts.Filters))
	for i, filter := range opts.Filters {
		filters[i] = &dtos.FilterCondition{
			Type:     string(filter.Type),
			Field:    filter.Field,
			Operator: string(filter.Operator),
			Value:    filter.Value,
		}
	}
	
	req := &dtos.FilteredConsumeRequest{
		Topic:       opts.Topic,
		Group:       opts.Group,
		MaxMessages: opts.MaxMessages,
		Filters:     filters,
		AndLogic:    opts.AndLogic,
	}
	
	resp, err := c.appService.ConsumeFiltered(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 转换消息
	messages := make([]*api.Message, len(resp.Messages))
	for i, msg := range resp.Messages {
		messages[i] = &api.Message{
			Topic:     msg.Topic,
			Key:       msg.Key,
			Value:     msg.Value,
			MessageID: msg.MessageID,
			Headers:   msg.Headers,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Timestamp: msg.Timestamp,
		}
	}
	
	return &api.FilteredConsumeResult{
		Messages:      messages,
		FilteredCount: resp.FilteredCount,
		TotalScanned:  resp.TotalScanned,
	}, nil
}

// ConsumeStream 流式消费
func (c *ConsumerAdapter) ConsumeStream(ctx context.Context, opts api.StreamConsumeOptions) (<-chan *api.StreamMessage, error) {
	if !*c.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	// 简化实现，返回空channel
	ch := make(chan *api.StreamMessage)
	close(ch)
	return ch, nil
}