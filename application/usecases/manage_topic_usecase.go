package usecases

import (
	"context"
	
	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/domain/services"
)

// ManageTopicUseCase 主题管理用例
type ManageTopicUseCase struct {
	topicRepo    repositories.TopicRepository
	topicService *services.TopicService
}

// NewManageTopicUseCase 创建主题管理用例
func NewManageTopicUseCase(
	topicRepo repositories.TopicRepository,
	topicService *services.TopicService,
) *ManageTopicUseCase {
	return &ManageTopicUseCase{
		topicRepo:    topicRepo,
		topicService: topicService,
	}
}

// CreateTopic 创建主题
func (uc *ManageTopicUseCase) CreateTopic(ctx context.Context, req *dtos.CreateTopicRequest) (*dtos.CreateTopicResponse, error) {
	// 验证主题名称
	if err := uc.topicService.ValidateTopicName(req.Name); err != nil {
		return &dtos.CreateTopicResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 创建主题实体
	topic := entities.NewTopic(req.Name, req.Partitions)
	
	if req.ReplicationFactor > 0 {
		topic.WithReplicationFactor(req.ReplicationFactor)
	}
	
	if req.RetentionMs > 0 {
		topic.WithRetention(req.RetentionMs)
	}
	
	if req.Config != nil {
		topic.WithConfig(req.Config)
	}
	
	topic.Description = req.Description
	
	// 验证主题配置
	if err := uc.topicService.ValidateTopicConfig(topic); err != nil {
		return &dtos.CreateTopicResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 创建主题
	if err := uc.topicRepo.Create(ctx, topic); err != nil {
		return &dtos.CreateTopicResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	return &dtos.CreateTopicResponse{
		Success: true,
	}, nil
}

// DeleteTopic 删除主题
func (uc *ManageTopicUseCase) DeleteTopic(ctx context.Context, req *dtos.DeleteTopicRequest) (*dtos.DeleteTopicResponse, error) {
	// 验证主题名称
	if err := uc.topicService.ValidateTopicName(req.Name); err != nil {
		return &dtos.DeleteTopicResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 删除主题
	if err := uc.topicRepo.Delete(ctx, req.Name); err != nil {
		return &dtos.DeleteTopicResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	return &dtos.DeleteTopicResponse{
		Success: true,
	}, nil
}

// ListTopics 列出主题
func (uc *ManageTopicUseCase) ListTopics(ctx context.Context) (*dtos.ListTopicsResponse, error) {
	topics, err := uc.topicRepo.List(ctx)
	if err != nil {
		return &dtos.ListTopicsResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	topicNames := make([]string, len(topics))
	for i, topic := range topics {
		topicNames[i] = topic.Name
	}
	
	return &dtos.ListTopicsResponse{
		Topics:  topicNames,
		Count:   len(topicNames),
		Success: true,
	}, nil
}

// GetTopicDetail 获取主题详情
func (uc *ManageTopicUseCase) GetTopicDetail(ctx context.Context, name string) (*dtos.TopicDetailResponse, error) {
	// 验证主题名称
	if err := uc.topicService.ValidateTopicName(name); err != nil {
		return &dtos.TopicDetailResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 获取主题详情
	topic, err := uc.topicRepo.GetDetail(ctx, name)
	if err != nil {
		return &dtos.TopicDetailResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 转换为DTO
	topicDTO := uc.entityToDTO(topic)
	
	return &dtos.TopicDetailResponse{
		Topic:   topicDTO,
		Success: true,
	}, nil
}

// GetTopicStats 获取主题统计
func (uc *ManageTopicUseCase) GetTopicStats(ctx context.Context, req *dtos.TopicStatsRequest) (*dtos.TopicStatsResponse, error) {
	// 验证主题名称
	if err := uc.topicService.ValidateTopicName(req.Topic); err != nil {
		return &dtos.TopicStatsResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 获取主题统计
	stats, err := uc.topicRepo.GetStats(ctx, req.Topic)
	if err != nil {
		return &dtos.TopicStatsResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	// 转换为DTO
	statsDTO := &dtos.TopicStatsDTO{
		Topic:              stats.Topic,
		TotalMessageCount:  stats.TotalMessageCount,
		TotalSizeBytes:     stats.TotalSizeBytes,
		PartitionCount:     stats.PartitionCount,
	}
	
	if req.IncludePartitions && len(stats.PartitionStats) > 0 {
		partitionDTOs := make([]*dtos.PartitionStatsDTO, len(stats.PartitionStats))
		for i, partStats := range stats.PartitionStats {
			partitionDTOs[i] = &dtos.PartitionStatsDTO{
				PartitionID:      partStats.PartitionID,
				MessageCount:     partStats.MessageCount,
				TotalSizeBytes:   partStats.TotalSizeBytes,
				HighWatermark:    partStats.HighWatermark,
				LowWatermark:     partStats.LowWatermark,
			}
		}
		statsDTO.Partitions = partitionDTOs
	}
	
	return &dtos.TopicStatsResponse{
		Topics:  []*dtos.TopicStatsDTO{statsDTO},
		Success: true,
	}, nil
}

// entityToDTO 将实体转换为DTO
func (uc *ManageTopicUseCase) entityToDTO(topic *entities.Topic) *dtos.TopicDTO {
	dto := &dtos.TopicDTO{
		Name:              topic.Name,
		Description:       topic.Description,
		Partitions:        topic.Partitions,
		ReplicationFactor: topic.ReplicationFactor,
		RetentionMs:       topic.RetentionMs,
		Config:            topic.Config,
		CreatedAt:         topic.CreatedAt,
		UpdatedAt:         topic.UpdatedAt,
	}
	
	if len(topic.PartitionDetails) > 0 {
		partitionDTOs := make([]*dtos.PartitionInfoDTO, len(topic.PartitionDetails))
		for i, partition := range topic.PartitionDetails {
			partitionDTOs[i] = &dtos.PartitionInfoDTO{
				PartitionID:    partition.PartitionID,
				LeaderID:       partition.LeaderID,
				ReplicaIDs:     partition.ReplicaIDs,
				HighWatermark:  partition.HighWatermark,
				LowWatermark:   partition.LowWatermark,
				MessageCount:   partition.MessageCount,
				TotalSizeBytes: partition.TotalSizeBytes,
			}
		}
		dto.PartitionDetails = partitionDTOs
	}
	
	return dto
}