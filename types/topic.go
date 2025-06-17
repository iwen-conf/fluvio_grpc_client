package types

// TopicInfo 主题信息
type TopicInfo struct {
	Name       string `json:"name"`
	Partitions int32  `json:"partitions"`
	Replicas   int32  `json:"replicas"`
}

// CreateTopicOptions 创建主题选项
type CreateTopicOptions struct {
	Name       string `json:"name"`
	Partitions int32  `json:"partitions"`
	Replicas   int32  `json:"replicas"`
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
