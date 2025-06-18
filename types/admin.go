package types

import "time"

// ClusterInfo 集群信息
type ClusterInfo struct {
	Status       string            `json:"status"`
	ControllerID int32             `json:"controller_id"`
	Brokers      []*BrokerInfo     `json:"brokers"`
	Metadata     map[string]string `json:"metadata"`
}

// BrokerInfo Broker信息
type BrokerInfo struct {
	ID       int64             `json:"id"`
	Addr     string            `json:"addr"`
	Status   string            `json:"status"`
	Metadata map[string]string `json:"metadata"`
}

// MetricInfo 指标信息
type MetricInfo struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
}

// DescribeClusterResult 描述集群结果
type DescribeClusterResult struct {
	Cluster *ClusterInfo `json:"cluster"`
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
}

// ListBrokersResult 列出Broker结果
type ListBrokersResult struct {
	Brokers []*BrokerInfo `json:"brokers"`
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
}

// GetMetricsOptions 获取指标选项
type GetMetricsOptions struct {
	MetricNames []string          `json:"metric_names"`
	Labels      map[string]string `json:"labels"`
}

// GetMetricsResult 获取指标结果
type GetMetricsResult struct {
	Metrics []*MetricInfo `json:"metrics"`
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
}

// SmartModuleInputType SmartModule输入类型
type SmartModuleInputType int32

const (
	SmartModuleInputUnknown SmartModuleInputType = 0
	SmartModuleInputStream  SmartModuleInputType = 1
	SmartModuleInputTable   SmartModuleInputType = 2
)

// SmartModuleOutputType SmartModule输出类型
type SmartModuleOutputType int32

const (
	SmartModuleOutputUnknown SmartModuleOutputType = 0
	SmartModuleOutputStream  SmartModuleOutputType = 1
	SmartModuleOutputTable   SmartModuleOutputType = 2
)

// SmartModuleParameter 参数定义
type SmartModuleParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Optional    bool   `json:"optional"`
}

// SmartModuleSpec SmartModule规格定义
type SmartModuleSpec struct {
	Name        string                  `json:"name"`
	InputKind   SmartModuleInputType    `json:"input_kind"`
	OutputKind  SmartModuleOutputType   `json:"output_kind"`
	Parameters  []*SmartModuleParameter `json:"parameters"`
	Description string                  `json:"description"`
	Version     string                  `json:"version"`
}

// SmartModuleInfo SmartModule信息
type SmartModuleInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
	Spec        *SmartModuleSpec  `json:"spec,omitempty"` // 新增：详细规格
}

// CreateSmartModuleOptions 创建SmartModule选项
type CreateSmartModuleOptions struct {
	Spec     *SmartModuleSpec `json:"spec"`
	WasmCode []byte           `json:"wasm_code"`
}

// CreateSmartModuleResult 创建SmartModule结果
type CreateSmartModuleResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// DeleteSmartModuleResult 删除SmartModule结果
type DeleteSmartModuleResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ListSmartModulesResult 列出SmartModule结果
type ListSmartModulesResult struct {
	SmartModules []*SmartModuleInfo `json:"smart_modules"`
	Success      bool               `json:"success"`
	Error        string             `json:"error,omitempty"`
}

// DescribeSmartModuleResult 描述SmartModule结果
type DescribeSmartModuleResult struct {
	SmartModule *SmartModuleInfo `json:"smart_module"`
	Success     bool             `json:"success"`
	Error       string           `json:"error,omitempty"`
}

// UpdateSmartModuleOptions 更新SmartModule选项
type UpdateSmartModuleOptions struct {
	Name     string           `json:"name"`
	Spec     *SmartModuleSpec `json:"spec,omitempty"`
	WasmCode []byte           `json:"wasm_code,omitempty"`
}

// UpdateSmartModuleResult 更新SmartModule结果
type UpdateSmartModuleResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// BulkDeleteOptions 批量删除选项
type BulkDeleteOptions struct {
	Topics         []string `json:"topics,omitempty"`
	ConsumerGroups []string `json:"consumer_groups,omitempty"`
	SmartModules   []string `json:"smart_modules,omitempty"`
	Force          bool     `json:"force"`
}

// BulkDeleteItemResult 批量删除单项结果
type BulkDeleteItemResult struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// BulkDeleteResult 批量删除结果
type BulkDeleteResult struct {
	Results           []*BulkDeleteItemResult `json:"results"`
	TotalRequested    int32                   `json:"total_requested"`
	SuccessfulDeletes int32                   `json:"successful_deletes"`
	FailedDeletes     int32                   `json:"failed_deletes"`
	Success           bool                    `json:"success"`
	Error             string                  `json:"error,omitempty"`
}

