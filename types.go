package fluvio

import "time"

// 重新导出核心类型，保持API一致性

// MessageHeader 消息头部
type MessageHeader map[string]string

// MessageMetadata 消息元数据
type MessageMetadata struct {
	Offset    int64     `json:"offset"`
	Partition int32     `json:"partition"`
	Topic     string    `json:"topic"`
	Timestamp time.Time `json:"timestamp"`
}

// ProduceOptions 生产选项（别名，保持兼容性）
type ProduceOptions = SendOptions

// ProduceResult 生产结果（别名，保持兼容性）
type ProduceResult = SendResult

// BatchProduceResult 批量生产结果（别名，保持兼容性）
type BatchProduceResult = BatchSendResult

// ConsumeOptions 消费选项（别名，保持兼容性）
type ConsumeOptions = ReceiveOptions

// ConsumeResult 消费结果
type ConsumeResult struct {
	Messages []*ConsumedMessage `json:"messages"`
	Count    int                `json:"count"`
}

// StreamConsumeOptions 流式消费选项（别名，保持兼容性）
type StreamConsumeOptions = StreamOptions

// TopicDescription 主题描述（别名，保持兼容性）
type TopicDescription = TopicInfo

// CreateTopicResult 创建主题结果
type CreateTopicResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// DeleteTopicResult 删除主题结果
type DeleteTopicResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// ListTopicsResult 主题列表结果
type ListTopicsResult struct {
	Topics []string `json:"topics"`
}

// DescribeTopicResult 主题描述结果
type DescribeTopicResult struct {
	Topic *TopicInfo `json:"topic"`
	Error string     `json:"error,omitempty"`
}

// DescribeClusterResult 集群描述结果
type DescribeClusterResult struct {
	Cluster *ClusterInfo `json:"cluster"`
	Error   string       `json:"error,omitempty"`
}

// ListBrokersResult Broker列表结果
type ListBrokersResult struct {
	Brokers []*BrokerInfo `json:"brokers"`
	Error   string        `json:"error,omitempty"`
}

// ListConsumerGroupsResult 消费者组列表结果
type ListConsumerGroupsResult struct {
	Groups []*ConsumerGroupInfo `json:"groups"`
	Error  string               `json:"error,omitempty"`
}

// DescribeConsumerGroupResult 消费者组描述结果
type DescribeConsumerGroupResult struct {
	Group *ConsumerGroupInfo `json:"group"`
	Error string             `json:"error,omitempty"`
}

// CreateSmartModuleOptions 创建SmartModule选项
type CreateSmartModuleOptions struct {
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	WasmCode    []byte `json:"wasm_code"`
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

// ListSmartModulesResult SmartModule列表结果
type ListSmartModulesResult struct {
	Modules []*SmartModuleInfo `json:"modules"`
	Error   string             `json:"error,omitempty"`
}

// DescribeSmartModuleResult SmartModule描述结果
type DescribeSmartModuleResult struct {
	Module *SmartModuleInfo `json:"module"`
	Error  string           `json:"error,omitempty"`
}

// 错误类型
const (
	ErrConnection      = "CONNECTION_ERROR"
	ErrInvalidArgument = "INVALID_ARGUMENT"
	ErrOperation       = "OPERATION_ERROR"
	ErrTimeout         = "TIMEOUT_ERROR"
	ErrNotFound        = "NOT_FOUND"
	ErrAlreadyExists   = "ALREADY_EXISTS"
	ErrPermission      = "PERMISSION_DENIED"
	ErrInternal        = "INTERNAL_ERROR"
)

// FluvioError Fluvio错误类型
type FluvioError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *FluvioError) Error() string {
	if e.Details != "" {
		return e.Code + ": " + e.Message + " (" + e.Details + ")"
	}
	return e.Code + ": " + e.Message
}

// NewError 创建新错误
func NewError(code, message string) *FluvioError {
	return &FluvioError{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithDetails 创建带详情的新错误
func NewErrorWithDetails(code, message, details string) *FluvioError {
	return &FluvioError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// 便捷函数

// IsConnectionError 检查是否为连接错误
func IsConnectionError(err error) bool {
	if fluvioErr, ok := err.(*FluvioError); ok {
		return fluvioErr.Code == ErrConnection
	}
	return false
}

// IsNotFoundError 检查是否为未找到错误
func IsNotFoundError(err error) bool {
	if fluvioErr, ok := err.(*FluvioError); ok {
		return fluvioErr.Code == ErrNotFound
	}
	return false
}

// IsTimeoutError 检查是否为超时错误
func IsTimeoutError(err error) bool {
	if fluvioErr, ok := err.(*FluvioError); ok {
		return fluvioErr.Code == ErrTimeout
	}
	return false
}