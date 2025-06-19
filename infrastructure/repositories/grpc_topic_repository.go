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
	// 简化实现
	return nil
}

// Delete 删除主题（实体接口）
func (r *GRPCTopicRepository) Delete(ctx context.Context, name string) error {
	r.logger.Debug("Deleting topic entity", logging.Field{Key: "name", Value: name})
	// 简化实现
	return nil
}

// Exists 检查主题是否存在
func (r *GRPCTopicRepository) Exists(ctx context.Context, name string) (bool, error) {
	r.logger.Debug("Checking topic existence", logging.Field{Key: "name", Value: name})
	// 简化实现：总是返回false
	return false, nil
}

// List 列出主题（实体接口）
func (r *GRPCTopicRepository) List(ctx context.Context) ([]*entities.Topic, error) {
	r.logger.Debug("Listing topic entities")
	// 简化实现：返回空列表
	return []*entities.Topic{}, nil
}

// GetByName 根据名称获取主题
func (r *GRPCTopicRepository) GetByName(ctx context.Context, name string) (*entities.Topic, error) {
	r.logger.Debug("Getting topic by name", logging.Field{Key: "name", Value: name})
	// 简化实现：返回nil
	return nil, nil
}

// GetDetail 获取主题详情
func (r *GRPCTopicRepository) GetDetail(ctx context.Context, name string) (*entities.Topic, error) {
	r.logger.Debug("Getting topic detail", logging.Field{Key: "name", Value: name})
	// 简化实现：返回nil
	return nil, nil
}

// GetStats 获取主题统计
func (r *GRPCTopicRepository) GetStats(ctx context.Context, name string) (*repositories.TopicStats, error) {
	r.logger.Debug("Getting topic stats", logging.Field{Key: "name", Value: name})
	// 简化实现：返回nil
	return nil, nil
}

// GetPartitionStats 获取分区统计
func (r *GRPCTopicRepository) GetPartitionStats(ctx context.Context, name string, partition int32) (*repositories.PartitionStats, error) {
	r.logger.Debug("Getting partition stats",
		logging.Field{Key: "name", Value: name},
		logging.Field{Key: "partition", Value: partition})
	// 简化实现：返回nil
	return nil, nil
}
