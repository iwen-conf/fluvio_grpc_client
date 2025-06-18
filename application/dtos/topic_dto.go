package dtos

import (
	"time"
)

// CreateTopicRequest 创建主题请求DTO
type CreateTopicRequest struct {
	Name              string            `json:"name"`
	Partitions        int32             `json:"partitions"`
	ReplicationFactor int32             `json:"replication_factor,omitempty"`
	RetentionMs       int64             `json:"retention_ms,omitempty"`
	Config            map[string]string `json:"config,omitempty"`
	Description       string            `json:"description,omitempty"`
}

// CreateTopicResponse 创建主题响应DTO
type CreateTopicResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// TopicDTO 主题DTO
type TopicDTO struct {
	Name              string               `json:"name"`
	Description       string               `json:"description,omitempty"`
	Partitions        int32                `json:"partitions"`
	ReplicationFactor int32                `json:"replication_factor"`
	RetentionMs       int64                `json:"retention_ms"`
	Config            map[string]string    `json:"config"`
	PartitionDetails  []*PartitionInfoDTO  `json:"partition_details,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

// PartitionInfoDTO 分区信息DTO
type PartitionInfoDTO struct {
	PartitionID    int32   `json:"partition_id"`
	LeaderID       int32   `json:"leader_id"`
	ReplicaIDs     []int32 `json:"replica_ids"`
	HighWatermark  int64   `json:"high_watermark"`
	LowWatermark   int64   `json:"low_watermark"`
	MessageCount   int64   `json:"message_count"`
	TotalSizeBytes int64   `json:"total_size_bytes"`
}

// ListTopicsResponse 列出主题响应DTO
type ListTopicsResponse struct {
	Topics  []string `json:"topics"`
	Count   int      `json:"count"`
	Success bool     `json:"success"`
	Error   string   `json:"error,omitempty"`
}

// TopicDetailResponse 主题详情响应DTO
type TopicDetailResponse struct {
	Topic   *TopicDTO `json:"topic"`
	Success bool      `json:"success"`
	Error   string    `json:"error,omitempty"`
}

// TopicStatsRequest 主题统计请求DTO
type TopicStatsRequest struct {
	Topic             string `json:"topic"`
	IncludePartitions bool   `json:"include_partitions"`
}

// TopicStatsResponse 主题统计响应DTO
type TopicStatsResponse struct {
	Topics  []*TopicStatsDTO `json:"topics"`
	Success bool             `json:"success"`
	Error   string           `json:"error,omitempty"`
}

// TopicStatsDTO 主题统计DTO
type TopicStatsDTO struct {
	Topic              string                `json:"topic"`
	TotalMessageCount  int64                 `json:"total_message_count"`
	TotalSizeBytes     int64                 `json:"total_size_bytes"`
	PartitionCount     int32                 `json:"partition_count"`
	Partitions         []*PartitionStatsDTO  `json:"partitions,omitempty"`
}

// PartitionStatsDTO 分区统计DTO
type PartitionStatsDTO struct {
	PartitionID      int32 `json:"partition_id"`
	MessageCount     int64 `json:"message_count"`
	TotalSizeBytes   int64 `json:"total_size_bytes"`
	HighWatermark    int64 `json:"high_watermark"`
	LowWatermark     int64 `json:"low_watermark"`
}

// DeleteTopicRequest 删除主题请求DTO
type DeleteTopicRequest struct {
	Name string `json:"name"`
}

// DeleteTopicResponse 删除主题响应DTO
type DeleteTopicResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}