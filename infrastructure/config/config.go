package config

import (
	"fmt"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
)

// Config 应用配置
type Config struct {
	// 连接配置
	Connection *valueobjects.ConnectionConfig

	// 日志配置
	Logging *LoggingConfig

	// 客户端配置
	Client *ClientConfig
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"`
	Output     string `json:"output" yaml:"output"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	UserAgent      string                `json:"user_agent" yaml:"user_agent"`
	RequestID      bool                  `json:"request_id" yaml:"request_id"`
	Metrics        bool                  `json:"metrics" yaml:"metrics"`
	Tracing        bool                  `json:"tracing" yaml:"tracing"`
	CircuitBreaker *CircuitBreakerConfig `json:"circuit_breaker" yaml:"circuit_breaker"`
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	Enabled          bool          `json:"enabled" yaml:"enabled"`
	FailureThreshold int           `json:"failure_threshold" yaml:"failure_threshold"`
	RecoveryTimeout  time.Duration `json:"recovery_timeout" yaml:"recovery_timeout"`
	MonitoringPeriod time.Duration `json:"monitoring_period" yaml:"monitoring_period"`
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Connection: valueobjects.NewConnectionConfig("localhost", 50051), // 开发环境默认值
		Logging: &LoggingConfig{
			Level:      "info",
			Format:     "json",
			Output:     "stdout",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
		Client: &ClientConfig{
			UserAgent: "fluvio-go-sdk/1.0.0",
			RequestID: true,
			Metrics:   false,
			Tracing:   false,
			CircuitBreaker: &CircuitBreakerConfig{
				Enabled:          false,
				FailureThreshold: 5,
				RecoveryTimeout:  30 * time.Second,
				MonitoringPeriod: 10 * time.Second,
			},
		},
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Connection == nil {
		return fmt.Errorf("connection config is required")
	}

	if !c.Connection.IsValid() {
		return fmt.Errorf("invalid connection config")
	}

	if c.Logging == nil {
		return fmt.Errorf("logging config is required")
	}

	if c.Client == nil {
		return fmt.Errorf("client config is required")
	}

	return nil
}
