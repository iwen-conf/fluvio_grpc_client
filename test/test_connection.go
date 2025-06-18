package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("🚀 开始测试 Fluvio Go SDK...")
	fmt.Println("📡 连接地址: 101.43.173.154:50051")
	fmt.Println()

	// 测试: 使用向后兼容API
	fmt.Println("=== 测试: 向后兼容API ===")
	testOldAPI()

	fmt.Println()
	fmt.Println("✅ 测试完成！")
}

func testOldAPI() {
	fmt.Println("📝 创建客户端（旧API）...")

	// 使用旧API创建客户端
	client, err := fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(5*time.Second, 10*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
		fluvio.WithMaxRetries(3),
	)
	if err != nil {
		log.Printf("❌ 创建客户端失败: %v", err)
		return
	}
	defer client.Close()

	fmt.Println("✅ 客户端创建成功")

	// 测试健康检查
	fmt.Println("🔍 执行健康检查...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.HealthCheck(ctx)
	if err != nil {
		log.Printf("❌ 健康检查失败: %v", err)
		return
	}
	fmt.Println("✅ 健康检查成功")

	// 测试Ping
	fmt.Println("🏓 测试Ping...")
	duration, err := client.Ping(ctx)
	if err != nil {
		log.Printf("❌ Ping失败: %v", err)
		return
	}
	fmt.Printf("✅ Ping成功，延迟: %v\n", duration)

	// 测试主题列表
	fmt.Println("📋 获取主题列表...")
	topics, err := client.Topic().List(ctx)
	if err != nil {
		log.Printf("❌ 获取主题列表失败: %v", err)
		return
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

	// 测试创建主题
	testTopicName := fmt.Sprintf("test-topic-%d", time.Now().Unix())
	fmt.Printf("🆕 创建测试主题: %s\n", testTopicName)

	createResult, err := client.Topic().Create(ctx, fluvio.CreateTopicOptions{
		Name:       testTopicName,
		Partitions: 1,
	})
	if err != nil {
		log.Printf("❌ 创建主题失败: %v", err)
	} else if createResult.Success {
		fmt.Println("✅ 主题创建成功")

		// 测试生产消息
		fmt.Println("📤 发送测试消息...")
		produceResult, err := client.Producer().Produce(ctx, "Hello from Old API!", fluvio.ProduceOptions{
			Topic: testTopicName,
			Key:   "test-key",
			Headers: map[string]string{
				"source": "old-api-test",
				"time":   time.Now().Format(time.RFC3339),
			},
		})
		if err != nil {
			log.Printf("❌ 发送消息失败: %v", err)
		} else {
			fmt.Printf("✅ 消息发送成功，ID: %s\n", produceResult.MessageID)
		}

		// 测试消费消息
		fmt.Println("📥 消费测试消息...")
		messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
			Topic:       testTopicName,
			Group:       "test-group",
			MaxMessages: 5,
		})
		if err != nil {
			log.Printf("❌ 消费消息失败: %v", err)
		} else {
			fmt.Printf("✅ 消费成功，收到 %d 条消息\n", len(messages))
			for i, msg := range messages {
				if i < 3 { // 只显示前3条
					fmt.Printf("   消息%d: [%s] %s\n", i+1, msg.Key, msg.Value)
				}
			}
		}

		// 清理：删除测试主题
		fmt.Printf("🗑️ 清理测试主题: %s\n", testTopicName)
		deleteResult, err := client.Topic().Delete(ctx, fluvio.DeleteTopicOptions{
			Name: testTopicName,
		})
		if err != nil {
			log.Printf("❌ 删除主题失败: %v", err)
		} else if deleteResult.Success {
			fmt.Println("✅ 主题删除成功")
		}
	} else {
		log.Printf("❌ 创建主题失败: %s", createResult.Error)
	}
}
