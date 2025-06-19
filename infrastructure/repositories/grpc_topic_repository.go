package repositories

import (
	"context"
	"fmt"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/grpc"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// GRPCTopicRepository gRPC主题仓储实现
type GRPCTopicRepository struct {
	client grpc.Client
	logger logging.Logger
}

// NewGRPCTopicRepository 创建gRPC主题仓储
func NewGRPCTopicRepository(client grpc.Client, logger logging.Logger) repositories.TopicRepository {
	return &GRPCTopicRepository{
		client: client,
		logger: logger,
	}
}

// CreateTopic 创建主题（DTO接口）
func (r *GRPCTopicRepository) CreateTopic(ctx context.Context, req *dtos.CreateTopicRequest) (*dtos.CreateTopicResponse, error) {
	r.logger.Debug("Creating topic", logging.Field{Key: "name", Value: req.Name})

	// 构建gRPC请求
	grpcReq := &pb.CreateTopicRequest{
		Topic:             req.Name,
		Partitions:        req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
		Config:            req.Config,
	}

	// 调用gRPC服务
	resp, err := r.client.CreateTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("创建主题失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: req.Name})
		return nil, fmt.Errorf("failed to create topic: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("主题创建被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "topic", Value: req.Name})
		return &dtos.CreateTopicResponse{
			Success: false,
			Error:   errMsg,
		}, nil
	}

	r.logger.Info("主题创建成功", logging.Field{Key: "topic", Value: req.Name})

	return &dtos.CreateTopicResponse{
		Success: true,
	}, nil
}

// DeleteTopic 删除主题（DTO接口）
func (r *GRPCTopicRepository) DeleteTopic(ctx context.Context, req *dtos.DeleteTopicRequest) (*dtos.DeleteTopicResponse, error) {
	r.logger.Debug("Deleting topic", logging.Field{Key: "name", Value: req.Name})

	// 构建gRPC请求
	grpcReq := &pb.DeleteTopicRequest{
		Topic: req.Name,
	}

	// 调用gRPC服务
	resp, err := r.client.DeleteTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("删除主题失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: req.Name})
		return nil, fmt.Errorf("failed to delete topic: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("主题删除被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "topic", Value: req.Name})
		return &dtos.DeleteTopicResponse{
			Success: false,
			Error:   errMsg,
		}, nil
	}

	r.logger.Info("主题删除成功", logging.Field{Key: "topic", Value: req.Name})

	return &dtos.DeleteTopicResponse{
		Success: true,
	}, nil
}

// ListTopics 列出主题（DTO接口）
func (r *GRPCTopicRepository) ListTopics(ctx context.Context, req *dtos.ListTopicsRequest) (*dtos.ListTopicsResponse, error) {
	r.logger.Debug("Listing topics")

	// 构建gRPC请求
	grpcReq := &pb.ListTopicsRequest{}

	// 调用gRPC服务
	resp, err := r.client.ListTopics(ctx, grpcReq)
	if err != nil {
		r.logger.Error("列出主题失败", logging.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to list topics: %w", err)
	}

	r.logger.Debug("列出主题成功", logging.Field{Key: "count", Value: len(resp.GetTopics())})

	return &dtos.ListTopicsResponse{
		Topics: resp.GetTopics(),
	}, nil
}

// DescribeTopic 描述主题（DTO接口）
func (r *GRPCTopicRepository) DescribeTopic(ctx context.Context, req *dtos.DescribeTopicRequest) (*dtos.DescribeTopicResponse, error) {
	r.logger.Debug("Describing topic", logging.Field{Key: "name", Value: req.Name})

	// 构建gRPC请求
	grpcReq := &pb.DescribeTopicRequest{
		Topic: req.Name,
	}

	// 调用gRPC服务
	resp, err := r.client.DescribeTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("描述主题失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: req.Name})
		return nil, fmt.Errorf("failed to describe topic: %w", err)
	}

	// 检查错误
	if resp.GetError() != "" {
		return &dtos.DescribeTopicResponse{
			Error: resp.GetError(),
		}, nil
	}

	r.logger.Debug("描述主题成功", logging.Field{Key: "topic", Value: req.Name})

	return &dtos.DescribeTopicResponse{
		Topic: &dtos.TopicDTO{
			Name:       resp.GetTopic(),
			Partitions: int32(len(resp.GetPartitions())), // 从分区列表计算分区数
			Config:     resp.GetConfig(),
		},
	}, nil
}

// 实体接口实现

// Create 创建主题（实体接口）
func (r *GRPCTopicRepository) Create(ctx context.Context, topic *entities.Topic) error {
	r.logger.Debug("Creating topic entity", logging.Field{Key: "name", Value: topic.Name})

	// 构建gRPC请求
	grpcReq := &pb.CreateTopicRequest{
		Topic:             topic.Name,
		Partitions:        topic.Partitions,
		ReplicationFactor: topic.ReplicationFactor,
		Config:            topic.Config,
	}

	// 调用gRPC服务
	resp, err := r.client.CreateTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("创建主题实体失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: topic.Name})
		return fmt.Errorf("failed to create topic entity: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("主题实体创建被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "topic", Value: topic.Name})
		return fmt.Errorf("create topic entity failed: %s", errMsg)
	}

	r.logger.Info("主题实体创建成功", logging.Field{Key: "topic", Value: topic.Name})
	return nil
}

