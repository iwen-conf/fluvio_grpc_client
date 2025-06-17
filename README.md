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
        fluvio.WithServer("localhost", 50051),
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

    // 生产消息
    result, err := client.Producer().Produce(ctx, "Hello, Fluvio!", fluvio.ProduceOptions{
        Topic: "my-topic",
        Key:   "key1",
    })
    if err != nil {
        log.Fatal("生产消息失败:", err)
    }
    fmt.Printf("消息发送成功: %+v\n", result)

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
}
```

## API 文档

### 客户端创建

```go
// 使用默认配置
client, err := fluvio.New()

// 使用自定义配置
client, err := fluvio.New(
    fluvio.WithServer("localhost", 50051),
    fluvio.WithTimeout(5*time.Second, 10*time.Second),
    fluvio.WithLogLevel(fluvio.LevelInfo),
    fluvio.WithMaxRetries(3),
    fluvio.WithPoolSize(5),
)

// 使用配置文件
cfg, err := fluvio.LoadConfigFromFile("config.json")
client, err := fluvio.NewWithConfig(cfg)

// 快速连接
client, err := fluvio.QuickStart("localhost", 50051)
```

### 消息生产

```go
// 基本生产
result, err := client.Producer().Produce(ctx, "Hello World", fluvio.ProduceOptions{
    Topic: "my-topic",
    Key:   "key1",
})

// 批量生产
messages := []fluvio.Message{
    {Topic: "my-topic", Key: "key1", Value: "message1"},
    {Topic: "my-topic", Key: "key2", Value: "message2"},
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

// 流式消费
stream, err := client.Consumer().ConsumeStream(ctx, fluvio.StreamConsumeOptions{
    Topic: "my-topic",
    Group: "my-group",
})

for msg := range stream {
    if msg.Error != nil {
        log.Printf("Error: %v", msg.Error)
        continue
    }
    fmt.Printf("Received: %s\n", msg.Message.Value)
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

// 创建主题
result, err := client.Topic().Create(ctx, fluvio.CreateTopicOptions{
    Name:       "new-topic",
    Partitions: 3,
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
// 集群信息
cluster, err := client.Admin().DescribeCluster(ctx)

// Broker列表
brokers, err := client.Admin().ListBrokers(ctx)

// 获取指标
metrics, err := client.Admin().GetMetrics(ctx, fluvio.GetMetricsOptions{
    MetricNames: []string{"cpu", "memory"},
})

// SmartModule管理
smartModules, err := client.Admin().ListSmartModules(ctx)
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

## 交互式命令行

本客户端提供了交互式命令行界面，支持以下命令：

- `help` - 显示帮助信息
- `produce <消息内容>` - 生产单条消息
- `batch_produce <消息1,消息2,...>` - 批量生产消息
- `consume` - 消费消息
- `health` - 健康检查
- `topics` - 列出所有主题
- `create_topic <主题名> <分区数>` - 创建主题
- `delete_topic <主题名>` - 删除主题
- `exit` 或 `quit` - 退出程序
