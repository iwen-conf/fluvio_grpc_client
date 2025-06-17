// Package fluvio provides a Go SDK for interacting with Fluvio streaming platform
package fluvio

import (
	"context"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/client"
	"github.com/iwen-conf/fluvio_grpc_client/config"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	"github.com/iwen-conf/fluvio_grpc_client/types"
)

// Client 是Fluvio SDK的主要客户端接口
type Client = client.Client

// 重新导出主要类型
type (
	// 消息相关类型
	Message              = types.Message
	ProduceOptions       = types.ProduceOptions
	ProduceResult        = types.ProduceResult
	BatchProduceResult   = types.BatchProduceResult
	ConsumeOptions       = types.ConsumeOptions
	ConsumeResult        = types.ConsumeResult
	StreamConsumeOptions = types.StreamConsumeOptions
	CommitOffsetOptions  = types.CommitOffsetOptions

	// 主题相关类型
	TopicInfo             = types.TopicInfo
	CreateTopicOptions    = types.CreateTopicOptions
	CreateTopicResult     = types.CreateTopicResult
	DeleteTopicOptions    = types.DeleteTopicOptions
	DeleteTopicResult     = types.DeleteTopicResult
	ListTopicsResult      = types.ListTopicsResult
	DescribeTopicResult   = types.DescribeTopicResult

	// 消费者组相关类型
	ConsumerGroupInfo           = types.ConsumerGroupInfo
	ConsumerGroupMember         = types.ConsumerGroupMember
	ListConsumerGroupsResult    = types.ListConsumerGroupsResult
	DescribeConsumerGroupResult = types.DescribeConsumerGroupResult

	// 管理相关类型
	ClusterInfo             = types.ClusterInfo
	BrokerInfo              = types.BrokerInfo
	MetricInfo              = types.MetricInfo
	DescribeClusterResult   = types.DescribeClusterResult
	ListBrokersResult       = types.ListBrokersResult
	GetMetricsOptions       = types.GetMetricsOptions
	GetMetricsResult        = types.GetMetricsResult

	// SmartModule相关类型
	SmartModuleInfo             = types.SmartModuleInfo
	CreateSmartModuleOptions    = types.CreateSmartModuleOptions
	CreateSmartModuleResult     = types.CreateSmartModuleResult
	DeleteSmartModuleResult     = types.DeleteSmartModuleResult
	ListSmartModulesResult      = types.ListSmartModulesResult
	DescribeSmartModuleResult   = types.DescribeSmartModuleResult

	// 配置相关类型
	Config           = config.Config
	ServerConfig     = config.ServerConfig
	ConnectionConfig = config.ConnectionConfig
	LoggingConfig    = config.LoggingConfig
	RetryConfig      = config.RetryConfig
	Option           = config.Option

	// 日志相关类型
	Logger = logger.Logger
	Level  = logger.Level
	Field  = logger.Field
)

// 重新导出配置选项函数
var (
	WithServer         = config.WithServer
	WithTimeout        = config.WithTimeout
	WithConnectTimeout = config.WithConnectTimeout
	WithCallTimeout    = config.WithCallTimeout
	WithLogger         = config.WithLogger
	WithLogLevel       = config.WithLogLevel
	WithRetry          = config.WithRetry
	WithMaxRetries     = config.WithMaxRetries
	WithConnectionPool = config.WithConnectionPool
	WithPoolSize       = config.WithPoolSize
	WithTLS            = config.WithTLS
	WithInsecureTLS    = config.WithInsecureTLS
	WithKeepAlive      = config.WithKeepAlive
	WithBackoff        = config.WithBackoff
)

// 重新导出日志级别常量
const (
	LevelDebug = logger.LevelDebug
	LevelInfo  = logger.LevelInfo
	LevelWarn  = logger.LevelWarn
	LevelError = logger.LevelError
)

// New 创建一个新的Fluvio客户端
// 这是创建客户端的推荐方式
func New(opts ...Option) (*Client, error) {
	return client.New(opts...)
}

// NewWithConfig 使用指定配置创建客户端
func NewWithConfig(cfg *Config) (*Client, error) {
	return client.NewWithConfig(cfg)
}

// Connect 连接到Fluvio服务器
// 这是一个便捷函数，等同于New()
func Connect(opts ...Option) (*Client, error) {
	return New(opts...)
}

// ConnectWithAddress 连接到指定地址的Fluvio服务器
func ConnectWithAddress(host string, port int, opts ...Option) (*Client, error) {
	allOpts := append([]Option{WithServer(host, port)}, opts...)
	return New(allOpts...)
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return config.DefaultConfig()
}

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(path string) (*Config, error) {
	return config.LoadFromFile(path)
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *Config {
	return config.LoadFromEnv()
}

// NewLogger 创建新的日志器
func NewLogger(level Level) Logger {
	return logger.NewDefaultLogger(level)
}

// NewNoopLogger 创建空日志器
func NewNoopLogger() Logger {
	return logger.NewNoopLogger()
}

// SetDefaultLogger 设置默认日志器
func SetDefaultLogger(log Logger) {
	logger.SetDefault(log)
}

// GetDefaultLogger 获取默认日志器
func GetDefaultLogger() Logger {
	return logger.GetDefault()
}

// QuickStart 快速开始示例
// 这个函数展示了如何快速连接和使用Fluvio
func QuickStart(host string, port int) (*Client, error) {
	return ConnectWithAddress(host, port,
		WithTimeout(5*time.Second, 10*time.Second),
		WithLogLevel(LevelInfo),
		WithMaxRetries(3),
	)
}

// SimpleProducer 创建一个简单的生产者客户端
func SimpleProducer(host string, port int) (*Client, error) {
	return ConnectWithAddress(host, port,
		WithTimeout(5*time.Second, 30*time.Second),
		WithLogLevel(LevelWarn),
		WithMaxRetries(5),
		WithPoolSize(1),
	)
}

// SimpleConsumer 创建一个简单的消费者客户端
func SimpleConsumer(host string, port int) (*Client, error) {
	return ConnectWithAddress(host, port,
		WithTimeout(10*time.Second, 60*time.Second),
		WithLogLevel(LevelWarn),
		WithMaxRetries(3),
		WithPoolSize(2),
	)
}

// HighThroughputClient 创建一个高吞吐量客户端
func HighThroughputClient(host string, port int) (*Client, error) {
	return ConnectWithAddress(host, port,
		WithTimeout(5*time.Second, 30*time.Second),
		WithLogLevel(LevelError),
		WithMaxRetries(5),
		WithPoolSize(10),
		WithKeepAlive(30*time.Second),
	)
}

// TestClient 创建一个用于测试的客户端
func TestClient(host string, port int) (*Client, error) {
	return ConnectWithAddress(host, port,
		WithTimeout(2*time.Second, 5*time.Second),
		WithLogLevel(LevelDebug),
		WithMaxRetries(1),
		WithPoolSize(1),
	)
}

// Ping 测试与Fluvio服务器的连接
func Ping(ctx context.Context, host string, port int) (time.Duration, error) {
	client, err := ConnectWithAddress(host, port,
		WithTimeout(5*time.Second, 10*time.Second),
		WithLogLevel(LevelError),
		WithMaxRetries(1),
	)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	return client.Ping(ctx)
}

// Version 返回SDK版本信息
func Version() string {
	return "1.0.0"
}

// UserAgent 返回用户代理字符串
func UserAgent() string {
	return "fluvio-go-sdk/" + Version()
}
