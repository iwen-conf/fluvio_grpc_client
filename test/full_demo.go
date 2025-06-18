package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("🚀 开始全面测试 Fluvio Go SDK...")
	fmt.Println("📡 连接地址: 101.43.173.154:50051")
	fmt.Println()

	// 创建客户端
	fmt.Println("📝 创建客户端...")
	client, err := fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(5*time.Second, 10*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
		fluvio.WithMaxRetries(3),
	)
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
	}
	defer client.Close()

	fmt.Println("✅ 客户端创建成功")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. 测试连接和健康检查
	fmt.Println("\n=== 1. 连接和健康检查 ===")
	testConnection(client, ctx)

	// 2. 测试主题管理
	fmt.Println("\n=== 2. 主题管理 ===")
	testTopicName := testTopicManagement(client, ctx)

	// 3. 测试消息生产和消费
	if testTopicName != "" {
		fmt.Println("\n=== 3. 消息生产和消费 ===")
		testMessaging(client, ctx, testTopicName)

		// 4. 清理测试主题
		fmt.Println("\n=== 4. 清理 ===")
		cleanupTopic(client, ctx, testTopicName)
	}

	// 5. 测试管理功能
	fmt.Println("\n=== 5. 管理功能 ===")
	testAdminFunctions(client, ctx)

	fmt.Println("\n🎉 全面测试完成！")
}

func testConnection(client *fluvio.Client, ctx context.Context) {
	// 健康检查
	fmt.Println("🔍 执行健康检查...")
	err := client.HealthCheck(ctx)
	if err != nil {
		log.Printf("❌ 健康检查失败: %v", err)
		return
	}
	fmt.Println("✅ 健康检查成功")

	// Ping测试
	fmt.Println("🏓 测试Ping...")
	duration, err := client.Ping(ctx)
	if err != nil {
		log.Printf("❌ Ping失败: %v", err)
		return
	}
	fmt.Printf("✅ Ping成功，延迟: %v\n", duration)
}

func testTopicManagement(client *fluvio.Client, ctx context.Context) string {
	// 列出现有主题
	fmt.Println("📋 获取现有主题列表...")
	topics, err := client.Topic().List(ctx)
	if err != nil {
		log.Printf("❌ 获取主题列表失败: %v", err)
		return ""
	}
	fmt.Printf("✅ 获取主题列表成功，共 %d 个主题\n", len(topics.Topics))
	for i, topic := range topics.Topics {
		if i < 5 { // 只显示前5个
			fmt.Printf("   - %s\n", topic)
		}
	}
	if len(topics.Topics) > 5 {
		fmt.Printf("   ... 还有 %d 个主题\n", len(topics.Topics)-5)
	}

	// 创建测试主题
	testTopicName := fmt.Sprintf("sdk-test-%d", time.Now().Unix())
	fmt.Printf("🆕 创建测试主题: %s\n", testTopicName)

	createResult, err := client.Topic().Create(ctx, fluvio.CreateTopicOptions{
		Name:              testTopicName,
		Partitions:        1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("❌ 创建主题失败: %v", err)
		return ""
	}
	if !createResult.Success {
		log.Printf("❌ 创建主题失败: %s", createResult.Error)
		return ""
	}
	fmt.Println("✅ 主题创建成功")

	// 验证主题是否存在
	fmt.Printf("🔍 验证主题是否存在: %s\n", testTopicName)
	exists, err := client.Topic().Exists(ctx, testTopicName)
	if err != nil {
		log.Printf("❌ 检查主题存在性失败: %v", err)
	} else if exists {
		fmt.Println("✅ 主题存在验证成功")
	} else {
		fmt.Println("❌ 主题不存在")
	}

	// 描述主题
	fmt.Printf("📖 描述主题: %s\n", testTopicName)
	topicInfo, err := client.Topic().Describe(ctx, testTopicName)
	if err != nil {
		log.Printf("❌ 描述主题失败: %v", err)
	} else {
		fmt.Printf("✅ 主题描述成功: %s (分区: %d)\n", topicInfo.Topic.Name, topicInfo.Topic.Partitions)
	}

	return testTopicName
}

func testMessaging(client *fluvio.Client, ctx context.Context, topicName string) {
	// 生产单条消息
	fmt.Println("📤 发送单条测试消息...")
	message1 := fmt.Sprintf("Hello from SDK test at %s", time.Now().Format(time.RFC3339))
	produceResult, err := client.Producer().Produce(ctx, message1, fluvio.ProduceOptions{
		Topic: topicName,
		Key:   "test-key-1",
		Headers: map[string]string{
			"source":    "sdk-test",
			"timestamp": time.Now().Format(time.RFC3339),
			"type":      "single",
		},
	})
	if err != nil {
		log.Printf("❌ 发送消息失败: %v", err)
		return
	}
	fmt.Printf("✅ 消息发送成功，ID: %s\n", produceResult.MessageID)

	// 批量生产消息
	fmt.Println("📤 批量发送测试消息...")
	messages := []fluvio.Message{
		{
			Topic: topicName,
			Key:   "batch-key-1",
			Value: "Batch message 1",
			Headers: map[string]string{
				"source": "sdk-batch-test",
				"index":  "1",
			},
		},
		{
			Topic: topicName,
			Key:   "batch-key-2",
			Value: "Batch message 2",
			Headers: map[string]string{
				"source": "sdk-batch-test",
				"index":  "2",
			},
		},
	}

	batchResult, err := client.Producer().ProduceBatch(ctx, messages)
	if err != nil {
		log.Printf("❌ 批量发送消息失败: %v", err)
	} else {
		fmt.Printf("✅ 批量消息发送成功，总数: %d，成功: %d\n", 
			batchResult.TotalCount, batchResult.TotalCount)
	}

	// 等待一下让消息被处理
	time.Sleep(2 * time.Second)

	// 消费消息
	fmt.Println("📥 消费测试消息...")
	consumedMessages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
		Topic:       topicName,
		Group:       "sdk-test-group",
		MaxMessages: 10,
	})
	if err != nil {
		log.Printf("❌ 消费消息失败: %v", err)
		return
	}
	fmt.Printf("✅ 消费成功，收到 %d 条消息\n", len(consumedMessages))
	for i, msg := range consumedMessages {
		if i < 5 { // 只显示前5条
			fmt.Printf("   消息%d: [%s] %s\n", i+1, msg.Key, msg.Value)
			if len(msg.Headers) > 0 {
				fmt.Printf("     头部: %v\n", msg.Headers)
			}
		}
	}
	if len(consumedMessages) > 5 {
		fmt.Printf("   ... 还有 %d 条消息\n", len(consumedMessages)-5)
	}
}

