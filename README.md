# Fluvio Go SDK

## 项目简介

Fluvio Go SDK 是一个基于 Go 语言的软件开发工具包，用于与 Fluvio 消息流处理系统进行交互。该SDK通过 gRPC 协议提供了丰富的功能，包括消息的生产和消费、主题管理、消费者组管理、SmartModule 管理以及集群管理等功能。SDK采用分层架构设计，提供简单易用的API接口。

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

## 项目结构

```
fluvio_grpc_client/
├── client/                 # 客户端API
│   ├── admin.go           # 管理功能
│   ├── consumer.go        # 消费者
│   ├── producer.go        # 生产者
│   └── topic.go           # 主题管理
├── config/                 # 配置管理
│   ├── config.go          # 配置定义
│   └── load.go            # 配置加载
├── errors/                 # 错误定义
│   └── errors.go          # 错误类型
├── examples/               # 使用示例
│   ├── basic/             # 基本示例
│   ├── advanced/          # 高级示例
│   └── integration/       # 集成测试
├── internal/               # 内部实现
│   ├── grpc/              # gRPC连接管理
│   ├── pool/              # 连接池
│   └── retry/             # 重试机制
├── logger/                 # 日志系统
│   └── logger.go          # 日志实现
├── proto/                  # 协议定义
│   └── fluvio_service/    # 生成的协议代码
├── types/                  # 类型定义
│   ├── admin.go           # 管理类型
│   ├── consumer.go        # 消费者类型
│   ├── producer.go        # 生产者类型
│   └── topic.go           # 主题类型
├── fluvio.go              # SDK主入口
├── go.mod                 # Go 模块定义
├── go.sum                 # 依赖校验和
└── README.md              # 项目说明文档
```

## 安装与使用

### 前置条件

- Go 1.18 或更高版本
- 正在运行的 Fluvio 服务实例

### 安装

```bash
go get github.com/iwen-conf/fluvio_grpc_client
```

### 基本使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 创建客户端
    client, err := fluvio.New(
        fluvio.WithServer("101.43.173.154", 50051),
        fluvio.WithTimeout(5*time.Second, 10*time.Second),
        fluvio.WithLogLevel(fluvio.LevelInfo),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 健康检查
    ctx := context.Background()
    err = client.HealthCheck(ctx)
    if err != nil {
        log.Fatal("健康检查失败:", err)
    }
    fmt.Println("连接成功!")

    // 生产消息（支持自定义消息ID）
    result, err := client.Producer().Produce(ctx, "Hello, Fluvio!", fluvio.ProduceOptions{
        Topic:     "my-topic",
        Key:       "key1",
        MessageID: "msg-001", // 🆕 自定义消息ID
        Headers: map[string]string{
            "source": "go-sdk",
            "type":   "greeting",
        },
    })
    if err != nil {
        log.Fatal("生产消息失败:", err)
    }
    fmt.Printf("消息发送成功! ID: %s\n", result.MessageID)

    // 消费消息
    messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
        Topic:       "my-topic",
        Group:       "my-group",
        MaxMessages: 10,
    })
    if err != nil {
        log.Fatal("消费消息失败:", err)
    }
    fmt.Printf("收到 %d 条消息\n", len(messages))
    for _, msg := range messages {
        fmt.Printf("消息: [%s] %s (ID: %s)\n", msg.Key, msg.Value, msg.MessageID)
        if len(msg.Headers) > 0 {
            fmt.Printf("  头部: %v\n", msg.Headers)
        }
    }
}
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

## API 文档

### 客户端创建

```go
// 使用默认配置
client, err := fluvio.New()

// 使用自定义配置
client, err := fluvio.New(
    fluvio.WithServer("101.43.173.154", 50051),
    fluvio.WithTimeout(5*time.Second, 10*time.Second),
    fluvio.WithLogLevel(fluvio.LevelInfo),
    fluvio.WithMaxRetries(3),
    fluvio.WithPoolSize(5),
)

// 使用配置文件
cfg, err := fluvio.LoadConfigFromFile("config.json")
client, err := fluvio.NewWithConfig(cfg)

// 快速连接
client, err := fluvio.QuickStart("101.43.173.154", 50051)
```

### 消息生产

