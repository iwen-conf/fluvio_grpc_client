# Fluvio Go SDK 快速入门

5分钟快速上手 Fluvio Go SDK！

## 🚀 第一步：安装

```bash
# 创建新项目
mkdir my-fluvio-app
cd my-fluvio-app
go mod init my-fluvio-app

# 安装SDK
go get github.com/iwen-conf/fluvio_grpc_client
```

## 📝 第二步：创建第一个应用

创建 `main.go` 文件：

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
    // 1. 创建客户端
    client, err := fluvio.New(
        fluvio.WithServer("101.43.173.154", 50051),
        fluvio.WithTimeout(5*time.Second, 10*time.Second),
    )
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()

    ctx := context.Background()

    // 2. 健康检查
    fmt.Println("🔍 检查连接...")
    err = client.HealthCheck(ctx)
    if err != nil {
        log.Fatal("连接失败:", err)
    }
    fmt.Println("✅ 连接成功!")

    // 3. 创建主题
    topicName := "quickstart-topic"
    fmt.Printf("📁 创建主题 '%s'...\n", topicName)
    _, err = client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
        Name:       topicName,
        Partitions: 1,
    })
    if err != nil {
        log.Fatal("创建主题失败:", err)
    }
    fmt.Println("✅ 主题已就绪!")

    // 4. 发送消息
    fmt.Println("📤 发送消息...")
    result, err := client.Producer().Produce(ctx, "Hello, Fluvio!", fluvio.ProduceOptions{
        Topic: topicName,
        Key:   "greeting",
    })
    if err != nil {
        log.Fatal("发送消息失败:", err)
    }
    fmt.Printf("✅ 消息发送成功! ID: %s\n", result.MessageID)

    // 5. 接收消息
    fmt.Println("📥 接收消息...")
    messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
        Topic:       topicName,
        Group:       "quickstart-group",
        MaxMessages: 1,
        Offset:      0,
    })
    if err != nil {
        log.Fatal("接收消息失败:", err)
    }

    if len(messages) > 0 {
        msg := messages[0]
        fmt.Printf("✅ 收到消息: [%s] %s\n", msg.Key, msg.Value)
    } else {
        fmt.Println("⚠️  没有收到消息")
    }

    fmt.Println("🎉 快速入门完成!")
}
```

## 🏃 第三步：运行应用

```bash
go run main.go
```

预期输出：
```
🔍 检查连接...
✅ 连接成功!
📁 创建主题 'quickstart-topic'...
✅ 主题已就绪!
📤 发送消息...
✅ 消息发送成功! ID: batch-0
📥 接收消息...
✅ 收到消息: [greeting] Hello, Fluvio!
🎉 快速入门完成!
```

## 🎯 下一步

### 1. 批量处理

```go
// 批量发送消息
messages := []fluvio.Message{
    {Topic: "my-topic", Key: "key1", Value: "消息1"},
    {Topic: "my-topic", Key: "key2", Value: "消息2"},
    {Topic: "my-topic", Key: "key3", Value: "消息3"},
}

batchResult, err := client.Producer().ProduceBatch(ctx, messages)
if err != nil {
    log.Fatal("批量发送失败:", err)
}

fmt.Printf("批量发送完成: %d 条消息\n", len(batchResult.Results))
```

### 2. 流式消费

```go
// 创建流式消费
stream, err := client.Consumer().ConsumeStream(ctx, fluvio.StreamConsumeOptions{
    Topic: "my-topic",
    Group: "stream-group",
})
if err != nil {
    log.Fatal("创建流式消费失败:", err)
}

// 持续接收消息
for msg := range stream {
    if msg.Error != nil {
        log.Printf("错误: %v", msg.Error)
        continue
    }
    fmt.Printf("流式消息: [%s] %s\n", msg.Message.Key, msg.Message.Value)
}
```

### 3. 错误处理

```go
// 带重试的操作
result, err := client.Producer().ProduceWithRetry(ctx, "重要消息", fluvio.ProduceOptions{
    Topic: "important-topic",
    Key:   "critical",
})
if err != nil {
    log.Printf("重试后仍失败: %v", err)
} else {
    fmt.Printf("重试成功: %s\n", result.MessageID)
}
```

## 📚 学习资源

- 📖 [完整使用指南](USAGE.md) - 详细的API文档和示例
- 🔧 [基本示例](examples/basic/) - 基础功能演示
- 🚀 [高级示例](examples/advanced/) - 高级功能和性能优化
- 🧪 [集成测试](examples/integration/) - 完整的功能测试

## 🆘 需要帮助？

1. **查看示例代码**: `examples/` 目录包含了各种使用场景
2. **阅读API文档**: 查看 `USAGE.md` 获取详细说明
3. **检查错误日志**: 启用调试日志查看详细信息

```go
// 启用调试日志
client, err := fluvio.New(
    fluvio.WithLogLevel(fluvio.LevelDebug),
)
```

## 🎊 恭喜！

你已经成功完成了 Fluvio Go SDK 的快速入门！现在你可以：

- ✅ 连接到 Fluvio 服务
- ✅ 创建和管理主题
- ✅ 发送和接收消息
- ✅ 处理错误和重试

继续探索更多高级功能，构建强大的流数据处理应用！🚀
