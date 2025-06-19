package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("🚀 测试新的Fluvio Go SDK v2.0 API...")
	fmt.Println("📡 连接地址: 101.43.173.154:50051")
	fmt.Println()

	// 创建客户端
	fmt.Println("📝 创建客户端...")
	client, err := fluvio.NewClient(
		fluvio.WithAddress("101.43.173.154", 50051),
		fluvio.WithTimeout(30*time.Second),
		fluvio.WithRetry(3, time.Second),
		fluvio.WithLogLevel(fluvio.LogLevelInfo),
		fluvio.WithConnectionPool(5, 5*time.Minute),
	)
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
	}
	defer client.Close()

	fmt.Printf("✅ 客户端创建成功，版本: %s\n", fluvio.Version())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 连接到服务器
	fmt.Println("🔗 连接到服务器...")
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("❌ 连接失败: %v", err)
	}
	fmt.Println("✅ 连接成功")

	// 健康检查
	fmt.Println("🔍 执行健康检查...")
	if err := client.HealthCheck(ctx); err != nil {
		log.Fatalf("❌ 健康检查失败: %v", err)
	}
	fmt.Println("✅ 健康检查成功")

	// Ping测试
	fmt.Println("🏓 测试Ping...")
	duration, err := client.Ping(ctx)
	if err != nil {
		log.Printf("❌ Ping失败: %v", err)
	} else {
		fmt.Printf("✅ Ping成功，延迟: %v\n", duration)
	}

	// 测试主题管理
	fmt.Println("\n=== 主题管理测试 ===")
	testTopicManagement(client, ctx)

	// 测试消息生产和消费
	fmt.Println("\n=== 消息生产和消费测试 ===")
	testMessaging(client, ctx)

	// 测试管理功能
	fmt.Println("\n=== 管理功能测试 ===")
	testAdminFunctions(client, ctx)

	fmt.Println("\n🎉 新API测试完成！")
}

