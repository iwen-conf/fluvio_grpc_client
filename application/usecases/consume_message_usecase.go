package usecases

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/domain/services"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
)

// ConsumeMessageUseCase 消费消息用例
type ConsumeMessageUseCase struct {
	messageRepo    repositories.MessageRepository
	messageService *services.MessageService
}

// NewConsumeMessageUseCase 创建消费消息用例
func NewConsumeMessageUseCase(
	messageRepo repositories.MessageRepository,
	messageService *services.MessageService,
) *ConsumeMessageUseCase {
	return &ConsumeMessageUseCase{
		messageRepo:    messageRepo,
		messageService: messageService,
	}
}

// Execute 执行消费消息
func (uc *ConsumeMessageUseCase) Execute(ctx context.Context, req *dtos.ConsumeMessageRequest) (*dtos.ConsumeMessageResponse, error) {
	// 消费消息
	messages, err := uc.messageRepo.Consume(ctx, req.Topic, req.Partition, req.Offset, req.MaxMessages, req.Group)
	if err != nil {
		return &dtos.ConsumeMessageResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// 转换为DTO
	messageDTOs := make([]*dtos.MessageDTO, len(messages))
	for i, message := range messages {
		messageDTOs[i] = uc.entityToDTO(message)
	}

	return &dtos.ConsumeMessageResponse{
		Messages: messageDTOs,
		Count:    len(messageDTOs),
		Success:  true,
	}, nil
}

// ExecuteFiltered 执行过滤消费
func (uc *ConsumeMessageUseCase) ExecuteFiltered(ctx context.Context, req *dtos.FilteredConsumeRequest) (*dtos.FilteredConsumeResponse, error) {
	// 转换过滤条件
	filters := make([]*valueobjects.FilterCondition, len(req.Filters))
	for i, filterDTO := range req.Filters {
		filters[i] = &valueobjects.FilterCondition{
			Type:     valueobjects.FilterType(filterDTO.Type),
			Field:    filterDTO.Field,
			Operator: valueobjects.FilterOperator(filterDTO.Operator),
			Value:    filterDTO.Value,
		}
	}

	// 过滤消费
	messages, err := uc.messageRepo.ConsumeFiltered(ctx, req.Topic, filters, req.MaxMessages)
	if err != nil {
		return &dtos.FilteredConsumeResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// 应用过滤逻辑
	filteredMessages := make([]*entities.Message, 0)
	totalScanned := len(messages)

	for _, message := range messages {
		if uc.messageService.ApplyFilters(message, filters, req.AndLogic) {
			filteredMessages = append(filteredMessages, message)
		}
	}

	// 转换为DTO
	messageDTOs := make([]*dtos.MessageDTO, len(filteredMessages))
	for i, message := range filteredMessages {
		messageDTOs[i] = uc.entityToDTO(message)
	}

	return &dtos.FilteredConsumeResponse{
		Messages:      messageDTOs,
		FilteredCount: len(filteredMessages),
		TotalScanned:  totalScanned,
		Success:       true,
	}, nil
}

// entityToDTO 将实体转换为DTO
func (uc *ConsumeMessageUseCase) entityToDTO(message *entities.Message) *dtos.MessageDTO {
	return &dtos.MessageDTO{
		ID:        message.ID,
		MessageID: message.MessageID,
		Topic:     message.Topic,
		Key:       message.Key,
		Value:     string(message.Value), // 转换为字符串
		Headers:   message.Headers,
		Partition: message.Partition,
		Offset:    message.Offset,
		Timestamp: message.Timestamp,
	}
}
