package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/logger"
)

// Config 定义SDK配置
type Config struct {
	Server     ServerConfig     `json:"server"`
	Connection ConnectionConfig `json:"connection"`
	Logging    LoggingConfig    `json:"logging"`
	Retry      RetryConfig      `json:"retry"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string    `json:"host"`
	Port int       `json:"port"`
	TLS  TLSConfig `json:"tls"`
}

// TLSConfig TLS配置
type TLSConfig struct {
	Enabled            bool   `json:"enabled"`
	CertFile           string `json:"cert_file"`
	KeyFile            string `json:"key_file"`
	CAFile             string `json:"ca_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

// ConnectionConfig 连接配置
type ConnectionConfig struct {
	ConnectTimeout time.Duration `json:"connect_timeout"`
	CallTimeout    time.Duration `json:"call_timeout"`
	MaxRetries     int           `json:"max_retries"`
	PoolSize       int           `json:"pool_size"`
	KeepAlive      time.Duration `json:"keep_alive"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  logger.Level `json:"level"`
	Format string       `json:"format"` // "text" or "json"
	Output string       `json:"output"` // "stdout", "stderr", or file path
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries      int           `json:"max_retries"`
	InitialBackoff  time.Duration `json:"initial_backoff"`
	MaxBackoff      time.Duration `json:"max_backoff"`
	BackoffMultiple float64       `json:"backoff_multiple"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 50051,
			TLS: TLSConfig{
				Enabled: false,
			},
		},
		Connection: ConnectionConfig{
			ConnectTimeout: 5 * time.Second,
			CallTimeout:    10 * time.Second,
			MaxRetries:     3,
			PoolSize:       1,
			KeepAlive:      30 * time.Second,
		},
		Logging: LoggingConfig{
			Level:  logger.LevelInfo,
			Format: "text",
			Output: "stdout",
		},
		Retry: RetryConfig{
			MaxRetries:      3,
			InitialBackoff:  100 * time.Millisecond,
			MaxBackoff:      5 * time.Second,
			BackoffMultiple: 2.0,
		},
	}
}

// LoadFromFile 从文件加载配置
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败 %s: %w", path, err)
	}

	cfg := DefaultConfig()
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败 %s: %w", path, err)
	}

	return cfg, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config {
	cfg := DefaultConfig()

	if host := os.Getenv("FLUVIO_HOST"); host != "" {
		cfg.Server.Host = host
	}

	if port := os.Getenv("FLUVIO_PORT"); port != "" {
		if p, err := parsePort(port); err == nil {
			cfg.Server.Port = p
		}
	}

	if timeout := os.Getenv("FLUVIO_CONNECT_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			cfg.Connection.ConnectTimeout = t
		}
	}

	if timeout := os.Getenv("FLUVIO_CALL_TIMEOUT"); timeout != "" {
		if t, err := time.ParseDuration(timeout); err == nil {
			cfg.Connection.CallTimeout = t
		}
	}

	if level := os.Getenv("FLUVIO_LOG_LEVEL"); level != "" {
		if l, err := parseLogLevel(level); err == nil {
			cfg.Logging.Level = l
		}
	}

	return cfg
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("服务器主机不能为空")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("服务器端口必须在1-65535之间")
	}

	if c.Connection.ConnectTimeout <= 0 {
		return fmt.Errorf("连接超时时间必须大于0")
	}

	if c.Connection.CallTimeout <= 0 {
		return fmt.Errorf("调用超时时间必须大于0")
	}

	if c.Connection.PoolSize <= 0 {
		return fmt.Errorf("连接池大小必须大于0")
	}

	if c.Retry.MaxRetries < 0 {
		return fmt.Errorf("最大重试次数不能为负数")
	}

	return nil
}

// Clone 克隆配置
func (c *Config) Clone() *Config {
	clone := *c
	return &clone
}

// parsePort 解析端口号
func parsePort(port string) (int, error) {
	var p int
	_, err := fmt.Sscanf(port, "%d", &p)
	if err != nil {
		return 0, err
	}
	if p <= 0 || p > 65535 {
		return 0, fmt.Errorf("端口号必须在1-65535之间")
	}
	return p, nil
}

// parseLogLevel 解析日志级别
func parseLogLevel(level string) (logger.Level, error) {
	switch level {
	case "debug", "DEBUG":
		return logger.LevelDebug, nil
	case "info", "INFO":
		return logger.LevelInfo, nil
	case "warn", "WARN":
		return logger.LevelWarn, nil
	case "error", "ERROR":
		return logger.LevelError, nil
	default:
		return logger.LevelInfo, fmt.Errorf("未知的日志级别: %s", level)
	}
}
