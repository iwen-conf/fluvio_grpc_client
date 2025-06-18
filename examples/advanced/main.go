package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
	"github.com/iwen-conf/fluvio_grpc_client/types"
)

func main() {
	fmt.Println("=== Fluvio Go SDK 高级示例 ===")

	// 创建高性能客户端
	client, err := fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(5*time.Second, 60*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
		fluvio.WithMaxRetries(5),
		fluvio.WithPoolSize(10), // 大连接池
	)
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 1. 演示过滤消费功能（新功能）
	fmt.Println("\n🔍 演示过滤消费功能...")
	demonstrateFilteredConsume(ctx, client)

	// 2. 演示流式消费增强功能
	fmt.Println("\n📡 演示流式消费增强功能...")
	demonstrateEnhancedStreamConsume(ctx, client)

	// 3. 演示SmartModule管理（新功能）
	fmt.Println("\n🧠 演示SmartModule管理...")
	demonstrateSmartModuleManagement(ctx, client)

	// 4. 演示存储管理功能（新功能）
	fmt.Println("\n💾 演示存储管理功能...")
	demonstrateStorageManagement(ctx, client)

	// 5. 演示批量删除功能（新功能）
	fmt.Println("\n🗑️  演示批量删除功能...")
	demonstrateBulkDelete(ctx, client)

	// 6. 演示并发处理
	fmt.Println("\n⚡ 演示并发处理...")
	demonstrateConcurrentProcessing(ctx, client)

	fmt.Println("\n🎉 高级示例完成!")
}

// 演示过滤消费功能
func demonstrateFilteredConsume(ctx context.Context, client *fluvio.Client) {
	topicName := "advanced-filter-topic"

	// 创建主题
	_, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 1,
	})
	if err != nil {
		log.Printf("创建主题失败: %v", err)
		return
	}

	// 生产一些测试消息
	testMessages := []fluvio.Message{
		{Topic: topicName, Key: "user-1", Value: "用户登录", Headers: map[string]string{"event": "login", "level": "info"}},
		{Topic: topicName, Key: "user-2", Value: "用户注册", Headers: map[string]string{"event": "register", "level": "info"}},
		{Topic: topicName, Key: "user-1", Value: "支付失败", Headers: map[string]string{"event": "payment", "level": "error"}},
		{Topic: topicName, Key: "user-3", Value: "用户登出", Headers: map[string]string{"event": "logout", "level": "info"}},
		{Topic: topicName, Key: "user-2", Value: "订单创建", Headers: map[string]string{"event": "order", "level": "info"}},
	}

	_, err = client.Producer().ProduceBatch(ctx, testMessages)
	if err != nil {
		log.Printf("生产测试消息失败: %v", err)
		return
	}

	// 过滤消费：只获取错误级别的消息
	fmt.Println("  🔍 过滤消费：只获取错误级别的消息")
	result, err := client.Consumer().ConsumeFiltered(ctx, types.FilteredConsumeOptions{
		Topic:       topicName,
		Group:       "filter-group-1",
		MaxMessages: 10,
		Filters: []types.FilterCondition{
			{
				Type:     types.FilterTypeHeader,
				Field:    "level",
				Operator: "eq",
				Value:    "error",
			},
		},
		AndLogic: true,
	})
	if err != nil {
		log.Printf("过滤消费失败: %v", err)
		return
	}

	fmt.Printf("  ✅ 过滤结果: 扫描了 %d 条消息，过滤出 %d 条消息\n",
		result.TotalScanned, result.FilteredCount)
	for i, msg := range result.Messages {
		fmt.Printf("    %d. [%s] %s (Headers: %v)\n", i+1, msg.Key, msg.Value, msg.Headers)
	}

	// 过滤消费：获取特定用户的消息
	fmt.Println("  🔍 过滤消费：只获取user-1的消息")
	result2, err := client.Consumer().ConsumeFiltered(ctx, types.FilteredConsumeOptions{
		Topic:       topicName,
		Group:       "filter-group-2",
		MaxMessages: 10,
		Filters: []types.FilterCondition{
			{
				Type:     types.FilterTypeKey,
				Operator: "eq",
				Value:    "user-1",
			},
		},
		AndLogic: true,
	})
	if err != nil {
		log.Printf("过滤消费失败: %v", err)
		return
	}

	fmt.Printf("  ✅ 过滤结果: 扫描了 %d 条消息，过滤出 %d 条消息\n",
		result2.TotalScanned, result2.FilteredCount)
	for i, msg := range result2.Messages {
		fmt.Printf("    %d. [%s] %s\n", i+1, msg.Key, msg.Value)
	}
}

