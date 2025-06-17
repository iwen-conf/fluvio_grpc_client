package pool

import (
	"context"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/config"
	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	grpcMgr "github.com/iwen-conf/fluvio_grpc_client/internal/grpc"

	"google.golang.org/grpc"
)

// ConnectionPool 连接池
type ConnectionPool struct {
	config    *config.Config
	logger    logger.Logger
	connMgr   *grpcMgr.ConnectionManager
	pool      chan *grpc.ClientConn
	mu        sync.RWMutex
	closed    bool
	activeConns int
}

// NewConnectionPool 创建连接池
func NewConnectionPool(cfg *config.Config, log logger.Logger) *ConnectionPool {
	connMgr := grpcMgr.NewConnectionManager(cfg, log)
	
	return &ConnectionPool{
		config:  cfg,
		logger:  log,
		connMgr: connMgr,
		pool:    make(chan *grpc.ClientConn, cfg.Connection.PoolSize),
	}
}

// Get 从连接池获取连接
func (p *ConnectionPool) Get(ctx context.Context) (*grpc.ClientConn, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, errors.New(errors.ErrConnection, "连接池已关闭")
	}
	p.mu.RUnlock()

	// 尝试从池中获取连接
	select {
	case conn := <-p.pool:
		if p.isConnectionValid(conn) {
			return conn, nil
		}
		// 连接无效，关闭它
		conn.Close()
		p.decrementActiveConns()
	default:
		// 池中没有可用连接
	}

	// 创建新连接
	return p.createNewConnection(ctx)
}

// Put 将连接放回连接池
func (p *ConnectionPool) Put(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}

	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		conn.Close()
		p.decrementActiveConns()
		return
	}
	p.mu.RUnlock()

	if !p.isConnectionValid(conn) {
		conn.Close()
		p.decrementActiveConns()
		return
	}

	// 尝试放回池中
	select {
	case p.pool <- conn:
		// 成功放回池中
	default:
		// 池已满，关闭连接
		conn.Close()
		p.decrementActiveConns()
	}
}

// createNewConnection 创建新连接
func (p *ConnectionPool) createNewConnection(ctx context.Context) (*grpc.ClientConn, error) {
	p.mu.Lock()
	if p.activeConns >= p.config.Connection.PoolSize {
		p.mu.Unlock()
		return nil, errors.New(errors.ErrResourceExhausted, "连接池已满")
	}
	p.activeConns++
	p.mu.Unlock()

	conn, err := p.connMgr.GetConnection(ctx)
	if err != nil {
		p.decrementActiveConns()
		return nil, err
	}

	return conn, nil
}

// isConnectionValid 检查连接是否有效
func (p *ConnectionPool) isConnectionValid(conn *grpc.ClientConn) bool {
	if conn == nil {
		return false
	}

	state := conn.GetState()
	return state.String() == "READY" || state.String() == "IDLE"
}

// decrementActiveConns 减少活跃连接数
func (p *ConnectionPool) decrementActiveConns() {
	p.mu.Lock()
	if p.activeConns > 0 {
		p.activeConns--
	}
	p.mu.Unlock()
}

// Close 关闭连接池
func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return nil
	}
	p.closed = true
	p.mu.Unlock()

	// 关闭池中的所有连接
	close(p.pool)
	for conn := range p.pool {
		conn.Close()
	}

	// 关闭连接管理器
	return p.connMgr.Close()
}

// Stats 连接池统计信息
type Stats struct {
	PoolSize    int `json:"pool_size"`
	ActiveConns int `json:"active_conns"`
	IdleConns   int `json:"idle_conns"`
}

// GetStats 获取连接池统计信息
func (p *ConnectionPool) GetStats() Stats {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return Stats{
		PoolSize:    p.config.Connection.PoolSize,
		ActiveConns: p.activeConns,
		IdleConns:   len(p.pool),
	}
}

// PooledConnection 池化连接包装器
type PooledConnection struct {
	conn *grpc.ClientConn
	pool *ConnectionPool
}

// NewPooledConnection 创建池化连接
func NewPooledConnection(conn *grpc.ClientConn, pool *ConnectionPool) *PooledConnection {
	return &PooledConnection{
		conn: conn,
		pool: pool,
	}
}

// GetConn 获取底层连接
func (pc *PooledConnection) GetConn() *grpc.ClientConn {
	return pc.conn
}

// Close 关闭连接（实际上是放回池中）
func (pc *PooledConnection) Close() error {
	if pc.conn != nil {
		pc.pool.Put(pc.conn)
		pc.conn = nil
	}
	return nil
}

// ConnectionFactory 连接工厂接口
type ConnectionFactory interface {
	GetConnection(ctx context.Context) (*PooledConnection, error)
	Close() error
	GetStats() Stats
}

// Factory 连接工厂实现
type Factory struct {
	pool *ConnectionPool
}

// NewFactory 创建连接工厂
func NewFactory(cfg *config.Config, log logger.Logger) *Factory {
	return &Factory{
		pool: NewConnectionPool(cfg, log),
	}
}

// GetConnection 获取连接
func (f *Factory) GetConnection(ctx context.Context) (*PooledConnection, error) {
	conn, err := f.pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	return NewPooledConnection(conn, f.pool), nil
}

// Close 关闭工厂
func (f *Factory) Close() error {
	return f.pool.Close()
}

// GetStats 获取统计信息
func (f *Factory) GetStats() Stats {
	return f.pool.GetStats()
}
