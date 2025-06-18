package types

// ConsumerGroupInfo 消费者组信息
type ConsumerGroupInfo struct {
	GroupID  string                 `json:"group_id"`
	Members  []*ConsumerGroupMember `json:"members"`  // 保持向后兼容
	Offsets  map[string]int64       `json:"offsets"`  // 保持向后兼容
	Metadata map[string]string      `json:"metadata"` // 保持向后兼容
}

// ConsumerGroupMember 消费者组成员
type ConsumerGroupMember struct {
	ID       string            `json:"id"`
	Host     string            `json:"host"`
	Topics   []string          `json:"topics"`
	Metadata map[string]string `json:"metadata"`
}

// ConsumerGroupOffsetInfo 消费组在特定分区的位点信息
type ConsumerGroupOffsetInfo struct {
	Topic           string `json:"topic"`
	Partition       int32  `json:"partition"`
	CommittedOffset int64  `json:"committed_offset"`
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

// DescribeConsumerGroupDetailResult 描述消费者组详细结果（新版本）
type DescribeConsumerGroupDetailResult struct {
	GroupID string                     `json:"group_id"`
	Offsets []*ConsumerGroupOffsetInfo `json:"offsets"`
	Success bool                       `json:"success"`
	Error   string                     `json:"error,omitempty"`
}
