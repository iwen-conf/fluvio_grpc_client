package grpc

import (
	"context"

	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// Client gRPC客户端接口（简化版本）
type Client interface {
	// 基本消息操作
	Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error)
	Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error)

	// 基本主题操作
	CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error)
	DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error)
	ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error)
	DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error)

	// 基本管理操作
	ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error)
	DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error)

	// SmartModule基本操作
	ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error)
	CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error)
	DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error)

	// 健康检查
	HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error)

	// 连接管理
	Connect() error
	Close() error
	IsConnected() bool
}
