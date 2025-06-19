package grpc

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// Client gRPC客户端接口（简化版本）
type Client interface {
	// 基本消息操作
	Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error)
	BatchProduce(ctx context.Context, req *pb.BatchProduceRequest) (*pb.BatchProduceReply, error)
	Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error)
	StreamConsume(ctx context.Context, req *pb.StreamConsumeRequest) (pb.FluvioService_StreamConsumeClient, error)

	// 基本主题操作
	CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error)
	DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error)
	ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error)
	DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error)

	// 基本管理操作
	ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error)
	DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error)
	CommitOffset(ctx context.Context, req *pb.CommitOffsetRequest) (*pb.CommitOffsetReply, error)

	// SmartModule基本操作
	ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error)
	CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error)
	DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error)
	DescribeSmartModule(ctx context.Context, req *pb.DescribeSmartModuleRequest) (*pb.DescribeSmartModuleReply, error)

	// 健康检查
	HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error)

	// 连接管理
	Connect() error
	Close() error
	IsConnected() bool
}

// DefaultClient 真实的gRPC客户端实现
type DefaultClient struct {
	connManager *ConnectionManager
	client      pb.FluvioServiceClient
	connected   bool
	mu          sync.RWMutex
}

// NewDefaultClient 创建新的gRPC客户端
func NewDefaultClient(connManager *ConnectionManager) *DefaultClient {
	return &DefaultClient{
		connManager: connManager,
		connected:   false,
	}
}

// Connect 连接到gRPC服务器
func (c *DefaultClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := c.connManager.GetConnection(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}

	c.client = pb.NewFluvioServiceClient(conn)
	c.connected = true
	return nil
}

// Close 关闭连接
func (c *DefaultClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	c.connected = false
	c.client = nil
	return c.connManager.Close()
}

// IsConnected 检查连接状态
func (c *DefaultClient) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// ensureConnected 确保客户端已连接
func (c *DefaultClient) ensureConnected() error {
	c.mu.RLock()
	connected := c.connected
	c.mu.RUnlock()

	if !connected {
		return fmt.Errorf("client not connected")
	}
	return nil
}

// 真实的gRPC方法实现

func (c *DefaultClient) Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.Produce(ctx, req)
}

func (c *DefaultClient) BatchProduce(ctx context.Context, req *pb.BatchProduceRequest) (*pb.BatchProduceReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.BatchProduce(ctx, req)
}

func (c *DefaultClient) Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.Consume(ctx, req)
}

func (c *DefaultClient) CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.CreateTopic(ctx, req)
}

func (c *DefaultClient) DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.DeleteTopic(ctx, req)
}

func (c *DefaultClient) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.ListTopics(ctx, req)
}

func (c *DefaultClient) DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.DescribeTopic(ctx, req)
}

func (c *DefaultClient) ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.ListConsumerGroups(ctx, req)
}

func (c *DefaultClient) DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.DescribeConsumerGroup(ctx, req)
}

func (c *DefaultClient) CommitOffset(ctx context.Context, req *pb.CommitOffsetRequest) (*pb.CommitOffsetReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.CommitOffset(ctx, req)
}

func (c *DefaultClient) ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.ListSmartModules(ctx, req)
}

func (c *DefaultClient) CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.CreateSmartModule(ctx, req)
}

func (c *DefaultClient) DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.DeleteSmartModule(ctx, req)
}

func (c *DefaultClient) DescribeSmartModule(ctx context.Context, req *pb.DescribeSmartModuleRequest) (*pb.DescribeSmartModuleReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.DescribeSmartModule(ctx, req)
}

func (c *DefaultClient) StreamConsume(ctx context.Context, req *pb.StreamConsumeRequest) (pb.FluvioService_StreamConsumeClient, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.StreamConsume(ctx, req)
}

func (c *DefaultClient) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error) {
	if err := c.ensureConnected(); err != nil {
		return nil, err
	}
	return c.client.HealthCheck(ctx, req)
}
