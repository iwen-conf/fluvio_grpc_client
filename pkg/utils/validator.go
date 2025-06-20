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

	// 验证主题名称
	if strings.TrimSpace(req.Topic) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic cannot be empty")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Topic); err != nil {
		return err
	}

	// 验证消息值
	if strings.TrimSpace(req.Value) == "" {
		return errors.New(errors.ErrInvalidArgument, "message value cannot be empty")
	}

	// 验证消息大小（假设最大1MB）
	const maxMessageSize = 1024 * 1024 // 1MB
	if len(req.Value) > maxMessageSize {
		return errors.New(errors.ErrInvalidArgument,
			fmt.Sprintf("message value too large: %d bytes (max %d bytes)", len(req.Value), maxMessageSize))
	}

	// 验证消息头部
	if req.Headers != nil {
		for key, value := range req.Headers {
			if strings.TrimSpace(key) == "" {
				return errors.New(errors.ErrInvalidArgument, "header key cannot be empty")
			}
			if len(key) > 255 {
				return errors.New(errors.ErrInvalidArgument, "header key too long (max 255 characters)")
			}
			if len(value) > 1024 {
				return errors.New(errors.ErrInvalidArgument, "header value too long (max 1024 characters)")
			}
		}
	}

	return nil
}

// ValidateConsumeMessageRequest 验证消费消息请求
func (v *Validator) ValidateConsumeMessageRequest(req *dtos.ConsumeMessageRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "consume message request cannot be nil")
	}

	// 验证主题名称
	if strings.TrimSpace(req.Topic) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic cannot be empty")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Topic); err != nil {
		return err
	}

	// 验证消费者组名称
	if strings.TrimSpace(req.Group) == "" {
		return errors.New(errors.ErrInvalidArgument, "consumer group cannot be empty")
	}

	if err := v.ValidateConsumerGroup(req.Group); err != nil {
		return err
	}

	// 验证分区号
	if req.Partition < 0 {
		return errors.New(errors.ErrInvalidArgument, "partition cannot be negative")
	}

	// 验证偏移量
	if req.Offset < 0 {
		return errors.New(errors.ErrInvalidArgument, "offset cannot be negative")
	}

	// 验证最大消息数
	if req.MaxMessages <= 0 {
		return errors.New(errors.ErrInvalidArgument, "max messages must be greater than 0")
	}

	if req.MaxMessages > 10000 {
		return errors.New(errors.ErrInvalidArgument, "max messages cannot exceed 10000")
	}

	return nil
}

// ValidateCreateTopicRequest 验证创建主题请求
func (v *Validator) ValidateCreateTopicRequest(req *dtos.CreateTopicRequest) error {
	if req == nil {
		return errors.New(errors.ErrInvalidArgument, "create topic request cannot be nil")
	}

	// 验证主题名称
	if strings.TrimSpace(req.Name) == "" {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot be empty")
	}

	// 验证主题名称格式
	if err := v.ValidateTopicName(req.Name); err != nil {
		return err
	}

	// 验证分区数
	if req.Partitions <= 0 {
		return errors.New(errors.ErrInvalidArgument, "partitions must be greater than 0")
	}

	if req.Partitions > 1000 {
		return errors.New(errors.ErrInvalidArgument, "partitions cannot exceed 1000")
	}

	// 验证复制因子
	if req.ReplicationFactor < 0 {
		return errors.New(errors.ErrInvalidArgument, "replication factor cannot be negative")
	}

	if req.ReplicationFactor > 10 {
		return errors.New(errors.ErrInvalidArgument, "replication factor cannot exceed 10")
	}

	// 验证保留时间
	if req.RetentionMs < 0 {
		return errors.New(errors.ErrInvalidArgument, "retention time cannot be negative")
	}

	// 验证配置项
	if req.Config != nil {
		for key, value := range req.Config {
			if strings.TrimSpace(key) == "" {
				return errors.New(errors.ErrInvalidArgument, "config key cannot be empty")
			}
			if strings.TrimSpace(value) == "" {
				return errors.New(errors.ErrInvalidArgument, "config value cannot be empty")
			}
		}
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

// ValidateTopicName 完整主题名称验证
func (v *Validator) ValidateTopicName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot be empty")
	}

	// 验证长度
	if len(name) > 249 {
		return errors.New(errors.ErrInvalidArgument, "topic name too long (max 249 characters)")
	}

	// 验证字符
	for i, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-') {
			return errors.New(errors.ErrInvalidArgument,
				fmt.Sprintf("invalid character '%c' at position %d in topic name", r, i))
		}
	}

	// 不能以点开头或结尾
	if name[0] == '.' || name[len(name)-1] == '.' {
		return errors.New(errors.ErrInvalidArgument, "topic name cannot start or end with dot")
	}

	// 检查保留名称
	reservedNames := []string{"__consumer_offsets", "__transaction_state"}
	for _, reserved := range reservedNames {
		if name == reserved {
			return errors.New(errors.ErrInvalidArgument, fmt.Sprintf("topic name is reserved: %s", name))
		}
	}

	return nil
}

// ValidateSmartModuleName 完整SmartModule名称验证
func (v *Validator) ValidateSmartModuleName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot be empty")
	}

	// 验证长度
	if len(name) > 100 {
		return errors.New(errors.ErrInvalidArgument, "smart module name too long (max 100 characters)")
	}

	// 验证字符（允许字母、数字、下划线、连字符）
	for i, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-') {
			return errors.New(errors.ErrInvalidArgument,
				fmt.Sprintf("invalid character '%c' at position %d in smart module name", r, i))
		}
	}

	// 不能以连字符或下划线开头
	if name[0] == '-' || name[0] == '_' {
		return errors.New(errors.ErrInvalidArgument, "smart module name cannot start with dash or underscore")
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

// 移除复杂的字符验证方法，由服务端处理

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
