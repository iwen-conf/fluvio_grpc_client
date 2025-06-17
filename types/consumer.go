package types

// ConsumerGroupInfo 消费者组信息
type ConsumerGroupInfo struct {
	GroupID  string                 `json:"group_id"`
	Members  []*ConsumerGroupMember `json:"members"`
	Offsets  map[string]int64       `json:"offsets"`
	Metadata map[string]string      `json:"metadata"`
}

// ConsumerGroupMember 消费者组成员
type ConsumerGroupMember struct {
	ID       string            `json:"id"`
	Host     string            `json:"host"`
	Topics   []string          `json:"topics"`
	Metadata map[string]string `json:"metadata"`
}

// ListConsumerGroupsResult 列出消费者组结果
type ListConsumerGroupsResult struct {
	Groups  []*ConsumerGroupInfo `json:"groups"`
	Success bool                 `json:"success"`
	Error   string               `json:"error,omitempty"`
}

// DescribeConsumerGroupResult 描述消费者组结果
type DescribeConsumerGroupResult struct {
	Group   *ConsumerGroupInfo `json:"group"`
	Success bool               `json:"success"`
	Error   string             `json:"error,omitempty"`
}
