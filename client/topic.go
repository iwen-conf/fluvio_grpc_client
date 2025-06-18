package client

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
	"github.com/iwen-conf/fluvio_grpc_client/types"

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
		Topic:             opts.Name,
		Partitions:        opts.Partitions,
		ReplicationFactor: opts.ReplicationFactor,
		RetentionMs:       opts.RetentionMs,
		Config:            opts.Config,
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
		if resp.GetTopic() != "" {
			topicInfo = &types.TopicInfo{
				Name:       resp.GetTopic(),
				Partitions: int32(len(resp.GetPartitions())), // 分区数量
				Replicas:   1,                                // 默认副本数
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

// DeleteIfExists 如果主题存在则删除
func (tm *TopicManager) DeleteIfExists(ctx context.Context, topicName string) (*types.DeleteTopicResult, error) {
	exists, err := tm.Exists(ctx, topicName)
	if err != nil {
		return nil, err
	}

	if !exists {
		tm.client.logger.Info("主题不存在，跳过删除", logger.Field{Key: "name", Value: topicName})
		return &types.DeleteTopicResult{
			Success: true,
		}, nil
	}

	return tm.Delete(ctx, types.DeleteTopicOptions{Name: topicName})
}

// WaitForTopic 等待主题创建完成
func (tm *TopicManager) WaitForTopic(ctx context.Context, topicName string, timeout time.Duration) error {
	if topicName == "" {
		return errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return errors.New(errors.ErrTimeout, "等待主题创建超时")
			}

			exists, err := tm.Exists(ctx, topicName)
			if err != nil {
				tm.client.logger.Debug("检查主题存在性失败",
					logger.Field{Key: "topic", Value: topicName},
					logger.Field{Key: "error", Value: err})
				continue
			}

			if exists {
				tm.client.logger.Info("主题已就绪", logger.Field{Key: "topic", Value: topicName})
				return nil
			}
		}
	}
}

// GetTopicInfo 获取主题详细信息（便民方法）
func (tm *TopicManager) GetTopicInfo(ctx context.Context, topicName string) (*types.TopicInfo, error) {
	result, err := tm.Describe(ctx, topicName)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, errors.New(errors.ErrTopicNotFound, result.Error)
	}

	return result.Topic, nil
}

// ListTopicNames 获取主题名称列表（便民方法）
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

// CreateTopicWithDefaults 使用默认配置创建主题
func (tm *TopicManager) CreateTopicWithDefaults(ctx context.Context, topicName string) (*types.CreateTopicResult, error) {
	return tm.Create(ctx, types.CreateTopicOptions{
		Name:              topicName,
		Partitions:        1,
		ReplicationFactor: 1,
	})
}

// TopicManager的高级功能
type AdvancedTopicManager struct {
	topicManager *TopicManager
	client       *Client
}

// NewAdvancedTopicManager 创建高级主题管理器
func NewAdvancedTopicManager(client *Client) *AdvancedTopicManager {
	return &AdvancedTopicManager{
		topicManager: NewTopicManager(client),
		client:       client,
	}
}

// BatchCreateTopics 批量创建主题
func (atm *AdvancedTopicManager) BatchCreateTopics(ctx context.Context, topics []types.CreateTopicOptions) (map[string]*types.CreateTopicResult, error) {
	results := make(map[string]*types.CreateTopicResult)

	for _, topic := range topics {
		result, err := atm.topicManager.Create(ctx, topic)
		if err != nil {
			results[topic.Name] = &types.CreateTopicResult{
				Success: false,
				Error:   err.Error(),
			}
		} else {
			results[topic.Name] = result
		}
	}

	return results, nil
}

// BatchDeleteTopics 批量删除主题
func (atm *AdvancedTopicManager) BatchDeleteTopics(ctx context.Context, topicNames []string) (map[string]*types.DeleteTopicResult, error) {
	results := make(map[string]*types.DeleteTopicResult)

	for _, topicName := range topicNames {
		result, err := atm.topicManager.Delete(ctx, types.DeleteTopicOptions{Name: topicName})
		if err != nil {
			results[topicName] = &types.DeleteTopicResult{
				Success: false,
				Error:   err.Error(),
			}
		} else {
			results[topicName] = result
		}
	}

	return results, nil
}

// EnsureTopics 确保主题存在，不存在则创建
func (atm *AdvancedTopicManager) EnsureTopics(ctx context.Context, topics []types.CreateTopicOptions) error {
	for _, topic := range topics {
		_, err := atm.topicManager.CreateIfNotExists(ctx, topic)
		if err != nil {
			return fmt.Errorf("确保主题 %s 存在失败: %w", topic.Name, err)
		}
	}
	return nil
}

// CleanupTopics 清理匹配模式的主题
func (atm *AdvancedTopicManager) CleanupTopics(ctx context.Context, pattern string) error {
	topics, err := atm.topicManager.ListTopicNames(ctx)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		if matched, _ := filepath.Match(pattern, topic); matched {
			_, err := atm.topicManager.Delete(ctx, types.DeleteTopicOptions{Name: topic})
			if err != nil {
				atm.client.logger.Error("删除主题失败",
					logger.Field{Key: "topic", Value: topic},
					logger.Field{Key: "error", Value: err})
			} else {
				atm.client.logger.Info("主题已删除", logger.Field{Key: "topic", Value: topic})
			}
		}
	}

	return nil
}

// GetTopicsWithPrefix 获取指定前缀的主题
func (atm *AdvancedTopicManager) GetTopicsWithPrefix(ctx context.Context, prefix string) ([]string, error) {
	topics, err := atm.topicManager.ListTopicNames(ctx)
	if err != nil {
		return nil, err
	}

	var filteredTopics []string
	for _, topic := range topics {
		if strings.HasPrefix(topic, prefix) {
			filteredTopics = append(filteredTopics, topic)
		}
	}

	return filteredTopics, nil
}