func cleanupTopic(client *fluvio.Client, ctx context.Context, topicName string) {
	fmt.Printf("🗑️ 删除测试主题: %s\n", topicName)
	deleteResult, err := client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{
		Name: topicName,
	})
	if err != nil {
		log.Printf("❌ 删除主题失败: %v", err)
	} else if deleteResult.Success {
		fmt.Println("✅ 主题删除成功")
	} else {
		log.Printf("❌ 删除主题失败: %s", deleteResult.Error)
	}
}

func testAdminFunctions(client *fluvio.Client, ctx context.Context) {
	// 获取集群信息
	fmt.Println("🏢 获取集群信息...")
	clusterInfo, err := client.Admin().DescribeCluster(ctx)
	if err != nil {
		log.Printf("❌ 获取集群信息失败: %v", err)
	} else {
		fmt.Printf("✅ 集群信息获取成功: 状态=%s, 控制器ID=%d\n", 
			clusterInfo.Cluster.Status, clusterInfo.Cluster.ControllerID)
	}

	// 获取Broker列表
	fmt.Println("🖥️ 获取Broker列表...")
	brokers, err := client.Admin().ListBrokers(ctx)
	if err != nil {
		log.Printf("❌ 获取Broker列表失败: %v", err)
	} else {
		fmt.Printf("✅ Broker列表获取成功，共 %d 个Broker\n", len(brokers.Brokers))
		for i, broker := range brokers.Brokers {
			if i < 3 { // 只显示前3个
				fmt.Printf("   - Broker %d: %s (%s)\n", broker.ID, broker.Addr, broker.Status)
			}
		}
	}

	// 获取消费者组列表
	fmt.Println("👥 获取消费者组列表...")
	groups, err := client.Admin().ListConsumerGroups(ctx)
	if err != nil {
		log.Printf("❌ 获取消费者组列表失败: %v", err)
	} else {
		fmt.Printf("✅ 消费者组列表获取成功，共 %d 个组\n", len(groups.Groups))
		for i, group := range groups.Groups {
			if i < 3 { // 只显示前3个
				fmt.Printf("   - 组: %s\n", group.GroupID)
			}
		}
	}
}
