package grpc

import (
	"context"
	
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// Client gRPC客户端接口
type Client interface {
	// 消息操作
	Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error)
	ProduceBatch(ctx context.Context, req *pb.ProduceBatchRequest) (*pb.ProduceBatchReply, error)
	Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error)
	ConsumeFiltered(ctx context.Context, req *pb.ConsumeFilteredRequest) (*pb.ConsumeFilteredReply, error)
	ConsumeStream(ctx context.Context, req *pb.ConsumeStreamRequest) (pb.FluvioService_ConsumeStreamClient, error)
	
	// 主题操作
	CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error)
	DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error)
	ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error)
	DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error)
	DescribeTopicDetail(ctx context.Context, req *pb.DescribeTopicDetailRequest) (*pb.DescribeTopicDetailReply, error)
	GetTopicStats(ctx context.Context, req *pb.GetTopicStatsRequest) (*pb.GetTopicStatsReply, error)
	
	// 消费组操作
	ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error)
	DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error)
	DeleteConsumerGroup(ctx context.Context, req *pb.DeleteConsumerGroupRequest) (*pb.DeleteConsumerGroupReply, error)
	
	// 管理操作
	ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error)
	CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error)
	UpdateSmartModule(ctx context.Context, req *pb.UpdateSmartModuleRequest) (*pb.UpdateSmartModuleReply, error)
	DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error)
	DescribeSmartModule(ctx context.Context, req *pb.DescribeSmartModuleRequest) (*pb.DescribeSmartModuleReply, error)
	
	// 存储管理
	GetStorageStatus(ctx context.Context, req *pb.GetStorageStatusRequest) (*pb.GetStorageStatusReply, error)
	MigrateStorage(ctx context.Context, req *pb.MigrateStorageRequest) (*pb.MigrateStorageReply, error)
	GetStorageMetrics(ctx context.Context, req *pb.GetStorageMetricsRequest) (*pb.GetStorageMetricsReply, error)
	
	// 批量操作
	BulkDelete(ctx context.Context, req *pb.BulkDeleteRequest) (*pb.BulkDeleteReply, error)
	
	// 健康检查
	HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error)
	Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error)
	
	// 连接管理
	Connect() error
	Close() error
	IsConnected() bool
}