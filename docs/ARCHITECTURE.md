# Fluvio Go SDK - Clean Architecture 设计文档

## 概述

Fluvio Go SDK v2.0 采用了 Clean Architecture（整洁架构）设计模式，这是由 Robert C. Martin（Uncle Bob）提出的软件架构模式。该架构强调关注点分离、依赖倒置和可测试性。

## 架构原则

### 1. 依赖规则（Dependency Rule）
- 内层不能依赖外层
- 外层可以依赖内层
- 依赖方向始终指向内部

### 2. 关注点分离（Separation of Concerns）
- 每一层都有明确的职责
- 业务逻辑与技术实现分离
- 接口与实现分离

### 3. 可测试性（Testability）
- 每一层都可以独立测试
- 支持依赖注入和模拟
- 业务逻辑不依赖外部系统

## 架构层次详解

### 🎯 Domain Layer（领域层）

**位置**: `domain/`
**职责**: 包含核心业务逻辑和规则
**依赖**: 无外部依赖

#### Entities（实体）
```go
// domain/entities/message.go
type Message struct {
    ID        string
    MessageID string
    Topic     string
    Key       string
    Value     string
    Headers   map[string]string
    Partition int32
    Offset    int64
    Timestamp time.Time
}

func NewMessage(key, value string) *Message {
    return &Message{
        ID:        generateID(),
        Key:       key,
        Value:     value,
        Headers:   make(map[string]string),
        Timestamp: time.Now(),
    }
}

func (m *Message) WithMessageID(id string) *Message {
    m.MessageID = id
    return m
}
```

#### Value Objects（值对象）
```go
// domain/valueobjects/connection_config.go
type ConnectionConfig struct {
    Host           string
    Port           int
    ConnectTimeout time.Duration
    RequestTimeout time.Duration
    // ... 其他配置
}

func (c *ConnectionConfig) Address() string {
    return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *ConnectionConfig) IsValid() bool {
    return c.Host != "" && c.Port > 0
}
```

#### Domain Services（领域服务）
```go
// domain/services/message_service.go
type MessageService struct{}

func NewMessageService() *MessageService {
    return &MessageService{}
}

func (s *MessageService) ValidateMessage(message *entities.Message) error {
    if message.Topic == "" {
        return errors.New(errors.ErrInvalidArgument, "topic cannot be empty")
    }
    if message.Value == "" {
        return errors.New(errors.ErrInvalidArgument, "message value cannot be empty")
    }
    return nil
}
```

#### Repository Interfaces（仓储接口）
```go
// domain/repositories/message_repository.go
type MessageRepository interface {
    Produce(ctx context.Context, message *entities.Message) error
    ProduceBatch(ctx context.Context, messages []*entities.Message) error
    Consume(ctx context.Context, topic string, partition int32, offset int64, maxMessages int) ([]*entities.Message, error)
    ConsumeFiltered(ctx context.Context, topic string, filters []*valueobjects.FilterCondition, maxMessages int) ([]*entities.Message, error)
}
```

### 🎮 Application Layer（应用层）

**位置**: `application/`
**职责**: 协调领域对象完成业务用例
**依赖**: Domain Layer

#### Use Cases（用例）
```go
// application/usecases/produce_message_usecase.go
type ProduceMessageUseCase struct {
    messageRepo    repositories.MessageRepository
    messageService *services.MessageService
}

func NewProduceMessageUseCase(
    messageRepo repositories.MessageRepository,
    messageService *services.MessageService,
) *ProduceMessageUseCase {
    return &ProduceMessageUseCase{
        messageRepo:    messageRepo,
        messageService: messageService,
    }
}

func (uc *ProduceMessageUseCase) Execute(ctx context.Context, req *dtos.ProduceMessageRequest) (*dtos.ProduceMessageResponse, error) {
    // 1. 创建领域实体
    message := entities.NewMessage(req.Key, req.Value)
    message.Topic = req.Topic
    message.MessageID = req.MessageID
    message.Headers = req.Headers
    
    // 2. 业务验证
    if err := uc.messageService.ValidateMessage(message); err != nil {
        return &dtos.ProduceMessageResponse{
            Success: false,
            Error:   err.Error(),
        }, err
    }
    
    // 3. 执行业务操作
    if err := uc.messageRepo.Produce(ctx, message); err != nil {
        return &dtos.ProduceMessageResponse{
            Success: false,
            Error:   err.Error(),
        }, err
    }
    
    // 4. 返回结果
    return &dtos.ProduceMessageResponse{
        Success:   true,
        MessageID: message.MessageID,
        Topic:     message.Topic,
        Partition: message.Partition,
        Offset:    message.Offset,
    }, nil
}
```

