package client

import (
	"context"
	
	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/interfaces/api"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// TopicAdapter 主题适配器
type TopicAdapter struct {
	appService *services.FluvioApplicationService
	connected  *bool
}

// Create 创建主题
func (t *TopicAdapter) Create(ctx context.Context, opts api.CreateTopicOptions) (*api.CreateTopicResult, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	req := &dtos.CreateTopicRequest{
		Name:              opts.Name,
		Partitions:        opts.Partitions,
		ReplicationFactor: opts.ReplicationFactor,
		RetentionMs:       opts.RetentionMs,
		Config:            opts.Config,
	}
	
	resp, err := t.appService.CreateTopic(ctx, req)
	if err != nil {
		return &api.CreateTopicResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	return &api.CreateTopicResult{
		Success: resp.Success,
		Error:   resp.Error,
	}, nil
}

// Delete 删除主题
func (t *TopicAdapter) Delete(ctx context.Context, opts api.DeleteTopicOptions) (*api.DeleteTopicResult, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	req := &dtos.DeleteTopicRequest{
		Name: opts.Name,
	}
	
	resp, err := t.appService.DeleteTopic(ctx, req)
	if err != nil {
		return &api.DeleteTopicResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	return &api.DeleteTopicResult{
		Success: resp.Success,
		Error:   resp.Error,
	}, nil
}

// Exists 检查主题是否存在
func (t *TopicAdapter) Exists(ctx context.Context, name string) (bool, error) {
	if !*t.connected {
		return false, errors.New(errors.ErrConnection, "client not connected")
	}
	
	// 通过列出主题来检查是否存在
	listResp, err := t.appService.ListTopics(ctx)
	if err != nil {
		return false, err
	}
	
	for _, topic := range listResp.Topics {
		if topic == name {
			return true, nil
		}
	}
	
	return false, nil
}

// CreateIfNotExists 如果不存在则创建主题
func (t *TopicAdapter) CreateIfNotExists(ctx context.Context, opts api.CreateTopicOptions) (*api.CreateTopicResult, error) {
	exists, err := t.Exists(ctx, opts.Name)
	if err != nil {
		return nil, err
	}
	
	if exists {
		return &api.CreateTopicResult{
			Success: true,
		}, nil
	}
	
	return t.Create(ctx, opts)
}

// List 列出主题
func (t *TopicAdapter) List(ctx context.Context) (*api.ListTopicsResult, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	resp, err := t.appService.ListTopics(ctx)
	if err != nil {
		return nil, err
	}
	
	return &api.ListTopicsResult{
		Topics: resp.Topics,
	}, nil
}

// Describe 描述主题
func (t *TopicAdapter) Describe(ctx context.Context, name string) (*api.TopicDescription, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	detail, err := t.DescribeTopicDetail(ctx, name)
	if err != nil {
		return nil, err
	}
	
	return &api.TopicDescription{
		Name:       detail.Topic,
		Partitions: int32(len(detail.Partitions)),
	}, nil
}

// DescribeTopicDetail 获取主题详情
func (t *TopicAdapter) DescribeTopicDetail(ctx context.Context, name string) (*api.TopicDetail, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	resp, err := t.appService.GetTopicDetail(ctx, name)
	if err != nil {
		return nil, err
	}
	
	if !resp.Success {
		return nil, errors.New(errors.ErrInternal, resp.Error)
	}
	
	// 转换分区信息
	partitions := make([]*api.PartitionInfo, len(resp.Topic.PartitionDetails))
	for i, partition := range resp.Topic.PartitionDetails {
		partitions[i] = &api.PartitionInfo{
			PartitionID:   partition.PartitionID,
			LeaderID:      partition.LeaderID,
			HighWatermark: partition.HighWatermark,
		}
	}
	
	return &api.TopicDetail{
		Topic:       resp.Topic.Name,
		Partitions:  partitions,
		RetentionMs: resp.Topic.RetentionMs,
		Config:      resp.Topic.Config,
	}, nil
}

// GetTopicStats 获取主题统计
func (t *TopicAdapter) GetTopicStats(ctx context.Context, opts api.GetTopicStatsOptions) (*api.GetTopicStatsResult, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}
	
	req := &dtos.TopicStatsRequest{
		Topic:             opts.Topic,
		IncludePartitions: opts.IncludePartitions,
	}
	
	resp, err := t.appService.GetTopicStats(ctx, req)
	if err != nil {
		return nil, err
	}
	
	if !resp.Success {
		return nil, errors.New(errors.ErrInternal, resp.Error)
	}
	
	// 转换统计信息
	topics := make([]*api.TopicStats, len(resp.Topics))
	for i, topicStats := range resp.Topics {
		// 转换分区统计
		var partitions []*api.PartitionStats
		if topicStats.Partitions != nil {
			partitions = make([]*api.PartitionStats, len(topicStats.Partitions))
			for j, partStats := range topicStats.Partitions {
				partitions[j] = &api.PartitionStats{
					PartitionID:    partStats.PartitionID,
					MessageCount:   partStats.MessageCount,
					TotalSizeBytes: partStats.TotalSizeBytes,
				}
			}
		}
		
		topics[i] = &api.TopicStats{
			Topic:              topicStats.Topic,
			TotalMessageCount:  topicStats.TotalMessageCount,
			TotalSizeBytes:     topicStats.TotalSizeBytes,
			PartitionCount:     topicStats.PartitionCount,
			Partitions:         partitions,
		}
	}
	
	return &api.GetTopicStatsResult{
		Topics: topics,
	}, nil
}