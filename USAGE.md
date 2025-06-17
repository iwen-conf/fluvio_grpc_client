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
