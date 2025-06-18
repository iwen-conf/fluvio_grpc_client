package api

import (
	"time"
)

// ProduceOptions 生产选项
type ProduceOptions struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key,omitempty"`
	MessageID string            `json:"message_id,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// ProduceResult 生产结果
type ProduceResult struct {
	MessageID string `json:"message_id"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
}

// Message 消息
type Message struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Value     string            `json:"value"`
	MessageID string            `json:"message_id,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	Partition int32             `json:"partition,omitempty"`
	Offset    int64             `json:"offset,omitempty"`
	Timestamp time.Time         `json:"timestamp,omitempty"`
}

// ProduceBatchResult 批量生产结果
type ProduceBatchResult struct {
	Results       []*ProduceResult `json:"results"`
	TotalMessages int              `json:"total_messages"`
	SuccessCount  int              `json:"success_count"`
	FailureCount  int              `json:"failure_count"`
}

// ConsumeOptions 消费选项
type ConsumeOptions struct {
	Topic       string `json:"topic"`
	Group       string `json:"group"`
	Partition   int32  `json:"partition,omitempty"`
	Offset      int64  `json:"offset,omitempty"`
	MaxMessages int    `json:"max_messages"`
}

// FilterCondition 过滤条件
type FilterCondition struct {
	Type     FilterType     `json:"type"`
	Field    string         `json:"field,omitempty"`
	Operator FilterOperator `json:"operator"`
	Value    string         `json:"value"`
}

// FilterType 过滤类型
type FilterType string

const (
	FilterTypeKey    FilterType = "key"
	FilterTypeValue  FilterType = "value"
	FilterTypeHeader FilterType = "header"
	FilterTypeOffset FilterType = "offset"
)

// FilterOperator 过滤操作符
type FilterOperator string

const (
	FilterOperatorEq       FilterOperator = "eq"
	FilterOperatorNe       FilterOperator = "ne"
	FilterOperatorGt       FilterOperator = "gt"
	FilterOperatorGte      FilterOperator = "gte"
	FilterOperatorLt       FilterOperator = "lt"
	FilterOperatorLte      FilterOperator = "lte"
	FilterOperatorContains FilterOperator = "contains"
	FilterOperatorRegex    FilterOperator = "regex"
)

// FilteredConsumeOptions 过滤消费选项
type FilteredConsumeOptions struct {
	Topic       string             `json:"topic"`
	Group       string             `json:"group"`
	MaxMessages int                `json:"max_messages"`
	Filters     []*FilterCondition `json:"filters"`
	AndLogic    bool               `json:"and_logic"`
}

// FilteredConsumeResult 过滤消费结果
type FilteredConsumeResult struct {
	Messages      []*Message `json:"messages"`
	FilteredCount int        `json:"filtered_count"`
	TotalScanned  int        `json:"total_scanned"`
}

// StreamConsumeOptions 流式消费选项
type StreamConsumeOptions struct {
	Topic        string `json:"topic"`
	Group        string `json:"group"`
	Partition    int32  `json:"partition,omitempty"`
	Offset       int64  `json:"offset,omitempty"`
	MaxBatchSize int    `json:"max_batch_size,omitempty"`
	MaxWaitMs    int    `json:"max_wait_ms,omitempty"`
}

// StreamMessage 流式消息
type StreamMessage struct {
	Message *Message `json:"message"`
	Error   error    `json:"error,omitempty"`
}

// CreateTopicOptions 创建主题选项
type CreateTopicOptions struct {
	Name              string            `json:"name"`
	Partitions        int32             `json:"partitions"`
	ReplicationFactor int32             `json:"replication_factor,omitempty"`
	RetentionMs       int64             `json:"retention_ms,omitempty"`
	Config            map[string]string `json:"config,omitempty"`
}

// CreateTopicResult 创建主题结果
type CreateTopicResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// DeleteTopicOptions 删除主题选项
type DeleteTopicOptions struct {
	Name string `json:"name"`
}

// DeleteTopicResult 删除主题结果
type DeleteTopicResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ListTopicsResult 列出主题结果
type ListTopicsResult struct {
	Topics []string `json:"topics"`
}

// TopicDescription 主题描述
type TopicDescription struct {
	Name       string `json:"name"`
	Partitions int32  `json:"partitions"`
}

// TopicDetail 主题详情
type TopicDetail struct {
	Topic       string               `json:"topic"`
	Partitions  []*PartitionInfo     `json:"partitions"`
	RetentionMs int64                `json:"retention_ms"`
	Config      map[string]string    `json:"config"`
}

// PartitionInfo 分区信息
type PartitionInfo struct {
	PartitionID   int32 `json:"partition_id"`
	LeaderID      int32 `json:"leader_id"`
	HighWatermark int64 `json:"high_watermark"`
}

// GetTopicStatsOptions 获取主题统计选项
type GetTopicStatsOptions struct {
	Topic             string `json:"topic"`
	IncludePartitions bool   `json:"include_partitions"`
}

// GetTopicStatsResult 获取主题统计结果
type GetTopicStatsResult struct {
	Topics []*TopicStats `json:"topics"`
}

// TopicStats 主题统计
type TopicStats struct {
	Topic              string            `json:"topic"`
	TotalMessageCount  int64             `json:"total_message_count"`
	TotalSizeBytes     int64             `json:"total_size_bytes"`
	PartitionCount     int32             `json:"partition_count"`
	Partitions         []*PartitionStats `json:"partitions,omitempty"`
}

// PartitionStats 分区统计
type PartitionStats struct {
	PartitionID    int32 `json:"partition_id"`
	MessageCount   int64 `json:"message_count"`
	TotalSizeBytes int64 `json:"total_size_bytes"`
}

// 其他类型定义（消费组、SmartModule、存储等）将在后续添加...

// ListConsumerGroupsResult 列出消费组结果
type ListConsumerGroupsResult struct {
	Groups []*ConsumerGroup `json:"groups"`
}

// ConsumerGroup 消费组
type ConsumerGroup struct {
	GroupID string `json:"group_id"`
	State   string `json:"state"`
}

// ConsumerGroupDetail 消费组详情
type ConsumerGroupDetail struct {
	Group   *ConsumerGroup `json:"group"`
	Members []*Member      `json:"members"`
	Offsets []*Offset      `json:"offsets"`
}

// Member 成员
type Member struct {
	MemberID string `json:"member_id"`
	ClientID string `json:"client_id"`
}

// Offset 偏移量
type Offset struct {
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
}

// DeleteConsumerGroupResult 删除消费组结果
type DeleteConsumerGroupResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// SmartModule相关类型
type ListSmartModulesResult struct {
	SmartModules []*SmartModule `json:"smart_modules"`
}

type SmartModule struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

type CreateSmartModuleOptions struct {
	Spec     *SmartModuleSpec `json:"spec"`
	WasmCode []byte           `json:"wasm_code"`
}

type SmartModuleSpec struct {
	Name        string                   `json:"name"`
	InputKind   SmartModuleInput         `json:"input_kind"`
	OutputKind  SmartModuleOutput        `json:"output_kind"`
	Description string                   `json:"description"`
	Version     string                   `json:"version"`
	Parameters  []*SmartModuleParameter  `json:"parameters"`
}

type SmartModuleInput int32
type SmartModuleOutput int32

const (
	SmartModuleInputStream  SmartModuleInput = 0
	SmartModuleOutputStream SmartModuleOutput = 0
)

type SmartModuleParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Optional    bool   `json:"optional"`
}

type CreateSmartModuleResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type UpdateSmartModuleOptions struct {
	Name     string           `json:"name"`
	Spec     *SmartModuleSpec `json:"spec,omitempty"`
	WasmCode []byte           `json:"wasm_code,omitempty"`
}

type UpdateSmartModuleResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type DeleteSmartModuleResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SmartModuleDetail struct {
	SmartModule *SmartModule `json:"smart_module"`
	Spec        *SmartModuleSpec `json:"spec"`
}

// 存储管理相关类型
type GetStorageStatusOptions struct {
	IncludeDetails bool `json:"include_details"`
}

type GetStorageStatusResult struct {
	PersistenceEnabled bool           `json:"persistence_enabled"`
	StorageStats       *StorageStats  `json:"storage_stats,omitempty"`
	CheckedAt          time.Time      `json:"checked_at"`
	Success            bool           `json:"success"`
	Error              string         `json:"error,omitempty"`
}

type StorageStats struct {
	StorageType       string                     `json:"storage_type"`
	ConsumerGroups    int32                      `json:"consumer_groups"`
	ConsumerOffsets   int32                      `json:"consumer_offsets"`
	SmartModules      int32                      `json:"smart_modules"`
	ConnectionStatus  string                     `json:"connection_status"`
	ConnectionStats   *StorageConnectionStats    `json:"connection_stats,omitempty"`
	DatabaseInfo      *StorageDatabaseInfo       `json:"database_info,omitempty"`
}

type StorageConnectionStats struct {
	CurrentConnections      int32 `json:"current_connections"`
	AvailableConnections    int32 `json:"available_connections"`
	TotalCreatedConnections int64 `json:"total_created_connections"`
}

type StorageDatabaseInfo struct {
	Name        string `json:"name"`
	Collections int32  `json:"collections"`
	DataSize    int64  `json:"data_size"`
	StorageSize int64  `json:"storage_size"`
	Indexes     int32  `json:"indexes"`
	IndexSize   int64  `json:"index_size"`
}

type MigrateStorageOptions struct {
	SourceType      string `json:"source_type"`
	TargetType      string `json:"target_type"`
	VerifyMigration bool   `json:"verify_migration"`
	ForceMigration  bool   `json:"force_migration"`
}

type MigrateStorageResult struct {
	Success            bool             `json:"success"`
	MigrationStats     *MigrationStats  `json:"migration_stats,omitempty"`
	VerificationPassed bool             `json:"verification_passed"`
	CompletedAt        time.Time        `json:"completed_at"`
	Error              string           `json:"error,omitempty"`
}

type MigrationStats struct {
	ConsumerGroupsMigrated  int32    `json:"consumer_groups_migrated"`
	ConsumerOffsetsMigrated int32    `json:"consumer_offsets_migrated"`
	SmartModulesMigrated    int32    `json:"smart_modules_migrated"`
	Errors                  []string `json:"errors"`
	TotalMigrated           int32    `json:"total_migrated"`
}

type GetStorageMetricsOptions struct {
	IncludeHistory bool  `json:"include_history"`
	HistoryLimit   int32 `json:"history_limit"`
}

type GetStorageMetricsResult struct {
	CurrentMetrics  *StorageMetrics            `json:"current_metrics,omitempty"`
	MetricsHistory  []*StorageMetrics          `json:"metrics_history,omitempty"`
	HealthStatus    *StorageHealthCheckResult  `json:"health_status,omitempty"`
	Alerts          []string                   `json:"alerts"`
	CollectedAt     time.Time                  `json:"collected_at"`
	Success         bool                       `json:"success"`
	Error           string                     `json:"error,omitempty"`
}

type StorageMetrics struct {
	StorageType           string    `json:"storage_type"`
	ResponseTimeMs        int32     `json:"response_time_ms"`
	OperationsPerSecond   float64   `json:"operations_per_second"`
	ErrorRate             float64   `json:"error_rate"`
	ConnectionPoolUsage   float64   `json:"connection_pool_usage"`
	MemoryUsageMB         int32     `json:"memory_usage_mb"`
	DiskUsageMB           int32     `json:"disk_usage_mb"`
	LastUpdated           time.Time `json:"last_updated"`
}

type StorageHealthCheckResult struct {
	Status         string    `json:"status"`
	ResponseTimeMs int32     `json:"response_time_ms"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

// 批量删除相关类型
type BulkDeleteOptions struct {
	Topics         []string `json:"topics,omitempty"`
	ConsumerGroups []string `json:"consumer_groups,omitempty"`
	SmartModules   []string `json:"smart_modules,omitempty"`
	Force          bool     `json:"force"`
}

type BulkDeleteResult struct {
	Results           []*BulkDeleteItemResult `json:"results"`
	TotalRequested    int32                   `json:"total_requested"`
	SuccessfulDeletes int32                   `json:"successful_deletes"`
	FailedDeletes     int32                   `json:"failed_deletes"`
	Success           bool                    `json:"success"`
	Error             string                  `json:"error,omitempty"`
}

type BulkDeleteItemResult struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}