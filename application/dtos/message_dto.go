package dtos

import (
	"time"
)

// ProduceMessageRequest 生产消息请求DTO
type ProduceMessageRequest struct {
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Value     string            `json:"value"`
	MessageID string            `json:"message_id,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// ProduceMessageResponse 生产消息响应DTO
type ProduceMessageResponse struct {
	MessageID string `json:"message_id"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// ProduceBatchRequest 批量生产请求DTO
type ProduceBatchRequest struct {
	Messages []*ProduceMessageRequest `json:"messages"`
}

// ProduceBatchResponse 批量生产响应DTO
type ProduceBatchResponse struct {
	Results       []*ProduceMessageResponse `json:"results"`
	TotalMessages int                       `json:"total_messages"`
	SuccessCount  int                       `json:"success_count"`
	FailureCount  int                       `json:"failure_count"`
}

// ConsumeMessageRequest 消费消息请求DTO
type ConsumeMessageRequest struct {
	Topic       string `json:"topic"`
	Group       string `json:"group"`
	Partition   int32  `json:"partition,omitempty"`
	Offset      int64  `json:"offset,omitempty"`
	MaxMessages int    `json:"max_messages"`
}

// ConsumeMessageResponse 消费消息响应DTO
type ConsumeMessageResponse struct {
	Messages []*MessageDTO `json:"messages"`
	Count    int           `json:"count"`
	Success  bool          `json:"success"`
	Error    string        `json:"error,omitempty"`
}

// MessageDTO 消息DTO
type MessageDTO struct {
	ID        string            `json:"id"`
	MessageID string            `json:"message_id"`
	Topic     string            `json:"topic"`
	Key       string            `json:"key"`
	Value     string            `json:"value"`
	Headers   map[string]string `json:"headers"`
	Partition int32             `json:"partition"`
	Offset    int64             `json:"offset"`
	Timestamp time.Time         `json:"timestamp"`
}

// FilteredConsumeRequest 过滤消费请求DTO
type FilteredConsumeRequest struct {
	Topic       string              `json:"topic"`
	Group       string              `json:"group"`
	MaxMessages int                 `json:"max_messages"`
	Filters     []*FilterCondition  `json:"filters"`
	AndLogic    bool                `json:"and_logic"`
}

// FilterCondition 过滤条件DTO
type FilterCondition struct {
	Type     string `json:"type"`
	Field    string `json:"field,omitempty"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// FilteredConsumeResponse 过滤消费响应DTO
type FilteredConsumeResponse struct {
	Messages       []*MessageDTO `json:"messages"`
	FilteredCount  int           `json:"filtered_count"`
	TotalScanned   int           `json:"total_scanned"`
	Success        bool          `json:"success"`
	Error          string        `json:"error,omitempty"`
}