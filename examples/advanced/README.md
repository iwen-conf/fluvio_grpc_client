# 高级示例

本示例展示了 Fluvio Go SDK 的高级功能和新特性。

## 功能展示

### 🔍 过滤消费功能（新功能）
- ✅ 按消息头部过滤
- ✅ 按消息键过滤
- ✅ 按消息内容过滤
- ✅ 多条件过滤（AND/OR逻辑）
- ✅ 过滤统计信息

### 📡 流式消费增强功能
- ✅ 批次大小控制
- ✅ 等待时间控制
- ✅ 实时消息处理
- ✅ 并发安全的流式处理

### 🧠 SmartModule管理（新功能）
- ✅ 列出SmartModules
- ✅ 创建SmartModule
- ✅ SmartModule规格定义
- ✅ 参数化SmartModule

### 💾 存储管理功能（新功能）
- ✅ 存储状态监控
- ✅ 存储指标获取
- ✅ 连接池状态
- ✅ 数据库信息
- ✅ 性能指标
- ✅ 健康检查

### 🗑️ 批量删除功能（新功能）
- ✅ 批量删除主题
- ✅ 批量删除消费组
- ✅ 批量删除SmartModules
- ✅ 详细的删除结果

### ⚡ 并发处理
- ✅ 多生产者并发
- ✅ 多消费者并发
- ✅ 分区并行处理
- ✅ 协程安全

## 运行示例

1. 确保 Fluvio 服务正在运行（默认在 101.43.173.154:50051）

2. 运行示例：
```bash
cd examples/advanced
go mod tidy
go run main.go
```

## 预期输出

```
=== Fluvio Go SDK 高级示例 ===

🔍 演示过滤消费功能...
  🔍 过滤消费：只获取错误级别的消息
  ✅ 过滤结果: 扫描了 5 条消息，过滤出 1 条消息
    1. [user-1] 支付失败 (Headers: map[event:payment level:error])
  🔍 过滤消费：只获取user-1的消息
  ✅ 过滤结果: 扫描了 5 条消息，过滤出 2 条消息
    1. [user-1] 用户登录
    2. [user-1] 支付失败

📡 演示流式消费增强功能...
  📡 开始流式消费（批次大小=3，等待时间=1秒）...
  📦 批次 1:
    1. [stream-key-1] 流式消息 1 (ID: stream-msg-001)
    2. [stream-key-2] 流式消息 2 (ID: stream-msg-002)
    3. [stream-key-3] 流式消息 3 (ID: stream-msg-003)
  📦 批次 2:
    4. [stream-key-4] 流式消息 4 (ID: stream-msg-004)
    5. [stream-key-5] 流式消息 5 (ID: stream-msg-005)
    6. [stream-key-6] 流式消息 6 (ID: stream-msg-006)
  ✅ 流式消费结束，共收到 10 条消息，4 个批次

🧠 演示SmartModule管理...
  📋 当前SmartModules数量: 2
    1. example-filter (版本: 1.0.0) - 示例过滤器
    2. data-transformer (版本: 2.1.0) - 数据转换器
  🧠 创建示例SmartModule...
  ⚠️  创建SmartModule失败（预期的，因为没有真实WASM代码）

💾 演示存储管理功能...
  💾 获取存储状态...
  ✅ 存储状态:
    - 持久化启用: true
    - 存储类型: MongoDB
    - 连接状态: Connected
    - 消费组数量: 15
    - 消费偏移量数量: 45
    - SmartModule数量: 8
    - 当前连接数: 5
    - 可用连接数: 15
    - 数据库: fluvio_metadata
    - 集合数: 3
    - 数据大小: 2048576 bytes
  📊 获取存储指标...
  ✅ 存储指标:
    - 存储类型: MongoDB
    - 响应时间: 15 ms
    - 每秒操作数: 1250.50
    - 错误率: 0.02%
    - 连接池使用率: 33.33%
    - 内存使用: 128 MB
    - 磁盘使用: 512 MB
    - 健康状态: Healthy

🗑️ 演示批量删除功能...
  🏗️  创建测试主题...
  🗑️  执行批量删除...
  ✅ 批量删除结果:
    - 总请求数: 3
    - 成功删除: 3
    - 删除失败: 0
    1. ✅ bulk-test-topic-1 (topic)
    2. ✅ bulk-test-topic-2 (topic)
    3. ✅ bulk-test-topic-3 (topic)

⚡ 演示并发处理...
  ⚡ 启动并发生产者...
    ✅ 生产者 0 完成
    ✅ 生产者 1 完成
    ✅ 生产者 2 完成
  ⚡ 启动并发消费者...
    ✅ 消费者 0 收到 8 条消息:
      1. [producer-0-msg-1] 并发消息 P0-M1 (ID: concurrent-p0-m1, Producer: 0)
      2. [producer-1-msg-1] 并发消息 P1-M1 (ID: concurrent-p1-m1, Producer: 1)
      ...
    ✅ 消费者 1 收到 7 条消息:
      1. [producer-2-msg-1] 并发消息 P2-M1 (ID: concurrent-p2-m1, Producer: 2)
      2. [producer-0-msg-2] 并发消息 P0-M2 (ID: concurrent-p0-m2, Producer: 0)
      ...
  ✅ 并发处理完成

🎉 高级示例完成!
```