```go
// 基本生产（支持消息ID和头部）
result, err := client.Producer().Produce(ctx, "Hello World", fluvio.ProduceOptions{
    Topic:     "my-topic",
    Key:       "key1",
    MessageID: "msg-001", // 🆕 自定义消息ID
    Headers: map[string]string{
        "source": "go-sdk",
        "type":   "greeting",
    },
})

// 批量生产
messages := []fluvio.Message{
    {Topic: "my-topic", Key: "key1", Value: "message1", MessageID: "batch-001"},
    {Topic: "my-topic", Key: "key2", Value: "message2", MessageID: "batch-002"},
}
batchResult, err := client.Producer().ProduceBatch(ctx, messages)

// 异步生产
resultChan := client.Producer().ProduceAsync(ctx, "Async message", fluvio.ProduceOptions{
    Topic: "my-topic",
})
result := <-resultChan
```

### 消息消费

```go
// 基本消费
messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
    Topic:       "my-topic",
    Group:       "my-group",
    MaxMessages: 10,
})

// 🆕 过滤消费
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

// 流式消费（增强功能）
stream, err := client.Consumer().ConsumeStream(ctx, fluvio.StreamConsumeOptions{
    Topic:        "my-topic",
    Group:        "my-group",
    MaxBatchSize: 10,   // 🆕 批次大小控制
    MaxWaitMs:    1000, // 🆕 等待时间控制
})

for msg := range stream {
    if msg.Error != nil {
        log.Printf("Error: %v", msg.Error)
        continue
    }
    fmt.Printf("Received: [%s] %s (ID: %s)\n",
        msg.Message.Key, msg.Message.Value, msg.Message.MessageID)
}

// 提交偏移量
err = client.Consumer().CommitOffset(ctx, fluvio.CommitOffsetOptions{
    Topic:  "my-topic",
    Group:  "my-group",
    Offset: 100,
})
```

### 主题管理

```go
// 列出主题
topics, err := client.Topic().List(ctx)

// 创建主题（增强配置）
result, err := client.Topic().Create(ctx, fluvio.CreateTopicOptions{
    Name:              "new-topic",
    Partitions:        3,
    ReplicationFactor: 1,                    // 🆕 复制因子
    RetentionMs:       24 * 60 * 60 * 1000, // 🆕 保留时间
    Config: map[string]string{               // 🆕 自定义配置
        "cleanup.policy": "delete",
        "segment.ms":     "3600000",
    },
})

// 🆕 获取主题详细信息
detail, err := client.Topic().DescribeTopicDetail(ctx, "my-topic")

// 🆕 获取主题统计信息
stats, err := client.Topic().GetTopicStats(ctx, fluvio.GetTopicStatsOptions{
    Topic:             "my-topic",
    IncludePartitions: true,
})

// 删除主题
result, err := client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{
    Name: "old-topic",
})

// 检查主题是否存在
exists, err := client.Topic().Exists(ctx, "my-topic")

// 如果不存在则创建
result, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
    Name:       "my-topic",
    Partitions: 1,
})
```

### 管理功能

```go
// 消费组管理
groups, err := client.Admin().ListConsumerGroups(ctx)
groupDetail, err := client.Admin().DescribeConsumerGroup(ctx, "my-group")

// 🆕 SmartModule管理
smartModules, err := client.Admin().ListSmartModules(ctx)
createResult, err := client.Admin().CreateSmartModule(ctx, fluvio.CreateSmartModuleOptions{
    Spec: &fluvio.SmartModuleSpec{
        Name:        "my-filter",
        InputKind:   fluvio.SmartModuleInputStream,
        OutputKind:  fluvio.SmartModuleOutputStream,
        Description: "自定义过滤器",
        Version:     "1.0.0",
    },
    WasmCode: wasmBytes,
})

// 🆕 存储管理
status, err := client.Admin().GetStorageStatus(ctx, fluvio.GetStorageStatusOptions{
    IncludeDetails: true,
})
metrics, err := client.Admin().GetStorageMetrics(ctx, fluvio.GetStorageMetricsOptions{
    IncludeHistory: true,
})

// 🆕 批量删除
bulkResult, err := client.Admin().BulkDelete(ctx, fluvio.BulkDeleteOptions{
    Topics:         []string{"topic1", "topic2"},
    ConsumerGroups: []string{"group1", "group2"},
    SmartModules:   []string{"module1", "module2"},
    Force:          false,
})
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

- 🚀 **[快速入门](QUICKSTART.md)** - 5分钟快速上手指南
- 📋 **[使用方法](HOW_TO_USE.md)** - 详细的导入、创建客户端和配置说明
- 📖 **[完整使用指南](USAGE.md)** - 详细的API文档和使用示例
- 🔧 **[配置示例](examples/config-example.json)** - 配置文件示例
- 💡 **[示例代码](examples/)** - 基本、高级和集成测试示例
