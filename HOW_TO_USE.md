# Fluvio Go SDK 使用方法

本文档详细说明如何使用 Fluvio Go SDK，包括导入、创建客户端、配置等。

## 📦 安装和导入

### 1. 安装 SDK

```bash
go get github.com/iwen-conf/fluvio_grpc_client
```

### 2. 在项目中导入

```go
import "github.com/iwen-conf/fluvio_grpc_client"
```

### 3. 创建 go.mod 文件

```bash
# 初始化新项目
go mod init your-project-name

# 添加 Fluvio SDK 依赖
go get github.com/iwen-conf/fluvio_grpc_client
```

## 🚀 创建客户端

### 方法1: 使用默认配置

```go
package main

import (
    "log"
    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 使用默认配置创建客户端
    client, err := fluvio.New()
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()
    
    // 客户端已就绪，可以使用
}
```

### 方法2: 使用自定义配置

```go
package main

import (
    "log"
    "time"
    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 使用配置选项创建客户端
    client, err := fluvio.New(
        fluvio.WithServer("101.43.173.154", 50051),  // 服务器地址
        fluvio.WithTimeout(5*time.Second, 30*time.Second), // 连接和操作超时
        fluvio.WithLogLevel(fluvio.LevelInfo),       // 日志级别
        fluvio.WithMaxRetries(3),                    // 最大重试次数
        fluvio.WithPoolSize(5),                      // 连接池大小
    )
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()
}
```

### 方法3: 使用配置文件

```go
package main

import (
    "log"
    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 从配置文件加载配置
    cfg, err := fluvio.LoadConfigFromFile("config.json")
    if err != nil {
        log.Fatal("加载配置文件失败:", err)
    }
    
    // 使用配置创建客户端
    client, err := fluvio.NewWithConfig(cfg)
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()
}
```

配置文件示例 (`config.json`):
```json
{
  "server": {
    "host": "101.43.173.154",
    "port": 50051,
    "tls": {
      "enabled": false
    }
  },
  "connection": {
    "connect_timeout": "5s",
    "call_timeout": "30s",
    "max_retries": 3,
    "pool_size": 5
  },
  "logging": {
    "level": "info",
    "format": "text",
    "output": "stdout"
  }
}
```

### 方法4: 使用环境变量

```bash
# 设置环境变量
export FLUVIO_HOST=101.43.173.154
export FLUVIO_PORT=50051
export FLUVIO_LOG_LEVEL=info
export FLUVIO_MAX_RETRIES=3
```

```go
package main

import (
    "log"
    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 从环境变量加载配置
    cfg := fluvio.LoadConfigFromEnv()
    
    client, err := fluvio.NewWithConfig(cfg)
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()
}
```

### 方法5: 快速连接

```go
package main

import (
    "log"
    "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
    // 快速连接到指定服务器
    client, err := fluvio.QuickStart("101.43.173.154", 50051)
    if err != nil {
        log.Fatal("快速连接失败:", err)
    }
    defer client.Close()
}
```

## ⚙️ 配置选项详解

### 服务器配置

```go
// 设置服务器地址和端口
fluvio.WithServer("101.43.173.154", 50051)

// 启用TLS
fluvio.WithTLS(true)

// 使用不安全的TLS（跳过证书验证）
fluvio.WithInsecureTLS(true)
```

### 超时配置

```go
// 设置连接和操作超时
fluvio.WithTimeout(5*time.Second, 30*time.Second)

// 只设置连接超时
fluvio.WithConnectTimeout(5*time.Second)

// 只设置操作超时
fluvio.WithCallTimeout(30*time.Second)
```

### 重试配置

```go
// 设置最大重试次数
fluvio.WithMaxRetries(5)

// 设置重试策略
fluvio.WithRetry(fluvio.RetryConfig{
    MaxRetries:      5,
    InitialBackoff:  100 * time.Millisecond,
    MaxBackoff:      10 * time.Second,
    BackoffMultiple: 2.0,
})
```

### 连接池配置

```go
// 设置连接池大小
fluvio.WithPoolSize(10)

// 设置Keep-Alive
fluvio.WithKeepAlive(30*time.Second)
```

