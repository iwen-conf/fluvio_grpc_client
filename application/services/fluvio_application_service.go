package services

import (
	"context"
	
	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/usecases"
)

// FluvioApplicationService Fluvio应用服务
type FluvioApplicationService struct {
	produceMessageUC *usecases.ProduceMessageUseCase
	consumeMessageUC *usecases.ConsumeMessageUseCase
	manageTopicUC    *usecases.ManageTopicUseCase
}

// NewFluvioApplicationService 创建Fluvio应用服务
func NewFluvioApplicationService(
	produceMessageUC *usecases.ProduceMessageUseCase,
	consumeMessageUC *usecases.ConsumeMessageUseCase,
	manageTopicUC *usecases.ManageTopicUseCase,
) *FluvioApplicationService {
	return &FluvioApplicationService{
		produceMessageUC: produceMessageUC,
		consumeMessageUC: consumeMessageUC,
		manageTopicUC:    manageTopicUC,
	}
}

// ProduceMessage 生产消息
func (s *FluvioApplicationService) ProduceMessage(ctx context.Context, req *dtos.ProduceMessageRequest) (*dtos.ProduceMessageResponse, error) {
	return s.produceMessageUC.Execute(ctx, req)
}

// ProduceBatch 批量生产消息
func (s *FluvioApplicationService) ProduceBatch(ctx context.Context, req *dtos.ProduceBatchRequest) (*dtos.ProduceBatchResponse, error) {
	return s.produceMessageUC.ExecuteBatch(ctx, req)
}

// ConsumeMessage 消费消息
func (s *FluvioApplicationService) ConsumeMessage(ctx context.Context, req *dtos.ConsumeMessageRequest) (*dtos.ConsumeMessageResponse, error) {
	return s.consumeMessageUC.Execute(ctx, req)
}

// ConsumeFiltered 过滤消费消息
func (s *FluvioApplicationService) ConsumeFiltered(ctx context.Context, req *dtos.FilteredConsumeRequest) (*dtos.FilteredConsumeResponse, error) {
	return s.consumeMessageUC.ExecuteFiltered(ctx, req)
}

// CreateTopic 创建主题
func (s *FluvioApplicationService) CreateTopic(ctx context.Context, req *dtos.CreateTopicRequest) (*dtos.CreateTopicResponse, error) {
	return s.manageTopicUC.CreateTopic(ctx, req)
}

// DeleteTopic 删除主题
func (s *FluvioApplicationService) DeleteTopic(ctx context.Context, req *dtos.DeleteTopicRequest) (*dtos.DeleteTopicResponse, error) {
	return s.manageTopicUC.DeleteTopic(ctx, req)
}

// ListTopics 列出主题
func (s *FluvioApplicationService) ListTopics(ctx context.Context) (*dtos.ListTopicsResponse, error) {
	return s.manageTopicUC.ListTopics(ctx)
}

// GetTopicDetail 获取主题详情
func (s *FluvioApplicationService) GetTopicDetail(ctx context.Context, name string) (*dtos.TopicDetailResponse, error) {
	return s.manageTopicUC.GetTopicDetail(ctx, name)
}

// GetTopicStats 获取主题统计
func (s *FluvioApplicationService) GetTopicStats(ctx context.Context, req *dtos.TopicStatsRequest) (*dtos.TopicStatsResponse, error) {
	return s.manageTopicUC.GetTopicStats(ctx, req)
}