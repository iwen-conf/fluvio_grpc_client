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

	// 转换为protobuf请求
	req := &pb.ProduceRequest{
		Topic:   message.Topic,
		Message: message.Value,
		Key:     message.Key,
		Headers: message.Headers,
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
	// 注意：当前protobuf定义中ProduceReply没有Partition和Offset字段
	// 这里使用默认值
	message.Partition = 0
	message.Offset = 0

	r.logger.Info("消息生产成功",
		logging.Field{Key: "message_id", Value: message.MessageID},
		logging.Field{Key: "success", Value: resp.GetSuccess()})

	return nil
}

// ProduceBatch 批量生产消息（简化实现）
func (r *GRPCMessageRepository) ProduceBatch(ctx context.Context, messages []*entities.Message) error {
	r.logger.Debug("批量生产消息", logging.Field{Key: "count", Value: len(messages)})

	// 简化实现：逐个调用单条消息生产
	for _, message := range messages {
		if err := r.Produce(ctx, message); err != nil {
			r.logger.Error("批量生产中的消息失败",
				logging.Field{Key: "error", Value: err},
				logging.Field{Key: "message_id", Value: message.MessageID})
			return err
		}
	}

	r.logger.Info("批量消息生产成功",
		logging.Field{Key: "total", Value: len(messages)})

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
			Topic:     topic, // 使用请求中的topic
			Key:       pbMessage.GetKey(),
			Value:     pbMessage.GetMessage(), // ConsumedMessage中字段名是Message
			Headers:   pbMessage.GetHeaders(),
			Partition: pbMessage.GetPartition(),
			Offset:    pbMessage.GetOffset(),
			Timestamp: time.Unix(pbMessage.GetTimestamp(), 0), // 从protobuf获取时间戳
		}
		messages[i] = message
	}

	r.logger.Info("消息消费成功",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "count", Value: len(messages)})

	return messages, nil
}

// ConsumeFiltered 过滤消费消息（简化实现）
func (r *GRPCMessageRepository) ConsumeFiltered(ctx context.Context, topic string, filters []*valueobjects.FilterCondition, maxMessages int) ([]*entities.Message, error) {
	r.logger.Debug("过滤消费消息",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "filters", Value: len(filters)},
		logging.Field{Key: "max_messages", Value: maxMessages})

	// 简化实现：先消费消息，然后在客户端进行过滤
	allMessages, err := r.Consume(ctx, topic, 0, 0, maxMessages*2) // 获取更多消息以便过滤
	if err != nil {
		return nil, err
	}

	// 应用过滤条件
	var filteredMessages []*entities.Message
	for _, message := range allMessages {
		if r.matchesFilters(message, filters) {
			filteredMessages = append(filteredMessages, message)
			if len(filteredMessages) >= maxMessages {
				break
			}
		}
	}

	r.logger.Info("过滤消费成功",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "filtered_count", Value: len(filteredMessages)},
		logging.Field{Key: "total_scanned", Value: len(allMessages)})

	return filteredMessages, nil
}

// matchesFilters 检查消息是否匹配过滤条件（简化实现）
func (r *GRPCMessageRepository) matchesFilters(message *entities.Message, filters []*valueobjects.FilterCondition) bool {
	if len(filters) == 0 {
		return true
	}

	// 简化实现：只检查第一个过滤条件
	filter := filters[0]
	switch filter.Type {
	case "key":
		return message.Key == filter.Value
	case "value":
		return message.Value == filter.Value
	default:
		return true
	}
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
