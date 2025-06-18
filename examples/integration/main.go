package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
	"github.com/iwen-conf/fluvio_grpc_client/types"
)

func main() {
	fmt.Println("=== Fluvio Go SDK 集成测试 ===")

	// 创建客户端
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

	// 运行所有测试
	tests := []struct {
		name string
		fn   func(context.Context, *fluvio.Client) error
	}{
		{"健康检查测试", testHealthCheck},
		{"主题管理测试", testTopicManagement},
		{"消息生产测试", testMessageProduction},
		{"消息消费测试", testMessageConsumption},
		{"过滤消费测试", testFilteredConsumption},
		{"流式消费测试", testStreamConsumption},
		{"主题统计测试", testTopicStats},
		{"消费组管理测试", testConsumerGroups},
		{"SmartModule管理测试", testSmartModuleManagement},
		{"存储管理测试", testStorageManagement},
		{"批量操作测试", testBulkOperations},
		{"错误处理测试", testErrorHandling},
	}

	passed := 0
	failed := 0

	for i, test := range tests {
		fmt.Printf("\n%d. 🧪 %s\n", i+1, test.name)

		err := test.fn(ctx, client)
		if err != nil {
			fmt.Printf("   ❌ 失败: %v\n", err)
			failed++
		} else {
			fmt.Printf("   ✅ 通过\n")
			passed++
		}
	}

	fmt.Printf("\n📊 测试结果: %d 通过, %d 失败, 总计 %d\n", passed, failed, len(tests))

	if failed > 0 {
		fmt.Printf("⚠️  有 %d 个测试失败，请检查Fluvio服务状态\n", failed)
	} else {
		fmt.Println("🎉 所有测试通过!")
	}
}

// 健康检查测试
func testHealthCheck(ctx context.Context, client *fluvio.Client) error {
	err := client.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("健康检查失败: %w", err)
	}

	// 测试Ping
	duration, err := client.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Ping失败: %w", err)
	}

	fmt.Printf("   响应时间: %v\n", duration)
	return nil
}

