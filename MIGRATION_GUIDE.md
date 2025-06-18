# Fluvio Go SDK 架构迁移指南

本文档说明了从旧架构到新Clean Architecture架构的迁移过程。

## 架构对比

### 旧架构
```
/
├── client/          # 客户端实现
├── types/           # 类型定义
├── config/          # 配置
├── logger/          # 日志
├── internal/        # 内部实现
├── errors/          # 错误处理
├── proto/           # protobuf
└── examples/        # 示例
```

### 新架构（Clean Architecture）
```
/
├── domain/                    # 领域层（核心业务逻辑）
│   ├── entities/             # 实体
│   ├── valueobjects/         # 值对象
│   ├── services/             # 领域服务
│   └── repositories/         # 仓储接口
├── application/              # 应用层（用例协调）
│   ├── usecases/            # 用例
│   ├── services/            # 应用服务
│   └── dtos/                # 数据传输对象
├── infrastructure/          # 基础设施层（技术实现）
│   ├── grpc/               # gRPC实现
│   ├── repositories/       # 仓储实现
│   ├── config/            # 配置
│   └── logging/           # 日志
├── interfaces/             # 接口层（对外API）
│   ├── api/               # 公共API
│   └── client/            # 客户端接口
├── pkg/                   # 共享包
│   ├── errors/           # 错误处理
│   └── utils/            # 工具函数
├── proto/                # protobuf定义
└── examples/             # 示例代码
```

## 迁移映射

### 类型定义迁移

| 旧位置 | 新位置 | 说明 |
|--------|--------|------|
| `types/message.go` | `domain/entities/message.go` | 消息实体 |
| `types/topic.go` | `domain/entities/topic.go` | 主题实体 |
| `types/consumer.go` | `domain/entities/consumer_group.go` | 消费组实体 |
| `types/admin.go` | `application/dtos/` | 管理相关DTOs |

### 客户端实现迁移

| 旧位置 | 新位置 | 说明 |
|--------|--------|------|
| `client/client.go` | `interfaces/client/fluvio_client_adapter.go` | 主客户端适配器 |
| `client/producer.go` | `interfaces/client/fluvio_client_adapter.go` | 生产者适配器 |
| `client/consumer.go` | `interfaces/client/fluvio_client_adapter.go` | 消费者适配器 |
| `client/topic.go` | `interfaces/client/topic_adapter.go` | 主题适配器 |
| `client/admin.go` | `interfaces/client/admin_adapter.go` | 管理适配器 |

### 基础设施迁移

| 旧位置 | 新位置 | 说明 |
|--------|--------|------|
| `internal/grpc/` | `infrastructure/grpc/` | gRPC实现 |
| `internal/pool/` | `infrastructure/grpc/connection_pool.go` | 连接池 |
| `internal/retry/` | `pkg/utils/retry.go` | 重试机制 |
| `config/` | `infrastructure/config/` | 配置管理 |
| `logger/` | `infrastructure/logging/` | 日志系统 |
| `errors/` | `pkg/errors/` | 错误处理 |

## API兼容性

### 旧API（仍然支持）
```go
// 使用旧的API
client, err := fluvio.New(
    fluvio.WithServer("localhost", 50051),
    fluvio.WithLogLevel(fluvio.LevelInfo),
)
```

### 新API（推荐使用）
```go
// 使用新的Clean Architecture API
client, err := fluvio.NewClient(
    fluvio.WithServerAddress("localhost", 50051),
    fluvio.WithLogLevelV2("info"),
)
```

## 迁移步骤

### 1. 立即可用（向后兼容）
现有代码无需修改，旧API仍然可用：
```go
// 这些代码仍然可以正常工作
client, err := fluvio.New(fluvio.WithServer("localhost", 50051))
result, err := client.Producer().Produce(ctx, "message", fluvio.ProduceOptions{
    Topic: "my-topic",
})
```

### 2. 渐进式迁移（推荐）
逐步迁移到新API：

#### 步骤1：更新客户端创建
```go
// 旧方式
client, err := fluvio.New(fluvio.WithServer("localhost", 50051))

// 新方式
client, err := fluvio.NewClient(fluvio.WithServerAddress("localhost", 50051))
```

