package valueobjects

import (
	"fmt"
	"time"
)

// ConnectionConfig 连接配置值对象
type ConnectionConfig struct {
	// 服务器配置
	Host string
	Port int

	// 超时配置
	ConnectTimeout time.Duration
	RequestTimeout time.Duration

	// 重试配置
	MaxRetries    int
	RetryInterval time.Duration

	// 连接池配置
	PoolSize         int
	MaxIdleTime      time.Duration
	KeepAliveTime    time.Duration
	KeepAliveTimeout time.Duration

	// TLS配置
	TLSEnabled bool
	CertFile   string
	KeyFile    string
	CAFile     string
	Insecure   bool // 跳过TLS验证
}

// NewConnectionConfig 创建默认连接配置
func NewConnectionConfig(host string, port int) *ConnectionConfig {
	return &ConnectionConfig{
		Host:             host,
		Port:             port,
		ConnectTimeout:   5 * time.Second,
		RequestTimeout:   30 * time.Second,
		MaxRetries:       3,
		RetryInterval:    1 * time.Second,
		PoolSize:         5,
		MaxIdleTime:      5 * time.Minute,
		KeepAliveTime:    30 * time.Second,
		KeepAliveTimeout: 5 * time.Second,
		TLSEnabled:       false,
	}
}

// WithTimeout 设置超时时间
func (cc *ConnectionConfig) WithTimeout(connect, request time.Duration) *ConnectionConfig {
	cc.ConnectTimeout = connect
	cc.RequestTimeout = request
	return cc
}

// WithRetry 设置重试配置
func (cc *ConnectionConfig) WithRetry(maxRetries int, interval time.Duration) *ConnectionConfig {
	cc.MaxRetries = maxRetries
	cc.RetryInterval = interval
	return cc
}

// WithPool 设置连接池配置
func (cc *ConnectionConfig) WithPool(size int, maxIdleTime time.Duration) *ConnectionConfig {
	cc.PoolSize = size
	cc.MaxIdleTime = maxIdleTime
	return cc
}

// WithTLS 启用TLS配置
func (cc *ConnectionConfig) WithTLS(certFile, keyFile, caFile string) *ConnectionConfig {
	cc.TLSEnabled = true
	cc.CertFile = certFile
	cc.KeyFile = keyFile
	cc.CAFile = caFile
	return cc
}

// IsValid 验证配置是否有效
func (cc *ConnectionConfig) IsValid() bool {
	if cc.Host == "" || cc.Port <= 0 || cc.Port > 65535 {
		return false
	}

	if cc.ConnectTimeout <= 0 || cc.RequestTimeout <= 0 {
		return false
	}

	if cc.MaxRetries < 0 || cc.PoolSize <= 0 {
		return false
	}

	// 如果启用TLS，检查证书文件
	if cc.TLSEnabled && (cc.CertFile == "" || cc.KeyFile == "") {
		return false
	}

	return true
}

// Address 返回完整的服务器地址
func (cc *ConnectionConfig) Address() string {
	return fmt.Sprintf("%s:%d", cc.Host, cc.Port)
}
