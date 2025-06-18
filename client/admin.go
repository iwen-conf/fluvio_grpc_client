package client

import (
	"context"
	"fmt"
	"time"

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
				Addr:     broker.GetAddr(),
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
			var timestamp time.Time
			if metric.GetTimestamp() != nil {
				timestamp = metric.GetTimestamp().AsTime()
			}
			metricInfo := &types.MetricInfo{
				Name:      metric.GetName(),
				Value:     metric.GetValue(),
				Labels:    metric.GetLabels(),
				Timestamp: timestamp,
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
				GroupID:  group.GetGroupId(),
				Members:  []*types.ConsumerGroupMember{},
				Offsets:  make(map[string]int64),
				Metadata: make(map[string]string),
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
		GroupId: groupName,
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
		groupInfo := &types.ConsumerGroupInfo{
			GroupID:  resp.GetGroupId(),
			Members:  []*types.ConsumerGroupMember{},
			Offsets:  make(map[string]int64),
			Metadata: make(map[string]string),
		}

		// 转换偏移量信息
		for _, offsetInfo := range resp.GetOffsets() {
			key := fmt.Sprintf("%s-%d", offsetInfo.GetTopic(), offsetInfo.GetPartition())
			groupInfo.Offsets[key] = offsetInfo.GetCommittedOffset()
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

	if opts.Spec == nil || opts.Spec.Name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "SmartModule规格不能为空")
	}

	// 构建SmartModule规格
	spec := &pb.SmartModuleSpec{
		Name:        opts.Spec.Name,
		InputKind:   pb.SmartModuleInput(opts.Spec.InputKind),
		OutputKind:  pb.SmartModuleOutput(opts.Spec.OutputKind),
		Description: opts.Spec.Description,
		Version:     opts.Spec.Version,
	}

	// 转换参数定义
	for _, param := range opts.Spec.Parameters {
		pbParam := &pb.SmartModuleParameter{
			Name:        param.Name,
			Description: param.Description,
			Optional:    param.Optional,
		}
		spec.Parameters = append(spec.Parameters, pbParam)
	}

	req := &pb.CreateSmartModuleRequest{
		Spec:     spec,
		WasmCode: opts.WasmCode,
	}

	am.client.logger.Debug("创建SmartModule",
		logger.Field{Key: "name", Value: opts.Spec.Name},
		logger.Field{Key: "description", Value: opts.Spec.Description})

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
			logger.Field{Key: "name", Value: opts.Spec.Name},
			logger.Field{Key: "error", Value: err})
		return &types.CreateSmartModuleResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("SmartModule创建成功", logger.Field{Key: "name", Value: opts.Spec.Name})

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

// UpdateSmartModule 更新SmartModule
func (am *AdminManager) UpdateSmartModule(ctx context.Context, opts types.UpdateSmartModuleOptions) (*types.UpdateSmartModuleResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Name == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "SmartModule名称不能为空")
	}

	req := &pb.UpdateSmartModuleRequest{
		Name: opts.Name,
	}

	// 如果提供了新的规格，设置规格
	if opts.Spec != nil {
		spec := &pb.SmartModuleSpec{
			Name:        opts.Spec.Name,
			InputKind:   pb.SmartModuleInput(opts.Spec.InputKind),
			OutputKind:  pb.SmartModuleOutput(opts.Spec.OutputKind),
			Description: opts.Spec.Description,
			Version:     opts.Spec.Version,
		}

		// 转换参数定义
		for _, param := range opts.Spec.Parameters {
			pbParam := &pb.SmartModuleParameter{
				Name:        param.Name,
				Description: param.Description,
				Optional:    param.Optional,
			}
			spec.Parameters = append(spec.Parameters, pbParam)
		}

		req.Spec = spec
	}

	// 如果提供了新的WASM代码，设置代码
	if len(opts.WasmCode) > 0 {
		req.WasmCode = opts.WasmCode
	}

	am.client.logger.Debug("更新SmartModule", logger.Field{Key: "name", Value: opts.Name})

	var result *types.UpdateSmartModuleResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.UpdateSmartModule(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "更新SmartModule失败", err)
		}

		result = &types.UpdateSmartModuleResult{
			Success: resp.GetSuccess(),
			Error:   resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("更新SmartModule失败",
			logger.Field{Key: "name", Value: opts.Name},
			logger.Field{Key: "error", Value: err})
		return &types.UpdateSmartModuleResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("SmartModule更新成功", logger.Field{Key: "name", Value: opts.Name})

	return result, nil
}

