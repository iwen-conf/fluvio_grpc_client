package entities

import (
	"time"
)

// Message 表示Fluvio中的消息实体
type Message struct {
	// 核心标识
	ID        string
	MessageID string // 用户自定义的消息ID

	// 消息内容
	Key     string
	Value   []byte // 支持二进制数据
	Headers map[string]string

	// 元数据
	Topic     string
	Partition int32
	Offset    int64

	// 时间戳
	Timestamp time.Time
	CreatedAt time.Time
}

// NewMessage 创建新的消息实体
func NewMessage(key, value string) *Message {
	now := time.Now()
	return &Message{
		Key:       key,
		Value:     []byte(value),
		Headers:   make(map[string]string),
		Timestamp: now,
		CreatedAt: now,
	}
}

// NewMessageBytes 创建新的消息实体（二进制数据）
func NewMessageBytes(key string, value []byte) *Message {
	now := time.Now()
	return &Message{
		Key:       key,
		Value:     value,
		Headers:   make(map[string]string),
		Timestamp: now,
		CreatedAt: now,
	}
}

// WithMessageID 设置自定义消息ID
func (m *Message) WithMessageID(messageID string) *Message {
	m.MessageID = messageID
	return m
}

// WithHeaders 设置消息头部
func (m *Message) WithHeaders(headers map[string]string) *Message {
	m.Headers = headers
	return m
}

// AddHeader 添加单个头部
func (m *Message) AddHeader(key, value string) *Message {
	if m.Headers == nil {
		m.Headers = make(map[string]string)
	}
	m.Headers[key] = value
	return m
}

// IsValid 验证消息是否有效
func (m *Message) IsValid() bool {
	return len(m.Value) > 0 && m.Topic != ""
}

// Size 计算消息大小（字节）
func (m *Message) Size() int {
	size := len(m.Key) + len(m.Value) + len(m.MessageID)
	for k, v := range m.Headers {
		size += len(k) + len(v)
	}
	return size
}
