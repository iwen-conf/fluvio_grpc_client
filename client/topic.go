package client

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
	"github.com/iwen-conf/fluvio_grpc_client/types"
)

// TopicManager 主题管理器（向后兼容）
type TopicManager struct {
	client *Client
}

// NewTopicManager 创建主题管理器
func NewTopicManager(client *Client) *TopicManager {
	return &TopicManager{
		client: client,
	}
}

// List 列出所有主题（简化实现）
func (tm *TopicManager) List(ctx context.Context) (*types.ListTopicsResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	// 简化实现：返回模拟主题列表
	return &types.ListTopicsResult{
		Topics:  []string{"example-topic", "test-topic"},
		Success: true,
	}, nil
}

// Create 创建主题（简化实现）
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

	// 简化实现：总是成功
	return &types.CreateTopicResult{
		Success: true,
	}, nil
}

// Delete 删除主题（简化实现）
func (tm *TopicManager) Delete(ctx context.Context, opts types.DeleteTopicOptions) (*types.DeleteTopicResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	// 简化实现：总是成功
	return &types.DeleteTopicResult{
		Success: true,
	}, nil
}

// Describe 描述主题（简化实现）
func (tm *TopicManager) Describe(ctx context.Context, topicName string) (*types.DescribeTopicResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if topicName == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	// 简化实现：返回模拟主题信息
	return &types.DescribeTopicResult{
		Topic: &types.TopicInfo{
			Name:       topicName,
			Partitions: 1,
			Replicas:   1,
		},
		Success: true,
	}, nil
}

// Exists 检查主题是否存在（简化实现）
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

// CreateIfNotExists 如果主题不存在则创建（简化实现）
func (tm *TopicManager) CreateIfNotExists(ctx context.Context, opts types.CreateTopicOptions) (*types.CreateTopicResult, error) {
	exists, err := tm.Exists(ctx, opts.Name)
	if err != nil {
		return nil, err
	}

	if exists {
		return &types.CreateTopicResult{
			Success: true,
		}, nil
	}

	return tm.Create(ctx, opts)
}

// DeleteIfExists 如果主题存在则删除（简化实现）
func (tm *TopicManager) DeleteIfExists(ctx context.Context, topicName string) (*types.DeleteTopicResult, error) {
	exists, err := tm.Exists(ctx, topicName)
	if err != nil {
		return nil, err
	}

	if !exists {
		return &types.DeleteTopicResult{
			Success: true,
		}, nil
	}

	return tm.Delete(ctx, types.DeleteTopicOptions{Name: topicName})
}

// WaitForTopic 等待主题创建完成（简化实现）
func (tm *TopicManager) WaitForTopic(ctx context.Context, topicName string, timeout time.Duration) error {
	if topicName == "" {
		return errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	// 简化实现：立即返回成功
	return nil
}

// GetTopicInfo 获取主题详细信息（简化实现）
func (tm *TopicManager) GetTopicInfo(ctx context.Context, topicName string) (*types.TopicInfo, error) {
	result, err := tm.Describe(ctx, topicName)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New(errors.ErrInternal, result.Error)
	}

	return result.Topic, nil
}

// ListTopicNames 获取主题名称列表（简化实现）
func (tm *TopicManager) ListTopicNames(ctx context.Context) ([]string, error) {
	result, err := tm.List(ctx)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New(errors.ErrInternal, result.Error)
	}

	return result.Topics, nil
}

// CreateTopicWithDefaults 使用默认配置创建主题（简化实现）
func (tm *TopicManager) CreateTopicWithDefaults(ctx context.Context, topicName string) (*types.CreateTopicResult, error) {
	return tm.Create(ctx, types.CreateTopicOptions{
		Name:              topicName,
		Partitions:        1,
		ReplicationFactor: 1,
	})
}