func testTopicManagement(client *fluvio.Client, ctx context.Context) {
	// 列出现有主题
	fmt.Println("📋 获取主题列表...")
	topics, err := client.Topics().List(ctx)
	if err != nil {
		log.Printf("❌ 获取主题列表失败: %v", err)
		return
	}
	fmt.Printf("✅ 获取主题列表成功，共 %d 个主题\n", len(topics))
	for i, topic := range topics {
		if i < 3 { // 只显示前3个
			fmt.Printf("   - %s\n", topic)
		}
	}
	if len(topics) > 3 {
		fmt.Printf("   ... 还有 %d 个主题\n", len(topics)-3)
	}

	// 创建测试主题
	testTopicName := fmt.Sprintf("new-api-test-%d", time.Now().Unix())
	fmt.Printf("🆕 创建测试主题: %s\n", testTopicName)

	err = client.Topics().Create(ctx, testTopicName, &fluvio.CreateTopicOptions{
		Partitions:        1,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Printf("❌ 创建主题失败: %v", err)
		return
	}
	fmt.Println("✅ 主题创建成功")

	// 获取主题信息
	fmt.Printf("📖 获取主题信息: %s\n", testTopicName)
	info, err := client.Topics().Info(ctx, testTopicName)
	if err != nil {
		log.Printf("❌ 获取主题信息失败: %v", err)
	} else {
		fmt.Printf("✅ 主题信息: %s (分区: %d)\n", info.Name, info.Partitions)
	}

	// 检查主题是否存在
	fmt.Printf("🔍 检查主题是否存在: %s\n", testTopicName)
	exists, err := client.Topics().Exists(ctx, testTopicName)
	if err != nil {
		log.Printf("❌ 检查主题存在性失败: %v", err)
	} else {
		fmt.Printf("✅ 主题存在性检查: %v\n", exists)
	}

	// 清理：删除测试主题
	fmt.Printf("🗑️ 删除测试主题: %s\n", testTopicName)
	err = client.Topics().Delete(ctx, testTopicName)
	if err != nil {
		log.Printf("❌ 删除主题失败: %v", err)
	} else {
		fmt.Println("✅ 主题删除成功")
	}
}

func testMessaging(client *fluvio.Client, ctx context.Context) {
	topicName := "new-api-messaging-test"

	// 确保主题存在
	created, err := client.Topics().CreateIfNotExists(ctx, topicName, &fluvio.CreateTopicOptions{
		Partitions: 1,
	})
	if err != nil {
		log.Printf("❌ 创建主题失败: %v", err)
		return
	}
	if created {
		fmt.Printf("✅ 创建主题: %s\n", topicName)
	} else {
		fmt.Printf("ℹ️ 主题已存在: %s\n", topicName)
	}

	// 发送单条消息
	fmt.Println("📤 发送单条消息...")
	result, err := client.Producer().Send(ctx, topicName, &fluvio.Message{
		Key:   "test-key-1",
		Value: []byte("Hello from new API!"),
		Headers: map[string]string{
			"source":    "new-api-test",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		log.Printf("❌ 发送消息失败: %v", err)
	} else {
		fmt.Printf("✅ 消息发送成功，ID: %s\n", result.MessageID)
	}

	// 使用便捷方法发送字符串消息
	fmt.Println("📤 发送字符串消息...")
	result2, err := client.Producer().SendString(ctx, topicName, "string-key", "Hello String!")
	if err != nil {
		log.Printf("❌ 发送字符串消息失败: %v", err)
	} else {
		fmt.Printf("✅ 字符串消息发送成功，ID: %s\n", result2.MessageID)
	}

	// 批量发送消息
	fmt.Println("📤 批量发送消息...")
	messages := []*fluvio.Message{
		{Key: "batch-1", Value: []byte("Batch message 1")},
		{Key: "batch-2", Value: []byte("Batch message 2")},
		{Key: "batch-3", Value: []byte("Batch message 3")},
	}
	batchResult, err := client.Producer().SendBatch(ctx, topicName, messages)
	if err != nil {
		log.Printf("❌ 批量发送失败: %v", err)
	} else {
		fmt.Printf("✅ 批量发送成功，成功: %d，失败: %d\n",
			batchResult.SuccessCount, batchResult.FailureCount)
	}

	// 等待一下让消息被处理
	time.Sleep(2 * time.Second)

	// 消费消息
	fmt.Println("📥 消费消息...")
	consumedMessages, err := client.Consumer().Receive(ctx, topicName, &fluvio.ReceiveOptions{
		Group:       "new-api-test-group",
		MaxMessages: 10,
	})
	if err != nil {
		log.Printf("❌ 消费消息失败: %v", err)
	} else {
		fmt.Printf("✅ 消费成功，收到 %d 条消息\n", len(consumedMessages))
		for i, msg := range consumedMessages {
			if i < 3 { // 只显示前3条
				fmt.Printf("   消息%d: [%s] %s\n", i+1, msg.Key, string(msg.Value))
			}
		}
		if len(consumedMessages) > 3 {
			fmt.Printf("   ... 还有 %d 条消息\n", len(consumedMessages)-3)
		}
	}

	// 接收单条消息
	fmt.Println("📥 接收单条消息...")
	singleMsg, err := client.Consumer().ReceiveOne(ctx, topicName, "single-msg-group")
	if err != nil {
		log.Printf("❌ 接收单条消息失败: %v", err)
	} else if singleMsg != nil {
		fmt.Printf("✅ 接收单条消息成功: [%s] %s\n", singleMsg.Key, string(singleMsg.Value))
	} else {
		fmt.Println("ℹ️ 没有可用的消息")
	}
}

func testAdminFunctions(client *fluvio.Client, ctx context.Context) {
	// 获取集群信息
	fmt.Println("🏢 获取集群信息...")
	clusterInfo, err := client.Admin().ClusterInfo(ctx)
	if err != nil {
		log.Printf("❌ 获取集群信息失败: %v", err)
	} else {
		fmt.Printf("✅ 集群信息: ID=%s, 状态=%s, 控制器ID=%d\n",
			clusterInfo.ID, clusterInfo.Status, clusterInfo.ControllerID)
	}

	// 获取Broker列表
	fmt.Println("🖥️ 获取Broker列表...")
	brokers, err := client.Admin().Brokers(ctx)
	if err != nil {
		log.Printf("❌ 获取Broker列表失败: %v", err)
	} else {
		fmt.Printf("✅ Broker列表，共 %d 个Broker\n", len(brokers))
		for i, broker := range brokers {
			if i < 2 { // 只显示前2个
				fmt.Printf("   - Broker %d: %s:%d (%s)\n", broker.ID, broker.Host, broker.Port, broker.Status)
			}
		}
	}

	// 获取消费者组列表
	fmt.Println("👥 获取消费者组列表...")
	groups, err := client.Admin().ConsumerGroups(ctx)
	if err != nil {
		log.Printf("❌ 获取消费者组列表失败: %v", err)
	} else {
		fmt.Printf("✅ 消费者组列表，共 %d 个组\n", len(groups))
		for i, group := range groups {
			if i < 3 { // 只显示前3个
				fmt.Printf("   - 组: %s (%s)\n", group.GroupID, group.State)
			}
		}
	}

	// 获取SmartModule列表
	fmt.Println("🧠 获取SmartModule列表...")
	modules, err := client.Admin().SmartModules().List(ctx)
	if err != nil {
		log.Printf("❌ 获取SmartModule列表失败: %v", err)
	} else {
		fmt.Printf("✅ SmartModule列表，共 %d 个模块\n", len(modules))
		for i, module := range modules {
			if i < 2 { // 只显示前2个
				fmt.Printf("   - 模块: %s (%s)\n", module.Name, module.Version)
			}
		}
	}
}