### 日志配置

```go
// 设置日志级别
fluvio.WithLogLevel(fluvio.LevelDebug) // Debug, Info, Warn, Error

// 使用自定义日志器
logger := fluvio.NewLogger(fluvio.LevelInfo)
fluvio.WithLogger(logger)
```

## 🔧 完整使用示例

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
        fluvio.WithTimeout(5*time.Second, 30*time.Second),
        fluvio.WithLogLevel(fluvio.LevelInfo),
        fluvio.WithMaxRetries(3),
    )
    if err != nil {
        log.Fatal("创建客户端失败:", err)
    }
    defer client.Close()

    ctx := context.Background()

    // 2. 健康检查
    err = client.HealthCheck(ctx)
    if err != nil {
        log.Fatal("健康检查失败:", err)
    }
    fmt.Println("✅ 连接成功")

    // 3. 创建主题
    _, err = client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
        Name:       "my-topic",
        Partitions: 1,
    })
    if err != nil {
        log.Fatal("创建主题失败:", err)
    }

    // 4. 生产消息
    result, err := client.Producer().Produce(ctx, "Hello, Fluvio!", fluvio.ProduceOptions{
        Topic: "my-topic",
        Key:   "greeting",
    })
    if err != nil {
        log.Fatal("生产消息失败:", err)
    }
    fmt.Printf("✅ 消息发送成功: %s\n", result.MessageID)

    // 5. 消费消息
    messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
        Topic:       "my-topic",
        Group:       "my-group",
        MaxMessages: 10,
    })
    if err != nil {
        log.Fatal("消费消息失败:", err)
    }
    
    fmt.Printf("✅ 收到 %d 条消息\n", len(messages))
    for i, msg := range messages {
        fmt.Printf("  %d. [%s] %s\n", i+1, msg.Key, msg.Value)
    }
}
```

## 🎯 最佳实践

### 1. 错误处理

```go
// 总是检查错误
client, err := fluvio.New()
if err != nil {
    log.Fatal("创建客户端失败:", err)
}

// 使用defer确保资源清理
defer client.Close()
```

### 2. 上下文管理

```go
// 使用带超时的上下文
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

err := client.HealthCheck(ctx)
```

### 3. 配置管理

```go
// 生产环境配置
client, err := fluvio.New(
    fluvio.WithServer("prod.fluvio.com", 50051),
    fluvio.WithTLS(true),
    fluvio.WithTimeout(10*time.Second, 60*time.Second),
    fluvio.WithMaxRetries(5),
    fluvio.WithLogLevel(fluvio.LevelWarn),
)

// 开发环境配置
client, err := fluvio.New(
    fluvio.WithServer("localhost", 50051),
    fluvio.WithTimeout(2*time.Second, 10*time.Second),
    fluvio.WithLogLevel(fluvio.LevelDebug),
)
```

### 4. 资源管理

```go
// 使用连接池提高性能
client, err := fluvio.New(
    fluvio.WithPoolSize(10),
    fluvio.WithKeepAlive(30*time.Second),
)

// 及时关闭客户端
defer client.Close()
```

## 🚨 常见错误

### 1. 忘记关闭客户端

```go
// ❌ 错误：没有关闭客户端
client, err := fluvio.New()
// 程序结束时可能泄露资源

// ✅ 正确：使用defer关闭
client, err := fluvio.New()
defer client.Close()
```

### 2. 没有处理错误

```go
// ❌ 错误：忽略错误
client, _ := fluvio.New()

// ✅ 正确：处理错误
client, err := fluvio.New()
if err != nil {
    log.Fatal("创建客户端失败:", err)
}
```

### 3. 超时设置不当

```go
// ❌ 错误：超时时间太短
fluvio.WithTimeout(100*time.Millisecond, 200*time.Millisecond)

// ✅ 正确：合理的超时时间
fluvio.WithTimeout(5*time.Second, 30*time.Second)
```

## 📚 更多资源

- 📖 [完整API文档](USAGE.md)
- 🚀 [快速入门指南](QUICKSTART.md)
- 💡 [示例代码](examples/)
- 🔧 [配置指南](examples/config-example.json)