// Delete 删除主题（实体接口）
func (r *GRPCTopicRepository) Delete(ctx context.Context, name string) error {
	r.logger.Debug("Deleting topic entity", logging.Field{Key: "name", Value: name})

	// 构建gRPC请求
	grpcReq := &pb.DeleteTopicRequest{
		Topic: name,
	}

	// 调用gRPC服务
	resp, err := r.client.DeleteTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("删除主题实体失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: name})
		return fmt.Errorf("failed to delete topic entity: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("主题实体删除被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "topic", Value: name})
		return fmt.Errorf("delete topic entity failed: %s", errMsg)
	}

	r.logger.Info("主题实体删除成功", logging.Field{Key: "topic", Value: name})
	return nil
}

// Exists 检查主题是否存在
func (r *GRPCTopicRepository) Exists(ctx context.Context, name string) (bool, error) {
	r.logger.Debug("Checking topic existence", logging.Field{Key: "name", Value: name})

	// 通过描述主题来检查是否存在
	grpcReq := &pb.DescribeTopicRequest{
		Topic: name,
	}

	resp, err := r.client.DescribeTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Debug("检查主题存在性失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: name})
		// 如果是网络错误等，返回错误
		return false, fmt.Errorf("failed to check topic existence: %w", err)
	}

	// 如果有错误信息，说明主题不存在
	if resp.GetError() != "" {
		r.logger.Debug("主题不存在", logging.Field{Key: "topic", Value: name})
		return false, nil
	}

	r.logger.Debug("主题存在", logging.Field{Key: "topic", Value: name})
	return true, nil
}

