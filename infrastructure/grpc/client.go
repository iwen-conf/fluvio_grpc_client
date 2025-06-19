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

// DefaultClient 默认gRPC客户端实现（简化版）
type DefaultClient struct {
	connected bool
}

// Connect 连接
func (c *DefaultClient) Connect() error {
	c.connected = true
	return nil
}

// Close 关闭连接
func (c *DefaultClient) Close() error {
	c.connected = false
	return nil
}

// IsConnected 检查连接状态
func (c *DefaultClient) IsConnected() bool {
	return c.connected
}

// 简化实现的方法（返回模拟数据）

func (c *DefaultClient) Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error) {
	return &pb.ProduceReply{Success: true, MessageId: "mock-msg-id"}, nil
}

func (c *DefaultClient) Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error) {
	return &pb.ConsumeReply{}, nil
}

func (c *DefaultClient) CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error) {
	return &pb.CreateTopicReply{Success: true}, nil
}

func (c *DefaultClient) DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error) {
	return &pb.DeleteTopicReply{Success: true}, nil
}

func (c *DefaultClient) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error) {
	return &pb.ListTopicsReply{Topics: []string{"example-topic", "test-topic"}}, nil
}

func (c *DefaultClient) DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error) {
	return &pb.DescribeTopicReply{}, nil
}

func (c *DefaultClient) ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error) {
	return &pb.ListConsumerGroupsReply{}, nil
}

func (c *DefaultClient) DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error) {
	return &pb.DescribeConsumerGroupReply{}, nil
}

func (c *DefaultClient) ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error) {
	return &pb.ListSmartModulesReply{}, nil
}

func (c *DefaultClient) CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error) {
	return &pb.CreateSmartModuleReply{Success: true}, nil
}

func (c *DefaultClient) DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error) {
	return &pb.DeleteSmartModuleReply{Success: true}, nil
}

func (c *DefaultClient) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error) {
	return &pb.HealthCheckReply{
		Ok:      true,
		Message: "Service is healthy",
		Status:  pb.HealthStatus_HEALTHY,
	}, nil
}
