# 集成测试

本集成测试全面验证 Fluvio Go SDK 的所有功能，包括新增的特性。

## 测试覆盖

### 🔍 基础功能测试
- ✅ 健康检查和连接测试
- ✅ Ping响应时间测试

### 📁 主题管理测试
- ✅ 主题创建（带新配置选项）
- ✅ 主题存在性检查
- ✅ 主题列表获取
- ✅ 主题详细信息获取（新功能）
- ✅ 分区信息验证

### 📤 消息生产测试
- ✅ 单条消息生产（带MessageID）
- ✅ 消息头部信息
- ✅ 批量消息生产
- ✅ 批量结果验证

### 📥 消息消费测试
- ✅ 基本消息消费
- ✅ 消息内容验证
- ✅ MessageID验证
- ✅ 消息头部验证

### 🔍 过滤消费测试（新功能）
- ✅ 按头部过滤消费
- ✅ 过滤条件验证
- ✅ 过滤统计信息

### 📡 流式消费测试（增强功能）
- ✅ 流式消费创建
- ✅ 批次大小控制
- ✅ 等待时间控制
- ✅ 实时消息处理

### 📊 主题统计测试（新功能）
- ✅ 主题统计信息获取
- ✅ 分区统计信息
- ✅ 消息计数验证

### 👥 消费组管理测试
- ✅ 消费组列表
- ✅ 消费组详细信息

### 🧠 SmartModule管理测试（新功能）
- ✅ SmartModule列表

### 💾 存储管理测试（新功能）
- ✅ 存储状态获取
- ✅ 存储指标获取
- ✅ 健康状态检查

### 🗑️ 批量操作测试（新功能）
- ✅ 批量删除主题
- ✅ 批量操作结果验证

### ❌ 错误处理测试
- ✅ 不存在资源的错误处理
- ✅ 无效参数的错误处理

## 运行测试

1. 确保 Fluvio 服务正在运行（默认在 101.43.173.154:50051）

2. 运行集成测试：
```bash
cd examples/integration
go mod tidy
go run main.go
```

## 预期输出

```
=== Fluvio Go SDK 集成测试 ===

1. 🧪 健康检查测试
   响应时间: 15.2ms
   ✅ 通过

2. 🧪 主题管理测试
   主题创建成功: integration-test-topic (分区: 2)
   ✅ 通过

3. 🧪 消息生产测试
   生产消息成功: 1条单独消息 + 2条批量消息
   ✅ 通过

4. 🧪 消息消费测试
   消费消息成功: 3条消息
   ✅ 通过

5. 🧪 过滤消费测试
   过滤消费成功: 扫描3条，过滤出1条
   ✅ 通过

6. 🧪 流式消费测试
   流式消费成功: 3条消息
   ✅ 通过

7. 🧪 主题统计测试
   主题统计成功: 3条消息, 2个分区
   ✅ 通过

8. 🧪 消费组管理测试
   消费组管理成功: 4个消费组
   ✅ 通过

9. 🧪 SmartModule管理测试
   SmartModule管理成功: 2个模块
   ✅ 通过

10. 🧪 存储管理测试
    存储状态: 持久化=true, 类型=MongoDB
    存储指标: 响应时间=12ms, 健康状态=Healthy
    ✅ 通过

11. 🧪 批量操作测试
    批量操作成功: 2个请求, 2个成功, 0个失败
    ✅ 通过

12. 🧪 错误处理测试
    错误处理正常: 正确捕获了预期错误
    ✅ 通过

📊 测试结果: 12 通过, 0 失败, 总计 12
🎉 所有测试通过!
```

## 测试说明

### 新功能验证

#### 1. 消息ID支持
```go
// 验证自定义消息ID
result, err := client.Producer().Produce(ctx, "测试消息", types.ProduceOptions{
    MessageID: "test-msg-001",
})
// 验证消费时能获取到MessageID
if msg.MessageID != "test-msg-001" {
    return fmt.Errorf("消息ID不匹配")
}
```

#### 2. 过滤消费
```go
// 验证过滤功能
result, err := client.Consumer().ConsumeFiltered(ctx, types.FilteredConsumeOptions{
    Filters: []types.FilterCondition{
        {
            Type:     types.FilterTypeHeader,
            Field:    "test",
            Operator: "eq",
            Value:    "true",
        },
    },
})
```

#### 3. 主题详细信息
```go
// 验证主题详细信息获取
detail, err := client.Topic().DescribeTopicDetail(ctx, topicName)
if len(detail.Partitions) != 2 {
    return fmt.Errorf("分区数不匹配")
}
```

#### 4. 主题统计信息
```go
// 验证主题统计
stats, err := client.Topic().GetTopicStats(ctx, types.GetTopicStatsOptions{
    IncludePartitions: true,
})
```

#### 5. 存储管理
```go
// 验证存储状态
status, err := client.Admin().GetStorageStatus(ctx, types.GetStorageStatusOptions{
    IncludeDetails: true,
})

// 验证存储指标
metrics, err := client.Admin().GetStorageMetrics(ctx, types.GetStorageMetricsOptions{})
```

#### 6. 批量删除
```go
// 验证批量删除
result, err := client.Admin().BulkDelete(ctx, types.BulkDeleteOptions{
    Topics: []string{"topic1", "topic2"},
})
```

### 错误处理验证

测试确保SDK能正确处理各种错误情况：
- 不存在的资源
- 无效的参数
- 网络连接问题
- 超时情况

## 故障排除

### 常见问题

1. **连接失败**
   - 检查Fluvio服务是否运行
   - 确认服务器地址和端口
   - 检查网络连接

2. **测试超时**
   - 增加超时时间
   - 检查服务器性能
   - 减少测试数据量

3. **权限错误**
   - 检查客户端权限
   - 确认服务器配置

4. **资源冲突**
   - 清理之前的测试数据
   - 使用唯一的资源名称

### 调试技巧

1. **启用详细日志**
```go
client, err := fluvio.New(
    fluvio.WithLogLevel(fluvio.LevelDebug),
)
```

2. **单独运行测试**
修改main函数只运行特定测试：
```go
tests := []struct{...}{
    {"特定测试", testSpecificFunction},
}
```

3. **检查服务器日志**
查看Fluvio服务器端的日志获取更多信息。

## 持续集成

这个集成测试可以集成到CI/CD流水线中：

```bash
#!/bin/bash
# 启动Fluvio服务
docker run -d --name fluvio-test -p 50051:50051 fluvio/fluvio

# 等待服务启动
sleep 10

# 运行集成测试
cd examples/integration
go test -v

# 清理
docker stop fluvio-test
docker rm fluvio-test
```
