package types

import (
	"time"
)

// Message 表示一条消息
type Message struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Value     string            `json:"value"`
	Headers   map[string]string `json:"headers"`
	Offset    int64             `json:"offset"`
	Partition int32             `json:"partition"`
	Timestamp time.Time         `json:"timestamp"`
	MessageID string            `json:"message_id"` // 新增：消息唯一ID
}

// ProduceResult 生产结果
type ProduceResult struct {
	MessageID string `json:"message_id"`
	Offset    int64  `json:"offset"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// BatchProduceResult 批量生产结果
type BatchProduceResult struct {
	Results    []*ProduceResult `json:"results"`
	TotalCount int              `json:"total_count"`
	Success    bool             `json:"success"`
	Errors     []string         `json:"errors,omitempty"`
}

// ConsumeResult 消费结果
type ConsumeResult struct {
	Messages   []*Message `json:"messages"`
	NextOffset int64      `json:"next_offset"`
	Success    bool       `json:"success"`
	Error      string     `json:"error,omitempty"`
}

// ProduceOptions 生产选项
type ProduceOptions struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Headers   map[string]string `json:"headers"`
	Timestamp *time.Time        `json:"timestamp,omitempty"`
	MessageID string            `json:"message_id,omitempty"` // 新增：可选的消息ID
}

// ConsumeOptions 消费选项
type ConsumeOptions struct {
	Topic       string `json:"topic"`
	Group       string `json:"group"`
	Offset      int64  `json:"offset"`
	Partition   int32  `json:"partition"`
	MaxMessages int32  `json:"max_messages"`
	AutoCommit  bool   `json:"auto_commit"`
}

// StreamConsumeOptions 流式消费选项
type StreamConsumeOptions struct {
	Topic        string `json:"topic"`
	Group        string `json:"group"`
	Offset       int64  `json:"offset"`
	Partition    int32  `json:"partition"`
	MaxBatchSize int32  `json:"max_batch_size,omitempty"` // 新增：每次流式响应的最大消息数
	MaxWaitMs    int32  `json:"max_wait_ms,omitempty"`    // 新增：服务器等待批次满足的最长时间
}

// CommitOffsetOptions 提交偏移量选项
type CommitOffsetOptions struct {
	Topic     string `json:"topic"`
	Group     string `json:"group"`
	Offset    int64  `json:"offset"`
	Partition int32  `json:"partition"`
}

// FilterType 过滤器类型
type FilterType int32

const (
	FilterTypeUnknown   FilterType = 0
	FilterTypeKey       FilterType = 1 // 按消息键过滤
	FilterTypeHeader    FilterType = 2 // 按消息头过滤
	FilterTypeContent   FilterType = 3 // 按消息内容过滤
	FilterTypeTimestamp FilterType = 4 // 按时间戳过滤
)

// FilterCondition 过滤条件
type FilterCondition struct {
	Type     FilterType `json:"type"`     // 过滤类型
	Field    string     `json:"field"`    // 过滤字段（对于header和content过滤）
	Operator string     `json:"operator"` // 操作符：eq, ne, contains, starts_with, ends_with, gt, lt, gte, lte
	Value    string     `json:"value"`    // 过滤值
}

// FilteredConsumeOptions 过滤消费选项
type FilteredConsumeOptions struct {
	Topic       string            `json:"topic"`
	Group       string            `json:"group"`
	Offset      int64             `json:"offset"`
	Partition   int32             `json:"partition"`
	MaxMessages int32             `json:"max_messages"`
	Filters     []FilterCondition `json:"filters"`   // 过滤条件列表
	AndLogic    bool              `json:"and_logic"` // true为AND逻辑，false为OR逻辑
}

// FilteredConsumeResult 过滤消费结果
type FilteredConsumeResult struct {
	Messages      []*Message `json:"messages"`       // 过滤后的消息列表
	NextOffset    int64      `json:"next_offset"`    // 建议下次消费的起始偏移量
	TotalScanned  int32      `json:"total_scanned"`  // 总共扫描的消息数量
	FilteredCount int32      `json:"filtered_count"` // 过滤后的消息数量
	Success       bool       `json:"success"`
	Error         string     `json:"error,omitempty"`
}