// BulkDelete 批量删除
func (am *AdminManager) BulkDelete(ctx context.Context, opts types.BulkDeleteOptions) (*types.BulkDeleteResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	req := &pb.BulkDeleteRequest{
		Topics:         opts.Topics,
		ConsumerGroups: opts.ConsumerGroups,
		SmartModules:   opts.SmartModules,
		Force:          opts.Force,
	}

	am.client.logger.Debug("批量删除",
		logger.Field{Key: "topics_count", Value: len(opts.Topics)},
		logger.Field{Key: "consumer_groups_count", Value: len(opts.ConsumerGroups)},
		logger.Field{Key: "smart_modules_count", Value: len(opts.SmartModules)},
		logger.Field{Key: "force", Value: opts.Force})

	var result *types.BulkDeleteResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.BulkDelete(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "批量删除失败", err)
		}

		// 转换删除结果
		var results []*types.BulkDeleteItemResult
		for _, item := range resp.GetResults() {
			itemResult := &types.BulkDeleteItemResult{
				Name:    item.GetName(),
				Type:    item.GetType(),
				Success: item.GetSuccess(),
				Error:   item.GetError(),
			}
			results = append(results, itemResult)
		}

		result = &types.BulkDeleteResult{
			Results:           results,
			TotalRequested:    resp.GetTotalRequested(),
			SuccessfulDeletes: resp.GetSuccessfulDeletes(),
			FailedDeletes:     resp.GetFailedDeletes(),
			Success:           resp.GetFailedDeletes() == 0,
			Error:             "",
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("批量删除失败", logger.Field{Key: "error", Value: err})
		return &types.BulkDeleteResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("批量删除完成",
		logger.Field{Key: "total_requested", Value: result.TotalRequested},
		logger.Field{Key: "successful_deletes", Value: result.SuccessfulDeletes},
		logger.Field{Key: "failed_deletes", Value: result.FailedDeletes})

	return result, nil
}

