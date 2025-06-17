package client

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	"github.com/iwen-conf/fluvio_grpc_client/types"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"

	"google.golang.org/grpc"
)

// TopicManager 主题管理器
type TopicManager struct {
	client *Client
}

// NewTopicManager 创建主题管理器
func NewTopicManager(client *Client) *TopicManager {
	return &TopicManager{
		client: client,
	}
}

// List 列出所有主题
func (tm *TopicManager) List(ctx context.Context) (*types.ListTopicsResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	tm.client.logger.Debug("获取主题列表")

	var result *types.ListTopicsResult
	err := tm.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.ListTopics(ctx, &pb.ListTopicsRequest{})
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取主题列表失败", err)
		}

		result = &types.ListTopicsResult{
			Topics:  resp.GetTopics(),
			Success: true,
		}

		return nil
	})

	if err != nil {
		tm.client.logger.Error("获取主题列表失败", logger.Field{Key: "error", Value: err})
		return &types.ListTopicsResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	tm.client.logger.Info("获取主题列表成功", 
		logger.Field{Key: "count", Value: len(result.Topics)})

	return result, nil
}

// Create 创建主题
func (tm *TopicManager) Create(ctx context.Context, opts types.CreateTopicOptions) (*types.CreateTopicResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Partitions <= 0 {
		opts.Partitions = 1
	}

	req := &pb.CreateTopicRequest{
		Topic:      opts.Name,
		Partitions: opts.Partitions,
	}

	tm.client.logger.Debug("创建主题", 
		logger.Field{Key: "name", Value: opts.Name},
		logger.Field{Key: "partitions", Value: opts.Partitions})

	var result *types.CreateTopicResult
	err := tm.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.CreateTopic(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "创建主题失败", err)
		}

		result = &types.CreateTopicResult{
			Success: resp.GetSuccess(),
			Error:   resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		tm.client.logger.Error("创建主题失败", 
			logger.Field{Key: "name", Value: opts.Name},
			logger.Field{Key: "error", Value: err})
		return &types.CreateTopicResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	tm.client.logger.Info("主题创建成功", 
		logger.Field{Key: "name", Value: opts.Name},
		logger.Field{Key: "partitions", Value: opts.Partitions})

	return result, nil
}

// Delete 删除主题
func (tm *TopicManager) Delete(ctx context.Context, opts types.DeleteTopicOptions) (*types.DeleteTopicResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	req := &pb.DeleteTopicRequest{
		Topic: opts.Name,
	}

	tm.client.logger.Debug("删除主题", logger.Field{Key: "name", Value: opts.Name})

	var result *types.DeleteTopicResult
	err := tm.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.DeleteTopic(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "删除主题失败", err)
		}

		result = &types.DeleteTopicResult{
			Success: resp.GetSuccess(),
			Error:   resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		tm.client.logger.Error("删除主题失败", 
			logger.Field{Key: "name", Value: opts.Name},
			logger.Field{Key: "error", Value: err})
		return &types.DeleteTopicResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	tm.client.logger.Info("主题删除成功", logger.Field{Key: "name", Value: opts.Name})

	return result, nil
}

// Describe 描述主题
func (tm *TopicManager) Describe(ctx context.Context, topicName string) (*types.DescribeTopicResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if topicName == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	req := &pb.DescribeTopicRequest{
		Topic: topicName,
	}

	tm.client.logger.Debug("描述主题", logger.Field{Key: "name", Value: topicName})

	var result *types.DescribeTopicResult
	err := tm.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)
		
		resp, err := client.DescribeTopic(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "描述主题失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换主题信息
		var topicInfo *types.TopicInfo
		if resp.GetTopic() != nil {
			topicInfo = &types.TopicInfo{
				Name:       resp.GetTopic().GetName(),
				Partitions: resp.GetTopic().GetPartitions(),
				Replicas:   resp.GetTopic().GetReplicas(),
			}
		}

		result = &types.DescribeTopicResult{
			Topic:   topicInfo,
			Success: true,
		}

		return nil
	})

	if err != nil {
		tm.client.logger.Error("描述主题失败", 
			logger.Field{Key: "name", Value: topicName},
			logger.Field{Key: "error", Value: err})
		return &types.DescribeTopicResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	tm.client.logger.Info("主题描述成功", logger.Field{Key: "name", Value: topicName})

	return result, nil
}

// Exists 检查主题是否存在
func (tm *TopicManager) Exists(ctx context.Context, topicName string) (bool, error) {
	if topicName == "" {
		return false, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	listResult, err := tm.List(ctx)
	if err != nil {
		return false, err
	}

	for _, topic := range listResult.Topics {
		if topic == topicName {
			return true, nil
		}
	}

	return false, nil
}

// CreateIfNotExists 如果主题不存在则创建
func (tm *TopicManager) CreateIfNotExists(ctx context.Context, opts types.CreateTopicOptions) (*types.CreateTopicResult, error) {
	exists, err := tm.Exists(ctx, opts.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		tm.client.logger.Info("主题已存在，跳过创建", logger.Field{Key: "name", Value: opts.Name})
		return &types.CreateTopicResult{
			Success: true,
		}, nil
	}

	return tm.Create(ctx, opts)
}