// List 列出主题（实体接口）
func (r *GRPCTopicRepository) List(ctx context.Context) ([]*entities.Topic, error) {
	r.logger.Debug("Listing topic entities")

	// 构建gRPC请求
	grpcReq := &pb.ListTopicsRequest{}

	// 调用gRPC服务
	resp, err := r.client.ListTopics(ctx, grpcReq)
	if err != nil {
		r.logger.Error("列出主题实体失败", logging.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to list topic entities: %w", err)
	}

	// 转换为实体
	topics := make([]*entities.Topic, len(resp.GetTopics()))
	for i, topicName := range resp.GetTopics() {
		// 为每个主题创建基本实体，详细信息需要单独查询
		topics[i] = &entities.Topic{
			Name: topicName,
		}
	}

	r.logger.Debug("列出主题实体成功", logging.Field{Key: "count", Value: len(topics)})
	return topics, nil
}

// GetByName 根据名称获取主题
func (r *GRPCTopicRepository) GetByName(ctx context.Context, name string) (*entities.Topic, error) {
	r.logger.Debug("Getting topic by name", logging.Field{Key: "name", Value: name})

	// 构建gRPC请求
	grpcReq := &pb.DescribeTopicRequest{
		Topic: name,
	}

	// 调用gRPC服务
	resp, err := r.client.DescribeTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("根据名称获取主题失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: name})
		return nil, fmt.Errorf("failed to get topic by name: %w", err)
	}

	// 检查错误
	if resp.GetError() != "" {
		r.logger.Debug("主题不存在", logging.Field{Key: "topic", Value: name})
		return nil, nil // 主题不存在时返回nil而不是错误
	}

	// 转换为实体
	topic := &entities.Topic{
		Name:              resp.GetTopic(),
		Partitions:        int32(len(resp.GetPartitions())),
		ReplicationFactor: 1, // 默认值，protobuf中没有这个字段
		Config:            resp.GetConfig(),
	}

	r.logger.Debug("根据名称获取主题成功", logging.Field{Key: "topic", Value: name})
	return topic, nil
}

// GetDetail 获取主题详情
func (r *GRPCTopicRepository) GetDetail(ctx context.Context, name string) (*entities.Topic, error) {
	r.logger.Debug("Getting topic detail", logging.Field{Key: "name", Value: name})

	// GetDetail和GetByName功能相同，直接调用
	return r.GetByName(ctx, name)
}

// GetStats 获取主题统计
func (r *GRPCTopicRepository) GetStats(ctx context.Context, name string) (*repositories.TopicStats, error) {
	r.logger.Debug("Getting topic stats", logging.Field{Key: "name", Value: name})

	// 首先获取主题详情
	topic, err := r.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic for stats: %w", err)
	}

	if topic == nil {
		return nil, fmt.Errorf("topic not found: %s", name)
	}

	// 构建主题统计信息
	stats := &repositories.TopicStats{
		Topic:          name,
		PartitionCount: topic.Partitions,
		PartitionStats: make([]*repositories.PartitionStats, topic.Partitions),
	}

	// 获取每个分区的统计信息
	var totalMessageCount int64
	var totalSizeBytes int64

	for i := int32(0); i < topic.Partitions; i++ {
		partitionStats, err := r.GetPartitionStats(ctx, name, i)
		if err != nil {
			r.logger.Error("获取分区统计失败",
				logging.Field{Key: "topic", Value: name},
				logging.Field{Key: "partition", Value: i},
				logging.Field{Key: "error", Value: err})
			// 继续处理其他分区，不中断整个过程
			partitionStats = &repositories.PartitionStats{
				PartitionID:    i,
				MessageCount:   0,
				TotalSizeBytes: 0,
				HighWatermark:  0,
				LowWatermark:   0,
			}
		}

		stats.PartitionStats[i] = partitionStats
		totalMessageCount += partitionStats.MessageCount
		totalSizeBytes += partitionStats.TotalSizeBytes
	}

	stats.TotalMessageCount = totalMessageCount
	stats.TotalSizeBytes = totalSizeBytes

	r.logger.Debug("获取主题统计成功",
		logging.Field{Key: "topic", Value: name},
		logging.Field{Key: "total_messages", Value: totalMessageCount})

	return stats, nil
}

// GetPartitionStats 获取分区统计
func (r *GRPCTopicRepository) GetPartitionStats(ctx context.Context, name string, partition int32) (*repositories.PartitionStats, error) {
	r.logger.Debug("Getting partition stats",
		logging.Field{Key: "name", Value: name},
		logging.Field{Key: "partition", Value: partition})

	// 注意：当前protobuf定义中没有专门的分区统计方法
	// 这里使用DescribeTopic来获取分区信息
	grpcReq := &pb.DescribeTopicRequest{
		Topic: name,
	}

	resp, err := r.client.DescribeTopic(ctx, grpcReq)
	if err != nil {
		r.logger.Error("获取分区统计失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: name},
			logging.Field{Key: "partition", Value: partition})
		return nil, fmt.Errorf("failed to get partition stats: %w", err)
	}

	// 检查错误
	if resp.GetError() != "" {
		return nil, fmt.Errorf("get partition stats failed: %s", resp.GetError())
	}

	// 从分区列表中查找对应的分区
	partitions := resp.GetPartitions()
	if int(partition) >= len(partitions) {
		return nil, fmt.Errorf("partition %d not found in topic %s", partition, name)
	}

	// 构建分区统计信息（基于可用的信息）
	stats := &repositories.PartitionStats{
		PartitionID:    partition,
		MessageCount:   0, // protobuf中没有这个信息，设为0
		TotalSizeBytes: 0, // protobuf中没有这个信息，设为0
		HighWatermark:  0, // protobuf中没有这个信息，设为0
		LowWatermark:   0, // protobuf中没有这个信息，设为0
	}

	r.logger.Debug("获取分区统计成功",
		logging.Field{Key: "topic", Value: name},
		logging.Field{Key: "partition", Value: partition})

	return stats, nil
}
