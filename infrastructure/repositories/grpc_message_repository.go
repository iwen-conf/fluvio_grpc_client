package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/grpc"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Topic:     message.Topic,
		Message:   string(message.Value), // 转换为字符串
		Key:       message.Key,
		Headers:   message.Headers,
		MessageId: message.MessageID,
	}

	// 设置时间戳
	if !message.Timestamp.IsZero() {
		req.Timestamp = timestamppb.New(message.Timestamp)
	}

	// 调用gRPC服务
	resp, err := r.client.Produce(ctx, req)
	if err != nil {
		r.logger.Error("生产消息失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: message.Topic})
		return err
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("生产消息被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "topic", Value: message.Topic})
		return fmt.Errorf("produce failed: %s", errMsg)
	}

	// 更新消息元数据
	message.ID = resp.GetMessageId()
	message.MessageID = resp.GetMessageId()
	// 注意：当前protobuf定义中ProduceReply没有Partition和Offset字段
	// 这里使用默认值，实际应该从服务器响应中获取
	message.Partition = 0
	message.Offset = 0

	r.logger.Info("消息生产成功",
		logging.Field{Key: "message_id", Value: resp.GetMessageId()},
		logging.Field{Key: "topic", Value: message.Topic})

	return nil
}

// ProduceBatch 批量生产消息
func (r *GRPCMessageRepository) ProduceBatch(ctx context.Context, messages []*entities.Message) error {
	r.logger.Debug("批量生产消息", logging.Field{Key: "count", Value: len(messages)})

	if len(messages) == 0 {
		return nil
	}

	// 转换为protobuf请求
	pbMessages := make([]*pb.ProduceRequest, len(messages))
	for i, message := range messages {
		pbMessage := &pb.ProduceRequest{
			Topic:     message.Topic,
			Message:   string(message.Value),
			Key:       message.Key,
			Headers:   message.Headers,
			MessageId: message.MessageID,
		}

		// 设置时间戳
		if !message.Timestamp.IsZero() {
			pbMessage.Timestamp = timestamppb.New(message.Timestamp)
		}

		pbMessages[i] = pbMessage
	}

	// 构建批量生产请求
	req := &pb.BatchProduceRequest{
		Topic:    messages[0].Topic, // 假设所有消息都是同一个主题
		Messages: pbMessages,
	}

	// 调用gRPC服务
	resp, err := r.client.BatchProduce(ctx, req)
	if err != nil {
		r.logger.Error("批量生产消息失败", logging.Field{Key: "error", Value: err})
		return fmt.Errorf("batch produce failed: %w", err)
	}

	// 检查响应状态
	successFlags := resp.GetSuccess()
	errorMessages := resp.GetError()

	successCount := 0
	failureCount := 0

	// 处理每个消息的结果
	for i, message := range messages {
		if i < len(successFlags) {
			if successFlags[i] {
				successCount++
				r.logger.Debug("消息生产成功",
					logging.Field{Key: "message_id", Value: message.MessageID})
			} else {
				failureCount++
				errMsg := "unknown error"
				if i < len(errorMessages) && errorMessages[i] != "" {
					errMsg = errorMessages[i]
				}
				r.logger.Error("消息生产失败",
					logging.Field{Key: "message_id", Value: message.MessageID},
					logging.Field{Key: "error", Value: errMsg})
			}
		}
	}

	// 如果有失败的消息，返回错误
	if failureCount > 0 {
		return fmt.Errorf("batch produce partially failed: %d success, %d failure", successCount, failureCount)
	}

	r.logger.Info("批量消息生产成功",
		logging.Field{Key: "total", Value: len(messages)},
		logging.Field{Key: "success_count", Value: successCount},
		logging.Field{Key: "failure_count", Value: failureCount})

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
			Value:     []byte(pbMessage.GetMessage()), // 转换为字节数组
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

// matchesFilters 检查消息是否匹配过滤条件
func (r *GRPCMessageRepository) matchesFilters(message *entities.Message, filters []*valueobjects.FilterCondition) bool {
	if len(filters) == 0 {
		return true
	}

	// 检查所有过滤条件（AND逻辑）
	for _, filter := range filters {
		if !r.matchesFilter(message, filter) {
			return false
		}
	}
	return true
}

// matchesFilter 检查消息是否匹配单个过滤条件
func (r *GRPCMessageRepository) matchesFilter(message *entities.Message, filter *valueobjects.FilterCondition) bool {
	var targetValue string

	switch filter.Type {
	case valueobjects.FilterTypeKey:
		targetValue = message.Key
	case valueobjects.FilterTypeValue:
		targetValue = string(message.Value)
	case valueobjects.FilterTypeHeader:
		if filter.Field == "" {
			return false
		}
		targetValue = message.Headers[filter.Field]
	default:
		return false
	}

	return r.compareValues(targetValue, filter.Operator, filter.Value)
}

// compareValues 比较值
func (r *GRPCMessageRepository) compareValues(target string, operator valueobjects.FilterOperator, value string) bool {
	switch operator {
	case valueobjects.FilterOperatorEq:
		return target == value
	case valueobjects.FilterOperatorNe:
		return target != value
	case valueobjects.FilterOperatorContains:
		return r.contains(target, value)
	case valueobjects.FilterOperatorGt:
		return target > value
	case valueobjects.FilterOperatorGte:
		return target >= value
	case valueobjects.FilterOperatorLt:
		return target < value
	case valueobjects.FilterOperatorLte:
		return target <= value
	default:
		return false
	}
}

// contains 简单的字符串包含检查
func (r *GRPCMessageRepository) contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ConsumeStream 流式消费消息
func (r *GRPCMessageRepository) ConsumeStream(ctx context.Context, topic string, partition int32, offset int64) (<-chan *entities.Message, error) {
	r.logger.Debug("开始流式消费",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition},
		logging.Field{Key: "offset", Value: offset})

	// 创建流式消费请求
	req := &pb.StreamConsumeRequest{
		Topic:     topic,
		Partition: partition,
		Offset:    offset,
	}

	// 建立gRPC流
	stream, err := r.client.StreamConsume(ctx, req)
	if err != nil {
		r.logger.Error("建立流式消费失败", logging.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	// 创建消息通道
	messageChan := make(chan *entities.Message, 100) // 缓冲通道，支持背压控制

	// 启动goroutine处理流式数据
	go func() {
		defer close(messageChan)
		defer func() {
			if err := stream.CloseSend(); err != nil {
				r.logger.Error("关闭流失败", logging.Field{Key: "error", Value: err})
			}
		}()

		for {
			select {
			case <-ctx.Done():
				r.logger.Debug("流式消费被取消")
				return
			default:
				// 接收消息
				pbMessage, err := stream.Recv()
				if err != nil {
					if err.Error() == "EOF" {
						r.logger.Debug("流式消费结束")
						return
					}
					r.logger.Error("接收流式消息失败", logging.Field{Key: "error", Value: err})
					return
				}

				// 转换为实体
				message := &entities.Message{
					ID:        pbMessage.GetMessageId(),
					MessageID: pbMessage.GetMessageId(),
					Topic:     topic,
					Key:       pbMessage.GetKey(),
					Value:     []byte(pbMessage.GetMessage()),
					Headers:   pbMessage.GetHeaders(),
					Partition: pbMessage.GetPartition(),
					Offset:    pbMessage.GetOffset(),
					Timestamp: time.Unix(pbMessage.GetTimestamp(), 0),
				}

				// 发送到通道（支持背压控制）
				select {
				case messageChan <- message:
					// 消息发送成功
				case <-ctx.Done():
					r.logger.Debug("流式消费被取消")
					return
				}
			}
		}
	}()

	r.logger.Info("流式消费已启动",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition})

	return messageChan, nil
}

// GetOffset 获取偏移量
func (r *GRPCMessageRepository) GetOffset(ctx context.Context, topic string, partition int32, consumerGroup string) (int64, error) {
	r.logger.Debug("获取偏移量",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition},
		logging.Field{Key: "consumer_group", Value: consumerGroup})

	// 注意：当前protobuf定义中没有GetOffset方法
	// 这里使用DescribeConsumerGroup来获取偏移量信息
	req := &pb.DescribeConsumerGroupRequest{
		GroupId: consumerGroup,
	}

	resp, err := r.client.DescribeConsumerGroup(ctx, req)
	if err != nil {
		r.logger.Error("获取消费者组信息失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "group", Value: consumerGroup})
		return 0, fmt.Errorf("failed to describe consumer group: %w", err)
	}

	// 从响应中查找对应主题和分区的偏移量
	for _, offsetInfo := range resp.GetOffsets() {
		if offsetInfo.GetTopic() == topic && offsetInfo.GetPartition() == partition {
			offset := offsetInfo.GetCommittedOffset()
			r.logger.Debug("获取偏移量成功",
				logging.Field{Key: "offset", Value: offset})
			return offset, nil
		}
	}

	// 如果没有找到，返回0（从头开始消费）
	r.logger.Debug("未找到偏移量信息，返回0")
	return 0, nil
}

// CommitOffset 提交偏移量
func (r *GRPCMessageRepository) CommitOffset(ctx context.Context, topic string, partition int32, consumerGroup string, offset int64) error {
	r.logger.Debug("提交偏移量",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition},
		logging.Field{Key: "consumer_group", Value: consumerGroup},
		logging.Field{Key: "offset", Value: offset})

	// 构建提交偏移量请求
	req := &pb.CommitOffsetRequest{
		Topic:     topic,
		Partition: partition,
		Group:     consumerGroup,
		Offset:    offset,
	}

	// 调用gRPC服务
	resp, err := r.client.CommitOffset(ctx, req)
	if err != nil {
		r.logger.Error("提交偏移量失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: topic},
			logging.Field{Key: "offset", Value: offset})
		return fmt.Errorf("failed to commit offset: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("偏移量提交被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "topic", Value: topic})
		return fmt.Errorf("commit offset failed: %s", errMsg)
	}

	r.logger.Info("偏移量提交成功",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "partition", Value: partition},
		logging.Field{Key: "offset", Value: offset})

	return nil
}
