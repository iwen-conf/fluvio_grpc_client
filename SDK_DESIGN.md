# Fluvio gRPC Client SDK 设计文档

## 概述

将现有的命令行工具重构为分层设计的SDK，支持外部引入和使用。SDK采用分层架构，提供清晰的API接口和灵活的配置选项。

## 设计原则

1. **分层架构**：核心层、服务层、接口层、工具层
2. **简单易用**：提供简单的默认配置和高级配置选项
3. **可扩展性**：支持插件化和自定义扩展
4. **线程安全**：所有公共API都是线程安全的
5. **向后兼容**：保持现有命令行工具的功能不变

## 包结构设计

```
fluvio_grpc_client/
├── client/                 # 公共SDK客户端包
│   ├── client.go          # 主要SDK入口点
│   ├── producer.go        # 消息生产者
│   ├── consumer.go        # 消息消费者
│   ├── topic.go           # 主题管理
│   ├── admin.go           # 管理功能
│   └── stream.go          # 流式处理
├── config/                # 配置管理包
│   ├── config.go          # 配置结构和加载
│   ├── options.go         # 配置选项
│   └── defaults.go        # 默认配置
├── types/                 # 公共类型定义包
│   ├── message.go         # 消息类型
│   ├── topic.go           # 主题类型
│   ├── consumer.go        # 消费者类型
│   └── admin.go           # 管理类型
├── errors/                # 错误处理包
│   ├── errors.go          # 错误定义
│   └── codes.go           # 错误代码
├── logger/                # 日志接口包
│   ├── logger.go          # 日志接口
│   └── noop.go            # 空日志实现
├── examples/              # 使用示例
│   ├── basic/             # 基本使用示例
│   ├── advanced/          # 高级使用示例
│   └── integration/       # 集成示例
├── internal/              # 内部实现（保持现有结构）
│   ├── grpc/              # gRPC连接管理
│   ├── pool/              # 连接池
│   └── retry/             # 重试机制
├── cmd/                   # 命令行工具（保持不变）
├── proto/                 # 协议定义（保持不变）
└── tests/                 # 测试（更新使用新SDK）
```

## 核心API设计

### 1. 主要入口点 (client/client.go)

```go
// Client 是Fluvio gRPC客户端的主要入口点
type Client struct {
    conn       *grpc.ClientConn
    config     *config.Config
    logger     logger.Logger
    producer   *Producer
    consumer   *Consumer
    topic      *TopicManager
    admin      *AdminManager
}

// New 创建一个新的Fluvio客户端
func New(opts ...config.Option) (*Client, error)

// NewWithConfig 使用指定配置创建客户端
func NewWithConfig(cfg *config.Config) (*Client, error)

// Close 关闭客户端连接
func (c *Client) Close() error

// Producer 返回消息生产者
func (c *Client) Producer() *Producer

// Consumer 返回消息消费者
func (c *Client) Consumer() *Consumer

// Topic 返回主题管理器
func (c *Client) Topic() *TopicManager

// Admin 返回管理功能
func (c *Client) Admin() *AdminManager

// HealthCheck 执行健康检查
func (c *Client) HealthCheck(ctx context.Context) error
```

### 2. 配置管理 (config/config.go)

```go
// Config 定义SDK配置
type Config struct {
    Server      ServerConfig
    Connection  ConnectionConfig
    Logging     LoggingConfig
    Retry       RetryConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Host string
    Port int
    TLS  TLSConfig
}

// ConnectionConfig 连接配置
type ConnectionConfig struct {
    ConnectTimeout time.Duration
    CallTimeout    time.Duration
    MaxRetries     int
    PoolSize       int
}

// Option 配置选项函数
type Option func(*Config)

// WithServer 设置服务器地址
func WithServer(host string, port int) Option

// WithTimeout 设置超时时间
func WithTimeout(connect, call time.Duration) Option

// WithLogger 设置日志器
func WithLogger(logger logger.Logger) Option

// WithRetry 设置重试配置
func WithRetry(maxRetries int, backoff time.Duration) Option
```

### 3. 消息生产者 (client/producer.go)

