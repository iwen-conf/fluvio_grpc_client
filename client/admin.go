package client

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
	"github.com/iwen-conf/fluvio_grpc_client/types"

	"google.golang.org/grpc"
)

// AdminManager 管理功能
type AdminManager struct {
	client *Client
}

// NewAdminManager 创建管理功能
func NewAdminManager(client *Client) *AdminManager {
	return &AdminManager{
		client: client,
	}
}

// DescribeCluster 描述集群
func (am *AdminManager) DescribeCluster(ctx context.Context) (*types.DescribeClusterResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	am.client.logger.Debug("获取集群信息")

	var result *types.DescribeClusterResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioAdminServiceClient(conn)

		resp, err := client.DescribeCluster(ctx, &pb.DescribeClusterRequest{})
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取集群信息失败", err)
		}

		// 转换集群信息
		clusterInfo := &types.ClusterInfo{
			Status:       resp.GetStatus(),
			ControllerID: resp.GetControllerId(),
			Metadata:     make(map[string]string),
		}

		result = &types.DescribeClusterResult{
			Cluster: clusterInfo,
			Success: true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取集群信息失败", logger.Field{Key: "error", Value: err})
		return &types.DescribeClusterResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("获取集群信息成功",
		logger.Field{Key: "status", Value: result.Cluster.Status},
		logger.Field{Key: "controller_id", Value: result.Cluster.ControllerID})

	return result, nil
}

// ListBrokers 列出Broker
func (am *AdminManager) ListBrokers(ctx context.Context) (*types.ListBrokersResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	am.client.logger.Debug("获取Broker列表")

	var result *types.ListBrokersResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioAdminServiceClient(conn)

		resp, err := client.ListBrokers(ctx, &pb.ListBrokersRequest{})
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取Broker列表失败", err)
		}

		// 转换Broker信息
		var brokers []*types.BrokerInfo
		for _, broker := range resp.GetBrokers() {
			brokerInfo := &types.BrokerInfo{
				ID:       broker.GetId(),
				Host:     broker.GetHost(),
				Port:     broker.GetPort(),
				Status:   broker.GetStatus(),
				Metadata: make(map[string]string),
			}
			brokers = append(brokers, brokerInfo)
		}

		result = &types.ListBrokersResult{
			Brokers: brokers,
			Success: true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取Broker列表失败", logger.Field{Key: "error", Value: err})
		return &types.ListBrokersResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("获取Broker列表成功",
		logger.Field{Key: "count", Value: len(result.Brokers)})

	return result, nil
}

// GetMetrics 获取指标
func (am *AdminManager) GetMetrics(ctx context.Context, opts types.GetMetricsOptions) (*types.GetMetricsResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	req := &pb.GetMetricsRequest{
		MetricNames: opts.MetricNames,
		Labels:      opts.Labels,
	}

	am.client.logger.Debug("获取指标",
		logger.Field{Key: "metric_names", Value: opts.MetricNames})

	var result *types.GetMetricsResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioAdminServiceClient(conn)

		resp, err := client.GetMetrics(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取指标失败", err)
		}

		// 转换指标信息
		var metrics []*types.MetricInfo
		for _, metric := range resp.GetMetrics() {
			metricInfo := &types.MetricInfo{
				Name:      metric.GetName(),
				Value:     metric.GetValue(),
				Unit:      metric.GetUnit(),
				Labels:    metric.GetLabels(),
				Timestamp: metric.GetTimestamp(),
			}
			metrics = append(metrics, metricInfo)
		}

		result = &types.GetMetricsResult{
			Metrics: metrics,
			Success: true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取指标失败", logger.Field{Key: "error", Value: err})
		return &types.GetMetricsResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("获取指标成功",
		logger.Field{Key: "count", Value: len(result.Metrics)})

	return result, nil
}

// ListConsumerGroups 列出消费者组
func (am *AdminManager) ListConsumerGroups(ctx context.Context) (*types.ListConsumerGroupsResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	am.client.logger.Debug("获取消费者组列表")

	var result *types.ListConsumerGroupsResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.ListConsumerGroups(ctx, &pb.ListConsumerGroupsRequest{})
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取消费者组列表失败", err)
		}

		// 转换消费者组信息
		var groups []*types.ConsumerGroupInfo
		for _, group := range resp.GetGroups() {
			groupInfo := &types.ConsumerGroupInfo{
				Name:     group.GetName(),
				Members:  []*types.ConsumerGroupMember{},
				Offsets:  make(map[string]int64),
				Metadata: make(map[string]string),
			}

			// 转换成员信息
			for _, member := range group.GetMembers() {
				memberInfo := &types.ConsumerGroupMember{
					ID:       member.GetId(),
					Host:     member.GetHost(),
					Topics:   member.GetTopics(),
					Metadata: make(map[string]string),
				}
				groupInfo.Members = append(groupInfo.Members, memberInfo)
			}

			// 转换偏移量信息
			for topic, offset := range group.GetOffsets() {
				groupInfo.Offsets[topic] = offset
			}

			groups = append(groups, groupInfo)
		}

		result = &types.ListConsumerGroupsResult{
			Groups:  groups,
			Success: true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取消费者组列表失败", logger.Field{Key: "error", Value: err})
		return &types.ListConsumerGroupsResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("获取消费者组列表成功",
		logger.Field{Key: "count", Value: len(result.Groups)})

	return result, nil
}

// DescribeConsumerGroup 描述消费者组
func (am *AdminManager) DescribeConsumerGroup(ctx context.Context, groupName string) (*types.DescribeConsumerGroupResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if groupName == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "消费者组名称不能为空")
	}

	req := &pb.DescribeConsumerGroupRequest{
		Group: groupName,
	}

	am.client.logger.Debug("描述消费者组", logger.Field{Key: "group", Value: groupName})

	var result *types.DescribeConsumerGroupResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.DescribeConsumerGroup(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "描述消费者组失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换消费者组信息
		var groupInfo *types.ConsumerGroupInfo
		if resp.GetGroup() != nil {
			group := resp.GetGroup()
			groupInfo = &types.ConsumerGroupInfo{
				Name:     group.GetName(),
				Members:  []*types.ConsumerGroupMember{},
				Offsets:  make(map[string]int64),
				Metadata: make(map[string]string),
			}

			// 转换成员信息
			for _, member := range group.GetMembers() {
				memberInfo := &types.ConsumerGroupMember{
					ID:       member.GetId(),
					Host:     member.GetHost(),
					Topics:   member.GetTopics(),
					Metadata: make(map[string]string),
				}
				groupInfo.Members = append(groupInfo.Members, memberInfo)
			}

			// 转换偏移量信息
			for topic, offset := range group.GetOffsets() {
				groupInfo.Offsets[topic] = offset
			}
		}

		result = &types.DescribeConsumerGroupResult{
			Group:   groupInfo,
			Success: true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("描述消费者组失败",
			logger.Field{Key: "group", Value: groupName},
			logger.Field{Key: "error", Value: err})
		return &types.DescribeConsumerGroupResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("描述消费者组成功", logger.Field{Key: "group", Value: groupName})

	return result, nil
}

// CreateSmartModule 创建SmartModule
func (am *AdminManager) CreateSmartModule(ctx context.Context, opts types.CreateSmartModuleOptions) (*types.CreateSmartModuleResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "SmartModule名称不能为空")
	}

	// 构建SmartModule规格
	spec := &pb.SmartModuleSpec{
		Name:        opts.Name,
		Description: opts.Description,
		Version:     "1.0.0", // 默认版本
	}

	req := &pb.CreateSmartModuleRequest{
		Spec: spec,
	}

	am.client.logger.Debug("创建SmartModule",
		logger.Field{Key: "name", Value: opts.Name},
		logger.Field{Key: "description", Value: opts.Description})

	var result *types.CreateSmartModuleResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.CreateSmartModule(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "创建SmartModule失败", err)
		}

		result = &types.CreateSmartModuleResult{
			Success: resp.GetSuccess(),
			Error:   resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("创建SmartModule失败",
			logger.Field{Key: "name", Value: opts.Name},
			logger.Field{Key: "error", Value: err})
		return &types.CreateSmartModuleResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("SmartModule创建成功", logger.Field{Key: "name", Value: opts.Name})

	return result, nil
}

