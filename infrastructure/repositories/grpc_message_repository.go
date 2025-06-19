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
	"github.com/iwen-conf/fluvio_grpc_client/pkg/utils"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// GRPCMessageRepository gRPC消息仓储实现
type GRPCMessageRepository struct {
	client    grpc.Client
	logger    logging.Logger
	handler   *utils.GRPCResponseHandler
	converter *utils.DTOConverter
	validator *utils.Validator
}

// NewGRPCMessageRepository 创建gRPC消息仓储
func NewGRPCMessageRepository(client grpc.Client, logger logging.Logger) repositories.MessageRepository {
	return &GRPCMessageRepository{
		client:    client,
		logger:    logger,
		handler:   utils.NewGRPCResponseHandler(logger),
		converter: utils.NewDTOConverter(),
		validator: utils.NewValidator(),
	}
}

// Produce 生产消息
func (r *GRPCMessageRepository) Produce(ctx context.Context, message *entities.Message) error {
	// 验证消息
	if err := r.validator.ValidateMessage(message); err != nil {
		return err
	}

	// 记录调试日志
	context := utils.NewContextBuilder().
		Add("topic", message.Topic).
		Add("key", message.Key).
		Add("message_id", message.MessageID).
		Build()
	r.handler.LogDebugOperation("生产消息", context)

	// 转换为gRPC请求
	req := r.converter.MessageEntityToProtoRequest(message)

	// 调用gRPC服务
	resp, err := r.client.Produce(ctx, req)
	if err != nil {
		return r.handler.HandleError(err, "生产消息", context)
	}

	// 验证响应
	if err := r.handler.ValidateResponse(resp.GetSuccess(), resp.GetError(), "生产消息", context); err != nil {
		return err
	}

	// 更新消息元数据
	message.ID = resp.GetMessageId()
	message.MessageID = resp.GetMessageId()
	// 注意：当前protobuf定义中ProduceReply没有Partition和Offset字段
	// 这里使用默认值，实际应该从服务器响应中获取
	message.Partition = 0
	message.Offset = 0

	// 记录成功日志
	successContext := utils.NewContextBuilder().
		Add("message_id", message.MessageID).
		Add("topic", message.Topic).
		Build()
	r.handler.HandleSuccessResponse("消息生产", successContext)

	return nil
}

// ProduceBatch 批量生产消息
func (r *GRPCMessageRepository) ProduceBatch(ctx context.Context, messages []*entities.Message) error {
	// 验证批量消息
	if err := r.validator.ValidateBatchMessages(messages); err != nil {
		return err
	}

	// 记录调试日志
	context := utils.NewContextBuilder().
		Add("count", len(messages)).
		Add("topic", messages[0].Topic).
		Build()
	r.handler.LogDebugOperation("批量生产消息", context)

	// 转换为protobuf请求
	pbMessages := r.converter.MessageEntitiesToProtoRequests(messages)

	// 构建批量生产请求
	req := &pb.BatchProduceRequest{
		Topic:    messages[0].Topic,
		Messages: pbMessages,
	}

	// 调用gRPC服务
	resp, err := r.client.BatchProduce(ctx, req)
	if err != nil {
		return r.handler.HandleError(err, "批量生产消息", context)
	}

	// 处理批量响应
	result := utils.NewBatchOperationResult()
	successFlags := resp.GetSuccess()
	errorMessages := resp.GetError()

	// 处理每个消息的结果
	for i, message := range messages {
		if i < len(successFlags) {
			if successFlags[i] {
				result.AddSuccess()
				r.logger.Debug("消息生产成功",
					logging.Field{Key: "message_id", Value: message.MessageID})
			} else {
				errMsg := "unknown error"
				if i < len(errorMessages) && errorMessages[i] != "" {
					errMsg = errorMessages[i]
				}
				result.AddFailure(fmt.Errorf("message %d failed: %s", i, errMsg))
				r.logger.Error("消息生产失败",
					logging.Field{Key: "message_id", Value: message.MessageID},
					logging.Field{Key: "error", Value: errMsg})
			}
		}
	}

	// 记录汇总日志
	result.LogSummary(r.handler, "批量消息生产", context)

	// 如果有失败的消息，返回错误
	return result.GetSummaryError()
}