// StorageConnectionStats 存储连接统计信息
type StorageConnectionStats struct {
	CurrentConnections      uint32 `json:"current_connections"`
	AvailableConnections    uint32 `json:"available_connections"`
	TotalCreatedConnections uint32 `json:"total_created_connections"`
}

// StorageDatabaseInfo 存储数据库信息
type StorageDatabaseInfo struct {
	Name        string `json:"name"`
	Collections uint32 `json:"collections"`
	DataSize    uint64 `json:"data_size"`
	StorageSize uint64 `json:"storage_size"`
	Indexes     uint32 `json:"indexes"`
	IndexSize   uint64 `json:"index_size"`
}

// StorageStats 存储统计信息
type StorageStats struct {
	StorageType      string                  `json:"storage_type"`
	ConsumerGroups   uint64                  `json:"consumer_groups"`
	ConsumerOffsets  uint64                  `json:"consumer_offsets"`
	SmartModules     uint64                  `json:"smart_modules"`
	ConnectionStatus string                  `json:"connection_status"`
	ConnectionStats  *StorageConnectionStats `json:"connection_stats,omitempty"`
	DatabaseInfo     *StorageDatabaseInfo    `json:"database_info,omitempty"`
}

// GetStorageStatusOptions 获取存储状态选项
type GetStorageStatusOptions struct {
	IncludeDetails bool `json:"include_details"`
}

// GetStorageStatusResult 获取存储状态结果
type GetStorageStatusResult struct {
	PersistenceEnabled bool          `json:"persistence_enabled"`
	StorageStats       *StorageStats `json:"storage_stats"`
	CheckedAt          time.Time     `json:"checked_at"`
	Success            bool          `json:"success"`
	Error              string        `json:"error,omitempty"`
}

// MigrateStorageOptions 存储迁移选项
type MigrateStorageOptions struct {
	SourceType      string `json:"source_type"`
	TargetType      string `json:"target_type"`
	VerifyMigration bool   `json:"verify_migration"`
	ForceMigration  bool   `json:"force_migration"`
}

// MigrationStats 迁移统计信息
type MigrationStats struct {
	ConsumerGroupsMigrated  uint64   `json:"consumer_groups_migrated"`
	ConsumerOffsetsMigrated uint64   `json:"consumer_offsets_migrated"`
	SmartModulesMigrated    uint64   `json:"smart_modules_migrated"`
	Errors                  []string `json:"errors"`
	TotalMigrated           uint64   `json:"total_migrated"`
}

// MigrateStorageResult 存储迁移结果
type MigrateStorageResult struct {
	Success            bool            `json:"success"`
	MigrationStats     *MigrationStats `json:"migration_stats"`
	VerificationPassed bool            `json:"verification_passed"`
	CompletedAt        time.Time       `json:"completed_at"`
	Error              string          `json:"error,omitempty"`
}

// GetStorageMetricsOptions 获取存储指标选项
type GetStorageMetricsOptions struct {
	IncludeHistory bool   `json:"include_history"`
	HistoryLimit   uint32 `json:"history_limit"`
}

// StorageMetrics 存储性能指标
type StorageMetrics struct {
	StorageType         string    `json:"storage_type"`
	ResponseTimeMs      uint64    `json:"response_time_ms"`
	OperationsPerSecond float64   `json:"operations_per_second"`
	ErrorRate           float64   `json:"error_rate"`
	ConnectionPoolUsage float64   `json:"connection_pool_usage"`
	MemoryUsageMB       uint64    `json:"memory_usage_mb"`
	DiskUsageMB         uint64    `json:"disk_usage_mb"`
	LastUpdated         time.Time `json:"last_updated"`
}

// StorageHealthCheckResult 存储健康检查结果
type StorageHealthCheckResult struct {
	Status         string    `json:"status"`
	ResponseTimeMs uint64    `json:"response_time_ms"`
	ErrorMessage   string    `json:"error_message,omitempty"`
	CheckedAt      time.Time `json:"checked_at"`
}

// GetStorageMetricsResult 获取存储指标结果
type GetStorageMetricsResult struct {
	CurrentMetrics *StorageMetrics           `json:"current_metrics"`
	MetricsHistory []*StorageMetrics         `json:"metrics_history,omitempty"`
	HealthStatus   *StorageHealthCheckResult `json:"health_status"`
	Alerts         []string                  `json:"alerts"`
	CollectedAt    time.Time                 `json:"collected_at"`
	Success        bool                      `json:"success"`
	Error          string                    `json:"error,omitempty"`
}
