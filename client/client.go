package client

import (
	"context"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/config"
	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	"github.com/iwen-conf/fluvio_grpc_client/internal/pool"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"

	"google.golang.org/grpc"
)

// Client 是Fluvio gRPC客户端的主要入口点
type Client struct {
	config     *config.Config
	logger     logger.Logger
	connFactory pool.ConnectionFactory
	
	// 服务客户端
	producer *Producer
	consumer *Consumer
	topic    *TopicManager
	admin    *AdminManager
	
	// 内部状态
	mu     sync.RWMutex
	closed bool
}

// New 创建一个新的Fluvio客户端
func New(opts ...config.Option) (*Client, error) {
	cfg := config.DefaultConfig()
	config.ApplyOptions(cfg, opts...)
	
	return NewWithConfig(cfg)
}

// NewWithConfig 使用指定配置创建客户端
func NewWithConfig(cfg *config.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "配置验证失败", err)
	}

	// 创建日志器
	log := logger.NewDefaultLogger(cfg.Logging.Level)

	// 创建连接工厂
	connFactory := pool.NewFactory(cfg, log)

	client := &Client{
		config:      cfg,
		logger:      log,
		connFactory: connFactory,
	}

	// 初始化服务客户端
	client.producer = NewProducer(client)
	client.consumer = NewConsumer(client)
	client.topic = NewTopicManager(client)
	client.admin = NewAdminManager(client)

	return client, nil
}

// Producer 返回消息生产者
func (c *Client) Producer() *Producer {
	return c.producer
}

// Consumer 返回消息消费者
func (c *Client) Consumer() *Consumer {
	return c.consumer
}

// Topic 返回主题管理器
func (c *Client) Topic() *TopicManager {
	return c.topic
}

// Admin 返回管理功能
func (c *Client) Admin() *AdminManager {
	return c.admin
}

// HealthCheck 执行健康检查
func (c *Client) HealthCheck(ctx context.Context) error {
	if c.isClosed() {
		return errors.New(errors.ErrConnection, "客户端已关闭")
	}

	conn, err := c.getConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewFluvioServiceClient(conn.GetConn())
	
	resp, err := client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		return errors.Wrap(errors.ErrServiceUnavailable, "健康检查失败", err)
	}

	if !resp.GetOk() {
		return errors.New(errors.ErrServiceUnavailable, resp.GetMessage())
	}

	return nil
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	
	if c.connFactory != nil {
		return c.connFactory.Close()
	}

	return nil
}

// GetConfig 获取配置
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// GetLogger 获取日志器
func (c *Client) GetLogger() logger.Logger {
	return c.logger
}

// GetStats 获取客户端统计信息
func (c *Client) GetStats() ClientStats {
	stats := c.connFactory.GetStats()
	
	return ClientStats{
		ConnectionPool: ConnectionPoolStats{
			PoolSize:    stats.PoolSize,
			ActiveConns: stats.ActiveConns,
			IdleConns:   stats.IdleConns,
		},
		Closed: c.isClosed(),
	}
}

// ClientStats 客户端统计信息
type ClientStats struct {
	ConnectionPool ConnectionPoolStats `json:"connection_pool"`
	Closed         bool                `json:"closed"`
}

// ConnectionPoolStats 连接池统计信息
type ConnectionPoolStats struct {
	PoolSize    int `json:"pool_size"`
	ActiveConns int `json:"active_conns"`
	IdleConns   int `json:"idle_conns"`
}

// getConnection 获取连接
func (c *Client) getConnection(ctx context.Context) (*pool.PooledConnection, error) {
	if c.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	return c.connFactory.GetConnection(ctx)
}

// isClosed 检查客户端是否已关闭
func (c *Client) isClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// withConnection 使用连接执行操作
func (c *Client) withConnection(ctx context.Context, fn func(*grpc.ClientConn) error) error {
	conn, err := c.getConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	return fn(conn.GetConn())
}

// withRetry 带重试执行操作
func (c *Client) withRetry(ctx context.Context, fn func(context.Context) error) error {
	// TODO: 实现重试逻辑
	return fn(ctx)
}

// SetLogger 设置日志器
func (c *Client) SetLogger(log logger.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = log
}

// Ping 测试连接
func (c *Client) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	err := c.HealthCheck(ctx)
	duration := time.Since(start)
	
	if err != nil {
		return 0, err
	}
	
	return duration, nil
}
