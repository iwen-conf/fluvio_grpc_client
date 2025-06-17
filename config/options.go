package config

import (
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/logger"
)

// Option 配置选项函数
type Option func(*Config)

// WithServer 设置服务器地址
func WithServer(host string, port int) Option {
	return func(c *Config) {
		c.Server.Host = host
		c.Server.Port = port
	}
}

// WithTimeout 设置超时时间
func WithTimeout(connect, call time.Duration) Option {
	return func(c *Config) {
		c.Connection.ConnectTimeout = connect
		c.Connection.CallTimeout = call
	}
}

// WithConnectTimeout 设置连接超时时间
func WithConnectTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Connection.ConnectTimeout = timeout
	}
}

// WithCallTimeout 设置调用超时时间
func WithCallTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.Connection.CallTimeout = timeout
	}
}

// WithLogger 设置日志器配置
func WithLogger(level logger.Level, format, output string) Option {
	return func(c *Config) {
		c.Logging.Level = level
		c.Logging.Format = format
		c.Logging.Output = output
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(level logger.Level) Option {
	return func(c *Config) {
		c.Logging.Level = level
	}
}

// WithRetry 设置重试配置
func WithRetry(maxRetries int, initialBackoff, maxBackoff time.Duration, multiple float64) Option {
	return func(c *Config) {
		c.Retry.MaxRetries = maxRetries
		c.Retry.InitialBackoff = initialBackoff
		c.Retry.MaxBackoff = maxBackoff
		c.Retry.BackoffMultiple = multiple
	}
}

// WithMaxRetries 设置最大重试次数
func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.Retry.MaxRetries = maxRetries
	}
}

// WithConnectionPool 设置连接池配置
func WithConnectionPool(poolSize int, keepAlive time.Duration) Option {
	return func(c *Config) {
		c.Connection.PoolSize = poolSize
		c.Connection.KeepAlive = keepAlive
	}
}

// WithPoolSize 设置连接池大小
func WithPoolSize(poolSize int) Option {
	return func(c *Config) {
		c.Connection.PoolSize = poolSize
	}
}

// WithTLS 设置TLS配置
func WithTLS(enabled bool, certFile, keyFile, caFile string) Option {
	return func(c *Config) {
		c.Server.TLS.Enabled = enabled
		c.Server.TLS.CertFile = certFile
		c.Server.TLS.KeyFile = keyFile
		c.Server.TLS.CAFile = caFile
	}
}

// WithInsecureTLS 设置不安全的TLS（跳过证书验证）
func WithInsecureTLS() Option {
	return func(c *Config) {
		c.Server.TLS.Enabled = true
		c.Server.TLS.InsecureSkipVerify = true
	}
}

// WithKeepAlive 设置保持连接时间
func WithKeepAlive(keepAlive time.Duration) Option {
	return func(c *Config) {
		c.Connection.KeepAlive = keepAlive
	}
}

// WithBackoff 设置退避策略
func WithBackoff(initial, max time.Duration, multiple float64) Option {
	return func(c *Config) {
		c.Retry.InitialBackoff = initial
		c.Retry.MaxBackoff = max
		c.Retry.BackoffMultiple = multiple
	}
}

// ApplyOptions 应用配置选项
func ApplyOptions(cfg *Config, opts ...Option) {
	for _, opt := range opts {
		opt(cfg)
	}
}
