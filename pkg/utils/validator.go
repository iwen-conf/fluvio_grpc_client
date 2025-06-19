package utils

import (
	"fmt"
	"strings"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// Validator 验证器
type Validator struct{}

// NewValidator 创建验证器
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateProduceMessageRequest 验证生产消息请求
func (v *Validator) ValidateProduceMessageRequest(req *dtos.ProduceMessageRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "produce message request cannot be nil")
	}

	if strings.TrimSpace(req.Topic) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic cannot be empty")
	}

	if strings.TrimSpace(req.Value) == "" {
		return errors.New(errors.ErrInvalidArgument, "message value cannot be empty")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Topic); err != nil {
		return err
	}

	return nil
}

// ValidateConsumeMessageRequest 验证消费消息请求
func (v *Validator) ValidateConsumeMessageRequest(req *dtos.ConsumeMessageRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "consume message request cannot be nil")
	}

	if strings.TrimSpace(req.Topic) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic cannot be empty")
	}

	if req.Partition < 0 {
		return errors.New(errors.ErrInvalidArgument, "partition cannot be negative")
	}

	if req.Offset < 0 {
		return errors.New(errors.ErrInvalidArgument, "offset cannot be negative")
	}

	if req.MaxMessages <= 0 {
		return errors.New(errors.ErrInvalidArgument, "max messages must be positive")
	}

	if req.MaxMessages > 10000 {
		return errors.New(errors.ErrInvalidArgument, "max messages cannot exceed 10000")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Topic); err != nil {
		return err
	}

	return nil
}

// ValidateCreateTopicRequest 验证创建主题请求
func (v *Validator) ValidateCreateTopicRequest(req *dtos.CreateTopicRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "create topic request cannot be nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot be empty")
	}

	if req.Partitions <= 0 {
		return errors.New(errors.ErrInvalidArgument, "partitions must be positive")
	}

	if req.Partitions > 1000 {
		return errors.New(errors.ErrInvalidArgument, "partitions cannot exceed 1000")
	}

	if req.ReplicationFactor <= 0 {
		return errors.New(errors.ErrInvalidArgument, "replication factor must be positive")
	}

	if req.ReplicationFactor > 10 {
		return errors.New(errors.ErrInvalidArgument, "replication factor cannot exceed 10")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Name); err != nil {
		return err
	}

	return nil
}

// ValidateDeleteTopicRequest 验证删除主题请求
func (v *Validator) ValidateDeleteTopicRequest(req *dtos.DeleteTopicRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "delete topic request cannot be nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot be empty")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Name); err != nil {
		return err
	}

	return nil
}

// ValidateCreateSmartModuleRequest 验证创建SmartModule请求
func (v *Validator) ValidateCreateSmartModuleRequest(req *dtos.CreateSmartModuleRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "create smart module request cannot be nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot be empty")
	}

	if len(req.WasmCode) == 0 {
		return errors.New(errors.ErrInvalidArgument, "wasm code cannot be empty")
	}

	// 验证SmartModule名称格式
	if err := v.ValidateSmartModuleName(req.Name); err != nil {
		return err
	}

	return nil
}

// ValidateDeleteSmartModuleRequest 验证删除SmartModule请求
func (v *Validator) ValidateDeleteSmartModuleRequest(req *dtos.DeleteSmartModuleRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "delete smart module request cannot be nil")
	}

	if strings.TrimSpace(req.Name) == "" {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot be empty")
	}

	// 验证SmartModule名称格式
	if err := v.ValidateSmartModuleName(req.Name); err != nil {
		return err
	}

	return nil
}

