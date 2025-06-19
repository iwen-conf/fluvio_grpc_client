# Fluvio Go SDK API 参考文档

本文档提供了 Fluvio Go SDK v2.0 的完整 API 参考。

## 目录

- [客户端 (Client)](#客户端-client)
- [配置选项 (Configuration)](#配置选项-configuration)
- [生产者 (Producer)](#生产者-producer)
- [消费者 (Consumer)](#消费者-consumer)
- [主题管理 (Topics)](#主题管理-topics)
- [管理功能 (Admin)](#管理功能-admin)
- [类型定义 (Types)](#类型定义-types)
- [错误处理 (Errors)](#错误处理-errors)

---

## 客户端 (Client)

### NewClient

创建新的 Fluvio 客户端实例。

```go
func NewClient(opts ...ClientOption) (*Client, error)
```

**参数:**
- `opts ...ClientOption` - 客户端配置选项

**返回值:**
- `*Client` - 客户端实例
- `error` - 错误信息

**示例:**
```go
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithTimeout(30*time.Second),
)
```

### Client 方法

#### Connect

连接到 Fluvio 服务器。

```go
func (c *Client) Connect(ctx context.Context) error
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `error` - 错误信息

#### Close

关闭客户端连接。

```go
func (c *Client) Close() error
```

**返回值:**
- `error` - 错误信息

#### IsConnected

检查客户端是否已连接。

```go
func (c *Client) IsConnected() bool
```

**返回值:**
- `bool` - 连接状态

#### HealthCheck

执行健康检查。

```go
func (c *Client) HealthCheck(ctx context.Context) error
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `error` - 错误信息

#### Ping

测试与服务器的连接延迟。

```go
func (c *Client) Ping(ctx context.Context) (time.Duration, error)
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `time.Duration` - 延迟时间
- `error` - 错误信息

#### Producer

获取生产者实例。

```go
func (c *Client) Producer() *Producer
```

**返回值:**
- `*Producer` - 生产者实例

#### Consumer

获取消费者实例。

```go
func (c *Client) Consumer() *Consumer
```

**返回值:**
- `*Consumer` - 消费者实例

#### Topics

获取主题管理器实例。

```go
func (c *Client) Topics() *TopicManager
```

**返回值:**
- `*TopicManager` - 主题管理器实例

#### Admin

获取管理器实例。

```go
func (c *Client) Admin() *AdminManager
```

**返回值:**
- `*AdminManager` - 管理器实例

#### Logger

获取日志器实例。

```go
func (c *Client) Logger() logging.Logger
```

**返回值:**
- `logging.Logger` - 日志器实例

---

## 配置选项 (Configuration)

### WithAddress

设置服务器地址和端口。

```go
func WithAddress(host string, port int) ClientOption
```

**参数:**
- `host string` - 服务器主机名或IP
- `port int` - 服务器端口

### WithTimeout

设置操作超时时间。

```go
func WithTimeout(timeout time.Duration) ClientOption
```

**参数:**
- `timeout time.Duration` - 超时时间

### WithRetry

设置重试配置。

```go
func WithRetry(maxAttempts int, baseDelay time.Duration) ClientOption
```

**参数:**
- `maxAttempts int` - 最大重试次数
- `baseDelay time.Duration` - 基础延迟时间

### WithLogLevel

设置日志级别。

```go
func WithLogLevel(level LogLevel) ClientOption
```

**参数:**
- `level LogLevel` - 日志级别

**日志级别常量:**
- `LogLevelDebug` - 调试级别
- `LogLevelInfo` - 信息级别
- `LogLevelWarn` - 警告级别
- `LogLevelError` - 错误级别
- `LogLevelFatal` - 致命错误级别

### WithConnectionPool

设置连接池配置。

```go
func WithConnectionPool(maxConnections int, maxIdleTime time.Duration) ClientOption
```

**参数:**
- `maxConnections int` - 最大连接数
- `maxIdleTime time.Duration` - 最大空闲时间

### WithKeepAlive

设置 Keep-Alive 配置。

```go
func WithKeepAlive(interval time.Duration) ClientOption
```

**参数:**
- `interval time.Duration` - Keep-Alive 间隔

### WithTLS

设置 TLS 配置。

```go
func WithTLS(certFile, keyFile, caFile string) ClientOption
```

**参数:**
- `certFile string` - 客户端证书文件路径
- `keyFile string` - 客户端私钥文件路径
- `caFile string` - CA 证书文件路径

### WithInsecure

禁用 TLS（仅用于开发环境）。

```go
func WithInsecure() ClientOption
```

---

## 生产者 (Producer)

### Send

发送单条消息。

```go
func (p *Producer) Send(ctx context.Context, topic string, message *Message) (*SendResult, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `message *Message` - 消息对象

**返回值:**
- `*SendResult` - 发送结果
- `error` - 错误信息

### SendString

发送字符串消息（便捷方法）。

```go
func (p *Producer) SendString(ctx context.Context, topic, key, value string) (*SendResult, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `key string` - 消息键
- `value string` - 消息值

**返回值:**
- `*SendResult` - 发送结果
- `error` - 错误信息

### SendBatch

批量发送消息。

```go
func (p *Producer) SendBatch(ctx context.Context, topic string, messages []*Message) (*BatchSendResult, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `messages []*Message` - 消息列表

**返回值:**
- `*BatchSendResult` - 批量发送结果
- `error` - 错误信息

### SendJSON

发送 JSON 消息（便捷方法）。

```go
func (p *Producer) SendJSON(ctx context.Context, topic, key string, value interface{}) (*SendResult, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `key string` - 消息键
- `value interface{}` - 要序列化为 JSON 的对象

**返回值:**
- `*SendResult` - 发送结果
- `error` - 错误信息

---

## 消费者 (Consumer)

### Receive

批量接收消息。

```go
func (c *Consumer) Receive(ctx context.Context, topic string, opts *ReceiveOptions) ([]*ConsumedMessage, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `opts *ReceiveOptions` - 接收选项

**返回值:**
- `[]*ConsumedMessage` - 消息列表
- `error` - 错误信息

### ReceiveOne

接收单条消息（便捷方法）。

```go
func (c *Consumer) ReceiveOne(ctx context.Context, topic string, group string) (*ConsumedMessage, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `group string` - 消费者组

**返回值:**
- `*ConsumedMessage` - 消息对象
- `error` - 错误信息

### ReceiveString

接收字符串消息（便捷方法）。

```go
func (c *Consumer) ReceiveString(ctx context.Context, topic string, opts *ReceiveOptions) ([]string, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `opts *ReceiveOptions` - 接收选项

**返回值:**
- `[]string` - 字符串消息列表
- `error` - 错误信息

### Stream

流式消费消息。

```go
func (c *Consumer) Stream(ctx context.Context, topic string, opts *StreamOptions) (<-chan *ConsumedMessage, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `opts *StreamOptions` - 流式选项

**返回值:**
- `<-chan *ConsumedMessage` - 消息通道
- `error` - 错误信息

### Commit

提交偏移量。

```go
func (c *Consumer) Commit(ctx context.Context, topic string, group string, offset int64) error
```

**参数:**
- `ctx context.Context` - 上下文
- `topic string` - 主题名称
- `group string` - 消费者组
- `offset int64` - 偏移量

**返回值:**
- `error` - 错误信息

---

## 主题管理 (Topics)

### Create

创建主题。

```go
func (t *TopicManager) Create(ctx context.Context, name string, opts *CreateTopicOptions) error
```

**参数:**
- `ctx context.Context` - 上下文
- `name string` - 主题名称
- `opts *CreateTopicOptions` - 创建选项

**返回值:**
- `error` - 错误信息

### CreateIfNotExists

创建主题（如果不存在）。

```go
func (t *TopicManager) CreateIfNotExists(ctx context.Context, name string, opts *CreateTopicOptions) (bool, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `name string` - 主题名称
- `opts *CreateTopicOptions` - 创建选项

**返回值:**
- `bool` - 是否创建了新主题
- `error` - 错误信息

### Delete

删除主题。

```go
func (t *TopicManager) Delete(ctx context.Context, name string) error
```

**参数:**
- `ctx context.Context` - 上下文
- `name string` - 主题名称

**返回值:**
- `error` - 错误信息

### List

列出所有主题。

```go
func (t *TopicManager) List(ctx context.Context) ([]string, error)
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `[]string` - 主题名称列表
- `error` - 错误信息

### Info

获取主题信息。

```go
func (t *TopicManager) Info(ctx context.Context, name string) (*TopicInfo, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `name string` - 主题名称

**返回值:**
- `*TopicInfo` - 主题信息
- `error` - 错误信息

### Exists

检查主题是否存在。

```go
func (t *TopicManager) Exists(ctx context.Context, name string) (bool, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `name string` - 主题名称

**返回值:**
- `bool` - 是否存在
- `error` - 错误信息

---

## 管理功能 (Admin)

### ClusterInfo

获取集群信息。

```go
func (a *AdminManager) ClusterInfo(ctx context.Context) (*ClusterInfo, error)
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `*ClusterInfo` - 集群信息
- `error` - 错误信息

### Brokers

获取 Broker 列表。

```go
func (a *AdminManager) Brokers(ctx context.Context) ([]*BrokerInfo, error)
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `[]*BrokerInfo` - Broker 信息列表
- `error` - 错误信息

### ConsumerGroups

获取消费者组列表。

```go
func (a *AdminManager) ConsumerGroups(ctx context.Context) ([]*ConsumerGroupInfo, error)
```

**参数:**
- `ctx context.Context` - 上下文

**返回值:**
- `[]*ConsumerGroupInfo` - 消费者组信息列表
- `error` - 错误信息

### ConsumerGroupDetail

获取消费者组详细信息。

```go
func (a *AdminManager) ConsumerGroupDetail(ctx context.Context, groupID string) (*ConsumerGroupDetail, error)
```

**参数:**
- `ctx context.Context` - 上下文
- `groupID string` - 消费者组ID

**返回值:**
- `*ConsumerGroupDetail` - 消费者组详细信息
- `error` - 错误信息

### SmartModules

获取 SmartModule 管理器。

```go
func (a *AdminManager) SmartModules() *SmartModuleManager
```

**返回值:**
- `*SmartModuleManager` - SmartModule 管理器

---

## 类型定义 (Types)

### Message

消息结构体。

```go
type Message struct {
    Key       string            `json:"key,omitempty"`
    Value     []byte            `json:"value"`
    Headers   map[string]string `json:"headers,omitempty"`
    Timestamp time.Time         `json:"timestamp,omitempty"`
}
```

**字段:**
- `Key` - 消息键
- `Value` - 消息值（二进制数据）
- `Headers` - 消息头
- `Timestamp` - 时间戳

### ConsumedMessage

消费的消息结构体。

```go
type ConsumedMessage struct {
    *Message
    Offset    int64     `json:"offset"`
    Partition int32     `json:"partition"`
    Topic     string    `json:"topic"`
    Timestamp time.Time `json:"timestamp"`
}
```

**字段:**
- `Message` - 嵌入的消息对象
- `Offset` - 偏移量
- `Partition` - 分区
- `Topic` - 主题
- `Timestamp` - 时间戳

### SendResult

发送结果结构体。

```go
type SendResult struct {
    MessageID string `json:"message_id"`
    Offset    int64  `json:"offset"`
    Partition int32  `json:"partition"`
}
```

**字段:**
- `MessageID` - 消息ID
- `Offset` - 偏移量
- `Partition` - 分区

### BatchSendResult

批量发送结果结构体。

```go
type BatchSendResult struct {
    Results      []*SendResult `json:"results"`
    SuccessCount int           `json:"success_count"`
    FailureCount int           `json:"failure_count"`
}
```

**字段:**
- `Results` - 发送结果列表
- `SuccessCount` - 成功数量
- `FailureCount` - 失败数量

### ReceiveOptions

接收选项结构体。

```go
type ReceiveOptions struct {
    Group       string `json:"group"`
    Offset      int64  `json:"offset"`
    MaxMessages int    `json:"max_messages"`
}
```

**字段:**
- `Group` - 消费者组
- `Offset` - 起始偏移量
- `MaxMessages` - 最大消息数

### StreamOptions

流式选项结构体。

```go
type StreamOptions struct {
    Group      string `json:"group"`
    Offset     int64  `json:"offset"`
    BufferSize int    `json:"buffer_size"`
}
```

**字段:**
- `Group` - 消费者组
- `Offset` - 起始偏移量
- `BufferSize` - 缓冲区大小

### CreateTopicOptions

创建主题选项结构体。

```go
type CreateTopicOptions struct {
    Partitions        int32             `json:"partitions"`
    ReplicationFactor int32             `json:"replication_factor"`
    Config            map[string]string `json:"config"`
}
```

**字段:**
- `Partitions` - 分区数
- `ReplicationFactor` - 副本因子
- `Config` - 主题配置

### TopicInfo

主题信息结构体。

```go
type TopicInfo struct {
    Name              string            `json:"name"`
    Partitions        int32             `json:"partitions"`
    ReplicationFactor int32             `json:"replication_factor"`
    Config            map[string]string `json:"config"`
}
```

**字段:**
- `Name` - 主题名称
- `Partitions` - 分区数
- `ReplicationFactor` - 副本因子
- `Config` - 主题配置

### ClusterInfo

集群信息结构体。

```go
type ClusterInfo struct {
    ID           string `json:"id"`
    Status       string `json:"status"`
    ControllerID int32  `json:"controller_id"`
}
```

**字段:**
- `ID` - 集群ID
- `Status` - 集群状态
- `ControllerID` - 控制器ID

### BrokerInfo

Broker 信息结构体。

```go
type BrokerInfo struct {
    ID     int32  `json:"id"`
    Host   string `json:"host"`
    Port   int32  `json:"port"`
    Status string `json:"status"`
    Addr   string `json:"addr"`
}
```

**字段:**
- `ID` - Broker ID
- `Host` - 主机名
- `Port` - 端口
- `Status` - 状态
- `Addr` - 地址

### ConsumerGroupInfo

消费者组信息结构体。

```go
type ConsumerGroupInfo struct {
    GroupID string `json:"group_id"`
    State   string `json:"state"`
}
```

**字段:**
- `GroupID` - 消费者组ID
- `State` - 状态

---

## 错误处理 (Errors)

### 错误类型

SDK 定义了以下错误类型：

```go
const (
    ErrConnection      ErrorCode = "CONNECTION_ERROR"
    ErrTimeout         ErrorCode = "TIMEOUT_ERROR"
    ErrInvalidArgument ErrorCode = "INVALID_ARGUMENT"
    ErrNotFound        ErrorCode = "NOT_FOUND"
    ErrAlreadyExists   ErrorCode = "ALREADY_EXISTS"
    ErrPermission      ErrorCode = "PERMISSION_DENIED"
    ErrAuthentication  ErrorCode = "AUTHENTICATION_ERROR"
    ErrInternal        ErrorCode = "INTERNAL_ERROR"
    ErrUnavailable     ErrorCode = "UNAVAILABLE"
    ErrCancelled       ErrorCode = "CANCELLED"
    ErrOperation       ErrorCode = "OPERATION_ERROR"
)
```

### 错误检查函数

#### IsConnectionError

检查是否为连接错误。

```go
func IsConnectionError(err error) bool
```

#### IsTimeoutError

检查是否为超时错误。

```go
func IsTimeoutError(err error) bool
```

#### IsValidationError

检查是否为验证错误。

```go
func IsValidationError(err error) bool
```

#### IsNotFoundError

检查是否为未找到错误。

```go
func IsNotFoundError(err error) bool
```

#### IsAlreadyExistsError

检查是否为已存在错误。

```go
func IsAlreadyExistsError(err error) bool
```

#### IsAuthenticationError

检查是否为认证错误。

```go
func IsAuthenticationError(err error) bool
```

### 错误处理示例

```go
import "github.com/iwen-conf/fluvio_grpc_client/pkg/errors"

result, err := client.Producer().SendString(ctx, "topic", "key", "value")
if err != nil {
    switch {
    case errors.IsConnectionError(err):
        // 处理连接错误
        log.Println("连接错误，尝试重新连接")
    case errors.IsTimeoutError(err):
        // 处理超时错误
        log.Println("操作超时，增加超时时间")
    case errors.IsValidationError(err):
        // 处理验证错误
        log.Println("参数验证失败")
    default:
        // 处理其他错误
        log.Printf("未知错误: %v", err)
    }
}
```

---

## 版本信息

### Version

获取 SDK 版本。

```go
func Version() string
```

**返回值:**
- `string` - 版本号

### UserAgent

获取用户代理字符串。

```go
func UserAgent() string
```

**返回值:**
- `string` - 用户代理字符串

---

## 注意事项

1. **上下文管理**: 所有异步操作都应该使用适当的上下文来控制超时和取消。

2. **资源清理**: 总是调用 `client.Close()` 来清理资源。

3. **错误处理**: 使用提供的错误检查函数来处理不同类型的错误。

4. **并发安全**: 客户端实例是并发安全的，可以在多个 goroutine 中使用。

5. **性能优化**: 对于高吞吐量场景，使用批量操作和流式消费。

6. **连接管理**: 使用连接池和 Keep-Alive 来优化连接性能。