// GetStorageStatus 获取存储状态
func (am *AdminManager) GetStorageStatus(ctx context.Context, opts types.GetStorageStatusOptions) (*types.GetStorageStatusResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	req := &pb.GetStorageStatusRequest{
		IncludeDetails: opts.IncludeDetails,
	}

	am.client.logger.Debug("获取存储状态",
		logger.Field{Key: "include_details", Value: opts.IncludeDetails})

	var result *types.GetStorageStatusResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.GetStorageStatus(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取存储状态失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换存储统计信息
		var storageStats *types.StorageStats
		if resp.GetStorageStats() != nil {
			stats := resp.GetStorageStats()
			storageStats = &types.StorageStats{
				StorageType:      stats.GetStorageType(),
				ConsumerGroups:   stats.GetConsumerGroups(),
				ConsumerOffsets:  stats.GetConsumerOffsets(),
				SmartModules:     stats.GetSmartModules(),
				ConnectionStatus: stats.GetConnectionStatus(),
			}

			// 转换连接统计信息
			if stats.GetConnectionStats() != nil {
				connStats := stats.GetConnectionStats()
				storageStats.ConnectionStats = &types.StorageConnectionStats{
					CurrentConnections:      connStats.GetCurrentConnections(),
					AvailableConnections:    connStats.GetAvailableConnections(),
					TotalCreatedConnections: connStats.GetTotalCreatedConnections(),
				}
			}

			// 转换数据库信息
			if stats.GetDatabaseInfo() != nil {
				dbInfo := stats.GetDatabaseInfo()
				storageStats.DatabaseInfo = &types.StorageDatabaseInfo{
					Name:        dbInfo.GetName(),
					Collections: dbInfo.GetCollections(),
					DataSize:    dbInfo.GetDataSize(),
					StorageSize: dbInfo.GetStorageSize(),
					Indexes:     dbInfo.GetIndexes(),
					IndexSize:   dbInfo.GetIndexSize(),
				}
			}
		}

		result = &types.GetStorageStatusResult{
			PersistenceEnabled: resp.GetPersistenceEnabled(),
			StorageStats:       storageStats,
			Success:            true,
		}

		if resp.GetCheckedAt() != nil {
			result.CheckedAt = resp.GetCheckedAt().AsTime()
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取存储状态失败", logger.Field{Key: "error", Value: err})
		return &types.GetStorageStatusResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("存储状态获取成功",
		logger.Field{Key: "persistence_enabled", Value: result.PersistenceEnabled})

	return result, nil
}

// MigrateStorage 存储迁移
func (am *AdminManager) MigrateStorage(ctx context.Context, opts types.MigrateStorageOptions) (*types.MigrateStorageResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	req := &pb.MigrateStorageRequest{
		SourceType:      opts.SourceType,
		TargetType:      opts.TargetType,
		VerifyMigration: opts.VerifyMigration,
		ForceMigration:  opts.ForceMigration,
	}

	am.client.logger.Debug("存储迁移",
		logger.Field{Key: "source_type", Value: opts.SourceType},
		logger.Field{Key: "target_type", Value: opts.TargetType},
		logger.Field{Key: "verify_migration", Value: opts.VerifyMigration},
		logger.Field{Key: "force_migration", Value: opts.ForceMigration})

	var result *types.MigrateStorageResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.MigrateStorage(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "存储迁移失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换迁移统计信息
		var migrationStats *types.MigrationStats
		if resp.GetMigrationStats() != nil {
			stats := resp.GetMigrationStats()
			migrationStats = &types.MigrationStats{
				ConsumerGroupsMigrated:  stats.GetConsumerGroupsMigrated(),
				ConsumerOffsetsMigrated: stats.GetConsumerOffsetsMigrated(),
				SmartModulesMigrated:    stats.GetSmartModulesMigrated(),
				Errors:                  stats.GetErrors(),
				TotalMigrated:           stats.GetTotalMigrated(),
			}
		}

		result = &types.MigrateStorageResult{
			Success:            resp.GetSuccess(),
			MigrationStats:     migrationStats,
			VerificationPassed: resp.GetVerificationPassed(),
		}

		if resp.GetCompletedAt() != nil {
			result.CompletedAt = resp.GetCompletedAt().AsTime()
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("存储迁移失败", logger.Field{Key: "error", Value: err})
		return &types.MigrateStorageResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("存储迁移完成",
		logger.Field{Key: "success", Value: result.Success},
		logger.Field{Key: "verification_passed", Value: result.VerificationPassed})

	return result, nil
}

// GetStorageMetrics 获取存储指标
func (am *AdminManager) GetStorageMetrics(ctx context.Context, opts types.GetStorageMetricsOptions) (*types.GetStorageMetricsResult, error) {
	if am.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	req := &pb.GetStorageMetricsRequest{
		IncludeHistory: opts.IncludeHistory,
		HistoryLimit:   opts.HistoryLimit,
	}

	am.client.logger.Debug("获取存储指标",
		logger.Field{Key: "include_history", Value: opts.IncludeHistory},
		logger.Field{Key: "history_limit", Value: opts.HistoryLimit})

	var result *types.GetStorageMetricsResult
	err := am.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.GetStorageMetrics(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "获取存储指标失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换当前指标
		var currentMetrics *types.StorageMetrics
		if resp.GetCurrentMetrics() != nil {
			metrics := resp.GetCurrentMetrics()
			currentMetrics = &types.StorageMetrics{
				StorageType:         metrics.GetStorageType(),
				ResponseTimeMs:      metrics.GetResponseTimeMs(),
				OperationsPerSecond: metrics.GetOperationsPerSecond(),
				ErrorRate:           metrics.GetErrorRate(),
				ConnectionPoolUsage: metrics.GetConnectionPoolUsage(),
				MemoryUsageMB:       metrics.GetMemoryUsageMb(),
				DiskUsageMB:         metrics.GetDiskUsageMb(),
			}

			if metrics.GetLastUpdated() != nil {
				currentMetrics.LastUpdated = metrics.GetLastUpdated().AsTime()
			}
		}

		// 转换历史指标
		var metricsHistory []*types.StorageMetrics
		if opts.IncludeHistory {
			for _, metrics := range resp.GetMetricsHistory() {
				historyMetrics := &types.StorageMetrics{
					StorageType:         metrics.GetStorageType(),
					ResponseTimeMs:      metrics.GetResponseTimeMs(),
					OperationsPerSecond: metrics.GetOperationsPerSecond(),
					ErrorRate:           metrics.GetErrorRate(),
					ConnectionPoolUsage: metrics.GetConnectionPoolUsage(),
					MemoryUsageMB:       metrics.GetMemoryUsageMb(),
					DiskUsageMB:         metrics.GetDiskUsageMb(),
				}

				if metrics.GetLastUpdated() != nil {
					historyMetrics.LastUpdated = metrics.GetLastUpdated().AsTime()
				}

				metricsHistory = append(metricsHistory, historyMetrics)
			}
		}

		// 转换健康状态
		var healthStatus *types.StorageHealthCheckResult
		if resp.GetHealthStatus() != nil {
			health := resp.GetHealthStatus()
			healthStatus = &types.StorageHealthCheckResult{
				Status:         health.GetStatus(),
				ResponseTimeMs: health.GetResponseTimeMs(),
				ErrorMessage:   health.GetErrorMessage(),
			}

			if health.GetCheckedAt() != nil {
				healthStatus.CheckedAt = health.GetCheckedAt().AsTime()
			}
		}

		result = &types.GetStorageMetricsResult{
			CurrentMetrics: currentMetrics,
			MetricsHistory: metricsHistory,
			HealthStatus:   healthStatus,
			Alerts:         resp.GetAlerts(),
			Success:        true,
		}

		if resp.GetCollectedAt() != nil {
			result.CollectedAt = resp.GetCollectedAt().AsTime()
		}

		return nil
	})

	if err != nil {
		am.client.logger.Error("获取存储指标失败", logger.Field{Key: "error", Value: err})
		return &types.GetStorageMetricsResult{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	am.client.logger.Info("存储指标获取成功",
		logger.Field{Key: "has_current_metrics", Value: result.CurrentMetrics != nil},
		logger.Field{Key: "history_count", Value: len(result.MetricsHistory)},
		logger.Field{Key: "alerts_count", Value: len(result.Alerts)})

	return result, nil
}
