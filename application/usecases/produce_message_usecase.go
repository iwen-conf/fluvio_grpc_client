package usecases

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/domain/services"
)

// ProduceMessageUseCase 生产消息用例
type ProduceMessageUseCase struct {
	messageRepo    repositories.MessageRepository
	messageService *services.MessageService
}

// NewProduceMessageUseCase 创建生产消息用例
func NewProduceMessageUseCase(
	messageRepo repositories.MessageRepository,
	messageService *services.MessageService,
) *ProduceMessageUseCase {
	return &ProduceMessageUseCase{
		messageRepo:    messageRepo,
		messageService: messageService,
	}
}

// Execute 执行生产消息
func (uc *ProduceMessageUseCase) Execute(ctx context.Context, req *dtos.ProduceMessageRequest) (*dtos.ProduceMessageResponse, error) {
	// 创建消息实体
	message := entities.NewMessage(req.Key, req.Value)
	message.Topic = req.Topic

	if req.MessageID != "" {
		message.WithMessageID(req.MessageID)
	}

	if req.Headers != nil {
		message.WithHeaders(req.Headers)
	}

	// 验证消息
	if err := uc.messageService.ValidateMessage(message); err != nil {
		return &dtos.ProduceMessageResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// 生产消息
	if err := uc.messageRepo.Produce(ctx, message); err != nil {
		return &dtos.ProduceMessageResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &dtos.ProduceMessageResponse{
		MessageID: message.MessageID,
		Topic:     message.Topic,
		Partition: message.Partition,
		Offset:    message.Offset,
		Success:   true,
	}, nil
}

// ExecuteBatch 执行批量生产消息
func (uc *ProduceMessageUseCase) ExecuteBatch(ctx context.Context, req *dtos.ProduceBatchRequest) (*dtos.ProduceBatchResponse, error) {
	// 转换为实体
	messages := make([]*entities.Message, len(req.Messages))
	for i, msgReq := range req.Messages {
		message := entities.NewMessage(msgReq.Key, msgReq.Value)
		message.Topic = msgReq.Topic

		if msgReq.MessageID != "" {
			message.WithMessageID(msgReq.MessageID)
		}

		if msgReq.Headers != nil {
			message.WithHeaders(msgReq.Headers)
		}

		messages[i] = message
	}

	// 验证批量消息
	if err := uc.messageService.ValidateBatch(messages); err != nil {
		return &dtos.ProduceBatchResponse{
			TotalMessages: len(messages),
			SuccessCount:  0,
			FailureCount:  len(messages),
			Results: []*dtos.ProduceMessageResponse{
				{
					Success: false,
					Error:   err.Error(),
				},
			},
		}, err
	}

	// 批量生产
	if err := uc.messageRepo.ProduceBatch(ctx, messages); err != nil {
		return &dtos.ProduceBatchResponse{
			TotalMessages: len(messages),
			SuccessCount:  0,
			FailureCount:  len(messages),
			Results: []*dtos.ProduceMessageResponse{
				{
					Success: false,
					Error:   err.Error(),
				},
			},
		}, err
	}

	// 构建响应
	results := make([]*dtos.ProduceMessageResponse, len(messages))
	successCount := 0

	for i, message := range messages {
		results[i] = &dtos.ProduceMessageResponse{
			MessageID: message.MessageID,
			Topic:     message.Topic,
			Partition: message.Partition,
			Offset:    message.Offset,
			Success:   true,
		}
		successCount++
	}

	return &dtos.ProduceBatchResponse{
		Results:       results,
		TotalMessages: len(messages),
		SuccessCount:  successCount,
		FailureCount:  len(messages) - successCount,
	}, nil
}
