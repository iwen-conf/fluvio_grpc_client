package repositories

import (
	"context"
	"time"
	
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/grpc"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// GRPCMessageRepository gRPC消息仓储实现
type GRPCMessageRepository struct {
	client grpc.Client
	logger logging.Logger
}

// NewGRPCMessageRepository 创建gRPC消息仓储
func NewGRPCMessageRepository(client grpc.Client, logger logging.Logger) repositories.MessageRepository {
	return &GRPCMessageRepository{
		client: client,
		logger: logger,
	}
}

// Produce 生产消息
func (r *GRPCMessageRepository) Produce(ctx context.Context, message *entities.Message) error {
	r.logger.Debug("生产消息", 
		logging.Field{Key: "topic", Value: message.Topic},
		logging.Field{Key: "key", Value: message.Key},
		logging.Field{Key: "message_id", Value: message.MessageID})
	
	// 转换为protobuf消息
	pbMessage := &pb.Message{
		Topic:     message.Topic,
		Key:       message.Key,
		Value:     message.Value,
		MessageId: message.MessageID,
		Headers:   message.Headers,
	}
	
	req := &pb.ProduceRequest{
		Message: pbMessage,
	}
	
	// 调用gRPC服务
	resp, err := r.client.Produce(ctx, req)
	if err != nil {
		r.logger.Error("生产消息失败", 
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: message.Topic})
		return err
	}
	
	// 更新消息元数据
	message.ID = resp.GetMessageId()
	message.Partition = resp.GetPartition()
	message.Offset = resp.GetOffset()
	
	r.logger.Info("消息生产成功",
		logging.Field{Key: "message_id", Value: message.MessageID},
		logging.Field{Key: "partition", Value: message.Partition},
		logging.Field{Key: "offset", Value: message.Offset})
	
	return nil
}

// ProduceBatch 批量生产消息
func (r *GRPCMessageRepository) ProduceBatch(ctx context.Context, messages []*entities.Message) error {
	r.logger.Debug("批量生产消息", logging.Field{Key: "count", Value: len(messages)})
	
	// 转换为protobuf消息
	pbMessages := make([]*pb.Message, len(messages))
	for i, message := range messages {
		pbMessages[i] = &pb.Message{
			Topic:     message.Topic,
			Key:       message.Key,
			Value:     message.Value,
			MessageId: message.MessageID,
			Headers:   message.Headers,
		}
	}
	
	req := &pb.ProduceBatchRequest{
		Messages: pbMessages,
	}
	
	// 调用gRPC服务
	resp, err := r.client.ProduceBatch(ctx, req)
	if err != nil {
		r.logger.Error("批量生产消息失败", logging.Field{Key: "error", Value: err})
		return err
	}
	
	// 更新消息元数据
	for i, result := range resp.GetResults() {
		if i < len(messages) {
			messages[i].ID = result.GetMessageId()
			messages[i].Partition = result.GetPartition()
			messages[i].Offset = result.GetOffset()
		}
	}
	
	r.logger.Info("批量消息生产成功", 
		logging.Field{Key: "total", Value: len(messages)},
		logging.Field{Key: "success", Value: resp.GetSuccessCount()})
	
	return nil
}

// Consume 消费消息
func (r *GRPCMessageRepository) Consume(ctx context.Context, topic string, partition int32, offset int64, maxMessages int) ([]*entities.Message, error) {
	r.logger.Debug("消费消息",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition},
		logging.Field{Key: "offset", Value: offset},
		logging.Field{Key: "max_messages", Value: maxMessages})
	
	req := &pb.ConsumeRequest{
		Topic:       topic,
		Partition:   partition,
		Offset:      offset,
		MaxMessages: int32(maxMessages),
	}
	
	// 调用gRPC服务
	resp, err := r.client.Consume(ctx, req)
	if err != nil {
		r.logger.Error("消费消息失败", logging.Field{Key: "error", Value: err})
		return nil, err
	}
	
	// 转换为实体
	messages := make([]*entities.Message, len(resp.GetMessages()))
	for i, pbMessage := range resp.GetMessages() {
		message := &entities.Message{
			ID:        pbMessage.GetMessageId(),
			MessageID: pbMessage.GetMessageId(),
			Topic:     pbMessage.GetTopic(),
			Key:       pbMessage.GetKey(),
			Value:     pbMessage.GetValue(),
			Headers:   pbMessage.GetHeaders(),
			Partition: pbMessage.GetPartition(),
			Offset:    pbMessage.GetOffset(),
			Timestamp: time.Now(), // 简化实现，实际应该从protobuf获取
		}
		messages[i] = message
	}
	
	r.logger.Info("消息消费成功", 
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "count", Value: len(messages)})
	
	return messages, nil
}

// ConsumeFiltered 过滤消费消息
func (r *GRPCMessageRepository) ConsumeFiltered(ctx context.Context, topic string, filters []*valueobjects.FilterCondition, maxMessages int) ([]*entities.Message, error) {
	r.logger.Debug("过滤消费消息",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "filters", Value: len(filters)},
		logging.Field{Key: "max_messages", Value: maxMessages})
	
	// 转换过滤条件
	pbFilters := make([]*pb.FilterCondition, len(filters))
	for i, filter := range filters {
		pbFilters[i] = &pb.FilterCondition{
			Type:     string(filter.Type),
			Field:    filter.Field,
			Operator: string(filter.Operator),
			Value:    filter.Value,
		}
	}
	
	req := &pb.ConsumeFilteredRequest{
		Topic:       topic,
		Filters:     pbFilters,
		MaxMessages: int32(maxMessages),
		AndLogic:    true, // 简化实现
	}
	
	// 调用gRPC服务
	resp, err := r.client.ConsumeFiltered(ctx, req)
	if err != nil {
		r.logger.Error("过滤消费失败", logging.Field{Key: "error", Value: err})
		return nil, err
	}
	
	// 转换为实体
	messages := make([]*entities.Message, len(resp.GetMessages()))
	for i, pbMessage := range resp.GetMessages() {
		message := &entities.Message{
			ID:        pbMessage.GetMessageId(),
			MessageID: pbMessage.GetMessageId(),
			Topic:     pbMessage.GetTopic(),
			Key:       pbMessage.GetKey(),
			Value:     pbMessage.GetValue(),
			Headers:   pbMessage.GetHeaders(),
			Partition: pbMessage.GetPartition(),
			Offset:    pbMessage.GetOffset(),
			Timestamp: time.Now(),
		}
		messages[i] = message
	}
	
	r.logger.Info("过滤消费成功",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "filtered_count", Value: len(messages)},
		logging.Field{Key: "total_scanned", Value: resp.GetTotalScanned()})
	
	return messages, nil
}

// ConsumeStream 流式消费消息
func (r *GRPCMessageRepository) ConsumeStream(ctx context.Context, topic string, partition int32, offset int64) (<-chan *entities.Message, error) {
	r.logger.Debug("开始流式消费",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition},
		logging.Field{Key: "offset", Value: offset})
	
	// 这里应该实现流式消费逻辑
	// 简化实现，返回一个空的channel
	ch := make(chan *entities.Message)
	close(ch)
	
	return ch, nil
}

// GetOffset 获取偏移量
func (r *GRPCMessageRepository) GetOffset(ctx context.Context, topic string, partition int32, consumerGroup string) (int64, error) {
	// 简化实现
	return 0, nil
}

// CommitOffset 提交偏移量
func (r *GRPCMessageRepository) CommitOffset(ctx context.Context, topic string, partition int32, consumerGroup string, offset int64) error {
	// 简化实现
	return nil
}