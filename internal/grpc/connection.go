package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/iwen-conf/fluvio_grpc_client/config"
	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ConnectionManager 连接管理器
type ConnectionManager struct {
	config *config.Config
	logger logger.Logger
	mu     sync.RWMutex
	conns  map[string]*grpc.ClientConn
}

// NewConnectionManager 创建连接管理器
func NewConnectionManager(cfg *config.Config, log logger.Logger) *ConnectionManager {
	return &ConnectionManager{
		config: cfg,
		logger: log,
		conns:  make(map[string]*grpc.ClientConn),
	}
}

// GetConnection 获取连接
func (cm *ConnectionManager) GetConnection(ctx context.Context) (*grpc.ClientConn, error) {
	serverAddr := fmt.Sprintf("%s:%d", cm.config.Server.Host, cm.config.Server.Port)

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

// createConnection 创建新连接
func (cm *ConnectionManager) createConnection(ctx context.Context, serverAddr string) (*grpc.ClientConn, error) {
	cm.logger.Info("正在创建gRPC连接", logger.Field{Key: "address", Value: serverAddr})

	opts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                cm.config.Connection.KeepAlive,
			Timeout:             cm.config.Connection.CallTimeout,
			PermitWithoutStream: true,
		}),
	}

	// 配置TLS
	if cm.config.Server.TLS.Enabled {
		creds, err := cm.createTLSCredentials()
		if err != nil {
			return nil, errors.Wrap(errors.ErrConnection, "创建TLS凭据失败", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// 创建连接
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, errors.Wrap(errors.ErrConnection, "创建gRPC客户端失败", err)
	}

	// 等待连接就绪
	connectCtx, cancel := context.WithTimeout(ctx, cm.config.Connection.ConnectTimeout)
	defer cancel()

	if err := cm.waitForConnection(connectCtx, conn); err != nil {
		conn.Close()
		return nil, err
	}

	cm.logger.Info("gRPC连接创建成功", logger.Field{Key: "address", Value: serverAddr})
	return conn, nil
}

// createTLSCredentials 创建TLS凭据
func (cm *ConnectionManager) createTLSCredentials() (credentials.TransportCredentials, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: cm.config.Server.TLS.InsecureSkipVerify,
	}

	if cm.config.Server.TLS.CertFile != "" && cm.config.Server.TLS.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cm.config.Server.TLS.CertFile, cm.config.Server.TLS.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("加载客户端证书失败: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return credentials.NewTLS(tlsConfig), nil
}

// waitForConnection 等待连接就绪
func (cm *ConnectionManager) waitForConnection(ctx context.Context, conn *grpc.ClientConn) error {
	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return nil
		}

		if state == connectivity.TransientFailure || state == connectivity.Shutdown {
			return errors.New(errors.ErrConnection, fmt.Sprintf("连接失败，状态: %v", state))
		}

		if !conn.WaitForStateChange(ctx, state) {
			return errors.New(errors.ErrTimeout, "等待连接就绪超时")
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
				logger.Field{Key: "address", Value: addr},
				logger.Field{Key: "error", Value: err})
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
