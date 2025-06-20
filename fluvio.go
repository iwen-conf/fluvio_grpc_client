// Package fluvio provides a modern Go SDK for interacting with Fluvio streaming platform
// Based on Clean Architecture principles for better maintainability and testability
package fluvio

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/config"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/grpc"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// Client 是Fluvio SDK的主要客户端
type Client struct {
	config     *config.Config
	grpcClient grpc.Client
	appService *services.FluvioApplicationService
	logger     logging.Logger
	connected  bool
}

// ClientOption 客户端配置选项函数
type ClientOption func(*config.Config) error

// NewClient 创建一个新的Fluvio客户端
func NewClient(opts ...ClientOption) (*Client, error) {
	// 创建默认配置
	cfg := config.NewDefaultConfig()

	// 应用配置选项
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, errors.Wrap(errors.ErrInvalidArgument, "failed to apply client option", err)
		}
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrInvalidArgument, "invalid configuration", err)
	}

	// 创建日志器
	logger := logging.NewDefaultLogger()
	if cfg.Logging.Level != "" {
		if level, err := logging.ParseLevel(cfg.Logging.Level); err == nil {
			logger.SetLevel(level)
		}
	}

	// 创建连接管理器
	connManager := grpc.NewConnectionManager(cfg.Connection, logger)

	// 创建真实的gRPC客户端
	grpcClient := grpc.NewDefaultClient(connManager)

	// 创建仓储
	messageRepo := repositories.NewGRPCMessageRepository(grpcClient, logger)
	topicRepo := repositories.NewGRPCTopicRepository(grpcClient, logger)
	adminRepo := repositories.NewGRPCAdminRepository(grpcClient, logger)

	// 创建应用服务
	appService := services.NewFluvioApplicationService(messageRepo, topicRepo, adminRepo, logger)

	return &Client{
		config:     cfg,
		grpcClient: grpcClient,
		appService: appService,
		logger:     logger,
		connected:  false,
	}, nil
}

// Connect 连接到Fluvio服务器
func (c *Client) Connect(ctx context.Context) error {
	if c.connected {
		return nil
	}

	c.logger.Info("Connecting to Fluvio server",
		logging.Field{Key: "host", Value: c.config.Connection.Host},
		logging.Field{Key: "port", Value: c.config.Connection.Port})

	if err := c.grpcClient.Connect(); err != nil {
		return errors.Wrap(errors.ErrConnection, "failed to connect to server", err)
	}

	c.connected = true
	c.logger.Info("Successfully connected to Fluvio server")
	return nil
}

// Close 关闭客户端连接
func (c *Client) Close() error {
	if !c.connected {
		return nil
	}

	c.logger.Info("Closing connection to Fluvio server")

	if err := c.grpcClient.Close(); err != nil {
		c.logger.Error("Error closing gRPC client", logging.Field{Key: "error", Value: err})
		return err
	}

	c.connected = false
	c.logger.Info("Connection closed successfully")
	return nil
}

// Ping 测试与服务器的连接
func (c *Client) Ping(ctx context.Context) (time.Duration, error) {
	if !c.connected {
		return 0, errors.New(errors.ErrConnection, "client not connected")
	}

	start := time.Now()

	// 执行健康检查作为ping
	if err := c.HealthCheck(ctx); err != nil {
		return 0, err
	}

	duration := time.Since(start)
	c.logger.Debug("Ping successful", logging.Field{Key: "duration", Value: duration})
	return duration, nil
}

// HealthCheck 执行健康检查
func (c *Client) HealthCheck(ctx context.Context) error {
	if !c.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}

	// 调用真实的健康检查gRPC方法
	req := &pb.HealthCheckRequest{}
	resp, err := c.grpcClient.HealthCheck(ctx, req)
	if err != nil {
		c.logger.Error("Health check failed", logging.Field{Key: "error", Value: err})
		return errors.Wrap(errors.ErrConnection, "health check failed", err)
	}

	if !resp.GetOk() {
		c.logger.Warn("Health check returned not ok", logging.Field{Key: "message", Value: resp.GetMessage()})
		return errors.New(errors.ErrConnection, "server health check failed: "+resp.GetMessage())
	}

	c.logger.Debug("Health check successful")
	return nil
}

// Producer 获取生产者实例
func (c *Client) Producer() *Producer {
	return &Producer{
		appService: c.appService,
		logger:     c.logger,
		connected:  &c.connected,
	}
}

// Consumer 获取消费者实例
func (c *Client) Consumer() *Consumer {
	return &Consumer{
		appService: c.appService,
		logger:     c.logger,
		connected:  &c.connected,
	}
}

// Topics 获取主题管理器实例
func (c *Client) Topics() *TopicManager {
	return &TopicManager{
		appService: c.appService,
		logger:     c.logger,
		connected:  &c.connected,
	}
}

// Admin 获取管理器实例
func (c *Client) Admin() *AdminManager {
	return &AdminManager{
		appService: c.appService,
		logger:     c.logger,
		connected:  &c.connected,
	}
}

// IsConnected 检查是否已连接
func (c *Client) IsConnected() bool {
	return c.connected
}

// Config 获取客户端配置
func (c *Client) Config() *config.Config {
	return c.config
}

// Logger 获取日志器
func (c *Client) Logger() logging.Logger {
	return c.logger
}

// LogLevel 日志级别类型
type LogLevel string

// 日志级别常量
const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// Version 返回SDK版本
func Version() string {
	return "2.0.0"
}
