package fluvio

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// AdminManager 管理器
type AdminManager struct {
	appService *services.FluvioApplicationService
	logger     logging.Logger
	connected  *bool
}

// ClusterInfo 集群信息
type ClusterInfo struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ControllerID int32  `json:"controller_id"`
}

// BrokerInfo Broker信息
type BrokerInfo struct {
	ID     int32  `json:"id"`
	Host   string `json:"host"`
	Port   int32  `json:"port"`
	Status string `json:"status"`
}

// ConsumerGroupInfo 消费者组信息
type ConsumerGroupInfo struct {
	GroupID string `json:"group_id"`
	State   string `json:"state"`
}

// SmartModuleInfo SmartModule信息
type SmartModuleInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

// ClusterInfo 获取集群信息
func (a *AdminManager) ClusterInfo(ctx context.Context) (*ClusterInfo, error) {
	if !*a.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	a.logger.Debug("Getting cluster info")

	req := &dtos.DescribeClusterRequest{}

	resp, err := a.appService.DescribeCluster(ctx, req)
	if err != nil {
		a.logger.Error("Failed to get cluster info", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(errors.ErrOperation, resp.Error)
	}

	info := &ClusterInfo{
		ID:           resp.Cluster.ID,
		Status:       resp.Cluster.Status,
		ControllerID: resp.Cluster.ControllerID,
	}

	a.logger.Info("Cluster info retrieved successfully")
	return info, nil
}

// Brokers 获取Broker列表
func (a *AdminManager) Brokers(ctx context.Context) ([]*BrokerInfo, error) {
	if !*a.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	a.logger.Debug("Getting broker list")

	req := &dtos.ListBrokersRequest{}

	resp, err := a.appService.ListBrokers(ctx, req)
	if err != nil {
		a.logger.Error("Failed to get broker list", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(errors.ErrOperation, resp.Error)
	}

	var brokers []*BrokerInfo
	for _, broker := range resp.Brokers {
		brokerInfo := &BrokerInfo{
			ID:     broker.ID,
			Host:   broker.Host,
			Port:   broker.Port,
			Status: broker.Status,
		}
		brokers = append(brokers, brokerInfo)
	}

	a.logger.Info("Broker list retrieved successfully", logging.Field{Key: "count", Value: len(brokers)})
	return brokers, nil
}

// ConsumerGroups 获取消费者组列表
func (a *AdminManager) ConsumerGroups(ctx context.Context) ([]*ConsumerGroupInfo, error) {
	if !*a.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	a.logger.Debug("Getting consumer group list")

	req := &dtos.ListConsumerGroupsRequest{}

	resp, err := a.appService.ListConsumerGroups(ctx, req)
	if err != nil {
		a.logger.Error("Failed to get consumer group list", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(errors.ErrOperation, resp.Error)
	}

	var groups []*ConsumerGroupInfo
	for _, group := range resp.Groups {
		groupInfo := &ConsumerGroupInfo{
			GroupID: group.GroupID,
			State:   group.State,
		}
		groups = append(groups, groupInfo)
	}

	a.logger.Info("Consumer group list retrieved successfully", logging.Field{Key: "count", Value: len(groups)})
	return groups, nil
}

// SmartModules 获取SmartModule管理器
func (a *AdminManager) SmartModules() *SmartModuleManager {
	return &SmartModuleManager{
		appService: a.appService,
		logger:     a.logger,
		connected:  a.connected,
	}
}

// SmartModuleManager SmartModule管理器
type SmartModuleManager struct {
	appService *services.FluvioApplicationService
	logger     logging.Logger
	connected  *bool
}

// List 列出SmartModule
func (s *SmartModuleManager) List(ctx context.Context) ([]*SmartModuleInfo, error) {
	if !*s.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	s.logger.Debug("Getting SmartModule list")

	req := &dtos.ListSmartModulesRequest{}

	resp, err := s.appService.ListSmartModules(ctx, req)
	if err != nil {
		s.logger.Error("Failed to get SmartModule list", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(errors.ErrOperation, resp.Error)
	}

	var modules []*SmartModuleInfo
	for _, module := range resp.Modules {
		moduleInfo := &SmartModuleInfo{
			Name:        module.Name,
			Version:     module.Version,
			Description: module.Description,
		}
		modules = append(modules, moduleInfo)
	}

	s.logger.Info("SmartModule list retrieved successfully", logging.Field{Key: "count", Value: len(modules)})
	return modules, nil
}

// Create 创建SmartModule
func (s *SmartModuleManager) Create(ctx context.Context, name string, wasmCode []byte) error {
	if !*s.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}

	s.logger.Debug("Creating SmartModule", logging.Field{Key: "name", Value: name})

	req := &dtos.CreateSmartModuleRequest{
		Name:     name,
		WasmCode: wasmCode,
	}

	resp, err := s.appService.CreateSmartModule(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create SmartModule", logging.Field{Key: "error", Value: err})
		return err
	}

	if !resp.Success {
		return errors.New(errors.ErrOperation, resp.Error)
	}

	s.logger.Info("SmartModule created successfully", logging.Field{Key: "name", Value: name})
	return nil
}

// Delete 删除SmartModule
func (s *SmartModuleManager) Delete(ctx context.Context, name string) error {
	if !*s.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}

	s.logger.Debug("Deleting SmartModule", logging.Field{Key: "name", Value: name})

	req := &dtos.DeleteSmartModuleRequest{
		Name: name,
	}

	resp, err := s.appService.DeleteSmartModule(ctx, req)
	if err != nil {
		s.logger.Error("Failed to delete SmartModule", logging.Field{Key: "error", Value: err})
		return err
	}

	if !resp.Success {
		return errors.New(errors.ErrOperation, resp.Error)
	}

	s.logger.Info("SmartModule deleted successfully", logging.Field{Key: "name", Value: name})
	return nil
}