```go
// Producer 消息生产者
type Producer struct {
    client *Client
}

// ProduceOptions 生产选项
type ProduceOptions struct {
    Topic     string
    Key       string
    Headers   map[string]string
    Timestamp time.Time
}

// Produce 发送单条消息
func (p *Producer) Produce(ctx context.Context, message string, opts ProduceOptions) (*types.ProduceResult, error)

// ProduceBatch 批量发送消息
func (p *Producer) ProduceBatch(ctx context.Context, messages []types.Message) (*types.BatchProduceResult, error)

// ProduceAsync 异步发送消息
func (p *Producer) ProduceAsync(ctx context.Context, message string, opts ProduceOptions) <-chan *types.ProduceResult
```

### 4. 消息消费者 (client/consumer.go)

```go
// Consumer 消息消费者
type Consumer struct {
    client *Client
}

// ConsumeOptions 消费选项
type ConsumeOptions struct {
    Topic       string
    Group       string
    Offset      int64
    MaxMessages int32
    AutoCommit  bool
}

// Consume 消费消息
func (c *Consumer) Consume(ctx context.Context, opts ConsumeOptions) ([]*types.Message, error)

// ConsumeStream 流式消费消息
func (c *Consumer) ConsumeStream(ctx context.Context, opts ConsumeOptions) (<-chan *types.Message, error)

// CommitOffset 提交偏移量
func (c *Consumer) CommitOffset(ctx context.Context, topic, group string, offset int64) error
```

## 类型定义 (types/)

### 消息类型 (types/message.go)

```go
// Message 表示一条消息
type Message struct {
    Topic     string
    Key       string
    Value     string
    Headers   map[string]string
    Offset    int64
    Partition int32
    Timestamp time.Time
}

// ProduceResult 生产结果
type ProduceResult struct {
    MessageID string
    Offset    int64
    Error     error
}

// BatchProduceResult 批量生产结果
type BatchProduceResult struct {
    Results []*ProduceResult
    Errors  []error
}
```

## 错误处理 (errors/)

```go
// Error SDK错误类型
type Error struct {
    Code    ErrorCode
    Message string
    Cause   error
}

// ErrorCode 错误代码
type ErrorCode int

const (
    ErrConnection ErrorCode = iota
    ErrTimeout
    ErrInvalidConfig
    ErrTopicNotFound
    ErrPermissionDenied
)
```

## 日志接口 (logger/)

```go
// Logger 日志接口
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
}

// Field 日志字段
type Field struct {
    Key   string
    Value interface{}
}
```

## 使用示例

### 基本使用

```go
// 创建客户端
client, err := fluvio.New(
    config.WithServer("localhost", 50051),
    config.WithTimeout(5*time.Second, 10*time.Second),
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 生产消息
result, err := client.Producer().Produce(ctx, "Hello, Fluvio!", fluvio.ProduceOptions{
    Topic: "my-topic",
    Key:   "key1",
})

// 消费消息
messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
    Topic:       "my-topic",
    Group:       "my-group",
    MaxMessages: 10,
})
```

### 高级使用

```go
// 自定义配置
cfg := &config.Config{
    Server: config.ServerConfig{
        Host: "localhost",
        Port: 50051,
    },
    Connection: config.ConnectionConfig{
        ConnectTimeout: 10 * time.Second,
        CallTimeout:    30 * time.Second,
        PoolSize:       5,
    },
}

client, err := fluvio.NewWithConfig(cfg)

// 流式消费
stream, err := client.Consumer().ConsumeStream(ctx, fluvio.ConsumeOptions{
    Topic:  "my-topic",
    Group:  "my-group",
    Offset: 0,
})

for message := range stream {
    // 处理消息
    fmt.Printf("Received: %s\n", message.Value)
}
```

## 迁移策略

1. **阶段1**：创建新的包结构，保持internal包不变
2. **阶段2**：实现新的公共API，内部调用existing代码
3. **阶段3**：重构命令行工具使用新SDK
4. **阶段4**：逐步优化内部实现
5. **阶段5**：添加高级功能和文档

这样的设计确保了向后兼容性，同时提供了清晰的升级路径。