// DeleteSmartModule 删除SmartModule
func (am *AdminManager) DeleteSmartModule(ctx context.Context, name string) (*types.DeleteSmartModuleResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "SmartModule名称不能为空")
	}

	req := &pb.DeleteSmartModuleRequest{
		Name: name,
	}

	am.client.logger.Debug("删除SmartModule", logger.Field{Key: "name", Value: name})

	var result *types.DeleteSmartModuleResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.DeleteSmartModule(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "删除SmartModule失败", err)
		}

		result = &types.DeleteSmartModuleResult{
			Success: resp.GetSuccess(),
			Error:   resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("删除SmartModule失败",
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "error", Value: err})
		return &types.DeleteSmartModuleResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("SmartModule删除成功", logger.Field{Key: "name", Value: name})

	return result, nil
}

// ListSmartModules 列出SmartModule
func (am *AdminManager) ListSmartModules(ctx context.Context) (*types.ListSmartModulesResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	am.client.logger.Debug("获取SmartModule列表")

	var result *types.ListSmartModulesResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.ListSmartModules(ctx, &pb.ListSmartModulesRequest{})
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取SmartModule列表失败", err)
		}

		// 转换SmartModule信息
		var smartModules []*types.SmartModuleInfo
		for _, sm := range resp.GetModules() {
			smInfo := &types.SmartModuleInfo{
				Name:        sm.GetName(),
				Version:     sm.GetVersion(),
				Description: sm.GetDescription(),
				Metadata:    make(map[string]string), // 初始化空的metadata
			}
			smartModules = append(smartModules, smInfo)
		}

		result = &types.ListSmartModulesResult{
			SmartModules: smartModules,
			Success:      true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取SmartModule列表失败", logger.Field{Key: "error", Value: err})
		return &types.ListSmartModulesResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("获取SmartModule列表成功",
		logger.Field{Key: "count", Value: len(result.SmartModules)})

	return result, nil
}

// DescribeSmartModule 描述SmartModule
func (am *AdminManager) DescribeSmartModule(ctx context.Context, name string) (*types.DescribeSmartModuleResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "SmartModule名称不能为空")
	}

	req := &pb.DescribeSmartModuleRequest{
		Name: name,
	}

	am.client.logger.Debug("描述SmartModule", logger.Field{Key: "name", Value: name})

	var result *types.DescribeSmartModuleResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.DescribeSmartModule(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "描述SmartModule失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换SmartModule信息
		var smInfo *types.SmartModuleInfo
		if resp.GetSpec() != nil {
			sm := resp.GetSpec()
			smInfo = &types.SmartModuleInfo{
				Name:        sm.GetName(),
				Version:     sm.GetVersion(),
				Description: sm.GetDescription(),
				Metadata:    make(map[string]string), // 初始化空的metadata
			}
		}

		result = &types.DescribeSmartModuleResult{
			SmartModule: smInfo,
			Success:     true,
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("描述SmartModule失败",
			logger.Field{Key: "name", Value: name},
			logger.Field{Key: "error", Value: err})
		return &types.DescribeSmartModuleResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("描述SmartModule成功", logger.Field{Key: "name", Value: name})

	return result, nil
}