#### Application Services（应用服务）
```go
// application/services/fluvio_application_service.go
type FluvioApplicationService struct {
    produceMessageUC *usecases.ProduceMessageUseCase
    consumeMessageUC *usecases.ConsumeMessageUseCase
    manageTopicUC    *usecases.ManageTopicUseCase
}

func (s *FluvioApplicationService) ProduceMessage(ctx context.Context, req *dtos.ProduceMessageRequest) (*dtos.ProduceMessageResponse, error) {
    return s.produceMessageUC.Execute(ctx, req)
}
```

#### DTOs（数据传输对象）
```go
// application/dtos/message_dtos.go
type ProduceMessageRequest struct {
    Topic     string            `json:"topic"`
    Key       string            `json:"key"`
    Value     string            `json:"value"`
    MessageID string            `json:"message_id,omitempty"`
    Headers   map[string]string `json:"headers,omitempty"`
}

type ProduceMessageResponse struct {
    Success   bool   `json:"success"`
    MessageID string `json:"message_id,omitempty"`
    Topic     string `json:"topic,omitempty"`
    Partition int32  `json:"partition,omitempty"`
    Offset    int64  `json:"offset,omitempty"`
    Error     string `json:"error,omitempty"`
}
```

### 🔧 Infrastructure Layer（基础设施层）

**位置**: `infrastructure/`
**职责**: 提供技术实现和外部系统集成
**依赖**: Domain Layer

#### Repository Implementations（仓储实现）
```go
// infrastructure/repositories/grpc_message_repository.go
type GRPCMessageRepository struct {
    client grpc.Client
    logger logging.Logger
}

func NewGRPCMessageRepository(client grpc.Client, logger logging.Logger) repositories.MessageRepository {
    return &GRPCMessageRepository{
        client: client,
        logger: logger,
    }
}

func (r *GRPCMessageRepository) Produce(ctx context.Context, message *entities.Message) error {
    // 转换为protobuf消息
    pbMessage := &pb.Message{
        Topic:     message.Topic,
        Key:       message.Key,
        Value:     message.Value,
        MessageId: message.MessageID,
        Headers:   message.Headers,
    }
    
    req := &pb.ProduceRequest{Message: pbMessage}
    
    // 调用gRPC服务
    resp, err := r.client.Produce(ctx, req)
    if err != nil {
        return err
    }
    
    // 更新实体状态
    message.ID = resp.GetMessageId()
    message.Partition = resp.GetPartition()
    message.Offset = resp.GetOffset()
    
    return nil
}
```

#### gRPC Client（gRPC客户端）
```go
// infrastructure/grpc/client.go
type Client interface {
    Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error)
    Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error)
    // ... 其他方法
}
```

#### Configuration（配置管理）
```go
// infrastructure/config/config.go
type Config struct {
    Connection *valueobjects.ConnectionConfig
    Logging    *LoggingConfig
    Client     *ClientConfig
}

func NewDefaultConfig() *Config {
    return &Config{
        Connection: valueobjects.NewConnectionConfig("localhost", 50051),
        Logging:    &LoggingConfig{Level: "info"},
        Client:     &ClientConfig{UserAgent: "fluvio-go-sdk/2.0.0"},
    }
}
```