#### 步骤2：更新类型引用
```go
// 旧方式
import "github.com/iwen-conf/fluvio_grpc_client/types"

// 新方式
import "github.com/iwen-conf/fluvio_grpc_client/interfaces/api"
```

#### 步骤3：更新错误处理
```go
// 旧方式
import "github.com/iwen-conf/fluvio_grpc_client/errors"

// 新方式
import "github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
```

### 3. 完全迁移（可选）
如果需要使用新架构的高级功能：

```go
// 直接使用应用服务
import (
    "github.com/iwen-conf/fluvio_grpc_client/application/services"
    "github.com/iwen-conf/fluvio_grpc_client/application/dtos"
)

// 使用用例
import (
    "github.com/iwen-conf/fluvio_grpc_client/application/usecases"
    "github.com/iwen-conf/fluvio_grpc_client/domain/services"
)
```

## 新功能

### 1. 依赖注入
新架构支持依赖注入，便于测试和扩展：
```go
// 可以注入自定义仓储实现
messageRepo := &CustomMessageRepository{}
useCase := usecases.NewProduceMessageUseCase(messageRepo, messageService)
```

### 2. 领域驱动设计
清晰的领域模型和业务逻辑分离：
```go
// 领域实体
message := entities.NewMessage("key", "value")
message.WithMessageID("custom-id").WithHeaders(headers)

// 领域服务
messageService := services.NewMessageService()
err := messageService.ValidateMessage(message)
```

### 3. 更好的测试支持
每一层都可以独立测试：
```go
// 测试用例
func TestProduceMessage(t *testing.T) {
    mockRepo := &MockMessageRepository{}
    useCase := usecases.NewProduceMessageUseCase(mockRepo, messageService)
    // 测试逻辑
}
```

### 4. 配置管理增强
更灵活的配置系统：
```go
config := config.NewDefaultConfig()
config.Connection.WithTLS("cert.pem", "key.pem", "ca.pem")
config.Client.CircuitBreaker.Enabled = true
```

## 性能改进

1. **连接池优化**: 新的连接池实现更高效
2. **重试机制增强**: 支持多种退避策略
3. **内存使用优化**: 更好的对象生命周期管理
4. **并发安全**: 全面的并发安全保证

## 故障排除

### 常见问题

1. **导入路径错误**
   ```go
   // 错误
   import "github.com/iwen-conf/fluvio_grpc_client/types"
   
   // 正确
   import "github.com/iwen-conf/fluvio_grpc_client/interfaces/api"
   ```

2. **函数名冲突**
   ```go
   // 如果遇到函数名冲突，使用V2版本
   fluvio.WithLogLevelV2("debug")  // 新架构
   fluvio.WithLogLevel("debug")    // 旧架构
   ```

3. **类型不匹配**
   ```go
   // 确保使用正确的类型
   var opts api.ProduceOptions  // 新架构
   var opts types.ProduceOptions // 旧架构
   ```

### 调试技巧

1. **启用详细日志**
   ```go
   client, err := fluvio.NewClient(
       fluvio.WithLogLevelV2("debug"),
   )
   ```

2. **检查连接状态**
   ```go
   duration, err := client.Ping(ctx)
   fmt.Printf("Ping: %v, Error: %v\n", duration, err)
   ```

3. **监控连接池**
   ```go
   // 在新架构中可以获取连接池统计
   stats := connectionFactory.GetStats()
   fmt.Printf("Pool stats: %+v\n", stats)
   ```

## 未来计划

1. **v3.0**: 完全移除旧API，只保留Clean Architecture
2. **插件系统**: 支持自定义插件和中间件
3. **监控集成**: 内置Prometheus指标和OpenTelemetry追踪
4. **配置热重载**: 支持运行时配置更新

## 支持

如果在迁移过程中遇到问题：

1. 查看示例代码：`examples/` 目录
2. 查看API文档：`interfaces/api/` 目录
3. 提交Issue：GitHub Issues
4. 查看测试用例：各层的测试文件

## 总结

新的Clean Architecture架构提供了：
- ✅ 更好的代码组织和可维护性
- ✅ 清晰的依赖关系和测试能力
- ✅ 向后兼容性保证
- ✅ 更强的扩展性和灵活性
- ✅ 更好的性能和稳定性

建议逐步迁移到新架构，享受Clean Architecture带来的好处！