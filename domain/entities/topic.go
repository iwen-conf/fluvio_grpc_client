package entities

import (
	"time"
)

// Topic 表示Fluvio中的主题实体
type Topic struct {
	// 基本信息
	Name        string
	Description string
	
	// 配置
	Partitions        int32
	ReplicationFactor int32
	RetentionMs       int64
	Config            map[string]string
	
	// 分区信息
	PartitionDetails []*PartitionInfo
	
	// 元数据
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PartitionInfo 分区信息
type PartitionInfo struct {
	PartitionID    int32
	LeaderID       int32
	ReplicaIDs     []int32
	HighWatermark  int64
	LowWatermark   int64
	MessageCount   int64
	TotalSizeBytes int64
}

// NewTopic 创建新的主题实体
func NewTopic(name string, partitions int32) *Topic {
	now := time.Now()
	return &Topic{
		Name:       name,
		Partitions: partitions,
		Config:     make(map[string]string),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// WithReplicationFactor 设置复制因子
func (t *Topic) WithReplicationFactor(factor int32) *Topic {
	t.ReplicationFactor = factor
	return t
}

// WithRetention 设置保留时间
func (t *Topic) WithRetention(retentionMs int64) *Topic {
	t.RetentionMs = retentionMs
	return t
}

// WithConfig 设置配置
func (t *Topic) WithConfig(config map[string]string) *Topic {
	t.Config = config
	return t
}

// AddConfig 添加单个配置项
func (t *Topic) AddConfig(key, value string) *Topic {
	if t.Config == nil {
		t.Config = make(map[string]string)
	}
	t.Config[key] = value
	return t
}

// IsValid 验证主题是否有效
func (t *Topic) IsValid() bool {
	return t.Name != "" && t.Partitions > 0
}

// TotalMessages 计算总消息数
func (t *Topic) TotalMessages() int64 {
	var total int64
	for _, partition := range t.PartitionDetails {
		total += partition.MessageCount
	}
	return total
}

// TotalSize 计算总大小
func (t *Topic) TotalSize() int64 {
	var total int64
	for _, partition := range t.PartitionDetails {
		total += partition.TotalSizeBytes
	}
	return total
}