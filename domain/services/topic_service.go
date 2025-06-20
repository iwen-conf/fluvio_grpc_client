package services

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
)

// TopicService 主题领域服务
type TopicService struct{}

// NewTopicService 创建主题服务
func NewTopicService() *TopicService {
	return &TopicService{}
}

// ValidateTopicName 验证主题名称
func (ts *TopicService) ValidateTopicName(name string) error {
	if name == "" {
		return fmt.Errorf("topic name cannot be empty")
	}

	// 主题名称长度限制
	if len(name) > 249 {
		return fmt.Errorf("topic name too long: %d characters (max 249)", len(name))
	}

	// 主题名称格式验证（字母、数字、下划线、连字符、点）
	validName := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid topic name format: %s", name)
	}

	// 不能以点开头或结尾
	if name[0] == '.' || name[len(name)-1] == '.' {
		return fmt.Errorf("topic name cannot start or end with dot: %s", name)
	}

	// 保留名称检查
	reservedNames := []string{"__consumer_offsets", "__transaction_state"}
	for _, reserved := range reservedNames {
		if name == reserved {
			return fmt.Errorf("topic name is reserved: %s", name)
		}
	}

	return nil
}

// ValidateTopicConfig 验证主题配置
func (ts *TopicService) ValidateTopicConfig(topic *entities.Topic) error {
	if !topic.IsValid() {
		return fmt.Errorf("invalid topic configuration")
	}

	// 验证分区数
	if topic.Partitions <= 0 {
		return fmt.Errorf("partitions must be greater than 0")
	}

	if topic.Partitions > 1000 {
		return fmt.Errorf("too many partitions: %d (max 1000)", topic.Partitions)
	}

	// 验证复制因子
	if topic.ReplicationFactor < 0 {
		return fmt.Errorf("replication factor cannot be negative")
	}

	// 验证保留时间
	if topic.RetentionMs < 0 {
		return fmt.Errorf("retention time cannot be negative")
	}

	// 验证配置项
	return ts.validateTopicConfigItems(topic.Config)
}

// validateTopicConfigItems 验证配置项
func (ts *TopicService) validateTopicConfigItems(config map[string]string) error {
	validConfigs := map[string]func(string) error{
		"cleanup.policy":      ts.validateCleanupPolicy,
		"compression.type":    ts.validateCompressionType,
		"delete.retention.ms": ts.validateDeleteRetention,
		"segment.ms":          ts.validateSegmentMs,
		"max.message.bytes":   ts.validateMaxMessageBytes,
	}

	for key, value := range config {
		if validator, exists := validConfigs[key]; exists {
			if err := validator(value); err != nil {
				return fmt.Errorf("invalid config %s=%s: %w", key, value, err)
			}
		}
		// 未知配置项会被忽略或传递给服务器验证
	}

	return nil
}

// validateCleanupPolicy 验证清理策略
func (ts *TopicService) validateCleanupPolicy(value string) error {
	validPolicies := []string{"delete", "compact", "compact,delete"}
	for _, policy := range validPolicies {
		if value == policy {
			return nil
		}
	}
	return fmt.Errorf("invalid cleanup policy: %s", value)
}

// validateCompressionType 验证压缩类型
func (ts *TopicService) validateCompressionType(value string) error {
	validTypes := []string{"uncompressed", "gzip", "snappy", "lz4", "zstd"}
	for _, compressionType := range validTypes {
		if value == compressionType {
			return nil
		}
	}
	return fmt.Errorf("invalid compression type: %s", value)
}

// validateDeleteRetention 验证删除保留时间
func (ts *TopicService) validateDeleteRetention(value string) error {
	if value == "" {
		return fmt.Errorf("delete retention cannot be empty")
	}

	// 解析数值并验证范围
	retentionMs, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("delete retention must be a valid number: %v", err)
	}

	// 验证范围：最小1分钟，最大30天
	const minRetentionMs = 60 * 1000                // 1分钟
	const maxRetentionMs = 30 * 24 * 60 * 60 * 1000 // 30天

	if retentionMs < minRetentionMs {
		return fmt.Errorf("delete retention must be at least %d ms (1 minute)", minRetentionMs)
	}

	if retentionMs > maxRetentionMs {
		return fmt.Errorf("delete retention must be at most %d ms (30 days)", maxRetentionMs)
	}

	return nil
}

// validateSegmentMs 验证段时间
func (ts *TopicService) validateSegmentMs(value string) error {
	if value == "" {
		return fmt.Errorf("segment ms cannot be empty")
	}
	return nil
}

// validateMaxMessageBytes 验证最大消息字节数
func (ts *TopicService) validateMaxMessageBytes(value string) error {
	if value == "" {
		return fmt.Errorf("max message bytes cannot be empty")
	}
	return nil
}

// CalculateOptimalPartitions 计算最优分区数
func (ts *TopicService) CalculateOptimalPartitions(expectedThroughput, targetLatency int64) int32 {
	// 完整的分区计算逻辑，考虑多种因素

	// 参数验证
	if expectedThroughput <= 0 {
		expectedThroughput = 1000 // 默认每秒1000条消息
	}
	if targetLatency <= 0 {
		targetLatency = 100 // 默认100ms延迟要求
	}

	// 基础分区数计算
	basePartitions := int32(1)

	// 1. 根据吞吐量计算分区数
	// 假设每个分区每秒可以处理1000条消息
	const messagesPerPartitionPerSecond = 1000
	throughputBasedPartitions := int32((expectedThroughput + messagesPerPartitionPerSecond - 1) / messagesPerPartitionPerSecond)

	// 2. 根据延迟要求调整
	// 延迟越低，需要更多分区来并行处理
	latencyFactor := float64(1.0)
	if targetLatency < 50 {
		latencyFactor = 2.0 // 极低延迟要求
	} else if targetLatency < 100 {
		latencyFactor = 1.5 // 低延迟要求
	} else if targetLatency < 500 {
		latencyFactor = 1.2 // 中等延迟要求
	}

	// 3. 考虑消费者并行度
	// 通常建议分区数不超过预期消费者数量的2倍
	maxConsumers := int32(10) // 假设最多10个消费者
	consumerBasedPartitions := maxConsumers * 2

	// 4. 综合计算
	basePartitions = int32(float64(throughputBasedPartitions) * latencyFactor)

	// 5. 应用约束条件
	// 最小分区数为1
	if basePartitions < 1 {
		basePartitions = 1
	}

	// 不超过基于消费者数量的限制
	if basePartitions > consumerBasedPartitions {
		basePartitions = consumerBasedPartitions
	}

	// 6. 考虑Kafka集群限制
	// 通常单个主题不建议超过100个分区
	const maxPartitionsPerTopic = 100
	if basePartitions > maxPartitionsPerTopic {
		basePartitions = maxPartitionsPerTopic
	}

	// 7. 考虑存储和网络开销
	// 分区数过多会增加存储和网络开销
	if basePartitions > 50 && expectedThroughput < 10000 {
		basePartitions = 50 // 中等吞吐量时限制分区数
	}

	return basePartitions
}
