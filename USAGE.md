# Fluvio Go SDK 使用指南

本文档详细介绍如何使用 Fluvio Go SDK 进行流数据处理。

## 📦 安装

### 使用 go get 安装

```bash
go get github.com/iwen-conf/fluvio_grpc_client
```

### 在项目中导入

```go
import "github.com/iwen-conf/fluvio_grpc_client"
```

## 🚀 快速开始

### 1. 创建客户端

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
    // 最简单的方式 - 使用默认配置
    client, err := fluvio.New()
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()

    // 或者使用自定义配置
    client, err = fluvio.New(
        fluvio.WithServer("101.43.173.154", 50051),
        fluvio.WithTimeout(5*time.Second, 10*time.Second),
        fluvio.WithLogLevel(fluvio.LevelInfo),
    )
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()

    fmt.Println("客户端创建成功!")
}
```

### 2. 健康检查

```go
func healthCheck(client *fluvio.Client) {
    ctx := context.Background()
    
    // 基本健康检查
    err := client.HealthCheck(ctx)
    if err != nil {
        log.Printf("健康检查失败: %v", err)
        return
    }
    fmt.Println("✅ 服务健康")

    // 带延迟测试的健康检查
    duration, err := client.Ping(ctx)
    if err != nil {
        log.Printf("Ping失败: %v", err)
        return
    }
    fmt.Printf("✅ 服务响应时间: %v\n", duration)
}
```

## 🏗️ 客户端配置

### 配置选项

```go
client, err := fluvio.New(
    // 服务器地址配置
    fluvio.WithServer("101.43.173.154", 50051),
    
    // 超时配置
    fluvio.WithTimeout(
        5*time.Second,  // 连接超时
        30*time.Second, // 操作超时
    ),
    
    // 日志级别
    fluvio.WithLogLevel(fluvio.LevelInfo), // Debug, Info, Warn, Error
    
    // 重试配置
    fluvio.WithMaxRetries(3),
    
    // 连接池配置
    fluvio.WithPoolSize(5),
    
    // TLS配置（如果需要）
    fluvio.WithTLS(true),
)
```

### 使用配置文件

```go
// 创建配置文件 config.json
{
    "server": {
        "host": "101.43.173.154",
        "port": 50051
    },
    "timeout": {
        "connect": "5s",
        "operation": "30s"
    },
    "log_level": "info",
    "max_retries": 3,
    "pool_size": 5
}

// 从配置文件加载
cfg, err := fluvio.LoadConfigFromFile("config.json")
if err != nil {
    log.Fatal("加载配置失败:", err)
}

