package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/retry"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ConnectionManager gRPC连接管理器
type ConnectionManager struct {
	config *valueobjects.ConnectionConfig
	logger logging.Logger
	mu     sync.RWMutex
	conns  map[string]*grpc.ClientConn
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(config *valueobjects.ConnectionConfig, logger logging.Logger) *ConnectionManager {
	return &ConnectionManager{
		config: config,
		logger: logger,
		conns:  make(map[string]*grpc.ClientConn),
	}
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(ctx context.Context) (*grpc.ClientConn, error) {
	serverAddr := cm.config.Address()

	cm.mu.RLock()
	conn, exists := cm.conns[serverAddr]
	cm.mu.RUnlock()

	if exists && cm.isConnectionHealthy(conn) {
		return conn, nil
	}

	// 需要创建新连接
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 双重检查
	if conn, exists := cm.conns[serverAddr]; exists && cm.isConnectionHealthy(conn) {
		return conn, nil
	}

	// 创建新连接
	newConn, err := cm.createConnection(ctx, serverAddr)
	if err != nil {
		return nil, err
	}

	// 关闭旧连接（如果存在）
	if conn != nil {
		conn.Close()
	}

	cm.conns[serverAddr] = newConn
	return newConn, nil
}

// createConnection 创建新连接（带重试机制）
func (cm *ConnectionManager) createConnection(ctx context.Context, serverAddr string) (*grpc.ClientConn, error) {
	cm.logger.Info("正在创建gRPC连接", logging.Field{Key: "address", Value: serverAddr})

	var conn *grpc.ClientConn

	// 使用重试机制创建连接
	retryConfig := &retry.RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}

	err := retry.Retry(ctx, retryConfig, retry.DefaultIsRetryableError, func() error {
		opts := []grpc.DialOption{
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                cm.config.KeepAliveTime,
				Timeout:             cm.config.KeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		}

		// 配置TLS
		if cm.config.TLSEnabled {
			creds, err := cm.createTLSCredentials()
			if err != nil {
				return errors.Wrap(errors.ErrConnection, "创建TLS凭据失败", err)
			}
			opts = append(opts, grpc.WithTransportCredentials(creds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		// 创建连接
		newConn, err := grpc.NewClient(serverAddr, opts...)
		if err != nil {
			return errors.Wrap(errors.ErrConnection, "创建gRPC客户端失败", err)
		}

		// 等待连接就绪
		connectCtx, cancel := context.WithTimeout(ctx, cm.config.ConnectTimeout)
		defer cancel()

		if err := cm.waitForConnection(connectCtx, newConn); err != nil {
			newConn.Close()
			return err
		}

		conn = newConn
		return nil
	}, cm.logger)

	if err != nil {
		return nil, err
	}

	cm.logger.Info("gRPC连接创建成功", logging.Field{Key: "address", Value: serverAddr})
	return conn, nil
}

// createTLSCredentials 创建TLS凭据
func (cm *ConnectionManager) createTLSCredentials() (credentials.TransportCredentials, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 默认安全
	}

	if cm.config.CertFile != "" && cm.config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cm.config.CertFile, cm.config.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("加载客户端证书失败: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return credentials.NewTLS(tlsConfig), nil
}

// waitForConnection 等待连接就绪
func (cm *ConnectionManager) waitForConnection(ctx context.Context, conn *grpc.ClientConn) error {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return errors.New(errors.ErrTimeout, "等待连接就绪超时")
		case <-ticker.C:
			state := conn.GetState()
			cm.logger.Debug("连接状态检查", logging.Field{Key: "state", Value: state.String()})
			
			if state == connectivity.Ready {
				return nil
			}
			
			if state == connectivity.TransientFailure || state == connectivity.Shutdown {
				return errors.New(errors.ErrConnection, fmt.Sprintf("连接失败，状态: %v", state))
			}
			
			// 尝试触发连接状态变化
			conn.Connect()
		}
	}
}

// isConnectionHealthy 检查连接是否健康
func (cm *ConnectionManager) isConnectionHealthy(conn *grpc.ClientConn) bool {
	if conn == nil {
		return false
	}

	state := conn.GetState()
	return state == connectivity.Ready || state == connectivity.Idle
}

// Close 关闭所有连接
func (cm *ConnectionManager) Close() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	var lastErr error
	for addr, conn := range cm.conns {
		if err := conn.Close(); err != nil {
			cm.logger.Error("关闭连接失败",
				logging.Field{Key: "address", Value: addr},
				logging.Field{Key: "error", Value: err})
			lastErr = err
		}
	}

	cm.conns = make(map[string]*grpc.ClientConn)
	return lastErr
}

// GetConnectionCount 获取连接数量
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.conns)
}

// GetConfig 获取连接配置
func (cm *ConnectionManager) GetConfig() *valueobjects.ConnectionConfig {
	return cm.config
}

// GetConnectionStates 获取所有连接状态
func (cm *ConnectionManager) GetConnectionStates() map[string]connectivity.State {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	states := make(map[string]connectivity.State)
	for addr, conn := range cm.conns {
		states[addr] = conn.GetState()
	}
	return states
}

// Ping 测试连接
func (cm *ConnectionManager) Ping(ctx context.Context) (time.Duration, error) {
	start := time.Now()

	conn, err := cm.GetConnection(ctx)
	if err != nil {
		return 0, err
	}

	// 简单的状态检查作为ping
	state := conn.GetState()
	if state != connectivity.Ready && state != connectivity.Idle {
		return 0, errors.New(errors.ErrConnection, fmt.Sprintf("连接状态异常: %v", state))
	}

	return time.Since(start), nil
}
