# 基本示例

本示例展示了 Fluvio Go SDK 的基本功能，包括新增的特性。

## 功能展示

### 基础功能
- ✅ 客户端创建和连接
- ✅ 健康检查
- ✅ 主题创建和管理

### 消息生产（新功能）
- ✅ 带自定义消息ID的消息生产
- ✅ 消息头部信息
- ✅ 批量消息生产
- ✅ 批量生产结果处理

### 消息消费（增强功能）
- ✅ 消费消息并显示MessageID
- ✅ 消息头部信息展示
- ✅ 偏移量信息

### 主题管理（新功能）
- ✅ 主题详细信息获取
- ✅ 分区信息查看
- ✅ 主题配置查看
- ✅ 主题统计信息
- ✅ 分区统计信息

## 运行示例

1. 确保 Fluvio 服务正在运行（默认在 101.43.173.154:50051）

2. 运行示例：
```bash
cd examples/basic
go mod tidy
go run main.go
```

## 预期输出

```
=== Fluvio Go SDK 基本示例 ===
🔍 检查连接...
✅ 连接成功!
📁 创建主题 'basic-example-topic'...
✅ 主题已就绪!
📤 生产消息...
✅ 消息发送成功! ID: msg-001
📤 批量生产消息...
  ✅ 批量消息 1 发送成功: batch-msg-001
  ✅ 批量消息 2 发送成功: batch-msg-002
  ✅ 批量消息 3 发送成功: batch-msg-003
✅ 批量发送完成: 3/3 成功
📥 消费消息...
✅ 收到 4 条消息:
  1. [greeting] Hello, Fluvio with MessageID! (MessageID: msg-001, Offset: 0)
     Headers: map[source:basic-example timestamp:2024-01-01T12:00:00Z version:1.0]
  2. [batch-1] 第一条批量消息 (MessageID: batch-msg-001, Offset: 1)
     Headers: map[batch:true index:1]
  3. [batch-2] 第二条批量消息 (MessageID: batch-msg-002, Offset: 2)
     Headers: map[batch:true index:2]
  4. [batch-3] 第三条批量消息 (MessageID: batch-msg-003, Offset: 3)
     Headers: map[batch:true index:3]
📊 获取主题详细信息...
✅ 主题详细信息:
  - 主题: basic-example-topic
  - 保留时间: 86400000 ms
  - 分区数: 2
  - 配置: map[cleanup.policy:delete segment.ms:3600000]
  - 分区 0: Leader=1, HighWatermark=2
  - 分区 1: Leader=1, HighWatermark=2
📈 获取主题统计信息...
✅ 主题统计信息:
  - 主题: basic-example-topic
  - 总消息数: 4
  - 总大小: 256 bytes
  - 分区数: 2
  - 分区统计:
    分区 0: 2 条消息, 128 bytes
    分区 1: 2 条消息, 128 bytes
🎉 基本示例完成!
```

## 新功能说明

### 1. 消息ID支持
- 可以为每条消息指定自定义ID
- 消费时可以获取消息ID
- 便于消息追踪和去重

### 2. 增强的主题配置
- 支持复制因子设置
- 支持保留时间配置
- 支持自定义主题配置参数

### 3. 主题详细信息
- 获取主题的详细配置
- 查看分区信息和状态
- 监控分区的高水位标记

### 4. 主题统计信息
- 查看主题和分区的消息统计
- 监控存储使用情况
- 性能分析和容量规划

### 5. 批量操作增强
- 批量生产支持消息ID
- 详细的批量操作结果
- 更好的错误处理

## 故障排除

如果遇到连接问题：
1. 检查 Fluvio 服务是否运行
2. 确认服务器地址和端口正确
3. 检查网络连接
4. 查看日志输出获取详细错误信息
