package services

import (
	"context"
	"fmt"
	
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
)

// MessageService 消息领域服务
type MessageService struct{}

// NewMessageService 创建消息服务
func NewMessageService() *MessageService {
	return &MessageService{}
}

// ValidateMessage 验证消息
func (ms *MessageService) ValidateMessage(message *entities.Message) error {
	if !message.IsValid() {
		return fmt.Errorf("invalid message: value and topic are required")
	}
	
	// 检查消息大小限制（例如：1MB）
	if message.Size() > 1024*1024 {
		return fmt.Errorf("message size exceeds limit: %d bytes", message.Size())
	}
	
	return nil
}

// ValidateBatch 验证批量消息
func (ms *MessageService) ValidateBatch(messages []*entities.Message) error {
	if len(messages) == 0 {
		return fmt.Errorf("batch cannot be empty")
	}
	
	// 检查批量大小限制
	if len(messages) > 1000 {
		return fmt.Errorf("batch size exceeds limit: %d messages", len(messages))
	}
	
	// 验证每条消息
	for i, message := range messages {
		if err := ms.ValidateMessage(message); err != nil {
			return fmt.Errorf("message %d in batch is invalid: %w", i, err)
		}
	}
	
	return nil
}

// ApplyFilters 应用过滤条件
func (ms *MessageService) ApplyFilters(message *entities.Message, filters []*valueobjects.FilterCondition, andLogic bool) bool {
	if len(filters) == 0 {
		return true
	}
	
	results := make([]bool, len(filters))
	
	for i, filter := range filters {
		results[i] = ms.applyFilter(message, filter)
	}
	
	// 应用逻辑运算
	if andLogic {
		// AND逻辑：所有条件都必须满足
		for _, result := range results {
			if !result {
				return false
			}
		}
		return true
	} else {
		// OR逻辑：任一条件满足即可
		for _, result := range results {
			if result {
				return true
			}
		}
		return false
	}
}

// applyFilter 应用单个过滤条件
func (ms *MessageService) applyFilter(message *entities.Message, filter *valueobjects.FilterCondition) bool {
	var targetValue string
	
	switch filter.Type {
	case valueobjects.FilterTypeKey:
		targetValue = message.Key
	case valueobjects.FilterTypeValue:
		targetValue = message.Value
	case valueobjects.FilterTypeHeader:
		if filter.Field == "" {
			return false
		}
		targetValue = message.Headers[filter.Field]
	default:
		return false
	}
	
	return ms.compareValues(targetValue, filter.Operator, filter.Value)
}

// compareValues 比较值
func (ms *MessageService) compareValues(target string, operator valueobjects.FilterOperator, expected string) bool {
	switch operator {
	case valueobjects.FilterOperatorEq:
		return target == expected
	case valueobjects.FilterOperatorNe:
		return target != expected
	case valueobjects.FilterOperatorContains:
		return len(target) > 0 && len(expected) > 0 && 
			   target != expected && 
			   (target == expected || 
			    (len(target) > len(expected) && 
			     target[:len(expected)] == expected || 
			     target[len(target)-len(expected):] == expected ||
			     ms.contains(target, expected)))
	default:
		return false
	}
}

// contains 简单的字符串包含检查
func (ms *MessageService) contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}