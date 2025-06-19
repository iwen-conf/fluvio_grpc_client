# Fluvio Go SDK

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-orange.svg)](docs/ARCHITECTURE.md)
[![Version](https://img.shields.io/badge/Version-2.0.0-green.svg)](https://github.com/iwen-conf/fluvio_grpc_client)

## 项目简介

Fluvio Go SDK 是一个现代化的 Go 语言软件开发工具包，用于与 Fluvio 消息流处理系统进行交互。该SDK基于 **Clean Architecture** 设计原则，通过 gRPC 协议提供了丰富的功能，包括消息的生产和消费、主题管理、消费者组管理、SmartModule 管理以及集群管理等功能。

🚀 **v2.0 全新设计**: 采用现代化的 API 设计，简洁易用，类型安全，高性能。

## ✨ 核心特性

- 🎯 **现代化API**: 简洁直观的API设计，类型安全
- 🏗️ **Clean Architecture**: 清晰的分层架构，遵循SOLID原则
- 🚀 **高性能**: 优化的连接池、重试机制和资源管理
- 🧪 **易于测试**: 每一层都可以独立测试，支持依赖注入
- 📦 **模块化设计**: 清晰的模块边界和职责分离
- 🔧 **函数式配置**: 使用函数式选项模式，配置灵活
- 🛡️ **错误处理**: 完善的错误类型和处理机制
- 📊 **可观测性**: 内置日志和指标支持

## 功能特性

### 核心服务 (FluvioService)

- **消息生产/消费**

  - 单条消息生产 (Produce)
  - 批量消息生产 (BatchProduce)
  - 消息消费 (Consume)
  - 流式消息消费 (StreamConsume)
  - 提交消费位点 (CommitOffset)

- **主题管理**

  - 创建主题 (CreateTopic)
  - 删除主题 (DeleteTopic)
  - 列出所有主题 (ListTopics)
  - 获取主题详情 (DescribeTopic)

- **消费者组管理**

  - 列出消费组 (ListConsumerGroups)
  - 获取消费组详情 (DescribeConsumerGroup)

- **SmartModule 管理**

  - 创建 SmartModule (CreateSmartModule)
  - 删除 SmartModule (DeleteSmartModule)
  - 列出 SmartModule (ListSmartModules)
  - 获取 SmartModule 详情 (DescribeSmartModule)
  - 更新 SmartModule (UpdateSmartModule)

- **其他功能**
  - 健康检查 (HealthCheck)

### 管理服务 (FluvioAdminService)

- **集群管理**
  - 获取集群状态 (DescribeCluster)
  - 列出 Broker 信息 (ListBrokers)
  - 获取系统指标 (GetMetrics)

## 🏗️ 项目结构（Clean Architecture）

```
fluvio_grpc_client/
├── domain/                    # 🎯 领域层（核心业务逻辑）
│   ├── entities/             # 业务实体
│   │   ├── message.go        # 消息实体
│   │   ├── topic.go          # 主题实体
│   │   └── consumer_group.go # 消费组实体
│   ├── valueobjects/         # 值对象
│   │   ├── connection_config.go # 连接配置
│   │   └── filter_condition.go # 过滤条件
│   ├── services/             # 领域服务
│   │   ├── message_service.go # 消息业务逻辑
│   │   └── topic_service.go   # 主题业务逻辑
│   └── repositories/         # 仓储接口
│       ├── message_repository.go # 消息仓储接口
│       └── topic_repository.go   # 主题仓储接口
├── application/              # 🎮 应用层（用例协调）
│   ├── usecases/            # 用例
│   │   ├── produce_message_usecase.go # 生产消息用例
│   │   ├── consume_message_usecase.go # 消费消息用例
│   │   └── manage_topic_usecase.go    # 主题管理用例
│   ├── services/            # 应用服务
│   │   └── fluvio_application_service.go # 应用服务
│   └── dtos/                # 数据传输对象
│       ├── message_dtos.go  # 消息DTOs
│       └── topic_dtos.go    # 主题DTOs
├── infrastructure/          # 🔧 基础设施层（技术实现）
│   ├── grpc/               # gRPC实现
│   │   ├── client.go       # gRPC客户端接口
│   │   ├── connection_manager.go # 连接管理
│   │   └── connection_pool.go    # 连接池
│   ├── repositories/       # 仓储实现
│   │   ├── grpc_message_repository.go # gRPC消息仓储
│   │   └── grpc_topic_repository.go   # gRPC主题仓储
│   ├── config/            # 配置管理
│   │   └── config.go      # 配置实现
│   └── logging/           # 日志系统
│       └── logger.go      # 日志实现
├── interfaces/             # 🌐 接口层（对外API）
│   ├── api/               # 公共API定义
│   │   ├── fluvio_api.go  # 主API接口
│   │   └── types.go       # API类型定义
│   └── client/            # 客户端适配器
│       ├── fluvio_client_adapter.go # 主客户端适配器
│       ├── topic_adapter.go         # 主题适配器
│       └── admin_adapter.go         # 管理适配器
├── pkg/                   # 📦 共享包
│   ├── errors/           # 错误处理
│   │   └── errors.go     # 错误类型定义
│   └── utils/            # 工具函数
│       └── retry.go      # 重试机制
├── proto/                # 📡 协议定义
│   └── fluvio_service/   # 生成的gRPC代码
├── examples/             # 📚 使用示例
│   ├── basic/           # 基本示例
│   ├── advanced/        # 高级示例
│   └── integration/     # 集成测试
├── client/              # 🔄 旧API（向后兼容）
├── types/               # 🔄 旧类型（向后兼容）
├── fluvio.go            # 🔄 旧SDK入口（向后兼容）
├── fluvio_new.go        # 🆕 新SDK入口（Clean Architecture）
├── MIGRATION_GUIDE.md   # 📖 迁移指南
├── go.mod               # Go 模块定义
├── go.sum               # 依赖校验和
└── README.md            # 项目说明文档
```

### 🎯 架构层次说明

| 层次 | 职责 | 依赖方向 |
|------|------|----------|
| **Domain** | 核心业务逻辑，不依赖任何外部技术 | 无外部依赖 |
| **Application** | 协调领域对象完成业务用例 | 依赖 Domain |
| **Infrastructure** | 技术实现（数据库、网络、文件等） | 依赖 Domain |
| **Interfaces** | 对外API和适配器 | 依赖 Application |

## 安装与使用

### 前置条件

- Go 1.18 或更高版本
- 正在运行的 Fluvio 服务实例

### 安装

```bash
go get github.com/iwen-conf/fluvio_grpc_client
```

### 🚀 快速开始

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 创建客户端
    client, err := fluvio.NewClient(
        fluvio.WithAddress("localhost", 50051),
        fluvio.WithTimeout(30*time.Second),
        fluvio.WithRetry(3, time.Second),
        fluvio.WithLogLevel(fluvio.LogLevelInfo),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()

    // 连接到服务器
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }

    // 健康检查
    if err := client.HealthCheck(ctx); err != nil {
        log.Fatal(err)
    }

    // 创建主题
    if err := client.Topics().Create(ctx, "my-topic", &fluvio.CreateTopicOptions{
        Partitions:        3,
        ReplicationFactor: 1,
    }); err != nil {
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
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Message sent: %s", result.MessageID)

    // 接收消息
    messages, err := client.Consumer().Receive(ctx, "my-topic", &fluvio.ReceiveOptions{
        Group:       "my-group",
        MaxMessages: 10,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, msg := range messages {
        log.Printf("Received: %s", string(msg.Value))
    }
}
```

### 🎯 便捷方法

```go
// 发送字符串消息
result, err := client.Producer().SendString(ctx, "my-topic", "key", "Hello World")

// 接收单条消息
message, err := client.Consumer().ReceiveOne(ctx, "my-topic", "my-group")

// 流式消费
stream, err := client.Consumer().Stream(ctx, "my-topic", &fluvio.StreamOptions{
    Group:      "stream-group",
    BufferSize: 100,
})

for message := range stream {
    fmt.Printf("Received: %s\n", string(message.Value))
}
```

## 🏗️ Clean Architecture 优势

### 1. 清晰的依赖关系
```go
// 领域层：纯业务逻辑，无外部依赖
type MessageService struct{}
func (s *MessageService) ValidateMessage(msg *entities.Message) error

// 应用层：协调业务用例
type ProduceMessageUseCase struct {
    messageRepo repositories.MessageRepository
    messageService *services.MessageService
}

// 基础设施层：技术实现
type GRPCMessageRepository struct {
    client grpc.Client
}

// 接口层：对外API
type FluvioClientAdapter struct {
    appService *services.FluvioApplicationService
}
```

### 2. 易于测试
```go
// 可以轻松模拟任何依赖进行单元测试
func TestProduceMessage(t *testing.T) {
    mockRepo := &MockMessageRepository{}
    mockService := &MockMessageService{}
    useCase := usecases.NewProduceMessageUseCase(mockRepo, mockService)

    // 测试业务逻辑
    err := useCase.Execute(ctx, request)
    assert.NoError(t, err)
}
```

### 3. 灵活的配置和扩展
```go
// 可以注入自定义实现
customRepo := &MyCustomMessageRepository{}
useCase := usecases.NewProduceMessageUseCase(customRepo, messageService)

// 支持多种配置方式
config := config.NewDefaultConfig()
config.Connection.WithTLS("cert.pem", "key.pem", "ca.pem")
config.Client.CircuitBreaker.Enabled = true
```

## 🆕 新功能示例

### 过滤消费
```go
// 按消息头部过滤
result, err := client.Consumer().ConsumeFiltered(ctx, fluvio.FilteredConsumeOptions{
    Topic: "my-topic",
    Group: "filter-group",
    Filters: []fluvio.FilterCondition{
        {
            Type:     fluvio.FilterTypeHeader,
            Field:    "level",
            Operator: "eq",
            Value:    "error",
        },
    },
    AndLogic: true,
})
```

### 主题详细信息
```go
// 获取主题详细信息
detail, err := client.Topic().DescribeTopicDetail(ctx, "my-topic")
if err == nil {
    fmt.Printf("主题: %s, 分区数: %d\n", detail.Topic, len(detail.Partitions))
    fmt.Printf("保留时间: %d ms\n", detail.RetentionMs)
    fmt.Printf("配置: %v\n", detail.Config)
}
```

### 主题统计信息
```go
// 获取主题统计
stats, err := client.Topic().GetTopicStats(ctx, fluvio.GetTopicStatsOptions{
    Topic:             "my-topic",
    IncludePartitions: true,
})
if err == nil {
    for _, topicStats := range stats.Topics {
        fmt.Printf("主题: %s, 消息数: %d, 大小: %d bytes\n",
            topicStats.Topic, topicStats.TotalMessageCount, topicStats.TotalSizeBytes)
    }
}
```

### 存储管理
```go
// 获取存储状态
status, err := client.Admin().GetStorageStatus(ctx, fluvio.GetStorageStatusOptions{
    IncludeDetails: true,
})
if err == nil {
    fmt.Printf("持久化: %v, 存储类型: %s\n",
        status.PersistenceEnabled, status.StorageStats.StorageType)
}

// 获取存储指标
metrics, err := client.Admin().GetStorageMetrics(ctx, fluvio.GetStorageMetricsOptions{})
if err == nil && metrics.CurrentMetrics != nil {
    fmt.Printf("响应时间: %d ms, 操作/秒: %.2f\n",
        metrics.CurrentMetrics.ResponseTimeMs, metrics.CurrentMetrics.OperationsPerSecond)
}
```

### 批量删除
```go
// 批量删除主题
result, err := client.Admin().BulkDelete(ctx, fluvio.BulkDeleteOptions{
    Topics: []string{"topic1", "topic2", "topic3"},
    Force:  false,
})
if err == nil {
    fmt.Printf("删除结果: %d成功, %d失败\n",
        result.SuccessfulDeletes, result.FailedDeletes)
}
```

## 📖 API 文档

### 客户端创建

```go
// 基本配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithTimeout(30*time.Second),
)

// 完整配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("101.43.173.154", 50051),
    fluvio.WithTimeouts(5*time.Second, 30*time.Second),
    fluvio.WithRetry(3, time.Second),
    fluvio.WithLogLevel(fluvio.LogLevelInfo),
    fluvio.WithConnectionPool(10, 5*time.Minute),
    fluvio.WithTLS("cert.pem", "key.pem", "ca.pem"),
    fluvio.WithKeepAlive(30*time.Second),
)

// 不安全连接（开发环境）
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithInsecure(),
)
```

### 消息生产

```go
// 发送消息
result, err := client.Producer().Send(ctx, "my-topic", &fluvio.Message{
    Key:   "user-123",
    Value: []byte("Hello World"),
    Headers: map[string]string{
        "source": "go-sdk",
        "type":   "greeting",
    },
})

// 使用选项发送
result, err := client.Producer().SendWithOptions(ctx, &fluvio.SendOptions{
    Topic:   "my-topic",
    Key:     "key1",
    Value:   []byte("Hello World"),
    Headers: map[string]string{"source": "app"},
})

// 批量发送
messages := []*fluvio.Message{
    {Key: "key1", Value: []byte("message1")},
    {Key: "key2", Value: []byte("message2")},
}
batchResult, err := client.Producer().SendBatch(ctx, "my-topic", messages)

// 便捷方法
result, err := client.Producer().SendString(ctx, "my-topic", "key", "Hello")
result, err := client.Producer().SendJSON(ctx, "my-topic", "key", map[string]string{"msg": "hello"})
```

### 消息消费

```go
// 基本消费
messages, err := client.Consumer().Receive(ctx, "my-topic", &fluvio.ReceiveOptions{
    Group:       "my-group",
    MaxMessages: 10,
    Offset:      0,
})

// 流式消费
stream, err := client.Consumer().Stream(ctx, "my-topic", &fluvio.StreamOptions{
    Group:      "my-group",
    BufferSize: 100,
    Offset:     0,
})

for msg := range stream {
    fmt.Printf("Received: [%s] %s\n", msg.Key, string(msg.Value))

    // 处理消息...

    // 可选：提交偏移量
    err := client.Consumer().Commit(ctx, "my-topic", "my-group", msg.Offset)
    if err != nil {
        log.Printf("Failed to commit offset: %v", err)
    }
}

// 便捷方法
message, err := client.Consumer().ReceiveOne(ctx, "my-topic", "my-group")
values, err := client.Consumer().ReceiveString(ctx, "my-topic", &fluvio.ReceiveOptions{
    Group:       "my-group",
    MaxMessages: 5,
})
```

### 主题管理

```go
// 列出主题
topics, err := client.Topics().List(ctx)

// 创建主题
err = client.Topics().Create(ctx, "new-topic", &fluvio.CreateTopicOptions{
    Partitions:        3,
    ReplicationFactor: 1,
    Config: map[string]string{
        "cleanup.policy": "delete",
        "segment.ms":     "3600000",
    },
})

// 获取主题信息
info, err := client.Topics().Info(ctx, "my-topic")
fmt.Printf("Topic: %s, Partitions: %d\n", info.Name, info.Partitions)

// 删除主题
err = client.Topics().Delete(ctx, "old-topic")

// 检查主题是否存在
exists, err := client.Topics().Exists(ctx, "my-topic")

// 如果不存在则创建
created, err := client.Topics().CreateIfNotExists(ctx, "my-topic", &fluvio.CreateTopicOptions{
    Partitions: 1,
})
if created {
    fmt.Println("Topic created")
} else {
    fmt.Println("Topic already exists")
}
```

### 管理功能

```go
// 集群信息
clusterInfo, err := client.Admin().ClusterInfo(ctx)
fmt.Printf("Cluster: %s, Status: %s\n", clusterInfo.ID, clusterInfo.Status)

// Broker管理
brokers, err := client.Admin().Brokers(ctx)
for _, broker := range brokers {
    fmt.Printf("Broker %d: %s:%d (%s)\n", broker.ID, broker.Host, broker.Port, broker.Status)
}

// 消费者组管理
groups, err := client.Admin().ConsumerGroups(ctx)
for _, group := range groups {
    fmt.Printf("Group: %s, State: %s\n", group.GroupID, group.State)
}

// SmartModule管理
smartModules, err := client.Admin().SmartModules().List(ctx)
for _, module := range smartModules {
    fmt.Printf("Module: %s, Version: %s\n", module.Name, module.Version)
}

// 创建SmartModule
err = client.Admin().SmartModules().Create(ctx, "my-filter", wasmBytes)

// 删除SmartModule
err = client.Admin().SmartModules().Delete(ctx, "my-filter")
```

或

```bash
quit
```

## 开发指南

### 生成 gRPC 代码

如需修改 proto 文件后重新生成代码，请执行：

```bash
protoc --go_out=. --go-grpc_out=. proto/fluvio_grpc.proto
```

生成的代码将保存在 `proto/fluvio_service/` 目录下。

### 运行测试

```bash
go test ./tests/...
```

测试文件包括健康检查测试和服务功能测试。

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件

## 联系方式

如有任何问题或建议，请通过 [issues](https://github.com/iwen-conf/fluvio_grpc_client/issues) 页面与我们联系。

## 特性

### 核心功能
- **简单易用**: 提供简洁的API接口，快速上手
- **高性能**: 内置连接池和重试机制，支持高并发
- **类型安全**: 完整的类型定义，编译时错误检查
- **可扩展**: 分层架构设计，支持自定义扩展
- **完整文档**: 丰富的示例和API文档

### 🆕 新增功能
- **消息ID支持**: 自定义消息ID，便于追踪和去重
- **过滤消费**: 服务端过滤，支持按键、头部、内容过滤
- **主题增强管理**: 详细配置、分区信息、统计数据
- **存储管理**: 状态监控、性能指标、健康检查
- **SmartModule管理**: 完整生命周期管理和参数化配置
- **批量操作**: 批量删除资源，提高管理效率
- **流式消费增强**: 批次大小控制、等待时间优化

## 📚 文档

### 🏗️ 架构文档
- 🔄 **[迁移指南](MIGRATION_GUIDE.md)** - 从旧架构到Clean Architecture的完整迁移指南
- 🎯 **[架构设计](docs/ARCHITECTURE.md)** - Clean Architecture设计原理和实现细节
- 🧪 **[测试指南](docs/TESTING.md)** - 如何在新架构中编写和运行测试

### 📖 使用文档
- 🚀 **[快速入门](QUICKSTART.md)** - 5分钟快速上手指南
- 📋 **[使用方法](HOW_TO_USE.md)** - 详细的导入、创建客户端和配置说明
- 📖 **[完整使用指南](USAGE.md)** - 详细的API文档和使用示例
- 🔧 **[配置示例](examples/config-example.json)** - 配置文件示例
- 💡 **[示例代码](examples/)** - 基本、高级和集成测试示例

### 🔄 兼容性
- ✅ **向后兼容**: 所有旧API仍然可用
- 🆕 **新API**: 推荐使用新的Clean Architecture API
- 📈 **渐进式迁移**: 可以逐步迁移到新架构
- 🛠️ **工具支持**: 提供迁移工具和检查脚本
