# 基本使用示例

这个示例展示了 Fluvio Go SDK 的基本使用方法。

## 功能演示

- 创建客户端连接
- 健康检查
- 主题管理（创建主题）
- 消息生产
- 消息消费
- 列出主题
- 获取集群信息

## 运行示例

1. 确保 Fluvio 服务正在运行（默认在 localhost:50051）

2. 运行示例：
```bash
cd examples/basic
go run main.go
```

## 预期输出

```
执行健康检查...
✓ 健康检查成功
确保主题 'example-topic' 存在...
✓ 主题已就绪
生产消息...
✓ 消息 1 发送成功: batch-0
✓ 消息 2 发送成功: batch-1
✓ 消息 3 发送成功: batch-2
✓ 消息 4 发送成功: batch-3
✓ 消息 5 发送成功: batch-4
消费消息...
✓ 收到 5 条消息:
  1. [key-1] Hello from Fluvio SDK! Message #1 (offset: 0)
  2. [key-2] Hello from Fluvio SDK! Message #2 (offset: 0)
  3. [key-3] Hello from Fluvio SDK! Message #3 (offset: 0)
  4. [key-4] Hello from Fluvio SDK! Message #4 (offset: 0)
  5. [key-5] Hello from Fluvio SDK! Message #5 (offset: 0)
列出所有主题...
✓ 找到 1 个主题: [example-topic]
获取集群信息...
✓ 集群状态: online, 控制器ID: 0
✓ 基本示例完成!
```
