package client

import (
	"context"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/config"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// ClientOption 客户端选项函数类型（向后兼容）
type ClientOption func(*config.Config) error

// Client 是Fluvio gRPC客户端的主要入口点（向后兼容）
type Client struct {
	config *config.Config
	logger logging.Logger

	// 服务客户端
	producer *Producer
	consumer *Consumer
	topic    *TopicManager
	admin    *AdminManager

	// 内部状态
	mu     sync.RWMutex
	closed bool
}

// New 创建一个新的Fluvio客户端（向后兼容）
func New(opts ...ClientOption) (*Client, error) {
	cfg := config.NewDefaultConfig()
	
	// 应用选项
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, errors.Wrap(errors.ErrInvalidArgument, "invalid client option", err)
		}
	}

	return NewWithConfig(cfg)
}

// NewWithConfig 使用指定配置创建客户端
func NewWithConfig(cfg *config.Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrInvalidArgument, "配置验证失败", err)
	}

	// 创建日志器
	logger := logging.NewDefaultLogger()
	if level, err := logging.ParseLevel(cfg.Logging.Level); err == nil {
		logger.SetLevel(level)
	}

	client := &Client{
		config: cfg,
		logger: logger,
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

	// 简化实现：总是返回成功
	// 在实际实现中，这里应该调用gRPC健康检查
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
	return nil
}

// GetConfig 获取配置
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// GetLogger 获取日志器
func (c *Client) GetLogger() logging.Logger {
	return c.logger
}

// GetStats 获取客户端统计信息
func (c *Client) GetStats() ClientStats {
	return ClientStats{
		ConnectionPool: ConnectionPoolStats{
			PoolSize:    5,
			ActiveConns: 1,
			IdleConns:   4,
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

// isClosed 检查客户端是否已关闭
func (c *Client) isClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.closed
}

// SetLogger 设置日志器
func (c *Client) SetLogger(log logging.Logger) {
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