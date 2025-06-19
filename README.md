# Fluvio Go SDK v2.0

一个现代化、生产就绪的 Go SDK，用于与 Fluvio 流处理平台交互。基于 Clean Architecture 设计，提供类型安全的 API 和强大的错误处理机制。

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](#)

## ✨ 特性

- 🚀 **现代化 API 设计** - 链式调用，类型安全
- 🏗️ **Clean Architecture** - 清晰的分层架构，易于测试和维护
- 🔄 **完整消息处理** - 支持生产、消费、流式处理
- 📊 **实时流处理** - 高性能流式消费，支持背压控制
- 🛡️ **类型安全** - 强类型接口，编译时错误检查
- 🔧 **智能重试机制** - 指数退避，可配置重试策略
- 📝 **完整错误处理** - 统一错误类型，详细错误信息
- 🔐 **安全连接** - 支持 TLS/SSL 和不安全连接
- 📈 **生产就绪** - 连接池、健康检查、监控支持

## 📦 安装

```bash
go get github.com/iwen-conf/fluvio_grpc_client
```

**系统要求:**
- Go 1.18 或更高版本
- Fluvio 服务器 0.9.0+

## 🚀 快速开始

### 基本示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    fluvio "github.com/iwen-conf/fluvio_grpc_client"
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
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 连接到服务器
    if err := client.Connect(ctx); err != nil {
        log.Fatal("连接失败:", err)
    }
    fmt.Println("✅ 连接成功")

    // 创建主题
    err = client.Topics().Create(ctx, "my-topic", &fluvio.CreateTopicOptions{
        Partitions: 1,
    })
    if err != nil {
        log.Printf("创建主题失败: %v", err)
    }

    // 发送消息
    result, err := client.Producer().SendString(ctx, "my-topic", "key1", "Hello, Fluvio!")
    if err != nil {
        log.Fatal("发送消息失败:", err)
    }
    fmt.Printf("✅ 消息发送成功: %s\n", result.MessageID)

    // 接收消息
    messages, err := client.Consumer().Receive(ctx, "my-topic", &fluvio.ReceiveOptions{
        Group:       "my-group",
        MaxMessages: 10,
    })
    if err != nil {
        log.Fatal("接收消息失败:", err)
    }

    fmt.Printf("✅ 接收到 %d 条消息:\n", len(messages))
    for i, msg := range messages {
        fmt.Printf("  %d. [%s] %s\n", i+1, msg.Key, string(msg.Value))
    }
}
```

### 流式消费示例

```go
func streamExample() {
    client, _ := fluvio.NewClient(
        fluvio.WithAddress("localhost", 50051),
    )
    defer client.Close()

    ctx := context.Background()
    client.Connect(ctx)

    // 启动流式消费
    stream, err := client.Consumer().Stream(ctx, "events", &fluvio.StreamOptions{
        Group:      "stream-group",
        BufferSize: 1000,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("🔄 开始流式消费...")
    for msg := range stream {
        fmt.Printf("📨 收到消息: [%s] %s\n", msg.Key, string(msg.Value))
        
        // 处理消息...
        
        // 提交偏移量
        client.Consumer().Commit(ctx, "events", "stream-group", msg.Offset)
    }
}
```

## ⚙️ 配置选项

### 完整配置示例

```go
client, err := fluvio.NewClient(
    // 🌐 服务器连接
    fluvio.WithAddress("fluvio.example.com", 50051),
    
    // ⏱️ 超时设置
    fluvio.WithTimeout(30*time.Second),
    
    // 🔄 重试策略
    fluvio.WithRetry(3, time.Second),
    
    // 📝 日志配置
    fluvio.WithLogLevel(fluvio.LogLevelInfo),
    
    // 🏊 连接池
    fluvio.WithConnectionPool(10, 10*time.Minute),
    
    // 💓 Keep-Alive
    fluvio.WithKeepAlive(30*time.Second),
    
    // 🔐 TLS 安全连接
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
    
    // 或者使用不安全连接（仅用于开发）
    // fluvio.WithInsecure(),
)
```

### 配置选项说明

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithAddress(host, port)` | 服务器地址和端口 | localhost:50051 |
| `WithTimeout(duration)` | 操作超时时间 | 30s |
| `WithRetry(attempts, delay)` | 重试次数和延迟 | 3次, 1s |
| `WithLogLevel(level)` | 日志级别 | Info |
| `WithConnectionPool(size, ttl)` | 连接池大小和TTL | 5, 5min |
| `WithKeepAlive(interval)` | Keep-Alive间隔 | 30s |
| `WithTLS(cert, key, ca)` | TLS证书配置 | - |
| `WithInsecure()` | 禁用TLS（不推荐生产环境） | false |

## 📖 主要功能

### 🔄 消息生产

```go
producer := client.Producer()

// 1. 发送字符串消息（最简单）
result, err := producer.SendString(ctx, "topic", "key", "Hello World")

// 2. 发送结构化消息
message := &fluvio.Message{
    Key:   "user-123",
    Value: []byte(`{"name": "Alice", "age": 30}`),
    Headers: map[string]string{
        "content-type": "application/json",
        "source":       "user-service",
    },
}
result, err := producer.Send(ctx, "user-events", message)

// 3. 批量发送（高性能）
var messages []*fluvio.Message
for i := 0; i < 1000; i++ {
    messages = append(messages, &fluvio.Message{
        Key:   fmt.Sprintf("batch-%d", i),
        Value: []byte(fmt.Sprintf("message-%d", i)),
    })
}
batchResult, err := producer.SendBatch(ctx, "batch-topic", messages)
fmt.Printf("批量发送: 成功 %d, 失败 %d\n", 
    batchResult.SuccessCount, batchResult.FailureCount)
```

### 📥 消息消费

```go
consumer := client.Consumer()

// 1. 批量接收
messages, err := consumer.Receive(ctx, "topic", &fluvio.ReceiveOptions{
    Group:       "my-group",
    MaxMessages: 100,
    Offset:      0, // 从头开始，-1 表示从最新开始
})

// 2. 接收单条消息
message, err := consumer.ReceiveOne(ctx, "topic", "my-group")
if message != nil {
    fmt.Printf("收到: %s\n", string(message.Value))
}

// 3. 流式消费（推荐用于实时处理）
stream, err := consumer.Stream(ctx, "events", &fluvio.StreamOptions{
    Group:      "stream-processor",
    BufferSize: 1000, // 缓冲区大小，支持背压控制
    Offset:     -1,   // 从最新消息开始
})

go func() {
    for msg := range stream {
        // 处理消息
        processMessage(msg)
        
        // 手动提交偏移量
        consumer.Commit(ctx, "events", "stream-processor", msg.Offset)
    }
}()

// 4. 便捷方法：接收字符串
values, err := consumer.ReceiveString(ctx, "text-topic", &fluvio.ReceiveOptions{
    Group: "text-processor",
})
```

### 🗂️ 主题管理

```go
topics := client.Topics()

// 创建主题
err := topics.Create(ctx, "new-topic", &fluvio.CreateTopicOptions{
    Partitions:        3,
    ReplicationFactor: 1,
    Config: map[string]string{
        "retention.ms": "86400000", // 1天
    },
})

// 列出所有主题
topicList, err := topics.List(ctx)
fmt.Printf("共有 %d 个主题\n", len(topicList))

// 获取主题详细信息
info, err := topics.Info(ctx, "my-topic")
fmt.Printf("主题 %s: %d 个分区\n", info.Name, info.Partitions)

// 检查主题是否存在
exists, err := topics.Exists(ctx, "my-topic")

// 便捷方法：创建主题（如果不存在）
created, err := topics.CreateIfNotExists(ctx, "auto-topic", &fluvio.CreateTopicOptions{
    Partitions: 1,
})

// 删除主题
err = topics.Delete(ctx, "old-topic")
```

### 🛠️ 集群管理

```go
admin := client.Admin()

// 获取集群信息
clusterInfo, err := admin.ClusterInfo(ctx)
fmt.Printf("集群状态: %s, 控制器: %d\n", 
    clusterInfo.Status, clusterInfo.ControllerID)

// 获取 Broker 列表
brokers, err := admin.Brokers(ctx)
for _, broker := range brokers {
    fmt.Printf("Broker %d: %s:%d (%s)\n", 
        broker.ID, broker.Host, broker.Port, broker.Status)
}

// 消费者组管理
groups, err := admin.ConsumerGroups(ctx)
for _, group := range groups {
    fmt.Printf("消费者组: %s (%s)\n", group.GroupID, group.State)
}

// 获取消费者组详情
groupDetail, err := admin.ConsumerGroupDetail(ctx, "my-group")

// SmartModule 管理
smartModules := admin.SmartModules()
modules, err := smartModules.List(ctx)
```

## 🔧 高级功能

### 错误处理

```go
import "github.com/iwen-conf/fluvio_grpc_client/pkg/errors"

result, err := client.Producer().SendString(ctx, "topic", "key", "value")
if err != nil {
    switch {
    case errors.IsConnectionError(err):
        log.Println("连接错误，检查网络和服务器状态")
    case errors.IsTimeoutError(err):
        log.Println("操作超时，可能需要增加超时时间")
    case errors.IsValidationError(err):
        log.Println("参数验证失败，检查输入参数")
    case errors.IsAuthenticationError(err):
        log.Println("认证失败，检查证书和权限")
    default:
        log.Printf("其他错误: %v", err)
    }
}
```

### 健康检查和监控

```go
// 健康检查
err := client.HealthCheck(ctx)
if err != nil {
    log.Printf("健康检查失败: %v", err)
}

// Ping 测试
duration, err := client.Ping(ctx)
if err == nil {
    fmt.Printf("Ping 延迟: %v\n", duration)
}

// 检查连接状态
if client.IsConnected() {
    fmt.Println("客户端已连接")
}
```

### 自定义日志

```go
// 使用自定义日志器
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithLogLevel(fluvio.LogLevelDebug),
)

// 获取内置日志器
logger := client.Logger()
logger.Info("自定义日志消息")
```

## 🎯 最佳实践

### 1. 连接管理

```go
// ✅ 推荐：使用连接池
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithConnectionPool(10, 10*time.Minute),
    fluvio.WithKeepAlive(30*time.Second),
)

// ✅ 总是关闭客户端
defer client.Close()

// ✅ 检查连接状态
if !client.IsConnected() {
    if err := client.Connect(ctx); err != nil {
        return err
    }
}
```

### 2. 错误处理和重试

```go
// ✅ 配置合适的重试策略
client, err := fluvio.NewClient(
    fluvio.WithRetry(3, time.Second),
    fluvio.WithTimeout(30*time.Second),
)

// ✅ 使用上下文控制超时
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

### 3. 性能优化

```go
// ✅ 批量操作提高吞吐量
var messages []*fluvio.Message
for i := 0; i < 1000; i++ {
    messages = append(messages, &fluvio.Message{
        Key:   fmt.Sprintf("key-%d", i),
        Value: []byte(fmt.Sprintf("data-%d", i)),
    })
}
result, err := client.Producer().SendBatch(ctx, "topic", messages)

// ✅ 流式消费处理大量数据
stream, err := client.Consumer().Stream(ctx, "topic", &fluvio.StreamOptions{
    Group:      "processor",
    BufferSize: 1000, // 适当的缓冲区大小
})
```

### 4. 生产环境配置

```go
// ✅ 生产环境推荐配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("fluvio-cluster.prod.com", 50051),
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
    fluvio.WithTimeout(30*time.Second),
    fluvio.WithRetry(5, 2*time.Second),
    fluvio.WithConnectionPool(20, 30*time.Minute),
    fluvio.WithKeepAlive(60*time.Second),
    fluvio.WithLogLevel(fluvio.LogLevelWarn), // 生产环境减少日志
)
```

## 📚 更多文档

- [API 参考文档](docs/API.md) - 完整的 API 文档
- [使用指南](docs/GUIDE.md) - 详细的使用指南和示例
- [故障排除](docs/TROUBLESHOOTING.md) - 常见问题和解决方案
- [更新日志](CHANGELOG.md) - 版本更新记录

## 🤝 贡献

我们欢迎社区贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解如何参与项目开发。

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/iwen-conf/fluvio_grpc_client.git
cd fluvio_grpc_client

# 安装依赖
go mod download

# 运行测试
go test ./...

# 构建
go build .
```

## 📄 许可证

本项目采用 [MIT 许可证](LICENSE)。

## 🆘 支持

- 📧 邮件支持: support@example.com
- 💬 社区讨论: [GitHub Discussions](https://github.com/iwen-conf/fluvio_grpc_client/discussions)
- 🐛 问题报告: [GitHub Issues](https://github.com/iwen-conf/fluvio_grpc_client/issues)

---

**Fluvio Go SDK v2.0** - 让流处理变得简单而强大 🚀