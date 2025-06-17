package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	// 集成测试场景
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

	// 运行集成测试
	tests := []struct {
		name string
		fn   func(context.Context, *fluvio.Client) error
	}{
		{"连接测试", testConnection},
		{"主题管理测试", testTopicManagement},
		{"消息生产消费测试", testProduceConsume},
		{"批量操作测试", testBatchOperations},
		{"消费组测试", testConsumerGroups},
		{"SmartModule测试", testSmartModules},
		{"管理功能测试", testAdminFunctions},
		{"错误处理测试", testErrorHandling},
	}

	passed := 0
	failed := 0

	for _, test := range tests {
		fmt.Printf("\n--- %s ---\n", test.name)
		err := test.fn(ctx, client)
		if err != nil {
			fmt.Printf("❌ %s 失败: %v\n", test.name, err)
			failed++
		} else {
			fmt.Printf("✅ %s 通过\n", test.name)
			passed++
		}
	}

	fmt.Printf("\n=== 测试结果 ===\n")
	fmt.Printf("通过: %d\n", passed)
	fmt.Printf("失败: %d\n", failed)
	fmt.Printf("总计: %d\n", passed+failed)

	if failed == 0 {
		fmt.Println("🎉 所有测试通过!")
	} else {
		fmt.Printf("⚠️  %d 个测试失败\n", failed)
	}
}

func testConnection(ctx context.Context, client *fluvio.Client) error {
	// 测试健康检查
	err := client.HealthCheck(ctx)
	if err != nil {
		return fmt.Errorf("健康检查失败: %w", err)
	}
	fmt.Println("✓ 健康检查通过")

	// 测试Ping
	duration, err := client.Ping(ctx)
	if err != nil {
		return fmt.Errorf("Ping失败: %w", err)
	}
	fmt.Printf("✓ Ping成功 (延迟: %v)\n", duration)

	return nil
}

func testTopicManagement(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-test-topic"

	// 创建主题
	_, err := client.Topic().Create(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 2,
	})
	if err != nil {
		return fmt.Errorf("创建主题失败: %w", err)
	}
	fmt.Printf("✓ 主题 '%s' 创建成功\n", topicName)

	// 检查主题是否存在
	exists, err := client.Topic().Exists(ctx, topicName)
	if err != nil {
		return fmt.Errorf("检查主题存在性失败: %w", err)
	}
	if !exists {
		return fmt.Errorf("主题应该存在但检查结果为不存在")
	}
	fmt.Println("✓ 主题存在性检查通过")

	// 列出主题
	result, err := client.Topic().List(ctx)
	if err != nil {
		return fmt.Errorf("列出主题失败: %w", err)
	}

	found := false
	for _, topic := range result.Topics {
		if topic == topicName {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("在主题列表中未找到创建的主题")
	}
	fmt.Printf("✓ 主题列表包含创建的主题 (共 %d 个主题)\n", len(result.Topics))

	// 删除主题
	_, err = client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{
		Name: topicName,
	})
	if err != nil {
		return fmt.Errorf("删除主题失败: %w", err)
	}
	fmt.Printf("✓ 主题 '%s' 删除成功\n", topicName)

	return nil
}

func testProduceConsume(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-produce-consume"

	// 创建测试主题
	_, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 1,
	})
	if err != nil {
		return fmt.Errorf("创建测试主题失败: %w", err)
	}

	// 生产消息
	testMessage := "集成测试消息 - " + time.Now().Format(time.RFC3339)
	result, err := client.Producer().Produce(ctx, testMessage, fluvio.ProduceOptions{
		Topic: topicName,
		Key:   "integration-test-key",
	})
	if err != nil {
		return fmt.Errorf("生产消息失败: %w", err)
	}
	fmt.Printf("✓ 消息生产成功: %s\n", result.MessageID)

	// 等待消息处理
	time.Sleep(1 * time.Second)

	// 消费消息
	messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
		Topic:       topicName,
		Group:       "integration-test-group",
		MaxMessages: 1,
		Offset:      0,
	})
	if err != nil {
		return fmt.Errorf("消费消息失败: %w", err)
	}

	if len(messages) == 0 {
		return fmt.Errorf("未消费到任何消息")
	}

	if messages[0].Value != testMessage {
		return fmt.Errorf("消费到的消息内容不匹配: 期望 '%s', 实际 '%s'",
			testMessage, messages[0].Value)
	}

	fmt.Printf("✓ 消息消费成功: %s\n", messages[0].Value)

	// 清理
	client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{Name: topicName})

	return nil
}

func testBatchOperations(ctx context.Context, client *fluvio.Client) error {
	topicName := "integration-batch-test"

	// 创建测试主题
	_, err := client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 1,
	})
	if err != nil {
		return fmt.Errorf("创建测试主题失败: %w", err)
	}

	// 批量生产消息
	messages := make([]fluvio.Message, 5)
	for i := 0; i < 5; i++ {
		messages[i] = fluvio.Message{
			Topic: topicName,
			Key:   fmt.Sprintf("batch-key-%d", i),
			Value: fmt.Sprintf("批量消息 #%d", i+1),
		}
	}

	batchResult, err := client.Producer().ProduceBatch(ctx, messages)
	if err != nil {
		return fmt.Errorf("批量生产失败: %w", err)
	}

	successCount := 0
	for _, result := range batchResult.Results {
		if result.Success {
			successCount++
		}
	}

	if successCount != len(messages) {
		return fmt.Errorf("批量生产部分失败: %d/%d 成功", successCount, len(messages))
	}

	fmt.Printf("✓ 批量生产成功: %d 条消息\n", successCount)

	// 清理
	client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{Name: topicName})

	return nil
}

func testConsumerGroups(ctx context.Context, client *fluvio.Client) error {
	// 列出消费组
	groups, err := client.Admin().ListConsumerGroups(ctx)
	if err != nil {
		return fmt.Errorf("列出消费组失败: %w", err)
	}
	fmt.Printf("✓ 消费组列表获取成功 (共 %d 个组)\n", len(groups.Groups))

	return nil
}

func testSmartModules(ctx context.Context, client *fluvio.Client) error {
	// 列出SmartModules
	modules, err := client.Admin().ListSmartModules(ctx)
	if err != nil {
		return fmt.Errorf("列出SmartModules失败: %w", err)
	}
	fmt.Printf("✓ SmartModules列表获取成功 (共 %d 个模块)\n", len(modules.SmartModules))

	return nil
}

func testAdminFunctions(ctx context.Context, client *fluvio.Client) error {
	// 获取集群信息
	cluster, err := client.Admin().DescribeCluster(ctx)
	if err != nil {
		return fmt.Errorf("获取集群信息失败: %w", err)
	}
	fmt.Printf("✓ 集群信息获取成功: 状态=%s, 控制器ID=%d\n",
		cluster.Cluster.Status, cluster.Cluster.ControllerID)

	// 列出Brokers
	brokers, err := client.Admin().ListBrokers(ctx)
	if err != nil {
		return fmt.Errorf("列出Brokers失败: %w", err)
	}
	fmt.Printf("✓ Brokers列表获取成功 (共 %d 个Broker)\n", len(brokers.Brokers))

	return nil
}

func testErrorHandling(ctx context.Context, client *fluvio.Client) error {
	// 测试操作不存在的主题
	_, err := client.Producer().Produce(ctx, "测试消息", fluvio.ProduceOptions{
		Topic: "non-existent-topic-12345",
	})
	if err == nil {
		return fmt.Errorf("操作不存在的主题应该返回错误")
	}
	fmt.Printf("✓ 错误处理正确: %v\n", err)

	return nil
}