// ValidateTopicName 验证主题名称
func (atm *AdvancedTopicManager) ValidateTopicName(topicName string) error {
	if topicName == "" {
		return errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if len(topicName) > 255 {
		return errors.New(errors.ErrInvalidArgument, "主题名称长度不能超过255个字符")
	}

	// 检查非法字符
	for _, char := range topicName {
		if char < 32 || char > 126 {
			return errors.New(errors.ErrInvalidArgument, "主题名称包含非法字符")
		}
	}

	// 检查保留名称
	reservedNames := []string{".", "..", "__consumer_offsets", "__transaction_state"}
	for _, reserved := range reservedNames {
		if topicName == reserved {
			return errors.New(errors.ErrInvalidArgument, "主题名称不能使用保留名称")
		}
	}

	return nil
}

// CreateTopicWithValidation 创建主题并验证名称
func (atm *AdvancedTopicManager) CreateTopicWithValidation(ctx context.Context, opts types.CreateTopicOptions) (*types.CreateTopicResult, error) {
	if err := atm.ValidateTopicName(opts.Name); err != nil {
		return &types.CreateTopicResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return atm.topicManager.Create(ctx, opts)
}

// DescribeTopicDetail 获取主题详细信息（新版本）
func (tm *TopicManager) DescribeTopicDetail(ctx context.Context, topicName string) (*types.DescribeTopicDetailResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if topicName == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	req := &pb.DescribeTopicRequest{
		Topic: topicName,
	}

	tm.client.logger.Debug("获取主题详细信息", logger.Field{Key: "name", Value: topicName})

	var result *types.DescribeTopicDetailResult
	err := tm.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.DescribeTopic(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取主题详细信息失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换分区信息
		var partitions []*types.PartitionInfo
		for _, p := range resp.GetPartitions() {
			partition := &types.PartitionInfo{
				PartitionID:    p.GetPartitionId(),
				LeaderID:       p.GetLeaderId(),
				ReplicaIDs:     p.GetReplicaIds(),
				ISRIDs:         p.GetIsrIds(),
				HighWatermark:  p.GetHighWatermark(),
				LogStartOffset: p.GetLogStartOffset(),
			}
			partitions = append(partitions, partition)
		}

		result = &types.DescribeTopicDetailResult{
			Topic:       resp.GetTopic(),
			RetentionMs: resp.GetRetentionMs(),
			Config:      resp.GetConfig(),
			Partitions:  partitions,
			Success:     true,
		}

		return nil
	})

	if err != nil {
		tm.client.logger.Error("获取主题详细信息失败",
			logger.Field{Key: "name", Value: topicName},
			logger.Field{Key: "error", Value: err})
		return &types.DescribeTopicDetailResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	tm.client.logger.Info("主题详细信息获取成功", logger.Field{Key: "name", Value: topicName})

	return result, nil
}

// GetTopicStats 获取主题统计信息
func (tm *TopicManager) GetTopicStats(ctx context.Context, opts types.GetTopicStatsOptions) (*types.GetTopicStatsResult, error) {
	if tm.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	req := &pb.GetTopicStatsRequest{
		Topic:             opts.Topic,
		IncludePartitions: opts.IncludePartitions,
	}

	tm.client.logger.Debug("获取主题统计信息",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "include_partitions", Value: opts.IncludePartitions})

	var result *types.GetTopicStatsResult
	err := tm.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.GetTopicStats(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取主题统计信息失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换主题统计信息
		var topics []*types.TopicStats
		for _, t := range resp.GetTopics() {
			topicStats := &types.TopicStats{
				Topic:             t.GetTopic(),
				PartitionCount:    t.GetPartitionCount(),
				ReplicationFactor: t.GetReplicationFactor(),
				TotalMessageCount: t.GetTotalMessageCount(),
				TotalSizeBytes:    t.GetTotalSizeBytes(),
			}

			if t.GetCreatedAt() != nil {
				topicStats.CreatedAt = t.GetCreatedAt().AsTime()
			}
			if t.GetLastUpdated() != nil {
				topicStats.LastUpdated = t.GetLastUpdated().AsTime()
			}

			// 转换分区统计信息
			if opts.IncludePartitions {
				var partitions []*types.PartitionStats
				for _, p := range t.GetPartitions() {
					partitionStats := &types.PartitionStats{
						PartitionID:    p.GetPartitionId(),
						MessageCount:   p.GetMessageCount(),
						TotalSizeBytes: p.GetTotalSizeBytes(),
						EarliestOffset: p.GetEarliestOffset(),
						LatestOffset:   p.GetLatestOffset(),
					}

					if p.GetLastUpdated() != nil {
						partitionStats.LastUpdated = p.GetLastUpdated().AsTime()
					}

					partitions = append(partitions, partitionStats)
				}
				topicStats.Partitions = partitions
			}

			topics = append(topics, topicStats)
		}

		result = &types.GetTopicStatsResult{
			Topics:  topics,
			Success: true,
		}

		if resp.GetCollectedAt() != nil {
			result.CollectedAt = resp.GetCollectedAt().AsTime()
		}

		return nil
	})

	if err != nil {
		tm.client.logger.Error("获取主题统计信息失败",
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "error", Value: err})
		return &types.GetTopicStatsResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	tm.client.logger.Info("主题统计信息获取成功",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "topics_count", Value: len(result.Topics)})

	return result, nil
}