## 新功能详解

### 1. 过滤消费
```go
// 按头部过滤
result, err := client.Consumer().ConsumeFiltered(ctx, types.FilteredConsumeOptions{
    Topic: "my-topic",
    Group: "filter-group",
    Filters: []types.FilterCondition{
        {
            Type:     types.FilterTypeHeader,
            Field:    "level",
            Operator: "eq",
            Value:    "error",
        },
    },
    AndLogic: true, // AND逻辑
})
```

### 2. 流式消费控制
```go
// 控制批次大小和等待时间
stream, err := client.Consumer().ConsumeStream(ctx, types.StreamConsumeOptions{
    Topic:        "my-topic",
    Group:        "stream-group",
    MaxBatchSize: 10,   // 每批最多10条消息
    MaxWaitMs:    1000, // 最多等待1秒
})
```

### 3. SmartModule管理
```go
// 创建SmartModule
spec := &types.SmartModuleSpec{
    Name:        "my-filter",
    InputKind:   types.SmartModuleInputStream,
    OutputKind:  types.SmartModuleOutputStream,
    Description: "自定义过滤器",
    Version:     "1.0.0",
}

result, err := client.Admin().CreateSmartModule(ctx, types.CreateSmartModuleOptions{
    Spec:     spec,
    WasmCode: wasmBytes,
})
```

### 4. 存储管理
```go
// 获取存储状态
status, err := client.Admin().GetStorageStatus(ctx, types.GetStorageStatusOptions{
    IncludeDetails: true,
})

// 获取存储指标
metrics, err := client.Admin().GetStorageMetrics(ctx, types.GetStorageMetricsOptions{
    IncludeHistory: true,
    HistoryLimit:   10,
})
```

### 5. 批量删除
```go
// 批量删除资源
result, err := client.Admin().BulkDelete(ctx, types.BulkDeleteOptions{
    Topics:         []string{"topic1", "topic2"},
    ConsumerGroups: []string{"group1", "group2"},
    SmartModules:   []string{"module1", "module2"},
    Force:          false,
})
```

## 性能优化建议

1. **连接池配置**: 使用大连接池提高并发性能
2. **批次处理**: 使用批量操作减少网络开销
3. **过滤消费**: 在服务端过滤减少网络传输
4. **流式消费**: 控制批次大小平衡延迟和吞吐量
5. **并发处理**: 利用多协程提高处理效率

## 故障排除

1. **过滤消费无结果**: 检查过滤条件是否正确
2. **SmartModule创建失败**: 确保WASM代码有效
3. **存储连接问题**: 检查存储服务状态
4. **并发冲突**: 使用不同的消费组避免冲突
