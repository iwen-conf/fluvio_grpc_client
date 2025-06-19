# 新API设计文档

## 设计原则

1. **简洁性**：用户只需几行代码就能开始使用
2. **类型安全**：充分利用Go的类型系统
3. **上下文感知**：所有操作都支持context.Context
4. **函数式选项**：使用函数式选项模式进行配置
5. **Clean Architecture**：基于现有的Clean Architecture结构

## 新API结构

### 主入口

```go
// 主包入口 - fluvio.go
package fluvio

// 创建客户端
func NewClient(opts ...ClientOption) (*Client, error)

// 配置选项
func WithAddress(host string, port int) ClientOption
func WithTimeout(timeout time.Duration) ClientOption
func WithRetry(maxRetries int, backoff time.Duration) ClientOption
func WithTLS(certFile, keyFile, caFile string) ClientOption
func WithLogger(logger Logger) ClientOption
func WithLogLevel(level LogLevel) ClientOption
func WithConnectionPool(size int, maxIdle time.Duration) ClientOption
```

### 核心客户端

```go
// Client 主客户端接口
type Client struct {
    // 内部实现基于Clean Architecture
}

// 连接管理
func (c *Client) Connect(ctx context.Context) error
func (c *Client) Close() error
func (c *Client) Ping(ctx context.Context) (time.Duration, error)
func (c *Client) HealthCheck(ctx context.Context) error

// 功能模块
func (c *Client) Producer() *Producer
func (c *Client) Consumer() *Consumer
func (c *Client) Topics() *TopicManager
func (c *Client) Admin() *AdminManager
```

### 生产者

```go
type Producer struct {}

// 生产单条消息
func (p *Producer) Send(ctx context.Context, topic string, message *Message) (*SendResult, error)

// 生产消息（带选项）
func (p *Producer) SendWithOptions(ctx context.Context, opts *SendOptions) (*SendResult, error)

// 批量生产
func (p *Producer) SendBatch(ctx context.Context, topic string, messages []*Message) (*BatchSendResult, error)
```

### 消费者

```go
type Consumer struct {}

// 消费消息
func (c *Consumer) Receive(ctx context.Context, topic string, opts *ReceiveOptions) ([]*Message, error)

// 流式消费
func (c *Consumer) Stream(ctx context.Context, topic string, opts *StreamOptions) (<-chan *Message, error)

// 提交偏移量
func (c *Consumer) Commit(ctx context.Context, topic string, group string, offset int64) error
```

### 主题管理

```go
type TopicManager struct {}

// 创建主题
func (t *TopicManager) Create(ctx context.Context, name string, opts *CreateTopicOptions) error

// 删除主题
func (t *TopicManager) Delete(ctx context.Context, name string) error

// 列出主题
func (t *TopicManager) List(ctx context.Context) ([]string, error)

// 获取主题信息
func (t *TopicManager) Info(ctx context.Context, name string) (*TopicInfo, error)

// 检查主题是否存在
func (t *TopicManager) Exists(ctx context.Context, name string) (bool, error)
```

### 管理功能

```go
type AdminManager struct {}

// 集群信息
func (a *AdminManager) ClusterInfo(ctx context.Context) (*ClusterInfo, error)

// Broker列表
func (a *AdminManager) Brokers(ctx context.Context) ([]*BrokerInfo, error)

// 消费者组
func (a *AdminManager) ConsumerGroups(ctx context.Context) ([]*ConsumerGroupInfo, error)

// SmartModule管理
func (a *AdminManager) SmartModules() *SmartModuleManager
```

## 类型定义

### 核心类型

```go
// Message 消息
type Message struct {
    Key       string            `json:"key,omitempty"`
    Value     []byte            `json:"value"`
    Headers   map[string]string `json:"headers,omitempty"`
    Timestamp time.Time         `json:"timestamp,omitempty"`
}

// SendOptions 发送选项
type SendOptions struct {
    Topic     string            `json:"topic"`
    Key       string            `json:"key,omitempty"`
    Value     []byte            `json:"value"`
    Headers   map[string]string `json:"headers,omitempty"`
    MessageID string            `json:"message_id,omitempty"`
}

// ReceiveOptions 接收选项
type ReceiveOptions struct {
    Group       string `json:"group,omitempty"`
    Offset      int64  `json:"offset,omitempty"`
    MaxMessages int    `json:"max_messages,omitempty"`
    Timeout     time.Duration `json:"timeout,omitempty"`
}

// StreamOptions 流式选项
type StreamOptions struct {
    Group       string        `json:"group,omitempty"`
    Offset      int64         `json:"offset,omitempty"`
    BufferSize  int           `json:"buffer_size,omitempty"`
    Timeout     time.Duration `json:"timeout,omitempty"`
}
```

### 结果类型

```go
// SendResult 发送结果
type SendResult struct {
    MessageID string `json:"message_id"`
    Offset    int64  `json:"offset"`
    Partition int32  `json:"partition"`
}

// BatchSendResult 批量发送结果
type BatchSendResult struct {
    Results      []*SendResult `json:"results"`
    SuccessCount int           `json:"success_count"`
    FailureCount int           `json:"failure_count"`
}

// TopicInfo 主题信息
type TopicInfo struct {
    Name              string            `json:"name"`
    Partitions        int32             `json:"partitions"`
    ReplicationFactor int32             `json:"replication_factor"`
    Config            map[string]string `json:"config,omitempty"`
}

// ClusterInfo 集群信息
type ClusterInfo struct {
    ID           string `json:"id"`
    Status       string `json:"status"`
    ControllerID int32  `json:"controller_id"`
}

// BrokerInfo Broker信息
type BrokerInfo struct {
    ID     int32  `json:"id"`
    Host   string `json:"host"`
    Port   int32  `json:"port"`
    Status string `json:"status"`
}
```

## 使用示例

### 基本使用

```go
// 创建客户端
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithTimeout(30*time.Second),
    fluvio.WithRetry(3, time.Second),
)
if err != nil {
    log.Fatal(err)
}
defer client.Close()

// 连接
if err := client.Connect(ctx); err != nil {
    log.Fatal(err)
}

// 发送消息
result, err := client.Producer().Send(ctx, "my-topic", &fluvio.Message{
    Key:   "user-123",
    Value: []byte("Hello, Fluvio!"),
    Headers: map[string]string{
        "source": "my-app",
    },
})

// 接收消息
messages, err := client.Consumer().Receive(ctx, "my-topic", &fluvio.ReceiveOptions{
    Group:       "my-group",
    MaxMessages: 10,
})

// 主题管理
err = client.Topics().Create(ctx, "new-topic", &fluvio.CreateTopicOptions{
    Partitions:        3,
    ReplicationFactor: 1,
})
```

### 高级使用

```go
// 流式消费
stream, err := client.Consumer().Stream(ctx, "my-topic", &fluvio.StreamOptions{
    Group:      "stream-group",
    BufferSize: 100,
})

go func() {
    for message := range stream {
        fmt.Printf("Received: %s\n", string(message.Value))
    }
}()

// 批量发送
messages := []*fluvio.Message{
    {Key: "key1", Value: []byte("message1")},
    {Key: "key2", Value: []byte("message2")},
}
batchResult, err := client.Producer().SendBatch(ctx, "my-topic", messages)
```

## 迁移策略

1. **保留Clean Architecture结构**：新API基于现有的Clean Architecture
2. **移除向后兼容层**：删除fluvio.go等向后兼容文件
3. **简化接口**：提供更直观的API
4. **保持功能完整性**：确保所有功能都可用
