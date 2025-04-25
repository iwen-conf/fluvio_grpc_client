# Fluvio gRPC 客户端

## 项目简介

Fluvio gRPC 客户端是一个基于 Go 语言的工具，用于与 Fluvio 消息流处理系统进行交互。该客户端通过 gRPC 协议提供了丰富的功能，包括消息的生产和消费、主题管理、消费者组管理、SmartModule 管理以及集群管理等功能。客户端提供了交互式命令行界面，方便用户进行操作。

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
├── cmd/                    # 命令行入口
│   └── client/             # 客户端命令
│       └── main.go         # 主程序入口
├── internal/               # 内部实现
│   ├── cli/                # 命令行工具
│   │   ├── handler.go      # 命令处理器
│   │   └── printer.go      # 输出格式化
│   ├── client/             # 客户端实现
│   │   ├── fluvio_admin.go # 管理服务客户端
│   │   ├── fluvio_service.go # 核心服务客户端
│   │   └── grpc_client.go  # gRPC 连接管理
│   └── config/             # 配置管理
│       ├── config.json     # 配置文件
│       └── load.go         # 配置加载
├── proto/                  # 协议定义
│   ├── fluvio_grpc.proto   # gRPC 协议定义
│   └── fluvio_service/     # 生成的协议代码
│       ├── fluvio_grpc.pb.go     # 消息定义
│       └── fluvio_grpc_grpc.pb.go # 服务定义
├── tests/                  # 测试代码
│   ├── grpc_health_test.go # 健康检查测试
│   └── grpc_service_test.go # 服务功能测试
├── go.mod                  # Go 模块定义
├── go.sum                  # 依赖校验和
└── README.md               # 项目说明文档
```

## 安装与使用

### 前置条件

- Go 1.18 或更高版本
- 正在运行的 Fluvio 服务实例
- 依赖库：github.com/iwen-conf/colorprint/clr

### 安装

```bash
git clone github.com/iwen-conf/fluvio_grpc_client
cd fluvio_grpc_client
go mod download
```

### 构建

```bash
go build -o fluvio-client ./cmd/client
```

### 运行

```bash
./fluvio-client
```

程序启动后会自动连接到配置文件中指定的服务器，执行健康检查，获取主题列表，并进入交互模式。

### 配置

编辑 `internal/config/config.json` 文件，设置 Fluvio 服务的连接信息：

```json
{
  "server": {
    "host": "localhost",
    "port": 50051
  }
}
```

### 使用示例

客户端启动后会自动进入交互模式，可以使用以下命令：

#### 生产消息

```bash
produce Hello, Fluvio!
```

#### 批量生产消息

```bash
batch_produce 消息1,消息2,消息3
```

#### 消费消息

```bash
consume
```

#### 健康检查

```bash
health
```

#### 列出主题

```bash
topics
```

#### 创建主题

```bash
create_topic new-topic 3
```

#### 删除主题

```bash
delete_topic topic-name
```

#### 退出程序

```bash
exit
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
