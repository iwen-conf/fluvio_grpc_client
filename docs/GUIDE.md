# Fluvio Go SDK 使用指南

本指南提供了 Fluvio Go SDK v2.0 的详细使用说明，包括高级功能、最佳实践和完整示例。

## 目录

- [快速入门](#快速入门)
- [客户端配置](#客户端配置)
- [消息生产](#消息生产)
- [消息消费](#消息消费)
- [主题管理](#主题管理)
- [集群管理](#集群管理)
- [高级功能](#高级功能)
- [最佳实践](#最佳实践)
- [性能优化](#性能优化)
- [故障排除](#故障排除)

---

## 快速入门

### 安装和导入

```bash
go get github.com/iwen-conf/fluvio_grpc_client
```

```go
import (
    "context"
    "log"
    "time"
    
    fluvio "github.com/iwen-conf/fluvio_grpc_client"
)
```

### 第一个程序

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
    // 1. 创建客户端
    client, err := fluvio.NewClient(
        fluvio.WithAddress("localhost", 50051),
        fluvio.WithTimeout(30*time.Second),
    )
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()

    // 2. 连接到服务器
    ctx := context.Background()
    if err := client.Connect(ctx); err != nil {
        log.Fatal("连接失败:", err)
    }

    // 3. 发送消息
    result, err := client.Producer().SendString(ctx, "hello-topic", "key1", "Hello, World!")
    if err != nil {
        log.Fatal("发送失败:", err)
    }
    fmt.Printf("消息发送成功: %s\n", result.MessageID)

    // 4. 接收消息
    messages, err := client.Consumer().Receive(ctx, "hello-topic", &fluvio.ReceiveOptions{
        Group: "hello-group",
        MaxMessages: 1,
    })
    if err != nil {
        log.Fatal("接收失败:", err)
    }

    for _, msg := range messages {
        fmt.Printf("收到消息: %s\n", string(msg.Value))
    }
}
```---

## 客户端配置

### 基础配置

最简单的客户端配置：

```go
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
)
```

### 完整配置示例

```go
client, err := fluvio.NewClient(
    // 服务器地址
    fluvio.WithAddress("fluvio.example.com", 50051),
    
    // 超时设置
    fluvio.WithTimeout(30*time.Second),
    
    // 重试策略
    fluvio.WithRetry(3, time.Second),
    
    // 日志级别
    fluvio.WithLogLevel(fluvio.LogLevelInfo),
    
    // 连接池配置
    fluvio.WithConnectionPool(10, 10*time.Minute),
    
    // Keep-Alive 设置
    fluvio.WithKeepAlive(30*time.Second),
    
    // TLS 安全连接
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
)
```

### 开发环境配置

```go
// 开发环境：使用不安全连接，详细日志
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithInsecure(),
    fluvio.WithLogLevel(fluvio.LogLevelDebug),
    fluvio.WithTimeout(10*time.Second),
)
```

### 生产环境配置

```go
// 生产环境：安全连接，优化性能
client, err := fluvio.NewClient(
    fluvio.WithAddress("fluvio-cluster.prod.com", 50051),
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
    fluvio.WithTimeout(30*time.Second),
    fluvio.WithRetry(5, 2*time.Second),
    fluvio.WithConnectionPool(20, 30*time.Minute),
    fluvio.WithKeepAlive(60*time.Second),
    fluvio.WithLogLevel(fluvio.LogLevelWarn),
)
```---

## 消息生产

### 基本消息发送

#### 发送字符串消息

```go
// 最简单的方式
result, err := client.Producer().SendString(ctx, "my-topic", "key1", "Hello World")
if err != nil {
    log.Printf("发送失败: %v", err)
    return
}
fmt.Printf("消息发送成功: %s\n", result.MessageID)
```

#### 发送二进制消息

```go
message := &fluvio.Message{
    Key:   "binary-key",
    Value: []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}, // "Hello" in bytes
    Headers: map[string]string{
        "content-type": "application/octet-stream",
        "encoding":     "binary",
    },
}

result, err := client.Producer().Send(ctx, "binary-topic", message)
if err != nil {
    log.Printf("发送失败: %v", err)
    return
}
```

#### 发送 JSON 消息

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{ID: 1, Name: "Alice", Age: 30}
result, err := client.Producer().SendJSON(ctx, "user-events", "user-1", user)
if err != nil {
    log.Printf("发送 JSON 失败: %v", err)
    return
}
```### 批量消息发送

批量发送可以显著提高吞吐量：

```go
// 准备批量消息
var messages []*fluvio.Message
for i := 0; i < 1000; i++ {
    message := &fluvio.Message{
        Key:   fmt.Sprintf("batch-key-%d", i),
        Value: []byte(fmt.Sprintf("batch message %d", i)),
        Headers: map[string]string{
            "batch_id": "batch-001",
            "index":    fmt.Sprintf("%d", i),
        },
    }
    messages = append(messages, message)
}

// 批量发送
result, err := client.Producer().SendBatch(ctx, "batch-topic", messages)
if err != nil {
    log.Printf("批量发送失败: %v", err)
    return
}

fmt.Printf("批量发送完成: 成功 %d, 失败 %d\n", 
    result.SuccessCount, result.FailureCount)

// 检查失败的消息
for i, res := range result.Results {
    if res == nil {
        fmt.Printf("消息 %d 发送失败\n", i)
    }
}
```

### 高级生产者配置

```go
// 创建带有自定义配置的生产者
producer := client.Producer()

// 发送带有时间戳的消息
message := &fluvio.Message{
    Key:       "timestamped-key",
    Value:     []byte("timestamped message"),
    Timestamp: time.Now(),
    Headers: map[string]string{
        "source":    "my-service",
        "version":   "1.0.0",
        "timestamp": time.Now().Format(time.RFC3339),
    },
}

result, err := producer.Send(ctx, "timestamped-topic", message)
```---

## 消息消费

### 基本消息消费

#### 批量接收消息

```go
// 接收最多 100 条消息
messages, err := client.Consumer().Receive(ctx, "my-topic", &fluvio.ReceiveOptions{
    Group:       "my-consumer-group",
    MaxMessages: 100,
    Offset:      0, // 从头开始，-1 表示从最新开始
})
if err != nil {
    log.Printf("接收失败: %v", err)
    return
}

fmt.Printf("接收到 %d 条消息\n", len(messages))
for i, msg := range messages {
    fmt.Printf("消息 %d: [%s] %s (偏移量: %d)\n", 
        i+1, msg.Key, string(msg.Value), msg.Offset)
}
```

#### 接收单条消息

```go
// 接收单条消息（便捷方法）
message, err := client.Consumer().ReceiveOne(ctx, "my-topic", "single-consumer")
if err != nil {
    log.Printf("接收失败: %v", err)
    return
}

if message != nil {
    fmt.Printf("收到消息: [%s] %s\n", message.Key, string(message.Value))
    
    // 手动提交偏移量
    err = client.Consumer().Commit(ctx, "my-topic", "single-consumer", message.Offset)
    if err != nil {
        log.Printf("提交偏移量失败: %v", err)
    }
} else {
    fmt.Println("没有可用的消息")
}
```### 流式消费

流式消费适合实时处理大量消息：

```go
// 启动流式消费
stream, err := client.Consumer().Stream(ctx, "events", &fluvio.StreamOptions{
    Group:      "stream-processor",
    BufferSize: 1000, // 缓冲区大小，支持背压控制
    Offset:     -1,   // 从最新消息开始
})
if err != nil {
    log.Printf("启动流式消费失败: %v", err)
    return
}

fmt.Println("开始流式消费...")

// 处理消息流
go func() {
    for msg := range stream {
        // 处理消息
        fmt.Printf("处理消息: [%s] %s\n", msg.Key, string(msg.Value))
        
        // 模拟处理时间
        time.Sleep(100 * time.Millisecond)
        
        // 提交偏移量
        if err := client.Consumer().Commit(ctx, "events", "stream-processor", msg.Offset); err != nil {
            log.Printf("提交偏移量失败: %v", err)
        }
    }
    fmt.Println("流式消费结束")
}()

// 等待一段时间或直到上下文取消
select {
case <-ctx.Done():
    fmt.Println("上下文取消，停止消费")
case <-time.After(30 * time.Second):
    fmt.Println("消费时间到，停止消费")
}
```

### 消费者组管理

```go
// 使用不同的消费者组并行处理
groups := []string{"processor-1", "processor-2", "processor-3"}

for _, group := range groups {
    go func(groupID string) {
        messages, err := client.Consumer().Receive(ctx, "parallel-topic", &fluvio.ReceiveOptions{
            Group:       groupID,
            MaxMessages: 50,
        })
        if err != nil {
            log.Printf("组 %s 接收失败: %v", groupID, err)
            return
        }
        
        fmt.Printf("组 %s 处理 %d 条消息\n", groupID, len(messages))
        // 处理消息...
    }(group)
}
```---

## 主题管理

### 创建和删除主题

```go
topics := client.Topics()

// 创建主题
err := topics.Create(ctx, "new-topic", &fluvio.CreateTopicOptions{
    Partitions:        3,
    ReplicationFactor: 1,
    Config: map[string]string{
        "retention.ms":     "86400000", // 1天保留期
        "cleanup.policy":   "delete",
        "compression.type": "gzip",
    },
})
if err != nil {
    log.Printf("创建主题失败: %v", err)
} else {
    fmt.Println("主题创建成功")
}

// 检查主题是否存在
exists, err := topics.Exists(ctx, "new-topic")
if err != nil {
    log.Printf("检查主题失败: %v", err)
} else {
    fmt.Printf("主题存在: %v\n", exists)
}

// 删除主题
err = topics.Delete(ctx, "old-topic")
if err != nil {
    log.Printf("删除主题失败: %v", err)
} else {
    fmt.Println("主题删除成功")
}
```

### 主题信息查询

```go
// 列出所有主题
topicList, err := topics.List(ctx)
if err != nil {
    log.Printf("列出主题失败: %v", err)
} else {
    fmt.Printf("共有 %d 个主题:\n", len(topicList))
    for i, topic := range topicList {
        fmt.Printf("  %d. %s\n", i+1, topic)
    }
}

// 获取主题详细信息
info, err := topics.Info(ctx, "my-topic")
if err != nil {
    log.Printf("获取主题信息失败: %v", err)
} else {
    fmt.Printf("主题信息:\n")
    fmt.Printf("  名称: %s\n", info.Name)
    fmt.Printf("  分区数: %d\n", info.Partitions)
    fmt.Printf("  副本因子: %d\n", info.ReplicationFactor)
    fmt.Printf("  配置:\n")
    for key, value := range info.Config {
        fmt.Printf("    %s: %s\n", key, value)
    }
}
```### 便捷方法

```go
// 创建主题（如果不存在）
created, err := topics.CreateIfNotExists(ctx, "auto-topic", &fluvio.CreateTopicOptions{
    Partitions: 1,
})
if err != nil {
    log.Printf("创建主题失败: %v", err)
} else if created {
    fmt.Println("主题已创建")
} else {
    fmt.Println("主题已存在")
}

// 批量创建主题
topicsToCreate := map[string]*fluvio.CreateTopicOptions{
    "events":    {Partitions: 3},
    "logs":      {Partitions: 1},
    "metrics":   {Partitions: 5},
}

for name, opts := range topicsToCreate {
    created, err := topics.CreateIfNotExists(ctx, name, opts)
    if err != nil {
        log.Printf("创建主题 %s 失败: %v", name, err)
    } else if created {
        fmt.Printf("主题 %s 已创建\n", name)
    } else {
        fmt.Printf("主题 %s 已存在\n", name)
    }
}
```---

## 集群管理

### 集群状态监控

```go
admin := client.Admin()

// 获取集群信息
clusterInfo, err := admin.ClusterInfo(ctx)
if err != nil {
    log.Printf("获取集群信息失败: %v", err)
} else {
    fmt.Printf("集群信息:\n")
    fmt.Printf("  ID: %s\n", clusterInfo.ID)
    fmt.Printf("  状态: %s\n", clusterInfo.Status)
    fmt.Printf("  控制器ID: %d\n", clusterInfo.ControllerID)
}

// 获取 Broker 列表
brokers, err := admin.Brokers(ctx)
if err != nil {
    log.Printf("获取 Broker 列表失败: %v", err)
} else {
    fmt.Printf("Broker 列表 (%d 个):\n", len(brokers))
    for _, broker := range brokers {
        fmt.Printf("  Broker %d: %s:%d (%s)\n", 
            broker.ID, broker.Host, broker.Port, broker.Status)
    }
}
```

### 消费者组管理

```go
// 获取所有消费者组
groups, err := admin.ConsumerGroups(ctx)
if err != nil {
    log.Printf("获取消费者组失败: %v", err)
} else {
    fmt.Printf("消费者组列表 (%d 个):\n", len(groups))
    for _, group := range groups {
        fmt.Printf("  组: %s (状态: %s)\n", group.GroupID, group.State)
    }
}

// 获取特定消费者组的详细信息
groupDetail, err := admin.ConsumerGroupDetail(ctx, "my-group")
if err != nil {
    log.Printf("获取消费者组详情失败: %v", err)
} else {
    fmt.Printf("消费者组详情:\n")
    fmt.Printf("  组ID: %s\n", groupDetail.GroupID)
    fmt.Printf("  状态: %s\n", groupDetail.State)
    fmt.Printf("  成员数: %d\n", len(groupDetail.Members))
    
    for i, member := range groupDetail.Members {
        fmt.Printf("  成员 %d:\n", i+1)
        fmt.Printf("    成员ID: %s\n", member.MemberID)
        fmt.Printf("    客户端ID: %s\n", member.ClientID)
        fmt.Printf("    客户端主机: %s\n", member.ClientHost)
    }
}
```---

## 高级功能

### 健康检查和监控

```go
// 健康检查
err := client.HealthCheck(ctx)
if err != nil {
    log.Printf("健康检查失败: %v", err)
    // 可能需要重新连接
    if err := client.Connect(ctx); err != nil {
        log.Printf("重新连接失败: %v", err)
    }
} else {
    fmt.Println("服务器健康状态正常")
}

// Ping 测试
duration, err := client.Ping(ctx)
if err != nil {
    log.Printf("Ping 失败: %v", err)
} else {
    fmt.Printf("Ping 延迟: %v\n", duration)
}

// 检查连接状态
if client.IsConnected() {
    fmt.Println("客户端已连接")
} else {
    fmt.Println("客户端未连接")
    // 尝试重新连接
    if err := client.Connect(ctx); err != nil {
        log.Printf("重新连接失败: %v", err)
    }
}
```

### 错误处理策略

```go
import "github.com/iwen-conf/fluvio_grpc_client/pkg/errors"

func handleError(err error) {
    switch {
    case errors.IsConnectionError(err):
        log.Println("连接错误 - 检查网络和服务器状态")
        // 实施重连逻辑
        
    case errors.IsTimeoutError(err):
        log.Println("超时错误 - 考虑增加超时时间")
        // 可能需要重试
        
    case errors.IsValidationError(err):
        log.Println("验证错误 - 检查输入参数")
        // 修正参数后重试
        
    case errors.IsNotFoundError(err):
        log.Println("资源未找到 - 检查主题或消费者组是否存在")
        // 可能需要创建资源
        
    case errors.IsAlreadyExistsError(err):
        log.Println("资源已存在 - 可以继续使用现有资源")
        // 通常可以忽略此错误
        
    case errors.IsAuthenticationError(err):
        log.Println("认证错误 - 检查证书和权限")
        // 检查 TLS 配置
        
    default:
        log.Printf("未知错误: %v", err)
        // 通用错误处理
    }
}
```---

## 最佳实践

### 1. 连接管理

```go
// ✅ 推荐：使用连接池和 Keep-Alive
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
        return fmt.Errorf("连接失败: %w", err)
    }
}

// ❌ 避免：频繁创建和销毁客户端
// 应该复用客户端实例
```

### 2. 上下文管理

```go
// ✅ 推荐：为每个操作设置合适的超时
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// ✅ 长时间运行的操作使用可取消的上下文
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// 在另一个 goroutine 中监听取消信号
go func() {
    <-sigChan // 等待信号
    cancel()  // 取消操作
}()

// ❌ 避免：使用 context.Background() 进行长时间操作
// 这会导致无法取消操作
```

### 3. 错误处理

```go
// ✅ 推荐：使用类型化错误检查
result, err := client.Producer().SendString(ctx, "topic", "key", "value")
if err != nil {
    if errors.IsConnectionError(err) {
        // 特定的连接错误处理
        return handleConnectionError(err)
    }
    return fmt.Errorf("发送失败: %w", err)
}

// ✅ 实施重试逻辑
func sendWithRetry(client *fluvio.Client, topic, key, value string) error {
    for attempt := 0; attempt < 3; attempt++ {
        _, err := client.Producer().SendString(ctx, topic, key, value)
        if err == nil {
            return nil
        }
        
        if !errors.IsConnectionError(err) && !errors.IsTimeoutError(err) {
            return err // 不可重试的错误
        }
        
        time.Sleep(time.Duration(attempt+1) * time.Second)
    }
    return fmt.Errorf("重试 3 次后仍然失败")
}
```---

## 性能优化

### 1. 批量操作

```go
// ✅ 高吞吐量：使用批量发送
var messages []*fluvio.Message
for i := 0; i < 1000; i++ {
    messages = append(messages, &fluvio.Message{
        Key:   fmt.Sprintf("key-%d", i),
        Value: []byte(fmt.Sprintf("data-%d", i)),
    })
}

// 批量发送比单条发送快 10-100 倍
result, err := client.Producer().SendBatch(ctx, "topic", messages)

// ❌ 低效：逐条发送
for i := 0; i < 1000; i++ {
    client.Producer().SendString(ctx, "topic", fmt.Sprintf("key-%d", i), fmt.Sprintf("data-%d", i))
}
```

### 2. 流式消费优化

```go
// ✅ 优化缓冲区大小
stream, err := client.Consumer().Stream(ctx, "topic", &fluvio.StreamOptions{
    Group:      "processor",
    BufferSize: 1000, // 根据处理能力调整
})

// ✅ 并行处理消息
const workerCount = 10
messageChan := make(chan *fluvio.ConsumedMessage, 100)

// 启动工作协程
for i := 0; i < workerCount; i++ {
    go func(workerID int) {
        for msg := range messageChan {
            processMessage(msg)
        }
    }(i)
}

// 分发消息
for msg := range stream {
    select {
    case messageChan <- msg:
    case <-ctx.Done():
        return
    }
}
```

### 3. 连接池优化

```go
// ✅ 根据负载调整连接池
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithConnectionPool(20, 30*time.Minute), // 高负载环境
    fluvio.WithKeepAlive(60*time.Second),
)

// ✅ 监控连接池使用情况
func monitorConnections(client *fluvio.Client) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // 这里可以添加连接池监控逻辑
        if !client.IsConnected() {
            log.Println("警告：客户端连接断开")
        }
    }
}
```---

## 故障排除

### 常见问题和解决方案

#### 1. 连接问题

**问题**: `connection refused` 或 `timeout` 错误

```go
// 解决方案：检查服务器状态和网络连接
func diagnoseConnection(client *fluvio.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // 测试连接
    if err := client.HealthCheck(ctx); err != nil {
        log.Printf("健康检查失败: %v", err)
        
        // 尝试 Ping
        if duration, pingErr := client.Ping(ctx); pingErr != nil {
            log.Printf("Ping 失败: %v", pingErr)
            log.Println("建议：检查服务器是否运行，网络是否可达")
        } else {
            log.Printf("Ping 成功，延迟: %v", duration)
        }
    }
}
```

#### 2. 认证问题

**问题**: `authentication failed` 或 TLS 错误

```go
// 解决方案：检查 TLS 配置
func checkTLSConfig() {
    // 验证证书文件是否存在
    certFiles := []string{"client.crt", "client.key", "ca.crt"}
    for _, file := range certFiles {
        if _, err := os.Stat(file); os.IsNotExist(err) {
            log.Printf("证书文件不存在: %s", file)
        }
    }
    
    // 使用正确的 TLS 配置
    client, err := fluvio.NewClient(
        fluvio.WithAddress("localhost", 50051),
        fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
    )
    if err != nil {
        log.Printf("TLS 配置错误: %v", err)
    }
}
```

#### 3. 性能问题

**问题**: 消息发送或接收缓慢

```go
// 解决方案：性能调优
func optimizePerformance() {
    client, err := fluvio.NewClient(
        fluvio.WithAddress("localhost", 50051),
        // 增加连接池大小
        fluvio.WithConnectionPool(20, 30*time.Minute),
        // 调整 Keep-Alive
        fluvio.WithKeepAlive(60*time.Second),
        // 增加超时时间
        fluvio.WithTimeout(60*time.Second),
    )
    
    // 使用批量操作
    // 调整流式消费缓冲区大小
    // 实施并行处理
}
```

### 调试技巧

#### 启用详细日志

```go
// 开发环境：启用调试日志
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithLogLevel(fluvio.LogLevelDebug),
)

// 自定义日志处理
logger := client.Logger()
logger.Debug("调试信息", logging.Field{Key: "key", Value: "value"})
```

#### 监控和指标

```go
// 实施基本监控
func monitorClient(client *fluvio.Client) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if client.IsConnected() {
            log.Println("✅ 客户端连接正常")
        } else {
            log.Println("❌ 客户端连接断开")
        }
        
        // 测试延迟
        if duration, err := client.Ping(context.Background()); err == nil {
            log.Printf("📊 延迟: %v", duration)
        }
    }
}
```

---

## 完整示例

### 生产者-消费者示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
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

    ctx := context.Background()
    
    // 连接
    if err := client.Connect(ctx); err != nil {
        log.Fatal("连接失败:", err)
    }

    // 创建主题
    _, err = client.Topics().CreateIfNotExists(ctx, "demo-topic", &fluvio.CreateTopicOptions{
        Partitions: 1,
    })
    if err != nil {
        log.Printf("创建主题失败: %v", err)
    }

    var wg sync.WaitGroup

    // 启动生产者
    wg.Add(1)
    go func() {
        defer wg.Done()
        producer(client)
    }()

    // 启动消费者
    wg.Add(1)
    go func() {
        defer wg.Done()
        consumer(client)
    }()

    wg.Wait()
}

func producer(client *fluvio.Client) {
    ctx := context.Background()
    
    for i := 0; i < 10; i++ {
        message := fmt.Sprintf("消息 %d", i)
        result, err := client.Producer().SendString(ctx, "demo-topic", fmt.Sprintf("key-%d", i), message)
        if err != nil {
            log.Printf("发送失败: %v", err)
            continue
        }
        fmt.Printf("✅ 发送成功: %s\n", result.MessageID)
        time.Sleep(time.Second)
    }
}

func consumer(client *fluvio.Client) {
    ctx := context.Background()
    
    stream, err := client.Consumer().Stream(ctx, "demo-topic", &fluvio.StreamOptions{
        Group:      "demo-group",
        BufferSize: 100,
    })
    if err != nil {
        log.Printf("启动消费失败: %v", err)
        return
    }

    for msg := range stream {
        fmt.Printf("📨 收到: [%s] %s\n", msg.Key, string(msg.Value))
        
        // 提交偏移量
        client.Consumer().Commit(ctx, "demo-topic", "demo-group", msg.Offset)
    }
}
```

这个完整的使用指南涵盖了 Fluvio Go SDK 的所有主要功能和最佳实践。每个部分都提供了实用的代码示例和详细的说明。