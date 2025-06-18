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
	fmt.Println("=== Fluvio Go SDK 基本示例 ===")

	// 1. 创建客户端
	client, err := fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(5*time.Second, 30*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
	)
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 2. 健康检查
	fmt.Println("🔍 检查连接...")
	err = client.HealthCheck(ctx)
	if err != nil {
		log.Printf("⚠️  健康检查失败: %v", err)
		fmt.Println("继续执行其他功能...")
	} else {
		fmt.Println("✅ 连接成功!")
	}

	// 3. 创建主题（使用新的配置选项）
	topicName := "basic-example-topic"
	fmt.Printf("📁 创建主题 '%s'...\n", topicName)
	_, err = client.Topic().CreateIfNotExists(ctx, types.CreateTopicOptions{
		Name:              topicName,
		Partitions:        2,                   // 多分区
		ReplicationFactor: 1,                   // 新字段：复制因子
		RetentionMs:       24 * 60 * 60 * 1000, // 新字段：保留时间（24小时）
		Config: map[string]string{ // 新字段：主题配置
			"cleanup.policy": "delete",
			"segment.ms":     "3600000",
		},
	})
	if err != nil {
		log.Fatal("创建主题失败:", err)
	}
	fmt.Println("✅ 主题已就绪!")

	// 4. 生产消息（展示新功能）
	fmt.Println("📤 生产消息...")

	// 生产带自定义消息ID的消息
	result, err := client.Producer().Produce(ctx, "Hello, Fluvio with MessageID!", types.ProduceOptions{
		Topic:     topicName,
		Key:       "greeting",
		MessageID: "msg-001", // 新功能：自定义消息ID
		Headers: map[string]string{
			"source":    "basic-example",
			"version":   "1.0",
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		log.Fatal("生产消息失败:", err)
	}
	fmt.Printf("✅ 消息发送成功! ID: %s\n", result.MessageID)

	// 5. 批量生产消息
	fmt.Println("📤 批量生产消息...")
	messages := []types.Message{
		{
			Topic:     topicName,
			Key:       "batch-1",
			Value:     "第一条批量消息",
			MessageID: "batch-msg-001",
			Headers: map[string]string{
				"batch": "true",
				"index": "1",
			},
		},
		{
			Topic:     topicName,
			Key:       "batch-2",
			Value:     "第二条批量消息",
			MessageID: "batch-msg-002",
			Headers: map[string]string{
				"batch": "true",
				"index": "2",
			},
		},
		{
			Topic:     topicName,
			Key:       "batch-3",
			Value:     "第三条批量消息",
			MessageID: "batch-msg-003",
			Headers: map[string]string{
				"batch": "true",
				"index": "3",
			},
		},
	}

	batchResult, err := client.Producer().ProduceBatch(ctx, messages)
	if err != nil {
		log.Fatal("批量生产失败:", err)
	}

	successCount := 0
	for i, result := range batchResult.Results {
		if result.Success {
			successCount++
			fmt.Printf("  ✅ 批量消息 %d 发送成功: %s\n", i+1, result.MessageID)
		} else {
			fmt.Printf("  ❌ 批量消息 %d 发送失败: %s\n", i+1, result.Error)
		}
	}
	fmt.Printf("✅ 批量发送完成: %d/%d 成功\n", successCount, len(messages))

	// 6. 消费消息（展示MessageID）
	fmt.Println("📥 消费消息...")
	consumedMessages, err := client.Consumer().Consume(ctx, types.ConsumeOptions{
		Topic:       topicName,
		Group:       "basic-example-group",
		MaxMessages: 10,
		Offset:      0,
	})
	if err != nil {
		log.Fatal("消费消息失败:", err)
	}

	fmt.Printf("✅ 收到 %d 条消息:\n", len(consumedMessages))
	for i, msg := range consumedMessages {
		fmt.Printf("  %d. [%s] %s (MessageID: %s, Offset: %d)\n",
			i+1, msg.Key, msg.Value, msg.MessageID, msg.Offset)
		if len(msg.Headers) > 0 {
			fmt.Printf("     Headers: %v\n", msg.Headers)
		}
	}

	// 7. 获取主题详细信息（新功能）
	fmt.Println("📊 获取主题详细信息...")
	topicDetail, err := client.Topic().DescribeTopicDetail(ctx, topicName)
	if err != nil {
		log.Printf("⚠️  获取主题详细信息失败: %v", err)
	} else {
		fmt.Printf("✅ 主题详细信息:\n")
		fmt.Printf("  - 主题: %s\n", topicDetail.Topic)
		fmt.Printf("  - 保留时间: %d ms\n", topicDetail.RetentionMs)
		fmt.Printf("  - 分区数: %d\n", len(topicDetail.Partitions))
		fmt.Printf("  - 配置: %v\n", topicDetail.Config)

		for _, partition := range topicDetail.Partitions {
			fmt.Printf("  - 分区 %d: Leader=%d, HighWatermark=%d\n",
				partition.PartitionID, partition.LeaderID, partition.HighWatermark)
		}
	}

	// 8. 获取主题统计信息（新功能）
	fmt.Println("📈 获取主题统计信息...")
	stats, err := client.Topic().GetTopicStats(ctx, types.GetTopicStatsOptions{
		Topic:             topicName,
		IncludePartitions: true,
	})
	if err != nil {
		log.Printf("⚠️  获取主题统计信息失败: %v", err)
	} else {
		fmt.Printf("✅ 主题统计信息:\n")
		for _, topicStats := range stats.Topics {
			fmt.Printf("  - 主题: %s\n", topicStats.Topic)
			fmt.Printf("  - 总消息数: %d\n", topicStats.TotalMessageCount)
			fmt.Printf("  - 总大小: %d bytes\n", topicStats.TotalSizeBytes)
			fmt.Printf("  - 分区数: %d\n", topicStats.PartitionCount)

			if len(topicStats.Partitions) > 0 {
				fmt.Printf("  - 分区统计:\n")
				for _, partStats := range topicStats.Partitions {
					fmt.Printf("    分区 %d: %d 条消息, %d bytes\n",
						partStats.PartitionID, partStats.MessageCount, partStats.TotalSizeBytes)
				}
			}
		}
	}

	fmt.Println("🎉 基本示例完成!")
}
