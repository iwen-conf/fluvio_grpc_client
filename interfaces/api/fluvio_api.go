package api

import (
	"context"
	"time"
)

// FluvioAPI Fluvio公共API接口
type FluvioAPI interface {
	// 连接管理
	Connect() error
	Close() error
	HealthCheck(ctx context.Context) error
	Ping(ctx context.Context) (time.Duration, error)

	// 消息操作
	Producer() ProducerAPI
	Consumer() ConsumerAPI

	// 主题管理
	Topic() TopicAPI

	// 管理功能
	Admin() AdminAPI
}

// ProducerAPI 生产者API接口
type ProducerAPI interface {
	// 生产单条消息
	Produce(ctx context.Context, value string, opts ProduceOptions) (*ProduceResult, error)

	// 批量生产消息
	ProduceBatch(ctx context.Context, messages []Message) (*ProduceBatchResult, error)
}

// ConsumerAPI 消费者API接口
type ConsumerAPI interface {
	// 基本消费
	Consume(ctx context.Context, opts ConsumeOptions) ([]*Message, error)

	// 过滤消费
	ConsumeFiltered(ctx context.Context, opts FilteredConsumeOptions) (*FilteredConsumeResult, error)

	// 流式消费
	ConsumeStream(ctx context.Context, opts StreamConsumeOptions) (<-chan *StreamMessage, error)
}

// TopicAPI 主题API接口
type TopicAPI interface {
	// 主题管理
	Create(ctx context.Context, opts CreateTopicOptions) (*CreateTopicResult, error)
	Delete(ctx context.Context, opts DeleteTopicOptions) (*DeleteTopicResult, error)
	Exists(ctx context.Context, name string) (bool, error)
	CreateIfNotExists(ctx context.Context, opts CreateTopicOptions) (*CreateTopicResult, error)

	// 主题查询
	List(ctx context.Context) (*ListTopicsResult, error)
	Describe(ctx context.Context, name string) (*TopicDescription, error)
	DescribeTopicDetail(ctx context.Context, name string) (*TopicDetail, error)

	// 主题统计
	GetTopicStats(ctx context.Context, opts GetTopicStatsOptions) (*GetTopicStatsResult, error)
}

// AdminAPI 管理API接口
type AdminAPI interface {
	// 消费组管理
	ListConsumerGroups(ctx context.Context) (*ListConsumerGroupsResult, error)
	DescribeConsumerGroup(ctx context.Context, groupID string) (*ConsumerGroupDetail, error)
	DeleteConsumerGroup(ctx context.Context, groupID string) (*DeleteConsumerGroupResult, error)

	// SmartModule管理
	ListSmartModules(ctx context.Context) (*ListSmartModulesResult, error)
	CreateSmartModule(ctx context.Context, opts CreateSmartModuleOptions) (*CreateSmartModuleResult, error)
	UpdateSmartModule(ctx context.Context, opts UpdateSmartModuleOptions) (*UpdateSmartModuleResult, error)
	DeleteSmartModule(ctx context.Context, name string) (*DeleteSmartModuleResult, error)
	DescribeSmartModule(ctx context.Context, name string) (*SmartModuleDetail, error)

	// 存储管理
	GetStorageStatus(ctx context.Context, opts GetStorageStatusOptions) (*GetStorageStatusResult, error)
	MigrateStorage(ctx context.Context, opts MigrateStorageOptions) (*MigrateStorageResult, error)
	GetStorageMetrics(ctx context.Context, opts GetStorageMetricsOptions) (*GetStorageMetricsResult, error)

	// 批量操作
	BulkDelete(ctx context.Context, opts BulkDeleteOptions) (*BulkDeleteResult, error)
}
