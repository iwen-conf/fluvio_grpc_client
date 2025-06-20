package fluvio

import (
	"fmt"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/config"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
)

// 配置选项函数

// WithAddress 设置服务器地址
func WithAddress(host string, port int) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.Host = host
		cfg.Connection.Port = port
		return nil
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithTimeout(timeout, timeout)
		return nil
	}
}

// WithTimeouts 设置连接和调用超时时间
func WithTimeouts(connectTimeout, callTimeout time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithTimeout(connectTimeout, callTimeout)
		return nil
	}
}

// WithRetry 设置重试配置
func WithRetry(maxRetries int, backoff time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithRetry(maxRetries, backoff)
		return nil
	}
}

// WithTLS 设置TLS配置
func WithTLS(certFile, keyFile, caFile string) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithTLS(certFile, keyFile, caFile)
		return nil
	}
}

// WithLogger 设置自定义日志器
func WithLogger(logger logging.Logger) ClientOption {
	return func(cfg *config.Config) error {
		if logger == nil {
			return fmt.Errorf("logger cannot be nil")
		}

		// 设置日志器配置
		cfg.Logging.Level = logger.GetLevel().String()

		// 将日志器实例存储在配置的扩展字段中
		// 这样客户端创建时可以直接使用这个日志器实例
		if cfg.Extensions == nil {
			cfg.Extensions = make(map[string]interface{})
		}
		cfg.Extensions["custom_logger"] = logger

		return nil
	}
}

// WithLogLevel 设置日志级别
func WithLogLevel(level LogLevel) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Logging.Level = string(level)
		return nil
	}
}

// WithConnectionPool 设置连接池配置
func WithConnectionPool(size int, maxIdle time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.WithPool(size, maxIdle)
		return nil
	}
}

// WithInsecure 设置不安全连接（跳过TLS验证）
func WithInsecure() ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.Insecure = true
		return nil
	}
}

// WithKeepAlive 设置保活配置
func WithKeepAlive(interval time.Duration) ClientOption {
	return func(cfg *config.Config) error {
		cfg.Connection.KeepAliveTime = interval
		return nil
	}
}
