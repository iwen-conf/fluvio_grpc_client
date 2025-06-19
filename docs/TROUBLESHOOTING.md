# Fluvio Go SDK 故障排除指南

本指南帮助您诊断和解决使用 Fluvio Go SDK 时遇到的常见问题。

## 目录

- [连接问题](#连接问题)
- [认证和安全问题](#认证和安全问题)
- [性能问题](#性能问题)
- [消息发送问题](#消息发送问题)
- [消息消费问题](#消息消费问题)
- [主题管理问题](#主题管理问题)
- [错误代码参考](#错误代码参考)
- [调试工具](#调试工具)
- [常见配置错误](#常见配置错误)

---

## 连接问题

### 问题 1: `connection refused` 错误

**症状:**
```
[CONNECTION_ERROR] failed to connect to server: dial tcp 127.0.0.1:50051: connect: connection refused
```

**可能原因:**
1. Fluvio 服务器未运行
2. 端口号错误
3. 防火墙阻止连接
4. 网络配置问题

**解决方案:**

```go
// 1. 检查服务器状态
func checkServerStatus() {
    // 使用系统命令检查端口
    cmd := exec.Command("nc", "-zv", "localhost", "50051")
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("端口检查失败: %v", err)
        log.Printf("输出: %s", output)
    } else {
        log.Println("端口可达")
    }
}

// 2. 使用正确的地址配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051), // 确认地址和端口
    fluvio.WithTimeout(10*time.Second),     // 增加超时时间
)

// 3. 实施连接重试
func connectWithRetry(client *fluvio.Client, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        err := client.Connect(ctx)
        cancel()
        
        if err == nil {
            return nil
        }
        
        log.Printf("连接尝试 %d/%d 失败: %v", i+1, maxRetries, err)
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return fmt.Errorf("连接重试 %d 次后失败", maxRetries)
}
```### 问题 2: `timeout` 错误

**症状:**
```
[TIMEOUT_ERROR] 等待连接就绪超时
```

**解决方案:**

```go
// 1. 增加超时时间
client, err := fluvio.NewClient(
    fluvio.WithAddress("remote-server.com", 50051),
    fluvio.WithTimeout(60*time.Second), // 增加到 60 秒
)

// 2. 检查网络延迟
func checkNetworkLatency(host string) {
    start := time.Now()
    conn, err := net.DialTimeout("tcp", host+":50051", 5*time.Second)
    if err != nil {
        log.Printf("网络连接失败: %v", err)
        return
    }
    defer conn.Close()
    
    latency := time.Since(start)
    log.Printf("网络延迟: %v", latency)
    
    if latency > 5*time.Second {
        log.Println("警告: 网络延迟较高，建议增加超时时间")
    }
}

// 3. 使用上下文控制超时
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := client.Connect(ctx); err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("连接超时，可能需要增加超时时间或检查网络")
    }
}
```

### 问题 3: 连接频繁断开

**症状:**
- 连接建立后很快断开
- `IsConnected()` 返回 false

**解决方案:**

```go
// 1. 启用 Keep-Alive
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithKeepAlive(30*time.Second), // 启用 Keep-Alive
)

// 2. 实施连接监控和自动重连
func monitorConnection(client *fluvio.Client) {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if !client.IsConnected() {
            log.Println("检测到连接断开，尝试重新连接...")
            
            ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
            if err := client.Connect(ctx); err != nil {
                log.Printf("重连失败: %v", err)
            } else {
                log.Println("重连成功")
            }
            cancel()
        }
    }
}

// 3. 使用连接池
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithConnectionPool(5, 10*time.Minute), // 连接池
)
```---

## 认证和安全问题

### 问题 4: TLS 证书错误

**症状:**
```
[AUTHENTICATION_ERROR] x509: certificate signed by unknown authority
```

**解决方案:**

```go
// 1. 检查证书文件
func validateCertificates() {
    certFiles := map[string]string{
        "客户端证书": "client.crt",
        "客户端私钥": "client.key",
        "CA证书":    "ca.crt",
    }
    
    for name, file := range certFiles {
        if _, err := os.Stat(file); os.IsNotExist(err) {
            log.Printf("❌ %s 文件不存在: %s", name, file)
        } else {
            log.Printf("✅ %s 文件存在: %s", name, file)
        }
    }
}

// 2. 正确配置 TLS
client, err := fluvio.NewClient(
    fluvio.WithAddress("secure-server.com", 50051),
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
)

// 3. 开发环境：使用不安全连接
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithInsecure(), // 仅用于开发环境
)
```

### 问题 5: 权限被拒绝

**症状:**
```
[PERMISSION_DENIED] insufficient permissions for operation
```

**解决方案:**

```go
// 1. 检查客户端权限
func checkPermissions(client *fluvio.Client) {
    ctx := context.Background()
    
    // 测试基本操作权限
    operations := map[string]func() error{
        "健康检查": func() error { return client.HealthCheck(ctx) },
        "列出主题": func() error { _, err := client.Topics().List(ctx); return err },
        "集群信息": func() error { _, err := client.Admin().ClusterInfo(ctx); return err },
    }
    
    for name, op := range operations {
        if err := op(); err != nil {
            log.Printf("❌ %s 权限检查失败: %v", name, err)
        } else {
            log.Printf("✅ %s 权限正常", name)
        }
    }
}

// 2. 使用正确的证书
// 确保客户端证书具有必要的权限
```

---

## 性能问题

### 问题 6: 消息发送缓慢

**症状:**
- 单条消息发送耗时过长
- 批量发送超时

**解决方案:**

```go
// 1. 使用批量发送
func optimizeSending() {
    var messages []*fluvio.Message
    
    // 收集消息到批次中
    for i := 0; i < 1000; i++ {
        messages = append(messages, &fluvio.Message{
            Key:   fmt.Sprintf("key-%d", i),
            Value: []byte(fmt.Sprintf("message-%d", i)),
        })
        
        // 每 100 条消息发送一次
        if len(messages) >= 100 {
            sendBatch(messages)
            messages = messages[:0] // 清空切片
        }
    }
    
    // 发送剩余消息
    if len(messages) > 0 {
        sendBatch(messages)
    }
}

func sendBatch(messages []*fluvio.Message) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    result, err := client.Producer().SendBatch(ctx, "topic", messages)
    if err != nil {
        log.Printf("批量发送失败: %v", err)
        return
    }
    
    log.Printf("批量发送完成: 成功 %d, 失败 %d", 
        result.SuccessCount, result.FailureCount)
}

// 2. 调整连接池大小
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithConnectionPool(20, 30*time.Minute), // 增加连接数
)

// 3. 并行发送
func parallelSend(messages []*fluvio.Message) {
    const workerCount = 10
    messageChan := make(chan *fluvio.Message, 100)
    
    // 启动工作协程
    var wg sync.WaitGroup
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for msg := range messageChan {
                _, err := client.Producer().Send(ctx, "topic", msg)
                if err != nil {
                    log.Printf("发送失败: %v", err)
                }
            }
        }()
    }
    
    // 分发消息
    for _, msg := range messages {
        messageChan <- msg
    }
    close(messageChan)
    
    wg.Wait()
}
```### 问题 7: 消息消费缓慢

**症状:**
- 消息接收延迟高
- 流式消费阻塞

**解决方案:**

```go
// 1. 优化流式消费
stream, err := client.Consumer().Stream(ctx, "topic", &fluvio.StreamOptions{
    Group:      "fast-consumer",
    BufferSize: 1000, // 增加缓冲区大小
})

// 2. 并行处理消息
func parallelConsume() {
    const workerCount = 10
    messageChan := make(chan *fluvio.ConsumedMessage, 1000)
    
    // 启动处理协程
    var wg sync.WaitGroup
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for msg := range messageChan {
                processMessage(msg)
            }
        }(i)
    }
    
    // 消费消息
    for msg := range stream {
        select {
        case messageChan <- msg:
        case <-ctx.Done():
            close(messageChan)
            wg.Wait()
            return
        }
    }
}

// 3. 批量提交偏移量
func batchCommitOffsets() {
    var offsets []int64
    const batchSize = 100
    
    for msg := range stream {
        processMessage(msg)
        offsets = append(offsets, msg.Offset)
        
        // 批量提交
        if len(offsets) >= batchSize {
            lastOffset := offsets[len(offsets)-1]
            if err := client.Consumer().Commit(ctx, "topic", "group", lastOffset); err != nil {
                log.Printf("提交偏移量失败: %v", err)
            }
            offsets = offsets[:0]
        }
    }
}
```

---

## 消息发送问题

### 问题 8: 消息发送失败

**症状:**
```
[OPERATION_ERROR] failed to send message: topic not found
```

**解决方案:**

```go
// 1. 检查主题是否存在
func ensureTopicExists(topicName string) error {
    exists, err := client.Topics().Exists(ctx, topicName)
    if err != nil {
        return fmt.Errorf("检查主题失败: %w", err)
    }
    
    if !exists {
        log.Printf("主题 %s 不存在，正在创建...", topicName)
        err = client.Topics().Create(ctx, topicName, &fluvio.CreateTopicOptions{
            Partitions: 1,
        })
        if err != nil {
            return fmt.Errorf("创建主题失败: %w", err)
        }
        log.Printf("主题 %s 创建成功", topicName)
    }
    
    return nil
}

// 2. 实施重试机制
func sendWithRetry(topic, key, value string, maxRetries int) error {
    for attempt := 0; attempt < maxRetries; attempt++ {
        _, err := client.Producer().SendString(ctx, topic, key, value)
        if err == nil {
            return nil
        }
        
        // 检查是否为可重试错误
        if errors.IsNotFoundError(err) {
            // 尝试创建主题
            if createErr := ensureTopicExists(topic); createErr != nil {
                return createErr
            }
            continue
        }
        
        if !errors.IsConnectionError(err) && !errors.IsTimeoutError(err) {
            return err // 不可重试的错误
        }
        
        log.Printf("发送失败 (尝试 %d/%d): %v", attempt+1, maxRetries, err)
        time.Sleep(time.Duration(attempt+1) * time.Second)
    }
    
    return fmt.Errorf("重试 %d 次后仍然失败", maxRetries)
}

// 3. 验证消息格式
func validateMessage(msg *fluvio.Message) error {
    if len(msg.Value) == 0 {
        return fmt.Errorf("消息值不能为空")
    }
    
    if len(msg.Value) > 1024*1024 { // 1MB
        return fmt.Errorf("消息过大: %d bytes", len(msg.Value))
    }
    
    return nil
}
```---

## 消息消费问题

### 问题 9: 消费者组偏移量问题

**症状:**
- 重复消费消息
- 丢失消息
- 偏移量提交失败

**解决方案:**

```go
// 1. 正确的偏移量管理
func properOffsetManagement() {
    messages, err := client.Consumer().Receive(ctx, "topic", &fluvio.ReceiveOptions{
        Group:       "my-group",
        Offset:      -1, // 从最新开始，0 表示从头开始
        MaxMessages: 100,
    })
    if err != nil {
        log.Printf("接收失败: %v", err)
        return
    }
    
    for _, msg := range messages {
        // 处理消息
        if err := processMessage(msg); err != nil {
            log.Printf("处理消息失败: %v", err)
            continue // 不提交此消息的偏移量
        }
        
        // 只有成功处理后才提交偏移量
        if err := client.Consumer().Commit(ctx, "topic", "my-group", msg.Offset); err != nil {
            log.Printf("提交偏移量失败: %v", err)
        }
    }
}

// 2. 检查消费者组状态
func checkConsumerGroupStatus(groupID string) {
    detail, err := client.Admin().ConsumerGroupDetail(ctx, groupID)
    if err != nil {
        log.Printf("获取消费者组详情失败: %v", err)
        return
    }
    
    log.Printf("消费者组状态:")
    log.Printf("  组ID: %s", detail.GroupID)
    log.Printf("  状态: %s", detail.State)
    log.Printf("  成员数: %d", len(detail.Members))
    
    for i, member := range detail.Members {
        log.Printf("  成员 %d: %s (%s)", i+1, member.MemberID, member.ClientHost)
    }
}

// 3. 处理重复消费
func handleDuplicateMessages() {
    processedMessages := make(map[string]bool)
    
    for msg := range stream {
        // 使用消息ID或组合键检查重复
        messageKey := fmt.Sprintf("%s-%d", msg.Key, msg.Offset)
        
        if processedMessages[messageKey] {
            log.Printf("跳过重复消息: %s", messageKey)
            continue
        }
        
        if err := processMessage(msg); err != nil {
            log.Printf("处理消息失败: %v", err)
            continue
        }
        
        processedMessages[messageKey] = true
        client.Consumer().Commit(ctx, msg.Topic, "my-group", msg.Offset)
    }
}
```

### 问题 10: 流式消费中断

**症状:**
- 流式消费突然停止
- 通道关闭
- 上下文取消

**解决方案:**

```go
// 1. 实施自动重启机制
func resilientStreamConsume(topicName string) {
    for {
        if err := startStreamConsume(topicName); err != nil {
            log.Printf("流式消费中断: %v", err)
            log.Println("5秒后重新启动...")
            time.Sleep(5 * time.Second)
            continue
        }
        break
    }
}

func startStreamConsume(topicName string) error {
    stream, err := client.Consumer().Stream(ctx, topicName, &fluvio.StreamOptions{
        Group:      "resilient-group",
        BufferSize: 1000,
    })
    if err != nil {
        return fmt.Errorf("启动流式消费失败: %w", err)
    }
    
    for msg := range stream {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := processMessage(msg); err != nil {
                log.Printf("处理消息失败: %v", err)
            }
        }
    }
    
    return fmt.Errorf("流式消费意外结束")
}

// 2. 优雅关闭
func gracefulShutdown() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    ctx, cancel := context.WithCancel(context.Background())
    
    go func() {
        <-sigChan
        log.Println("收到关闭信号，正在优雅关闭...")
        cancel()
    }()
    
    // 启动消费
    resilientStreamConsume("my-topic")
}
```

---

## 错误代码参考

### 常见错误代码及解决方案

| 错误代码 | 描述 | 常见原因 | 解决方案 |
|---------|------|---------|----------|
| `CONNECTION_ERROR` | 连接错误 | 服务器不可达、网络问题 | 检查服务器状态、网络配置 |
| `TIMEOUT_ERROR` | 超时错误 | 操作超时 | 增加超时时间、检查网络延迟 |
| `AUTHENTICATION_ERROR` | 认证错误 | TLS证书问题、权限不足 | 检查证书配置、权限设置 |
| `NOT_FOUND` | 资源未找到 | 主题不存在、消费者组不存在 | 创建相应资源 |
| `ALREADY_EXISTS` | 资源已存在 | 重复创建主题 | 使用 `CreateIfNotExists` |
| `INVALID_ARGUMENT` | 参数无效 | 配置参数错误 | 检查参数格式和范围 |
| `OPERATION_ERROR` | 操作错误 | 业务逻辑错误 | 检查操作逻辑 |

### 错误处理模板

```go
func handleError(err error) {
    switch {
    case errors.IsConnectionError(err):
        log.Println("🔌 连接错误 - 检查网络和服务器")
        // 实施重连逻辑
        
    case errors.IsTimeoutError(err):
        log.Println("⏰ 超时错误 - 考虑增加超时时间")
        // 可能需要重试
        
    case errors.IsAuthenticationError(err):
        log.Println("🔐 认证错误 - 检查证书和权限")
        // 检查 TLS 配置
        
    case errors.IsNotFoundError(err):
        log.Println("🔍 资源未找到 - 检查主题或组是否存在")
        // 创建缺失的资源
        
    case errors.IsAlreadyExistsError(err):
        log.Println("✅ 资源已存在 - 可以继续使用")
        // 通常可以忽略
        
    case errors.IsValidationError(err):
        log.Println("❌ 参数验证失败 - 检查输入参数")
        // 修正参数
        
    default:
        log.Printf("❓ 未知错误: %v", err)
        // 通用处理
    }
}
```---

## 调试工具

### 启用详细日志

```go
// 开发环境：启用调试日志
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithLogLevel(fluvio.LogLevelDebug),
)

// 自定义日志处理
logger := client.Logger()
logger.Debug("调试信息", logging.Field{Key: "operation", Value: "connect"})
logger.Info("信息日志", logging.Field{Key: "status", Value: "success"})
logger.Warn("警告日志", logging.Field{Key: "issue", Value: "high_latency"})
logger.Error("错误日志", logging.Field{Key: "error", Value: err})
```

### 连接诊断工具

```go
func diagnoseConnection(client *fluvio.Client) {
    fmt.Println("=== 连接诊断 ===")
    
    // 1. 检查连接状态
    if client.IsConnected() {
        fmt.Println("✅ 客户端已连接")
    } else {
        fmt.Println("❌ 客户端未连接")
    }
    
    // 2. 健康检查
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.HealthCheck(ctx); err != nil {
        fmt.Printf("❌ 健康检查失败: %v\n", err)
    } else {
        fmt.Println("✅ 健康检查通过")
    }
    
    // 3. Ping 测试
    if duration, err := client.Ping(ctx); err != nil {
        fmt.Printf("❌ Ping 失败: %v\n", err)
    } else {
        fmt.Printf("✅ Ping 成功，延迟: %v\n", duration)
        
        if duration > 100*time.Millisecond {
            fmt.Println("⚠️ 延迟较高，可能影响性能")
        }
    }
    
    // 4. 基本操作测试
    if topics, err := client.Topics().List(ctx); err != nil {
        fmt.Printf("❌ 列出主题失败: %v\n", err)
    } else {
        fmt.Printf("✅ 成功列出 %d 个主题\n", len(topics))
    }
}
```

### 性能监控工具

```go
func monitorPerformance(client *fluvio.Client) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        start := time.Now()
        
        // 测试发送性能
        _, err := client.Producer().SendString(context.Background(), 
            "test-topic", "test-key", "test-value")
        sendDuration := time.Since(start)
        
        if err != nil {
            log.Printf("📊 发送测试失败: %v", err)
        } else {
            log.Printf("📊 发送延迟: %v", sendDuration)
            
            if sendDuration > 1*time.Second {
                log.Println("⚠️ 发送延迟过高")
            }
        }
        
        // 测试接收性能
        start = time.Now()
        _, err = client.Consumer().Receive(context.Background(), 
            "test-topic", &fluvio.ReceiveOptions{
                Group: "monitor-group",
                MaxMessages: 1,
            })
        receiveDuration := time.Since(start)
        
        if err != nil {
            log.Printf("📊 接收测试失败: %v", err)
        } else {
            log.Printf("📊 接收延迟: %v", receiveDuration)
        }
    }
}
```

---

## 常见配置错误

### 错误 1: 地址配置错误

```go
// ❌ 错误的配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost:50051", 0), // 端口应该是数字
)

// ❌ 错误的配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("", 50051), // 主机名不能为空
)

// ✅ 正确的配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
)
```

### 错误 2: 超时配置错误

```go
// ❌ 超时时间过短
client, err := fluvio.NewClient(
    fluvio.WithTimeout(100*time.Millisecond), // 太短
)

// ❌ 超时时间过长
client, err := fluvio.NewClient(
    fluvio.WithTimeout(10*time.Minute), // 太长
)

// ✅ 合理的超时时间
client, err := fluvio.NewClient(
    fluvio.WithTimeout(30*time.Second), // 适中
)
```

### 错误 3: TLS 配置错误

```go
// ❌ 证书路径错误
client, err := fluvio.NewClient(
    fluvio.WithTLS("wrong-path.crt", "wrong-path.key", "wrong-ca.crt"),
)

// ❌ 混合使用安全和不安全配置
client, err := fluvio.NewClient(
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
    fluvio.WithInsecure(), // 冲突的配置
)

// ✅ 正确的 TLS 配置
client, err := fluvio.NewClient(
    fluvio.WithAddress("secure-server.com", 50051),
    fluvio.WithTLS("client.crt", "client.key", "ca.crt"),
)

// ✅ 正确的不安全配置（仅开发环境）
client, err := fluvio.NewClient(
    fluvio.WithAddress("localhost", 50051),
    fluvio.WithInsecure(),
)
```

---

## 快速诊断清单

当遇到问题时，按以下顺序检查：

### 🔍 基础检查
- [ ] Fluvio 服务器是否运行？
- [ ] 网络连接是否正常？
- [ ] 端口是否正确？
- [ ] 防火墙是否阻止连接？

### 🔧 配置检查
- [ ] 客户端地址配置是否正确？
- [ ] 超时时间是否合理？
- [ ] TLS 证书路径是否正确？
- [ ] 日志级别是否适当？

### 📊 性能检查
- [ ] 网络延迟是否过高？
- [ ] 连接池大小是否足够？
- [ ] 是否使用了批量操作？
- [ ] 缓冲区大小是否合适？

### 🛡️ 安全检查
- [ ] 证书是否有效？
- [ ] 权限是否足够？
- [ ] TLS 配置是否正确？

### 📝 代码检查
- [ ] 错误处理是否完整？
- [ ] 上下文是否正确使用？
- [ ] 资源是否正确释放？
- [ ] 重试逻辑是否合理？

---

## 获取帮助

如果以上解决方案都无法解决您的问题，请：

1. **收集诊断信息**：
   - 启用调试日志
   - 记录错误消息
   - 收集网络诊断信息

2. **查看文档**：
   - [API 参考文档](API.md)
   - [使用指南](GUIDE.md)

3. **社区支持**：
   - GitHub Issues
   - 社区论坛
   - 技术支持邮箱

4. **提供信息**：
   - Go 版本
   - SDK 版本
   - 操作系统
   - 错误日志
   - 复现步骤

记住：详细的错误信息和环境描述有助于快速定位和解决问题！