package types

// ClusterInfo 集群信息
type ClusterInfo struct {
	Status       string            `json:"status"`
	ControllerID int32             `json:"controller_id"`
	Brokers      []*BrokerInfo     `json:"brokers"`
	Metadata     map[string]string `json:"metadata"`
}

// BrokerInfo Broker信息
type BrokerInfo struct {
	ID       int32             `json:"id"`
	Host     string            `json:"host"`
	Port     int32             `json:"port"`
	Status   string            `json:"status"`
	Metadata map[string]string `json:"metadata"`
}

// MetricInfo 指标信息
type MetricInfo struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Unit      string            `json:"unit"`
	Labels    map[string]string `json:"labels"`
	Timestamp int64             `json:"timestamp"`
}

// DescribeClusterResult 描述集群结果
type DescribeClusterResult struct {
	Cluster *ClusterInfo `json:"cluster"`
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
}

// ListBrokersResult 列出Broker结果
type ListBrokersResult struct {
	Brokers []*BrokerInfo `json:"brokers"`
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
}

// GetMetricsOptions 获取指标选项
type GetMetricsOptions struct {
	MetricNames []string          `json:"metric_names"`
	Labels      map[string]string `json:"labels"`
}

// GetMetricsResult 获取指标结果
type GetMetricsResult struct {
	Metrics []*MetricInfo `json:"metrics"`
	Success bool          `json:"success"`
	Error   string        `json:"error,omitempty"`
}

// SmartModuleInfo SmartModule信息
type SmartModuleInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
}

// CreateSmartModuleOptions 创建SmartModule选项
type CreateSmartModuleOptions struct {
	Name        string            `json:"name"`
	WasmPath    string            `json:"wasm_path"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata"`
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

// ListSmartModulesResult 列出SmartModule结果
type ListSmartModulesResult struct {
	SmartModules []*SmartModuleInfo `json:"smart_modules"`
	Success      bool               `json:"success"`
	Error        string             `json:"error,omitempty"`
}

// DescribeSmartModuleResult 描述SmartModule结果
type DescribeSmartModuleResult struct {
	SmartModule *SmartModuleInfo `json:"smart_module"`
	Success     bool             `json:"success"`
	Error       string           `json:"error,omitempty"`
}