### 🌐 Interfaces Layer（接口层）

**位置**: `interfaces/`
**职责**: 对外提供API和适配器
**依赖**: Application Layer

#### API Interfaces（API接口）
```go
// interfaces/api/fluvio_api.go
type FluvioAPI interface {
    Connect() error
    Close() error
    HealthCheck(ctx context.Context) error
    
    Producer() ProducerAPI
    Consumer() ConsumerAPI
    Topic() TopicAPI
    Admin() AdminAPI
}

type ProducerAPI interface {
    Produce(ctx context.Context, value string, opts ProduceOptions) (*ProduceResult, error)
    ProduceBatch(ctx context.Context, messages []Message) (*ProduceBatchResult, error)
}
```

#### Client Adapters（客户端适配器）
```go
// interfaces/client/fluvio_client_adapter.go
type FluvioClientAdapter struct {
    appService *services.FluvioApplicationService
    connected  bool
}

func NewFluvioClientAdapter(appService *services.FluvioApplicationService) api.FluvioAPI {
    return &FluvioClientAdapter{
        appService: appService,
        connected:  false,
    }
}

func (c *FluvioClientAdapter) Producer() api.ProducerAPI {
    return &ProducerAdapter{
        appService: c.appService,
        connected:  &c.connected,
    }
}
```

## 依赖注入

### 手动依赖注入
```go
// 创建依赖
logger := logging.NewDefaultLogger()
config := config.NewDefaultConfig()
connectionManager := grpc.NewConnectionManager(config.Connection, logger)

// 创建仓储
messageRepo := repositories.NewGRPCMessageRepository(connectionManager, logger)
topicRepo := repositories.NewGRPCTopicRepository(connectionManager, logger)

// 创建领域服务
messageService := services.NewMessageService()
topicService := services.NewTopicService()

// 创建用例
produceUC := usecases.NewProduceMessageUseCase(messageRepo, messageService)
consumeUC := usecases.NewConsumeMessageUseCase(messageRepo, messageService)
topicUC := usecases.NewManageTopicUseCase(topicRepo, topicService)

// 创建应用服务
appService := services.NewFluvioApplicationService(produceUC, consumeUC, topicUC)

// 创建客户端适配器
client := client.NewFluvioClientAdapter(appService)
```

### 工厂模式
```go
// fluvio_new.go
func NewClient(opts ...ClientOption) (api.FluvioAPI, error) {
    // 创建配置
    cfg := config.NewDefaultConfig()
    for _, opt := range opts {
        if err := opt(cfg); err != nil {
            return nil, err
        }
    }
    
    // 创建所有依赖
    logger := logging.NewDefaultLogger()
    // ... 创建其他依赖
    
    // 组装并返回客户端
    return client.NewFluvioClientAdapter(appService), nil
}
```

## 测试策略

### 单元测试
```go
// 测试领域服务
func TestMessageService_ValidateMessage(t *testing.T) {
    service := services.NewMessageService()
    
    // 测试有效消息
    message := entities.NewMessage("key", "value")
    message.Topic = "test-topic"
    err := service.ValidateMessage(message)
    assert.NoError(t, err)
    
    // 测试无效消息
    invalidMessage := entities.NewMessage("", "")
    err = service.ValidateMessage(invalidMessage)
    assert.Error(t, err)
}

// 测试用例
func TestProduceMessageUseCase_Execute(t *testing.T) {
    // 创建模拟依赖
    mockRepo := &MockMessageRepository{}
    mockService := &MockMessageService{}
    
    useCase := usecases.NewProduceMessageUseCase(mockRepo, mockService)
    
    // 设置期望
    mockService.On("ValidateMessage", mock.Anything).Return(nil)
    mockRepo.On("Produce", mock.Anything, mock.Anything).Return(nil)
    
    // 执行测试
    req := &dtos.ProduceMessageRequest{
        Topic: "test-topic",
        Key:   "test-key",
        Value: "test-value",
    }
    
    resp, err := useCase.Execute(context.Background(), req)
    
    // 验证结果
    assert.NoError(t, err)
    assert.True(t, resp.Success)
    mockRepo.AssertExpectations(t)
    mockService.AssertExpectations(t)
}
```

