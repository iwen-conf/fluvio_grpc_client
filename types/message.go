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
	Topic     string `json:"topic"`
	Group     string `json:"group"`
	Offset    int64  `json:"offset"`
	Partition int32  `json:"partition"`
}

// CommitOffsetOptions 提交偏移量选项
type CommitOffsetOptions struct {
	Topic     string `json:"topic"`
	Group     string `json:"group"`
	Offset    int64  `json:"offset"`
	Partition int32  `json:"partition"`
}
