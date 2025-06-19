package repositories

import (
	"context"
	"fmt"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/grpc"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// GRPCAdminRepository gRPC管理仓储实现
type GRPCAdminRepository struct {
	client grpc.Client
	logger logging.Logger
}

// NewGRPCAdminRepository 创建gRPC管理仓储
func NewGRPCAdminRepository(client grpc.Client, logger logging.Logger) repositories.AdminRepository {
	return &GRPCAdminRepository{
		client: client,
		logger: logger,
	}
}

// DescribeCluster 描述集群
func (r *GRPCAdminRepository) DescribeCluster(ctx context.Context, req *dtos.DescribeClusterRequest) (*dtos.DescribeClusterResponse, error) {
	r.logger.Debug("Describing cluster")

	// 注意：当前gRPC客户端接口中没有DescribeCluster方法
	// 这个方法在FluvioAdminService中定义，但客户端接口只包含FluvioService
	// 这里使用健康检查来判断集群状态
	healthReq := &pb.HealthCheckRequest{}
	healthResp, err := r.client.HealthCheck(ctx, healthReq)
	if err != nil {
		r.logger.Error("健康检查失败", logging.Field{Key: "error", Value: err})
		return &dtos.DescribeClusterResponse{
			Cluster: &dtos.ClusterDTO{
				ID:           "fluvio-cluster",
				Status:       "Unhealthy",
				ControllerID: 0,
			},
			Error: fmt.Sprintf("cluster health check failed: %v", err),
		}, nil
	}

	// 根据健康检查结果判断集群状态
	status := "Healthy"
	if !healthResp.GetOk() {
		status = "Unhealthy"
	}

	r.logger.Debug("集群状态检查成功", logging.Field{Key: "status", Value: status})

	return &dtos.DescribeClusterResponse{
		Cluster: &dtos.ClusterDTO{
			ID:           "fluvio-cluster",
			Status:       status,
			ControllerID: 1, // 默认控制器ID
		},
	}, nil
}

// ListBrokers 列出Broker
func (r *GRPCAdminRepository) ListBrokers(ctx context.Context, req *dtos.ListBrokersRequest) (*dtos.ListBrokersResponse, error) {
	r.logger.Debug("Listing brokers")

	// 由于当前protobuf定义中没有ListBrokers方法，我们使用健康检查来模拟
	// 在实际实现中，应该有专门的ListBrokers gRPC方法
	healthReq := &pb.HealthCheckRequest{}
	_, err := r.client.HealthCheck(ctx, healthReq)
	if err != nil {
		r.logger.Error("健康检查失败", logging.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to check broker health: %w", err)
	}

	// 简化实现：返回当前连接的broker信息
	// 在真实实现中，应该从gRPC响应中获取broker列表
	return &dtos.ListBrokersResponse{
		Brokers: []*dtos.BrokerDTO{
			{
				ID:     1,
				Host:   "101.43.173.154", // 使用实际的服务器地址
				Port:   50051,
				Status: "Running",
				Addr:   "101.43.173.154:50051",
			},
		},
	}, nil
}

// ListConsumerGroups 列出消费者组
func (r *GRPCAdminRepository) ListConsumerGroups(ctx context.Context, req *dtos.ListConsumerGroupsRequest) (*dtos.ListConsumerGroupsResponse, error) {
	r.logger.Debug("Listing consumer groups")

	// 构建gRPC请求
	grpcReq := &pb.ListConsumerGroupsRequest{}

	// 调用gRPC服务
	resp, err := r.client.ListConsumerGroups(ctx, grpcReq)
	if err != nil {
		r.logger.Error("列出消费者组失败", logging.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to list consumer groups: %w", err)
	}

	// 转换响应
	groups := make([]*dtos.ConsumerGroupDTO, len(resp.GetGroups()))
	for i, group := range resp.GetGroups() {
		groups[i] = &dtos.ConsumerGroupDTO{
			GroupID: group.GetGroupId(),
			State:   "Active", // 简化实现：假设所有组都是活跃的
		}
	}

	r.logger.Debug("列出消费者组成功", logging.Field{Key: "count", Value: len(groups)})

	return &dtos.ListConsumerGroupsResponse{
		Groups: groups,
	}, nil
}

// DescribeConsumerGroup 描述消费者组
func (r *GRPCAdminRepository) DescribeConsumerGroup(ctx context.Context, req *dtos.DescribeConsumerGroupRequest) (*dtos.DescribeConsumerGroupResponse, error) {
	r.logger.Debug("Describing consumer group", logging.Field{Key: "group_id", Value: req.GroupID})

	// 构建gRPC请求
	grpcReq := &pb.DescribeConsumerGroupRequest{
		GroupId: req.GroupID,
	}

	// 调用gRPC服务
	resp, err := r.client.DescribeConsumerGroup(ctx, grpcReq)
	if err != nil {
		r.logger.Error("描述消费者组失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "group_id", Value: req.GroupID})
		return nil, fmt.Errorf("failed to describe consumer group: %w", err)
	}

	// 检查错误
	if resp.GetError() != "" {
		return &dtos.DescribeConsumerGroupResponse{
			Error: resp.GetError(),
		}, nil
	}

	r.logger.Debug("描述消费者组成功",
		logging.Field{Key: "group_id", Value: req.GroupID},
		logging.Field{Key: "offsets_count", Value: len(resp.GetOffsets())})

	// 简化实现：由于protobuf定义中没有成员信息，我们返回空的成员列表
	return &dtos.DescribeConsumerGroupResponse{
		Group: &dtos.ConsumerGroupDTO{
			GroupID: req.GroupID,
			State:   "Active",                         // 简化实现
			Members: []*dtos.ConsumerGroupMemberDTO{}, // 空成员列表
		},
	}, nil
}

// ListSmartModules 列出SmartModule
func (r *GRPCAdminRepository) ListSmartModules(ctx context.Context, req *dtos.ListSmartModulesRequest) (*dtos.ListSmartModulesResponse, error) {
	r.logger.Debug("Listing SmartModules")

	// 构建gRPC请求
	grpcReq := &pb.ListSmartModulesRequest{}

	// 调用gRPC服务
	resp, err := r.client.ListSmartModules(ctx, grpcReq)
	if err != nil {
		r.logger.Error("列出SmartModule失败", logging.Field{Key: "error", Value: err})
		return nil, fmt.Errorf("failed to list smart modules: %w", err)
	}

	// 检查错误
	if resp.GetError() != "" {
		return &dtos.ListSmartModulesResponse{
			Error: resp.GetError(),
		}, nil
	}

	// 转换响应
	modules := make([]*dtos.SmartModuleDTO, len(resp.GetModules()))
	for i, module := range resp.GetModules() {
		modules[i] = &dtos.SmartModuleDTO{
			Name:        module.GetName(),
			Version:     module.GetVersion(),
			Description: module.GetDescription(),
		}
	}

	r.logger.Debug("列出SmartModule成功", logging.Field{Key: "count", Value: len(modules)})

	return &dtos.ListSmartModulesResponse{
		Modules: modules,
	}, nil
}

// CreateSmartModule 创建SmartModule
func (r *GRPCAdminRepository) CreateSmartModule(ctx context.Context, req *dtos.CreateSmartModuleRequest) (*dtos.CreateSmartModuleResponse, error) {
	r.logger.Debug("Creating SmartModule", logging.Field{Key: "name", Value: req.Name})

	// 构建gRPC请求
	grpcReq := &pb.CreateSmartModuleRequest{
		Spec: &pb.SmartModuleSpec{
			Name:        req.Name,
			InputKind:   pb.SmartModuleInput_SMART_MODULE_INPUT_STREAM,   // 默认输入类型
			OutputKind:  pb.SmartModuleOutput_SMART_MODULE_OUTPUT_STREAM, // 默认输出类型
			Description: "SmartModule created via API",
			Version:     "1.0.0",
		},
		WasmCode: req.WasmCode,
	}

	// 调用gRPC服务
	resp, err := r.client.CreateSmartModule(ctx, grpcReq)
	if err != nil {
		r.logger.Error("创建SmartModule失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "name", Value: req.Name})
		return nil, fmt.Errorf("failed to create smart module: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("SmartModule创建被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "name", Value: req.Name})
		return &dtos.CreateSmartModuleResponse{
			Success: false,
			Error:   errMsg,
		}, nil
	}

	r.logger.Info("SmartModule创建成功", logging.Field{Key: "name", Value: req.Name})

	return &dtos.CreateSmartModuleResponse{
		Success: true,
	}, nil
}

// DeleteSmartModule 删除SmartModule
func (r *GRPCAdminRepository) DeleteSmartModule(ctx context.Context, req *dtos.DeleteSmartModuleRequest) (*dtos.DeleteSmartModuleResponse, error) {
	r.logger.Debug("Deleting SmartModule", logging.Field{Key: "name", Value: req.Name})

	// 构建gRPC请求
	grpcReq := &pb.DeleteSmartModuleRequest{
		Name: req.Name,
	}

	// 调用gRPC服务
	resp, err := r.client.DeleteSmartModule(ctx, grpcReq)
	if err != nil {
		r.logger.Error("删除SmartModule失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "name", Value: req.Name})
		return nil, fmt.Errorf("failed to delete smart module: %w", err)
	}

	// 检查响应状态
	if !resp.GetSuccess() {
		errMsg := resp.GetError()
		if errMsg == "" {
			errMsg = "unknown error"
		}
		r.logger.Error("SmartModule删除被服务器拒绝",
			logging.Field{Key: "error", Value: errMsg},
			logging.Field{Key: "name", Value: req.Name})
		return &dtos.DeleteSmartModuleResponse{
			Success: false,
			Error:   errMsg,
		}, nil
	}

	r.logger.Info("SmartModule删除成功", logging.Field{Key: "name", Value: req.Name})

	return &dtos.DeleteSmartModuleResponse{
		Success: true,
	}, nil
}

// DescribeSmartModule 描述SmartModule
func (r *GRPCAdminRepository) DescribeSmartModule(ctx context.Context, req *dtos.DescribeSmartModuleRequest) (*dtos.DescribeSmartModuleResponse, error) {
	r.logger.Debug("Describing SmartModule", logging.Field{Key: "name", Value: req.Name})

	// 构建gRPC请求
	grpcReq := &pb.DescribeSmartModuleRequest{
		Name: req.Name,
	}

	// 调用gRPC服务
	resp, err := r.client.DescribeSmartModule(ctx, grpcReq)
	if err != nil {
		r.logger.Error("描述SmartModule失败",
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "name", Value: req.Name})
		return nil, fmt.Errorf("failed to describe smart module: %w", err)
	}

	// 检查错误
	if resp.GetError() != "" {
		return &dtos.DescribeSmartModuleResponse{
			Error: resp.GetError(),
		}, nil
	}

	// 转换响应
	spec := resp.GetSpec()
	if spec == nil {
		return &dtos.DescribeSmartModuleResponse{
			Error: "SmartModule not found",
		}, nil
	}

	r.logger.Debug("描述SmartModule成功", logging.Field{Key: "name", Value: req.Name})

	return &dtos.DescribeSmartModuleResponse{
		Module: &dtos.SmartModuleDTO{
			Name:        spec.GetName(),
			Version:     spec.GetVersion(),
			Description: spec.GetDescription(),
		},
	}, nil
}
