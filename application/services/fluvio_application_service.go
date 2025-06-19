package services

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
)

// FluvioApplicationService Fluvio应用服务
type FluvioApplicationService struct {
	messageRepo repositories.MessageRepository
	topicRepo   repositories.TopicRepository
	adminRepo   repositories.AdminRepository
	logger      logging.Logger
}

// NewFluvioApplicationService 创建Fluvio应用服务
func NewFluvioApplicationService(
	messageRepo repositories.MessageRepository,
	topicRepo repositories.TopicRepository,
	adminRepo repositories.AdminRepository,
	logger logging.Logger,
) *FluvioApplicationService {
	return &FluvioApplicationService{
		messageRepo: messageRepo,
		topicRepo:   topicRepo,
		adminRepo:   adminRepo,
		logger:      logger,
	}
}

// ProduceMessage 生产消息
func (s *FluvioApplicationService) ProduceMessage(ctx context.Context, req *dtos.ProduceMessageRequest) (*dtos.ProduceMessageResponse, error) {
	s.logger.Debug("Producing message", logging.Field{Key: "topic", Value: req.Topic})

	// 创建消息实体
	message := entities.NewMessage(req.Key, req.Value)
	message.Topic = req.Topic

	if req.MessageID != "" {
		message.WithMessageID(req.MessageID)
	}

	if req.Headers != nil {
		message.WithHeaders(req.Headers)
	}

	// 调用仓储层进行实际的消息生产
	if err := s.messageRepo.Produce(ctx, message); err != nil {
		s.logger.Error("Failed to produce message",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "topic", Value: req.Topic})
		return &dtos.ProduceMessageResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	s.logger.Info("Message produced successfully",
		logging.Field{Key: "topic", Value: message.Topic},
		logging.Field{Key: "message_id", Value: message.MessageID})

	return &dtos.ProduceMessageResponse{
		MessageID: message.MessageID,
		Topic:     message.Topic,
		Partition: message.Partition,
		Offset:    message.Offset,
		Success:   true,
	}, nil
}

// ConsumeMessage 消费消息
func (s *FluvioApplicationService) ConsumeMessage(ctx context.Context, req *dtos.ConsumeMessageRequest) (*dtos.ConsumeMessageResponse, error) {
	s.logger.Debug("Consuming messages", logging.Field{Key: "topic", Value: req.Topic})

	// 简化实现：返回模拟消息
	return &dtos.ConsumeMessageResponse{
		Messages: []*dtos.MessageDTO{
			{
				Key:       "test-key",
				Value:     "Hello from " + req.Topic,
				Headers:   map[string]string{"source": "mock"},
				Offset:    0,
				Partition: 0,
			},
		},
	}, nil
}

// CreateTopic 创建主题
func (s *FluvioApplicationService) CreateTopic(ctx context.Context, req *dtos.CreateTopicRequest) (*dtos.CreateTopicResponse, error) {
	return s.topicRepo.CreateTopic(ctx, req)
}

// DeleteTopic 删除主题
func (s *FluvioApplicationService) DeleteTopic(ctx context.Context, req *dtos.DeleteTopicRequest) (*dtos.DeleteTopicResponse, error) {
	return s.topicRepo.DeleteTopic(ctx, req)
}

// ListTopics 列出主题
func (s *FluvioApplicationService) ListTopics(ctx context.Context, req *dtos.ListTopicsRequest) (*dtos.ListTopicsResponse, error) {
	return s.topicRepo.ListTopics(ctx, req)
}

// ListTopicsSimple 简单列出主题（向后兼容）
func (s *FluvioApplicationService) ListTopicsSimple(ctx context.Context) (*dtos.ListTopicsResponse, error) {
	return s.topicRepo.ListTopics(ctx, &dtos.ListTopicsRequest{})
}

// DescribeTopic 描述主题
func (s *FluvioApplicationService) DescribeTopic(ctx context.Context, req *dtos.DescribeTopicRequest) (*dtos.DescribeTopicResponse, error) {
	return s.topicRepo.DescribeTopic(ctx, req)
}

// 管理功能

// DescribeCluster 描述集群
func (s *FluvioApplicationService) DescribeCluster(ctx context.Context, req *dtos.DescribeClusterRequest) (*dtos.DescribeClusterResponse, error) {
	return s.adminRepo.DescribeCluster(ctx, req)
}

// ListBrokers 列出Broker
func (s *FluvioApplicationService) ListBrokers(ctx context.Context, req *dtos.ListBrokersRequest) (*dtos.ListBrokersResponse, error) {
	return s.adminRepo.ListBrokers(ctx, req)
}

// ListConsumerGroups 列出消费者组
func (s *FluvioApplicationService) ListConsumerGroups(ctx context.Context, req *dtos.ListConsumerGroupsRequest) (*dtos.ListConsumerGroupsResponse, error) {
	return s.adminRepo.ListConsumerGroups(ctx, req)
}

// DescribeConsumerGroup 描述消费者组
func (s *FluvioApplicationService) DescribeConsumerGroup(ctx context.Context, req *dtos.DescribeConsumerGroupRequest) (*dtos.DescribeConsumerGroupResponse, error) {
	return s.adminRepo.DescribeConsumerGroup(ctx, req)
}

// ListSmartModules 列出SmartModule
func (s *FluvioApplicationService) ListSmartModules(ctx context.Context, req *dtos.ListSmartModulesRequest) (*dtos.ListSmartModulesResponse, error) {
	return s.adminRepo.ListSmartModules(ctx, req)
}

// CreateSmartModule 创建SmartModule
func (s *FluvioApplicationService) CreateSmartModule(ctx context.Context, req *dtos.CreateSmartModuleRequest) (*dtos.CreateSmartModuleResponse, error) {
	return s.adminRepo.CreateSmartModule(ctx, req)
}

// DeleteSmartModule 删除SmartModule
func (s *FluvioApplicationService) DeleteSmartModule(ctx context.Context, req *dtos.DeleteSmartModuleRequest) (*dtos.DeleteSmartModuleResponse, error) {
	return s.adminRepo.DeleteSmartModule(ctx, req)
}

// DescribeSmartModule 描述SmartModule
func (s *FluvioApplicationService) DescribeSmartModule(ctx context.Context, req *dtos.DescribeSmartModuleRequest) (*dtos.DescribeSmartModuleResponse, error) {
	return s.adminRepo.DescribeSmartModule(ctx, req)
}