client, err := fluvio.NewWithConfig(cfg)
if err != nil {
    log.Fatal("创建客户端失败:", err)
}
```

### 环境变量配置

```bash
export FLUVIO_HOST=101.43.173.154
export FLUVIO_PORT=50051
export FLUVIO_LOG_LEVEL=info
export FLUVIO_MAX_RETRIES=3
```

```go
// 从环境变量加载配置
cfg := fluvio.LoadConfigFromEnv()
client, err := fluvio.NewWithConfig(cfg)
```

## 📝 主题管理

### 创建主题

```go
func createTopic(client *fluvio.Client) {
    ctx := context.Background()
    
    // 创建基本主题
    result, err := client.Topic().Create(ctx, fluvio.CreateTopicOptions{
        Name:       "my-topic",
        Partitions: 3,
    })
    if err != nil {
        log.Printf("创建主题失败: %v", err)
        return
    }
    fmt.Printf("主题创建成功: %+v\n", result)

    // 如果不存在则创建
    result, err = client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
        Name:       "my-topic",
        Partitions: 1,
    })
    if err != nil {
        log.Printf("创建主题失败: %v", err)
        return
    }
    fmt.Printf("主题已就绪: %+v\n", result)
}
```

### 列出和管理主题

```go
func manageTopic(client *fluvio.Client) {
    ctx := context.Background()
    
    // 列出所有主题
    topics, err := client.Topic().List(ctx)
    if err != nil {
        log.Printf("列出主题失败: %v", err)
        return
    }
    fmt.Printf("找到 %d 个主题: %v\n", len(topics.Topics), topics.Topics)

    // 检查主题是否存在
    exists, err := client.Topic().Exists(ctx, "my-topic")
    if err != nil {
        log.Printf("检查主题失败: %v", err)
        return
    }
    fmt.Printf("主题 'my-topic' 存在: %v\n", exists)

    // 删除主题
    result, err := client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{
        Name: "old-topic",
    })
    if err != nil {
        log.Printf("删除主题失败: %v", err)
        return
    }
    fmt.Printf("主题删除结果: %+v\n", result)
}
```

## 📤 消息生产

### 基本生产

```go
func basicProduce(client *fluvio.Client) {
    ctx := context.Background()
    
    // 生产单条消息
    result, err := client.Producer().Produce(ctx, "Hello, Fluvio!", fluvio.ProduceOptions{
        Topic: "my-topic",
        Key:   "greeting",
    })
    if err != nil {
        log.Printf("生产消息失败: %v", err)
        return
    }
    fmt.Printf("消息发送成功: %s\n", result.MessageID)

    // 带头部信息的消息
    result, err = client.Producer().Produce(ctx, "带头部的消息", fluvio.ProduceOptions{
        Topic: "my-topic",
        Key:   "with-headers",
        Headers: map[string]string{
            "source":    "go-sdk",
            "version":   "1.0",
            "timestamp": time.Now().Format(time.RFC3339),
        },
    })
    if err != nil {
        log.Printf("生产消息失败: %v", err)
        return
    }
    fmt.Printf("带头部消息发送成功: %s\n", result.MessageID)
}
```

### 批量生产

```go
func batchProduce(client *fluvio.Client) {
    ctx := context.Background()
    
    // 准备批量消息
    messages := []fluvio.Message{
        {
            Topic: "my-topic",
            Key:   "batch-1",
            Value: "第一条批量消息",
            Headers: map[string]string{"batch": "true"},
        },
        {
            Topic: "my-topic",
            Key:   "batch-2", 
            Value: "第二条批量消息",
            Headers: map[string]string{"batch": "true"},
        },
        {
            Topic: "my-topic",
            Key:   "batch-3",
            Value: "第三条批量消息", 
            Headers: map[string]string{"batch": "true"},
        },
    }

    // 批量发送
    batchResult, err := client.Producer().ProduceBatch(ctx, messages)
    if err != nil {
        log.Printf("批量生产失败: %v", err)
        return
    }

    // 检查结果
    successCount := 0
    for i, result := range batchResult.Results {
        if result.Success {
            successCount++
            fmt.Printf("消息 %d 发送成功: %s\n", i+1, result.MessageID)
        } else {
            fmt.Printf("消息 %d 发送失败: %s\n", i+1, result.Error)
        }
    }
    fmt.Printf("批量发送完成: %d/%d 成功\n", successCount, len(messages))
}
```

### 异步生产

```go
func asyncProduce(client *fluvio.Client) {
    ctx := context.Background()
    
    // 异步发送消息
    resultChan := client.Producer().ProduceAsync(ctx, "异步消息", fluvio.ProduceOptions{
        Topic: "my-topic",
        Key:   "async",
    })

    // 处理结果
    go func() {
        result := <-resultChan
        if result.Error != nil {
            log.Printf("异步发送失败: %v", result.Error)
        } else {
            fmt.Printf("异步发送成功: %s\n", result.Result.MessageID)
        }
    }()

    // 继续其他工作...
    fmt.Println("异步发送已启动，继续其他工作...")
}
```

## 📥 消息消费

### 基本消费

```go
func basicConsume(client *fluvio.Client) {
    ctx := context.Background()

    // 消费消息
    messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
        Topic:       "my-topic",
        Group:       "my-group",
        MaxMessages: 10,
        Offset:      0, // 从头开始
    })
    if err != nil {
        log.Printf("消费消息失败: %v", err)
        return
    }

    fmt.Printf("收到 %d 条消息:\n", len(messages))
    for i, msg := range messages {
        fmt.Printf("  %d. [%s] %s (offset: %d)\n",
            i+1, msg.Key, msg.Value, msg.Offset)

        // 处理头部信息
        if len(msg.Headers) > 0 {
            fmt.Printf("     Headers: %v\n", msg.Headers)
        }
    }
}
```

### 流式消费

```go
func streamConsume(client *fluvio.Client) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 创建流式消费
    stream, err := client.Consumer().ConsumeStream(ctx, fluvio.StreamConsumeOptions{
        Topic: "my-topic",
        Group: "stream-group",
    })
    if err != nil {
        log.Printf("创建流式消费失败: %v", err)
        return
    }

    fmt.Println("开始流式消费...")
    messageCount := 0

    for {
        select {
        case msg, ok := <-stream:
            if !ok {
                fmt.Printf("流式消费结束，共收到 %d 条消息\n", messageCount)
                return
            }

            if msg.Error != nil {
                log.Printf("流式消费错误: %v", msg.Error)
                continue
            }

            messageCount++
            fmt.Printf("流式消息 %d: [%s] %s\n",
                messageCount, msg.Message.Key, msg.Message.Value)

        case <-ctx.Done():
            fmt.Printf("流式消费超时，共收到 %d 条消息\n", messageCount)
            return
        }
    }
}
```

### 手动偏移量管理

```go
func manualOffsetManagement(client *fluvio.Client) {
    ctx := context.Background()

    // 消费消息但不自动提交偏移量
    messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
        Topic:      "my-topic",
        Group:      "manual-group",
        MaxMessages: 5,
        AutoCommit: false, // 禁用自动提交
    })
    if err != nil {
        log.Printf("消费消息失败: %v", err)
        return
    }

    // 处理消息
    for _, msg := range messages {
        fmt.Printf("处理消息: [%s] %s\n", msg.Key, msg.Value)

        // 这里进行业务处理...
        // 如果处理成功，手动提交偏移量
    }

    // 手动提交偏移量
    if len(messages) > 0 {
        lastMessage := messages[len(messages)-1]
        err = client.Consumer().CommitOffset(ctx, fluvio.CommitOffsetOptions{
            Topic:  "my-topic",
            Group:  "manual-group",
            Offset: lastMessage.Offset + 1,
        })
        if err != nil {
            log.Printf("提交偏移量失败: %v", err)
        } else {
            fmt.Printf("偏移量提交成功: %d\n", lastMessage.Offset+1)
        }
    }
}
```

## 🔧 管理功能

### 集群管理

```go
func clusterManagement(client *fluvio.Client) {
    ctx := context.Background()

    // 获取集群信息
    cluster, err := client.Admin().DescribeCluster(ctx)
    if err != nil {
        log.Printf("获取集群信息失败: %v", err)
        return
    }
    fmt.Printf("集群状态: %s, 控制器ID: %d\n",
        cluster.Cluster.Status, cluster.Cluster.ControllerID)

    // 列出Brokers
    brokers, err := client.Admin().ListBrokers(ctx)
    if err != nil {
        log.Printf("列出Brokers失败: %v", err)
        return
    }
    fmt.Printf("找到 %d 个Broker:\n", len(brokers.Brokers))
    for i, broker := range brokers.Brokers {
        fmt.Printf("  %d. ID: %d, 地址: %s, 状态: %s\n",
            i+1, broker.ID, broker.Addr, broker.Status)
    }
}
```

### 消费组管理

```go
func consumerGroupManagement(client *fluvio.Client) {
    ctx := context.Background()

    // 列出消费组
    groups, err := client.Admin().ListConsumerGroups(ctx)
    if err != nil {
        log.Printf("列出消费组失败: %v", err)
        return
    }
    fmt.Printf("找到 %d 个消费组:\n", len(groups.Groups))
    for i, group := range groups.Groups {
        fmt.Printf("  %d. %s\n", i+1, group.GroupID)
    }

    // 获取消费组详情
    if len(groups.Groups) > 0 {
        groupName := groups.Groups[0].GroupID
        groupDetail, err := client.Admin().DescribeConsumerGroup(ctx, groupName)
        if err != nil {
            log.Printf("获取消费组详情失败: %v", err)
            return
        }

        fmt.Printf("消费组 '%s' 详情:\n", groupDetail.Group.GroupID)
        fmt.Printf("  偏移量信息: %v\n", groupDetail.Group.Offsets)
    }
}
```

### SmartModule管理

```go
func smartModuleManagement(client *fluvio.Client) {
    ctx := context.Background()

    // 列出SmartModules
    modules, err := client.Admin().ListSmartModules(ctx)
    if err != nil {
        log.Printf("列出SmartModules失败: %v", err)
        return
    }

    fmt.Printf("找到 %d 个SmartModule:\n", len(modules.SmartModules))
    for i, module := range modules.SmartModules {
        fmt.Printf("  %d. 名称: %s, 版本: %s\n",
            i+1, module.Name, module.Version)
        if module.Description != "" {
            fmt.Printf("     描述: %s\n", module.Description)
        }
    }
}
```

## 🚀 高级用法

### 错误处理和重试

```go
func errorHandlingAndRetry(client *fluvio.Client) {
    ctx := context.Background()

    // 带重试的生产
    result, err := client.Producer().ProduceWithRetry(ctx, "重试消息", fluvio.ProduceOptions{
        Topic: "my-topic",
        Key:   "retry-test",
    })
    if err != nil {
        log.Printf("重试后仍然失败: %v", err)
        return
    }
    fmt.Printf("重试成功: %s\n", result.MessageID)

    // 带重试的消费
    messages, err := client.Consumer().ConsumeWithRetry(ctx, fluvio.ConsumeOptions{
        Topic:       "my-topic",
        Group:       "retry-group",
        MaxMessages: 5,
    })
    if err != nil {
        log.Printf("重试消费失败: %v", err)
        return
    }
    fmt.Printf("重试消费成功，收到 %d 条消息\n", len(messages))
}
```

### 并发处理

```go
func concurrentProcessing(client *fluvio.Client) {
    ctx := context.Background()
    var wg sync.WaitGroup

    // 并发生产者
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(producerID int) {
            defer wg.Done()

            for j := 0; j < 5; j++ {
                message := fmt.Sprintf("并发消息 P%d-M%d", producerID, j+1)
                _, err := client.Producer().Produce(ctx, message, fluvio.ProduceOptions{
                    Topic: "concurrent-topic",
                    Key:   fmt.Sprintf("producer-%d-msg-%d", producerID, j+1),
                })
                if err != nil {
                    log.Printf("生产者 %d 消息 %d 失败: %v", producerID, j+1, err)
                }
            }
            fmt.Printf("生产者 %d 完成\n", producerID)
        }(i)
    }

    // 并发消费者
    for i := 0; i < 2; i++ {
        wg.Add(1)
        go func(consumerID int) {
            defer wg.Done()

            messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
                Topic:       "concurrent-topic",
                Group:       fmt.Sprintf("concurrent-group-%d", consumerID),
                MaxMessages: 10,
            })
            if err != nil {
                log.Printf("消费者 %d 失败: %v", consumerID, err)
                return
            }
            fmt.Printf("消费者 %d 收到 %d 条消息\n", consumerID, len(messages))
        }(i)
    }

    wg.Wait()
    fmt.Println("并发处理完成")
}
```

### 性能优化

```go
func performanceOptimization() {
    // 高性能客户端配置
    client, err := fluvio.New(
        fluvio.WithServer("101.43.173.154", 50051),
        fluvio.WithPoolSize(10),           // 增加连接池大小
        fluvio.WithMaxRetries(5),          // 增加重试次数
        fluvio.WithTimeout(2*time.Second, 30*time.Second), // 优化超时
        fluvio.WithLogLevel(fluvio.LevelWarn), // 减少日志输出
    )
    if err != nil {
        log.Fatal("创建高性能客户端失败:", err)
    }
    defer client.Close()

    // 或者使用预设的高性能配置
    highPerfClient, err := fluvio.HighThroughputClient("101.43.173.154", 50051)
    if err != nil {
        log.Fatal("创建高性能客户端失败:", err)
    }
    defer highPerfClient.Close()

    fmt.Println("高性能客户端创建成功")
}
```

## 📚 完整示例

查看 `examples/` 目录下的完整示例：

- `examples/basic/` - 基本使用示例
- `examples/advanced/` - 高级功能示例
- `examples/integration/` - 集成测试示例

## 🔍 故障排除

### 常见问题

1. **连接超时**
   ```
   [TIMEOUT] 等待连接就绪超时
   ```
   - 检查服务器地址和端口是否正确
   - 确认Fluvio服务正在运行
   - 检查网络连接

2. **认证失败**
   - 检查TLS配置
   - 确认服务器证书

3. **主题不存在**
   - 使用 `CreateIfNotExists` 自动创建主题
   - 检查主题名称拼写

### 调试技巧

```go
// 启用详细日志
client, err := fluvio.New(
    fluvio.WithLogLevel(fluvio.LevelDebug),
)

// 使用自定义日志器
logger := fluvio.NewLogger(fluvio.LevelDebug)
fluvio.SetDefaultLogger(logger)
```

## 📖 API参考

详细的API文档请参考项目的README.md文件和代码注释。
```