### 集成测试
```go
func TestFluvioClient_Integration(t *testing.T) {
    // 使用真实的依赖进行集成测试
    client, err := fluvio.NewClient(
        fluvio.WithServerAddress("localhost", 50051),
        fluvio.WithLogLevelV2("debug"),
    )
    require.NoError(t, err)
    defer client.Close()
    
    ctx := context.Background()
    
    // 测试健康检查
    err = client.HealthCheck(ctx)
    assert.NoError(t, err)
    
    // 测试生产和消费
    result, err := client.Producer().Produce(ctx, "test message", api.ProduceOptions{
        Topic: "test-topic",
    })
    assert.NoError(t, err)
    assert.NotEmpty(t, result.MessageID)
}
```

## 性能优化

### 连接池
```go
// infrastructure/grpc/connection_pool.go
type ConnectionPool struct {
    config      *valueobjects.ConnectionConfig
    pool        chan *grpc.ClientConn
    mu          sync.RWMutex
    activeConns int
}

func (p *ConnectionPool) Get(ctx context.Context) (*grpc.ClientConn, error) {
    // 从池中获取连接或创建新连接
}

func (p *ConnectionPool) Put(conn *grpc.ClientConn) {
    // 将连接放回池中
}
```

### 重试机制
```go
// pkg/utils/retry.go
type Retryer struct {
    config *RetryConfig
    logger logging.Logger
}

func (r *Retryer) RetryWithContext(ctx context.Context, fn RetryableContextFunc) error {
    // 实现指数退避重试
}
```

### 批量操作
```go
// 批量生产消息
func (uc *ProduceMessageUseCase) ExecuteBatch(ctx context.Context, req *dtos.ProduceBatchRequest) (*dtos.ProduceBatchResponse, error) {
    // 批量处理消息
}
```

## 扩展性

### 添加新功能
1. **添加新实体**: 在 `domain/entities/` 中定义
2. **添加新用例**: 在 `application/usecases/` 中实现
3. **添加新仓储**: 在 `domain/repositories/` 定义接口，在 `infrastructure/repositories/` 实现
4. **添加新API**: 在 `interfaces/api/` 定义，在 `interfaces/client/` 实现适配器

### 替换实现
```go
// 可以轻松替换任何层的实现
type CustomMessageRepository struct {
    // 自定义实现
}

func (r *CustomMessageRepository) Produce(ctx context.Context, message *entities.Message) error {
    // 自定义逻辑
}

// 注入自定义实现
useCase := usecases.NewProduceMessageUseCase(&CustomMessageRepository{}, messageService)
```

## 最佳实践

### 1. 保持依赖方向
- 内层不依赖外层
- 使用接口定义依赖
- 通过依赖注入提供实现

### 2. 单一职责原则
- 每个类只有一个变化的理由
- 分离业务逻辑和技术实现
- 保持方法简洁

### 3. 接口隔离原则
- 定义小而专注的接口
- 客户端不应依赖它不使用的方法
- 使用组合而非继承

### 4. 开闭原则
- 对扩展开放，对修改关闭
- 通过接口和依赖注入实现扩展
- 避免修改现有代码

### 5. 依赖倒置原则
- 高层模块不依赖低层模块
- 抽象不依赖细节
- 细节依赖抽象

## 总结

Clean Architecture 为 Fluvio Go SDK 提供了：

1. **清晰的代码组织**: 每一层都有明确的职责
2. **高可测试性**: 每一层都可以独立测试
3. **灵活性**: 可以轻松替换任何层的实现
4. **可维护性**: 业务逻辑与技术实现分离
5. **可扩展性**: 易于添加新功能和修改现有功能

这种架构确保了代码的长期可维护性和可扩展性，同时保持了向后兼容性。