// 演示流式消费增强功能
func demonstrateEnhancedStreamConsume(ctx context.Context, client *fluvio.Client) {
	topicName := "advanced-stream-topic"

	// 创建主题
	_, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 1,
	})
	if err != nil {
		log.Printf("创建主题失败: %v", err)
		return
	}

	// 启动生产者协程
	go func() {
		for i := 0; i < 10; i++ {
			_, err := client.Producer().Produce(ctx, fmt.Sprintf("流式消息 %d", i+1), fluvio.ProduceOptions{
				Topic:     topicName,
				Key:       fmt.Sprintf("stream-key-%d", i+1),
				MessageID: fmt.Sprintf("stream-msg-%03d", i+1),
			})
			if err != nil {
				log.Printf("生产流式消息失败: %v", err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// 流式消费（使用新的批次控制参数）
	streamCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	stream, err := client.Consumer().ConsumeStream(streamCtx, types.StreamConsumeOptions{
		Topic:        topicName,
		Group:        "stream-group",
		MaxBatchSize: 3,    // 新功能：每批最多3条消息
		MaxWaitMs:    1000, // 新功能：最多等待1秒
	})
	if err != nil {
		log.Printf("创建流式消费失败: %v", err)
		return
	}

	fmt.Println("  📡 开始流式消费（批次大小=3，等待时间=1秒）...")
	messageCount := 0
	batchCount := 0

	for {
		select {
		case msg, ok := <-stream:
			if !ok {
				fmt.Printf("  ✅ 流式消费结束，共收到 %d 条消息，%d 个批次\n", messageCount, batchCount)
				return
			}

			if msg.Error != nil {
				log.Printf("流式消费错误: %v", msg.Error)
				continue
			}

			messageCount++
			if messageCount%3 == 1 {
				batchCount++
				fmt.Printf("  📦 批次 %d:\n", batchCount)
			}
			fmt.Printf("    %d. [%s] %s (ID: %s)\n",
				messageCount, msg.Message.Key, msg.Message.Value, msg.Message.MessageID)

		case <-streamCtx.Done():
			fmt.Printf("  ⏰ 流式消费超时，共收到 %d 条消息，%d 个批次\n", messageCount, batchCount)
			return
		}
	}
}

// 演示SmartModule管理
func demonstrateSmartModuleManagement(ctx context.Context, client *fluvio.Client) {
	// 列出现有的SmartModules
	modules, err := client.Admin().ListSmartModules(ctx)
	if err != nil {
		log.Printf("列出SmartModules失败: %v", err)
		return
	}

	fmt.Printf("  📋 当前SmartModules数量: %d\n", len(modules.SmartModules))
	for i, module := range modules.SmartModules {
		fmt.Printf("    %d. %s (版本: %s) - %s\n",
			i+1, module.Name, module.Version, module.Description)
	}

	// 创建一个示例SmartModule（注意：这需要实际的WASM代码）
	fmt.Println("  🧠 创建示例SmartModule...")

	// 示例SmartModule规格
	spec := &types.SmartModuleSpec{
		Name:        "example-filter",
		InputKind:   types.SmartModuleInputStream,
		OutputKind:  types.SmartModuleOutputStream,
		Description: "示例过滤器SmartModule",
		Version:     "1.0.0",
		Parameters: []*types.SmartModuleParameter{
			{
				Name:        "filter_key",
				Description: "要过滤的键值",
				Optional:    false,
			},
		},
	}

	// 注意：这里使用空的WASM代码作为示例
	// 在实际使用中，需要提供真实的WASM字节码
	createResult, err := client.Admin().CreateSmartModule(ctx, types.CreateSmartModuleOptions{
		Spec:     spec,
		WasmCode: []byte{}, // 实际使用时需要真实的WASM代码
	})
	if err != nil {
		log.Printf("⚠️  创建SmartModule失败（预期的，因为没有真实WASM代码）: %v", err)
	} else {
		fmt.Printf("  ✅ SmartModule创建成功: %+v\n", createResult)
	}
}

// 演示存储管理功能
func demonstrateStorageManagement(ctx context.Context, client *fluvio.Client) {
	// 获取存储状态
	fmt.Println("  💾 获取存储状态...")
	status, err := client.Admin().GetStorageStatus(ctx, types.GetStorageStatusOptions{
		IncludeDetails: true,
	})
	if err != nil {
		log.Printf("获取存储状态失败: %v", err)
		return
	}

	fmt.Printf("  ✅ 存储状态:\n")
	fmt.Printf("    - 持久化启用: %v\n", status.PersistenceEnabled)
	if status.StorageStats != nil {
		stats := status.StorageStats
		fmt.Printf("    - 存储类型: %s\n", stats.StorageType)
		fmt.Printf("    - 连接状态: %s\n", stats.ConnectionStatus)
		fmt.Printf("    - 消费组数量: %d\n", stats.ConsumerGroups)
		fmt.Printf("    - 消费偏移量数量: %d\n", stats.ConsumerOffsets)
		fmt.Printf("    - SmartModule数量: %d\n", stats.SmartModules)

		if stats.ConnectionStats != nil {
			fmt.Printf("    - 当前连接数: %d\n", stats.ConnectionStats.CurrentConnections)
			fmt.Printf("    - 可用连接数: %d\n", stats.ConnectionStats.AvailableConnections)
		}

		if stats.DatabaseInfo != nil {
			fmt.Printf("    - 数据库: %s\n", stats.DatabaseInfo.Name)
			fmt.Printf("    - 集合数: %d\n", stats.DatabaseInfo.Collections)
			fmt.Printf("    - 数据大小: %d bytes\n", stats.DatabaseInfo.DataSize)
		}
	}

	// 获取存储指标
	fmt.Println("  📊 获取存储指标...")
	metrics, err := client.Admin().GetStorageMetrics(ctx, types.GetStorageMetricsOptions{
		IncludeHistory: false,
		HistoryLimit:   10,
	})
	if err != nil {
		log.Printf("获取存储指标失败: %v", err)
		return
	}

	fmt.Printf("  ✅ 存储指标:\n")
	if metrics.CurrentMetrics != nil {
		m := metrics.CurrentMetrics
		fmt.Printf("    - 存储类型: %s\n", m.StorageType)
		fmt.Printf("    - 响应时间: %d ms\n", m.ResponseTimeMs)
		fmt.Printf("    - 每秒操作数: %.2f\n", m.OperationsPerSecond)
		fmt.Printf("    - 错误率: %.2f%%\n", m.ErrorRate*100)
		fmt.Printf("    - 连接池使用率: %.2f%%\n", m.ConnectionPoolUsage*100)
		fmt.Printf("    - 内存使用: %d MB\n", m.MemoryUsageMB)
		fmt.Printf("    - 磁盘使用: %d MB\n", m.DiskUsageMB)
	}

	if metrics.HealthStatus != nil {
		fmt.Printf("    - 健康状态: %s\n", metrics.HealthStatus.Status)
		if metrics.HealthStatus.ErrorMessage != "" {
			fmt.Printf("    - 错误信息: %s\n", metrics.HealthStatus.ErrorMessage)
		}
	}

	if len(metrics.Alerts) > 0 {
		fmt.Printf("    - 告警: %v\n", metrics.Alerts)
	}
}

// 演示批量删除功能
func demonstrateBulkDelete(ctx context.Context, client *fluvio.Client) {
	// 创建一些测试资源
	testTopics := []string{"bulk-test-topic-1", "bulk-test-topic-2", "bulk-test-topic-3"}

	fmt.Println("  🏗️  创建测试主题...")
	for _, topic := range testTopics {
		_, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
			Name:       topic,
			Partitions: 1,
		})
		if err != nil {
			log.Printf("创建测试主题 %s 失败: %v", topic, err)
		}
	}

	// 等待一下确保主题创建完成
	time.Sleep(2 * time.Second)

	// 批量删除
	fmt.Println("  🗑️  执行批量删除...")
	result, err := client.Admin().BulkDelete(ctx, types.BulkDeleteOptions{
		Topics: testTopics,
		Force:  false, // 非强制删除
	})
	if err != nil {
		log.Printf("批量删除失败: %v", err)
		return
	}

	fmt.Printf("  ✅ 批量删除结果:\n")
	fmt.Printf("    - 总请求数: %d\n", result.TotalRequested)
	fmt.Printf("    - 成功删除: %d\n", result.SuccessfulDeletes)
	fmt.Printf("    - 删除失败: %d\n", result.FailedDeletes)

	for i, itemResult := range result.Results {
		status := "✅"
		if !itemResult.Success {
			status = "❌"
		}
		fmt.Printf("    %d. %s %s (%s)", i+1, status, itemResult.Name, itemResult.Type)
		if itemResult.Error != "" {
			fmt.Printf(" - 错误: %s", itemResult.Error)
		}
		fmt.Println()
	}
}

// 演示并发处理
func demonstrateConcurrentProcessing(ctx context.Context, client *fluvio.Client) {
	topicName := "concurrent-topic"

	// 创建主题
	_, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 3, // 多分区支持并发
	})
	if err != nil {
		log.Printf("创建并发测试主题失败: %v", err)
		return
	}

	var wg sync.WaitGroup

	// 启动多个并发生产者
	fmt.Println("  ⚡ 启动并发生产者...")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				message := fmt.Sprintf("并发消息 P%d-M%d", producerID, j+1)
				_, err := client.Producer().Produce(ctx, message, fluvio.ProduceOptions{
					Topic:     topicName,
					Key:       fmt.Sprintf("producer-%d-msg-%d", producerID, j+1),
					MessageID: fmt.Sprintf("concurrent-p%d-m%d", producerID, j+1),
					Headers: map[string]string{
						"producer_id": fmt.Sprintf("%d", producerID),
						"message_seq": fmt.Sprintf("%d", j+1),
						"timestamp":   time.Now().Format(time.RFC3339),
					},
				})
				if err != nil {
					log.Printf("生产者 %d 消息 %d 失败: %v", producerID, j+1, err)
				}
				time.Sleep(100 * time.Millisecond)
			}
			fmt.Printf("    ✅ 生产者 %d 完成\n", producerID)
		}(i)
	}

	// 启动多个并发消费者
	fmt.Println("  ⚡ 启动并发消费者...")
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()

			// 等待一下让生产者先产生一些消息
			time.Sleep(1 * time.Second)

			messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
				Topic:       topicName,
				Group:       fmt.Sprintf("concurrent-group-%d", consumerID),
				MaxMessages: 10,
			})
			if err != nil {
				log.Printf("消费者 %d 失败: %v", consumerID, err)
				return
			}

			fmt.Printf("    ✅ 消费者 %d 收到 %d 条消息:\n", consumerID, len(messages))
			for j, msg := range messages {
				fmt.Printf("      %d. [%s] %s (ID: %s, Producer: %s)\n",
					j+1, msg.Key, msg.Value, msg.MessageID, msg.Headers["producer_id"])
			}
		}(i)
	}

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("  ✅ 并发处理完成")
}