// ValidateTopicName 验证主题名称格式
func (v *Validator) ValidateTopicName(name string) error {
	name = strings.TrimSpace(name)
	
	if len(name) == 0 {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot be empty")
	}

	if len(name) > 255 {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot exceed 255 characters")
	}

	// 检查是否包含非法字符
	for _, char := range name {
		if !v.isValidTopicChar(char) {
			return errors.New(errors.ErrInvalidArgument, 
				fmt.Sprintf("topic name contains invalid character: %c", char))
		}
	}

	// 不能以点开头或结尾
	if strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot start or end with '.'")
	}

	// 不能包含连续的点
	if strings.Contains(name, "..") {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot contain consecutive dots")
	}

	return nil
}

// ValidateSmartModuleName 验证SmartModule名称格式
func (v *Validator) ValidateSmartModuleName(name string) error {
	name = strings.TrimSpace(name)
	
	if len(name) == 0 {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot be empty")
	}

	if len(name) > 100 {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot exceed 100 characters")
	}

	// 检查是否包含非法字符（只允许字母、数字、连字符和下划线）
	for _, char := range name {
		if !v.isValidSmartModuleChar(char) {
			return errors.New(errors.ErrInvalidArgument, 
				fmt.Sprintf("smart module name contains invalid character: %c", char))
		}
	}

	// 不能以连字符开头或结尾
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot start or end with '-'")
	}

	return nil
}

// ValidateMessage 验证消息实体
func (v *Validator) ValidateMessage(message *entities.Message) error {
	if message == nil {
		return errors.New(errors.ErrInvalidArgument, "message cannot be nil")
	}

	if strings.TrimSpace(message.Topic) == "" {
		return errors.New(errors.ErrInvalidArgument, "message topic cannot be empty")
	}

	if len(message.Value) == 0 {
		return errors.New(errors.ErrInvalidArgument, "message value cannot be empty")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(message.Topic); err != nil {
		return err
	}

	return nil
}

// ValidateBatchMessages 验证批量消息
func (v *Validator) ValidateBatchMessages(messages []*entities.Message) error {
	if len(messages) == 0 {
		return errors.New(errors.ErrInvalidArgument, "message batch cannot be empty")
	}

	if len(messages) > 1000 {
		return errors.New(errors.ErrInvalidArgument, "message batch cannot exceed 1000 messages")
	}

	// 验证所有消息都是同一个主题
	firstTopic := messages[0].Topic
	for i, message := range messages {
		if err := v.ValidateMessage(message); err != nil {
			return errors.Wrap(errors.ErrInvalidArgument, 
				fmt.Sprintf("message at index %d is invalid", i), err)
		}

		if message.Topic != firstTopic {
			return errors.New(errors.ErrInvalidArgument, 
				fmt.Sprintf("all messages must have the same topic, got %s and %s", 
					firstTopic, message.Topic))
		}
	}

	return nil
}

// isValidTopicChar 检查字符是否为有效的主题字符
func (v *Validator) isValidTopicChar(char rune) bool {
	// 允许字母、数字、点、连字符和下划线
	return (char >= 'a' && char <= 'z') ||
		   (char >= 'A' && char <= 'Z') ||
		   (char >= '0' && char <= '9') ||
		   char == '.' ||
		   char == '-' ||
		   char == '_'
}

// isValidSmartModuleChar 检查字符是否为有效的SmartModule字符
func (v *Validator) isValidSmartModuleChar(char rune) bool {
	// 允许字母、数字、连字符和下划线
	return (char >= 'a' && char <= 'z') ||
		   (char >= 'A' && char <= 'Z') ||
		   (char >= '0' && char <= '9') ||
		   char == '-' ||
		   char == '_'
}

// ValidateConsumerGroup 验证消费者组名称
func (v *Validator) ValidateConsumerGroup(groupName string) error {
	groupName = strings.TrimSpace(groupName)
	
	if len(groupName) == 0 {
		return errors.New(errors.ErrInvalidArgument, "consumer group name cannot be empty")
	}

	if len(groupName) > 255 {
		return errors.New(errors.ErrInvalidArgument, "consumer group name cannot exceed 255 characters")
	}

	// 使用与主题名称相同的验证规则
	return v.ValidateTopicName(groupName)
}
