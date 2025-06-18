package types

import "time"

// TopicInfo 主题信息
type TopicInfo struct {
	Name       string `json:"name"`
	Partitions int32  `json:"partitions"`
	Replicas   int32  `json:"replicas"`
}

// PartitionInfo 分区信息
type PartitionInfo struct {
	PartitionID    int32   `json:"partition_id"`
	LeaderID       int64   `json:"leader_id"`
	ReplicaIDs     []int64 `json:"replica_ids"`
	ISRIDs         []int64 `json:"isr_ids"`
	HighWatermark  int64   `json:"high_watermark"`
	LogStartOffset int64   `json:"log_start_offset"`
}

// CreateTopicOptions 创建主题选项
type CreateTopicOptions struct {
	Name              string            `json:"name"`
	Partitions        int32             `json:"partitions"`
	ReplicationFactor int32             `json:"replication_factor"`
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
	Topics  []string `json:"topics"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}

// DescribeTopicResult 描述主题结果
type DescribeTopicResult struct {
	Topic   *TopicInfo `json:"topic"`
	Success bool       `json:"success"`
	Error   string     `json:"error,omitempty"`
}

// DescribeTopicDetailResult 描述主题详细结果（新版本）
type DescribeTopicDetailResult struct {
	Topic       string            `json:"topic"`
	RetentionMs int64             `json:"retention_ms"`
	Config      map[string]string `json:"config"`
	Partitions  []*PartitionInfo  `json:"partitions"`
	Success     bool              `json:"success"`
	Error       string            `json:"error,omitempty"`
}

// PartitionStats 分区统计信息
type PartitionStats struct {
	PartitionID    int32     `json:"partition_id"`
	MessageCount   int64     `json:"message_count"`
	TotalSizeBytes int64     `json:"total_size_bytes"`
	EarliestOffset int64     `json:"earliest_offset"`
	LatestOffset   int64     `json:"latest_offset"`
	LastUpdated    time.Time `json:"last_updated"`
}

// TopicStats 主题统计信息
type TopicStats struct {
	Topic             string            `json:"topic"`
	PartitionCount    int32             `json:"partition_count"`
	ReplicationFactor int32             `json:"replication_factor"`
	TotalMessageCount int64             `json:"total_message_count"`
	TotalSizeBytes    int64             `json:"total_size_bytes"`
	Partitions        []*PartitionStats `json:"partitions,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	LastUpdated       time.Time         `json:"last_updated"`
}

// GetTopicStatsOptions 获取主题统计选项
type GetTopicStatsOptions struct {
	Topic             string `json:"topic,omitempty"`
	IncludePartitions bool   `json:"include_partitions"`
}

// GetTopicStatsResult 获取主题统计结果
type GetTopicStatsResult struct {
	Topics      []*TopicStats `json:"topics"`
	CollectedAt time.Time     `json:"collected_at"`
	Success     bool          `json:"success"`
	Error       string        `json:"error,omitempty"`
}