// 主题管理测试
func testTopicManagement(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 创建主题
	_, err := client.Topic().CreateIfNotExists(ctx, types.CreateTopicOptions{
		Name:              topicName,
		Partitions:        2,
		ReplicationFactor: 1,
		RetentionMs:       3600000, // 1小时
		Config: map[string]string{
			"cleanup.policy": "delete",
		},
	})
	if err != nil {
		return fmt.Errorf("创建主题失败: %w", err)
	}

	// 检查主题是否存在
	exists, err := client.Topic().Exists(ctx, topicName)
	if err != nil {
		return fmt.Errorf("检查主题存在性失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("主题应该存在但检查结果为不存在")
	}

	// 列出主题
	topics, err := client.Topic().List(ctx)
	if err != nil {
		return fmt.Errorf("列出主题失败: %w", err)
	}

	found := false
	for _, topic := range topics.Topics {
		if topic == topicName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("在主题列表中未找到创建的主题")
	}

	// 获取主题详细信息
	detail, err := client.Topic().DescribeTopicDetail(ctx, topicName)
	if err != nil {
		return fmt.Errorf("获取主题详细信息失败: %w", err)
	}

	if detail.Topic != topicName {
		return fmt.Errorf("主题详细信息中的名称不匹配")
	}

	if len(detail.Partitions) != 2 {
		return fmt.Errorf("期望2个分区，实际得到%d个", len(detail.Partitions))
	}

	fmt.Printf("   主题创建成功: %s (分区: %d)\n", topicName, len(detail.Partitions))
	return nil
}

// 消息生产测试
func testMessageProduction(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 单条消息生产
	result, err := client.Producer().Produce(ctx, "测试消息", types.ProduceOptions{
		Topic:     topicName,
		Key:       "test-key",
		MessageID: "test-msg-001",
		Headers: map[string]string{
			"test": "true",
			"type": "integration",
		},
	})
	if err != nil {
		return fmt.Errorf("生产单条消息失败: %w", err)
	}

	if result.MessageID != "test-msg-001" {
		return fmt.Errorf("消息ID不匹配，期望: test-msg-001, 实际: %s", result.MessageID)
	}

	// 批量消息生产
	messages := []types.Message{
		{
			Topic:     topicName,
			Key:       "batch-1",
			Value:     "批量消息1",
			MessageID: "batch-msg-001",
		},
		{
			Topic:     topicName,
			Key:       "batch-2",
			Value:     "批量消息2",
			MessageID: "batch-msg-002",
		},
	}

	batchResult, err := client.Producer().ProduceBatch(ctx, messages)
	if err != nil {
		return fmt.Errorf("批量生产失败: %w", err)
	}

	if len(batchResult.Results) != 2 {
		return fmt.Errorf("期望2个批量结果，实际得到%d个", len(batchResult.Results))
	}

	for i, result := range batchResult.Results {
		if !result.Success {
			return fmt.Errorf("批量消息%d生产失败: %s", i+1, result.Error)
		}
	}

	fmt.Printf("   生产消息成功: 1条单独消息 + 2条批量消息\n")
	return nil
}

// 消息消费测试
func testMessageConsumption(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 消费消息
	messages, err := client.Consumer().Consume(ctx, types.ConsumeOptions{
		Topic:       topicName,
		Group:       "integration-test-group",
		MaxMessages: 10,
		Offset:      0,
	})
	if err != nil {
		return fmt.Errorf("消费消息失败: %w", err)
	}

	if len(messages) < 3 {
		return fmt.Errorf("期望至少3条消息，实际得到%d条", len(messages))
	}

	// 验证消息内容
	foundTestMessage := false
	for _, msg := range messages {
		if msg.MessageID == "test-msg-001" {
			foundTestMessage = true
			if msg.Key != "test-key" || msg.Value != "测试消息" {
				return fmt.Errorf("测试消息内容不匹配")
			}
			if msg.Headers["test"] != "true" {
				return fmt.Errorf("消息头部信息不匹配")
			}
		}
	}

	if !foundTestMessage {
		return fmt.Errorf("未找到测试消息")
	}

	fmt.Printf("   消费消息成功: %d条消息\n", len(messages))
	return nil
}

// 过滤消费测试
func testFilteredConsumption(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 过滤消费：只获取带有特定头部的消息
	result, err := client.Consumer().ConsumeFiltered(ctx, types.FilteredConsumeOptions{
		Topic:       topicName,
		Group:       "filter-test-group",
		MaxMessages: 10,
		Filters: []types.FilterCondition{
			{
				Type:     types.FilterTypeHeader,
				Field:    "test",
				Operator: "eq",
				Value:    "true",
			},
		},
		AndLogic: true,
	})
	if err != nil {
		return fmt.Errorf("过滤消费失败: %w", err)
	}

	if result.FilteredCount == 0 {
		return fmt.Errorf("过滤消费应该返回至少1条消息")
	}

	// 验证过滤结果
	for _, msg := range result.Messages {
		if msg.Headers["test"] != "true" {
			return fmt.Errorf("过滤结果包含不符合条件的消息")
		}
	}

	fmt.Printf("   过滤消费成功: 扫描%d条，过滤出%d条\n", result.TotalScanned, result.FilteredCount)
	return nil
}

// 流式消费测试
func testStreamConsumption(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 创建带超时的上下文
	streamCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 开始流式消费
	stream, err := client.Consumer().ConsumeStream(streamCtx, types.StreamConsumeOptions{
		Topic:        topicName,
		Group:        "stream-test-group",
		MaxBatchSize: 2,
		MaxWaitMs:    1000,
	})
	if err != nil {
		return fmt.Errorf("创建流式消费失败: %w", err)
	}

	messageCount := 0
	for {
		select {
		case msg, ok := <-stream:
			if !ok {
				if messageCount == 0 {
					return fmt.Errorf("流式消费未收到任何消息")
				}
				fmt.Printf("   流式消费成功: %d条消息\n", messageCount)
				return nil
			}

			if msg.Error != nil {
				return fmt.Errorf("流式消费错误: %w", msg.Error)
			}

			messageCount++
			if messageCount >= 3 { // 收到足够消息后退出
				fmt.Printf("   流式消费成功: %d条消息\n", messageCount)
				return nil
			}

		case <-streamCtx.Done():
			if messageCount == 0 {
				return fmt.Errorf("流式消费超时，未收到消息")
			}
			fmt.Printf("   流式消费成功: %d条消息\n", messageCount)
			return nil
		}
	}
}

// 主题统计测试
func testTopicStats(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 获取主题统计信息
	stats, err := client.Topic().GetTopicStats(ctx, types.GetTopicStatsOptions{
		Topic:             topicName,
		IncludePartitions: true,
	})
	if err != nil {
		return fmt.Errorf("获取主题统计失败: %w", err)
	}

	if len(stats.Topics) == 0 {
		return fmt.Errorf("统计结果中没有主题信息")
	}

	topicStats := stats.Topics[0]
	if topicStats.Topic != topicName {
		return fmt.Errorf("主题名称不匹配")
	}

	if topicStats.TotalMessageCount == 0 {
		return fmt.Errorf("主题应该包含消息但统计显示为0")
	}

	if len(topicStats.Partitions) == 0 {
		return fmt.Errorf("应该包含分区统计信息")
	}

	fmt.Printf("   主题统计成功: %d条消息, %d个分区\n",
		topicStats.TotalMessageCount, len(topicStats.Partitions))
	return nil
}

// 消费组管理测试
func testConsumerGroups(ctx context.Context, client *fluvio.Client) error {
	// 列出消费组
	groups, err := client.Admin().ListConsumerGroups(ctx)
	if err != nil {
		return fmt.Errorf("列出消费组失败: %w", err)
	}

	if len(groups.Groups) == 0 {
		return fmt.Errorf("应该至少有一个消费组")
	}

	// 获取第一个消费组的详细信息
	groupName := groups.Groups[0].GroupID
	detail, err := client.Admin().DescribeConsumerGroup(ctx, groupName)
	if err != nil {
		return fmt.Errorf("获取消费组详情失败: %w", err)
	}

	if detail.Group.GroupID != groupName {
		return fmt.Errorf("消费组名称不匹配")
	}

	fmt.Printf("   消费组管理成功: %d个消费组\n", len(groups.Groups))
	return nil
}

// SmartModule管理测试
func testSmartModuleManagement(ctx context.Context, client *fluvio.Client) error {
	// 列出SmartModules
	modules, err := client.Admin().ListSmartModules(ctx)
	if err != nil {
		return fmt.Errorf("列出SmartModules失败: %w", err)
	}

	fmt.Printf("   SmartModule管理成功: %d个模块\n", len(modules.SmartModules))

	// 注意：创建SmartModule需要真实的WASM代码，这里只测试列出功能
	return nil
}

// 存储管理测试
func testStorageManagement(ctx context.Context, client *fluvio.Client) error {
	// 获取存储状态
	status, err := client.Admin().GetStorageStatus(ctx, types.GetStorageStatusOptions{
		IncludeDetails: true,
	})
	if err != nil {
		return fmt.Errorf("获取存储状态失败: %w", err)
	}

	fmt.Printf("   存储状态: 持久化=%v", status.PersistenceEnabled)
	if status.StorageStats != nil {
		fmt.Printf(", 类型=%s", status.StorageStats.StorageType)
	}
	fmt.Println()

	// 获取存储指标
	metrics, err := client.Admin().GetStorageMetrics(ctx, types.GetStorageMetricsOptions{
		IncludeHistory: false,
	})
	if err != nil {
		return fmt.Errorf("获取存储指标失败: %w", err)
	}

	if metrics.CurrentMetrics != nil {
		fmt.Printf("   存储指标: 响应时间=%dms", metrics.CurrentMetrics.ResponseTimeMs)
		if metrics.HealthStatus != nil {
			fmt.Printf(", 健康状态=%s", metrics.HealthStatus.Status)
		}
		fmt.Println()
	}

	return nil
}

// 批量操作测试
func testBulkOperations(ctx context.Context, client *fluvio.Client) error {
	// 创建测试主题
	testTopics := []string{"bulk-test-1", "bulk-test-2"}

	for _, topic := range testTopics {
		_, err := client.Topic().CreateIfNotExists(ctx, types.CreateTopicOptions{
			Name:       topic,
			Partitions: 1,
		})
		if err != nil {
			return fmt.Errorf("创建测试主题失败: %w", err)
		}
	}

	// 等待主题创建完成
	time.Sleep(2 * time.Second)

	// 批量删除
	result, err := client.Admin().BulkDelete(ctx, types.BulkDeleteOptions{
		Topics: testTopics,
		Force:  false,
	})
	if err != nil {
		return fmt.Errorf("批量删除失败: %w", err)
	}

	if result.TotalRequested != int32(len(testTopics)) {
		return fmt.Errorf("批量删除请求数不匹配")
	}

	fmt.Printf("   批量操作成功: %d个请求, %d个成功, %d个失败\n",
		result.TotalRequested, result.SuccessfulDeletes, result.FailedDeletes)
	return nil
}

// 错误处理测试
func testErrorHandling(ctx context.Context, client *fluvio.Client) error {
	// 测试不存在的主题
	_, err := client.Topic().Describe(ctx, "non-existent-topic-12345")
	if err == nil {
		return fmt.Errorf("应该返回错误但没有")
	}

	// 测试无效的消费选项
	_, err = client.Consumer().Consume(ctx, types.ConsumeOptions{
		Topic:       "", // 空主题名
		Group:       "test-group",
		MaxMessages: 10,
	})
	if err == nil {
		return fmt.Errorf("空主题名应该返回错误")
	}

	fmt.Printf("   错误处理正常: 正确捕获了预期错误\n")
	return nil
}
