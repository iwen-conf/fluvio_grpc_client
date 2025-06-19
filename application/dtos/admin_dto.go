package dtos

import "time"

// 集群相关DTO

// DescribeClusterRequest 描述集群请求
type DescribeClusterRequest struct{}

// DescribeClusterResponse 描述集群响应
type DescribeClusterResponse struct {
	Cluster *ClusterDTO `json:"cluster"`
	Error   string      `json:"error,omitempty"`
}

// ClusterDTO 集群信息
type ClusterDTO struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ControllerID int32  `json:"controller_id"`
}

// Broker相关DTO

// ListBrokersRequest 列出Broker请求
type ListBrokersRequest struct{}

// ListBrokersResponse 列出Broker响应
type ListBrokersResponse struct {
	Brokers []*BrokerDTO `json:"brokers"`
	Error   string       `json:"error,omitempty"`
}

// BrokerDTO Broker信息
type BrokerDTO struct {
	ID     int32  `json:"id"`
	Host   string `json:"host"`
	Port   int32  `json:"port"`
	Status string `json:"status"`
	Addr   string `json:"addr"`
}

// 消费者组相关DTO

// ListConsumerGroupsRequest 列出消费者组请求
type ListConsumerGroupsRequest struct{}

// ListConsumerGroupsResponse 列出消费者组响应
type ListConsumerGroupsResponse struct {
	Groups []*ConsumerGroupDTO `json:"groups"`
	Error  string              `json:"error,omitempty"`
}

// DescribeConsumerGroupRequest 描述消费者组请求
type DescribeConsumerGroupRequest struct {
	GroupID string `json:"group_id"`
}

// DescribeConsumerGroupResponse 描述消费者组响应
type DescribeConsumerGroupResponse struct {
	Group *ConsumerGroupDTO `json:"group"`
	Error string            `json:"error,omitempty"`
}

// ConsumerGroupDTO 消费者组信息
type ConsumerGroupDTO struct {
	GroupID string                     `json:"group_id"`
	State   string                     `json:"state"`
	Members []*ConsumerGroupMemberDTO `json:"members,omitempty"`
}

// ConsumerGroupMemberDTO 消费者组成员信息
type ConsumerGroupMemberDTO struct {
	MemberID   string `json:"member_id"`
	ClientID   string `json:"client_id"`
	ClientHost string `json:"client_host"`
}

// SmartModule相关DTO

// ListSmartModulesRequest 列出SmartModule请求
type ListSmartModulesRequest struct{}

// ListSmartModulesResponse 列出SmartModule响应
type ListSmartModulesResponse struct {
	Modules []*SmartModuleDTO `json:"modules"`
	Error   string            `json:"error,omitempty"`
}

// CreateSmartModuleRequest 创建SmartModule请求
type CreateSmartModuleRequest struct {
	Name     string `json:"name"`
	WasmCode []byte `json:"wasm_code"`
}

// CreateSmartModuleResponse 创建SmartModule响应
type CreateSmartModuleResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// DeleteSmartModuleRequest 删除SmartModule请求
type DeleteSmartModuleRequest struct {
	Name string `json:"name"`
}

// DeleteSmartModuleResponse 删除SmartModule响应
type DeleteSmartModuleResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// DescribeSmartModuleRequest 描述SmartModule请求
type DescribeSmartModuleRequest struct {
	Name string `json:"name"`
}

// DescribeSmartModuleResponse 描述SmartModule响应
type DescribeSmartModuleResponse struct {
	Module *SmartModuleDTO `json:"module"`
	Error  string          `json:"error,omitempty"`
}

// SmartModuleDTO SmartModule信息
type SmartModuleDTO struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}