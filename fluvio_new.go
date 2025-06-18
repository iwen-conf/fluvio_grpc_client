// Package fluvio provides a Go SDK for interacting with Fluvio streaming platform
// This is the new architecture implementation with Clean Architecture principles
package fluvio

import (
	"context"
	"time"

	appservices "github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/application/usecases"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	domainservices "github.com/iwen-conf/fluvio_grpc_client/domain/services"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/config"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/interfaces/api"
	"github.com/iwen-conf/fluvio_grpc_client/interfaces/client"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// NewClient 创建新架构的Fluvio客户端
// 这是使用Clean Architecture的新实现
func NewClient(opts ...ClientOption) (api.FluvioAPI, error) {
	// 创建配置
	cfg := config.NewDefaultConfig()

	// 应用选项
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return nil, errors.Wrap(errors.ErrInvalidArgument, "invalid client option", err)
		}
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(errors.ErrInvalidArgument, "invalid configuration", err)
	}

	// 创建日志器
	logger := logging.NewDefaultLogger()
	if level, err := logging.ParseLevel(cfg.Logging.Level); err == nil {
		logger.SetLevel(level)
	}

	// 创建领域服务
	messageService := domainservices.NewMessageService()
	topicService := domainservices.NewTopicService()

	// 创建仓储（这里需要实际的gRPC客户端实现）
	// 简化实现，实际应该创建真实的gRPC连接
	// 这里应该注入真实的仓储实现
	var messageRepo repositories.MessageRepository = nil
	var topicRepo repositories.TopicRepository = nil

	// 创建用例
	produceMessageUC := usecases.NewProduceMessageUseCase(messageRepo, messageService)
	consumeMessageUC := usecases.NewConsumeMessageUseCase(messageRepo, messageService)
	manageTopicUC := usecases.NewManageTopicUseCase(topicRepo, topicService)

	// 创建应用服务
	appService := appservices.NewFluvioApplicationService(
		produceMessageUC,
		consumeMessageUC,
		manageTopicUC,
	)

	// 创建客户端适配器
	clientAdapter := client.NewFluvioClientAdapter(appService)

	return clientAdapter, nil
}

// ClientOption 客户端选项函数类型
type ClientOption func(*config.Config) error

// WithServerAddress 设置服务器地址
func WithServerAddress(host string, port int) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.Host = host
		cfg.Connection.Port = port
		return nil
	}
}

// WithTimeouts 设置超时时间
func WithTimeouts(connect, request time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithTimeout(connect, request)
		return nil
	}
}

// WithRetries 设置重试配置
func WithRetries(maxRetries int, interval time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithRetry(maxRetries, interval)
		return nil
	}
}

// WithConnectionPoolV2 设置连接池（新架构）
func WithConnectionPoolV2(size int, maxIdleTime time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithPool(size, maxIdleTime)
		return nil
	}
}

// WithLogLevelV2 设置日志级别（新架构）
func WithLogLevelV2(level string) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Logging.Level = level
		return nil
	}
}

// WithTLSV2 启用TLS（新架构）
func WithTLSV2(certFile, keyFile, caFile string) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithTLS(certFile, keyFile, caFile)
		return nil
	}
}

// WithCircuitBreaker 启用熔断器
func WithCircuitBreaker(enabled bool, failureThreshold int, recoveryTimeout time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Client.CircuitBreaker.Enabled = enabled
		cfg.Client.CircuitBreaker.FailureThreshold = failureThreshold
		cfg.Client.CircuitBreaker.RecoveryTimeout = recoveryTimeout
		return nil
	}
}

// WithMetrics 启用指标收集
func WithMetrics(enabled bool) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Client.Metrics = enabled
		return nil
	}
}

// WithTracing 启用链路追踪
func WithTracing(enabled bool) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Client.Tracing = enabled
		return nil
	}
}

// QuickConnect 快速连接到Fluvio服务器
func QuickConnect(host string, port int) (api.FluvioAPI, error) {
	return NewClient(
		WithServerAddress(host, port),
		WithTimeouts(5*time.Second, 30*time.Second),
		WithLogLevelV2("info"),
		WithRetries(3, 1*time.Second),
		WithConnectionPoolV2(5, 5*time.Minute),
	)
}

// ProductionClient 创建生产环境客户端
func ProductionClient(host string, port int) (api.FluvioAPI, error) {
	return NewClient(
		WithServerAddress(host, port),
		WithTimeouts(5*time.Second, 30*time.Second),
		WithLogLevelV2("warn"),
		WithRetries(5, 2*time.Second),
		WithConnectionPoolV2(10, 10*time.Minute),
		WithCircuitBreaker(true, 5, 30*time.Second),
		WithMetrics(true),
	)
}

// DevelopmentClient 创建开发环境客户端
func DevelopmentClient(host string, port int) (api.FluvioAPI, error) {
	return NewClient(
		WithServerAddress(host, port),
		WithTimeouts(10*time.Second, 60*time.Second),
		WithLogLevelV2("debug"),
		WithRetries(3, 1*time.Second),
		WithConnectionPoolV2(3, 5*time.Minute),
		WithMetrics(false),
		WithTracing(true),
	)
}

// TestClientV2 创建测试环境客户端（新架构）
func TestClientV2(host string, port int) (api.FluvioAPI, error) {
	return NewClient(
		WithServerAddress(host, port),
		WithTimeouts(2*time.Second, 10*time.Second),
		WithLogLevelV2("debug"),
		WithRetries(1, 500*time.Millisecond),
		WithConnectionPoolV2(1, 1*time.Minute),
	)
}

// PingServer 测试服务器连接
func PingServer(ctx context.Context, host string, port int) (time.Duration, error) {
	client, err := TestClientV2(host, port)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	return client.Ping(ctx)
}

// HealthCheckServer 检查服务器健康状态
func HealthCheckServer(ctx context.Context, host string, port int) error {
	client, err := TestClientV2(host, port)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.HealthCheck(ctx)
}

// VersionV2 返回新架构SDK版本
func VersionV2() string {
	return "2.0.0-clean-architecture"
}

// UserAgentV2 返回新架构用户代理字符串
func UserAgentV2() string {
	return "fluvio-go-sdk-v2/" + VersionV2()
}