// Consume 消费消息
func (r *GRPCMessageRepository) Consume(ctx context.Context, topic string, partition int32, offset int64, maxMessages int) ([]*entities.Message, error) {
	// 记录调试日志
	context := utils.NewContextBuilder().
		Add("topic", topic).
		Add("partition", partition).
		Add("offset", offset).
		Add("max_messages", maxMessages).
		Build()
	r.handler.LogDebugOperation("消费消息", context)

	req := &pb.ConsumeRequest{
		Topic:       topic,
		Partition:   partition,
		Offset:      offset,
		MaxMessages: int32(maxMessages),
	}

	// 调用gRPC服务
	resp, err := r.client.Consume(ctx, req)
	if err != nil {
		return nil, r.handler.HandleError(err, "消费消息", context)
	}

	// 转换为实体
	messages := r.converter.ConsumedMessagesToEntities(resp.GetMessages())

	// 设置主题（因为protobuf消息中可能没有主题信息）
	for _, message := range messages {
		if message.Topic == "" {
			message.Topic = topic
		}
	}

	// 记录成功日志
	successContext := utils.NewContextBuilder().
		Add("topic", topic).
		Add("count", len(messages)).
		Build()
	r.handler.HandleSuccessResponse("消息消费", successContext)

	return messages, nil
}

// ConsumeFiltered 过滤消费消息
func (r *GRPCMessageRepository) ConsumeFiltered(ctx context.Context, topic string, filters []*valueobjects.FilterCondition, maxMessages int) ([]*entities.Message, error) {
	r.logger.Debug("过滤消费消息",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "filters", Value: len(filters)},
		logging.Field{Key: "max_messages", Value: maxMessages})

	// 如果没有过滤条件，直接调用普通消费
	if len(filters) == 0 {
		return r.Consume(ctx, topic, 0, 0, maxMessages)
	}

	// 客户端过滤实现：分批获取消息直到满足过滤条件
	var filteredMessages []*entities.Message
	batchSize := maxMessages * 2 // 初始批次大小
	offset := int64(0)
	totalScanned := 0
	maxScanned := maxMessages * 10 // 最大扫描消息数，避免无限循环

	for len(filteredMessages) < maxMessages && totalScanned < maxScanned {
		// 获取一批消息
		batch, err := r.Consume(ctx, topic, 0, offset, batchSize)
		if err != nil {
			return nil, fmt.Errorf("failed to consume batch for filtering: %w", err)
		}

		// 如果没有更多消息，退出循环
		if len(batch) == 0 {
			break
		}

		// 应用过滤条件
		for _, message := range batch {
			totalScanned++
			if r.matchesFilters(message, filters) {
				filteredMessages = append(filteredMessages, message)
				if len(filteredMessages) >= maxMessages {
					break
				}
			}
		}

		// 更新偏移量
		if len(batch) > 0 {
			offset = batch[len(batch)-1].Offset + 1
		}

		// 如果这批消息少于请求的数量，说明已经到达末尾
		if len(batch) < batchSize {
			break
		}
	}

	r.logger.Info("过滤消费完成",
		logging.Field{Key: "topic", Value: topic},
		logging.Field{Key: "filtered_count", Value: len(filteredMessages)},
		logging.Field{Key: "total_scanned", Value: totalScanned},
		logging.Field{Key: "filter_efficiency", Value: float64(len(filteredMessages)) / float64(totalScanned) * 100})

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
	case valueobjects.FilterOperatorRegex:
		return r.matchesRegex(target, value)
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

// matchesRegex 正则表达式匹配
func (r *GRPCMessageRepository) matchesRegex(target, pattern string) bool {
	// 简化的正则表达式实现，只支持基本的模式匹配
	// 在生产环境中，应该使用完整的正则表达式库

	// 支持通配符 * 和 ?
	return r.simplePatternMatch(target, pattern)
}

// simplePatternMatch 简单的模式匹配（支持 * 和 ?）
func (r *GRPCMessageRepository) simplePatternMatch(text, pattern string) bool {
	if pattern == "*" {
		return true
	}

	// 递归匹配
	return r.matchPattern(text, pattern, 0, 0)
}

// matchPattern 递归模式匹配
func (r *GRPCMessageRepository) matchPattern(text, pattern string, textIdx, patternIdx int) bool {
	// 如果模式已经匹配完
	if patternIdx >= len(pattern) {
		return textIdx >= len(text)
	}

	// 如果文本已经匹配完但模式还有非*字符
	if textIdx >= len(text) {
		for i := patternIdx; i < len(pattern); i++ {
			if pattern[i] != '*' {
				return false
			}
		}
		return true
	}

	// 处理通配符
	if pattern[patternIdx] == '*' {
		// 尝试匹配0个或多个字符
		for i := textIdx; i <= len(text); i++ {
			if r.matchPattern(text, pattern, i, patternIdx+1) {
				return true
			}
		}
		return false
	}

	// 处理单字符通配符或精确匹配
	if pattern[patternIdx] == '?' || pattern[patternIdx] == text[textIdx] {
		return r.matchPattern(text, pattern, textIdx+1, patternIdx+1)
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
