package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	// 创建高性能客户端
	client, err := fluvio.HighThroughputClient("localhost", 50051)
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 健康检查
	fmt.Println("执行健康检查...")
	duration, err := client.Ping(ctx)
	if err != nil {
		log.Fatal("健康检查失败:", err)
	}
	fmt.Printf("✓ 健康检查成功 (延迟: %v)\n", duration)

	// 批量创建主题
	fmt.Println("批量创建主题...")
	topics := []fluvio.CreateTopicOptions{
		{Name: "advanced-topic-1", Partitions: 3},
		{Name: "advanced-topic-2", Partitions: 5},
		{Name: "advanced-topic-3", Partitions: 2},
	}

	for _, topic := range topics {
		_, err := client.Topic().CreateIfNotExists(ctx, topic)
		if err != nil {
			log.Printf("创建主题 %s 失败: %v", topic.Name, err)
		} else {
			fmt.Printf("✓ 主题 %s 已就绪\n", topic.Name)
		}
	}

	// 演示高级生产者功能
	fmt.Println("\n=== 高级生产者演示 ===")
	demonstrateAdvancedProducer(ctx, client)

	// 演示高级消费者功能
	fmt.Println("\n=== 高级消费者演示 ===")
	demonstrateAdvancedConsumer(ctx, client)

	// 演示流式消费
	fmt.Println("\n=== 流式消费演示 ===")
	demonstrateStreamConsumer(ctx, client)

	// 演示并发操作
	fmt.Println("\n=== 并发操作演示 ===")
	demonstrateConcurrentOperations(ctx, client)

	// 演示错误处理和重试
	fmt.Println("\n=== 错误处理和重试演示 ===")
	demonstrateErrorHandling(ctx, client)

	fmt.Println("\n✓ 高级示例完成!")
}

func demonstrateAdvancedProducer(ctx context.Context, client *fluvio.Client) {
	// 批量生产消息
	fmt.Println("批量生产消息...")
	messages := make([]fluvio.Message, 10)
	for i := 0; i < 10; i++ {
		messages[i] = fluvio.Message{
			Topic:     "advanced-topic-1",
			Key:       fmt.Sprintf("batch-key-%d", i),
			Value:     fmt.Sprintf("批量消息 #%d - %s", i+1, time.Now().Format(time.RFC3339)),
			Headers:   map[string]string{"source": "advanced-example", "batch": "true"},
			Timestamp: time.Now(),
		}
	}

	batchResult, err := client.Producer().ProduceBatch(ctx, messages)
	if err != nil {
		log.Printf("批量生产失败: %v", err)
		return
	}

	successCount := 0
	for _, result := range batchResult.Results {
		if result.Success {
			successCount++
		}
	}
	fmt.Printf("✓ 批量生产完成: %d/%d 成功\n", successCount, len(messages))

	// 异步生产消息
	fmt.Println("异步生产消息...")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			resultChan := client.Producer().ProduceAsync(ctx, 
				fmt.Sprintf("异步消息 #%d", id+1), 
				fluvio.ProduceOptions{
					Topic: "advanced-topic-2",
					Key:   fmt.Sprintf("async-key-%d", id),
					Headers: map[string]string{"type": "async"},
				})
			
			result := <-resultChan
			if result.Error != nil {
				log.Printf("异步生产消息 %d 失败: %v", id+1, result.Error)
			} else {
				fmt.Printf("✓ 异步消息 %d 发送成功\n", id+1)
			}
		}(i)
	}
	wg.Wait()
}

func demonstrateAdvancedConsumer(ctx context.Context, client *fluvio.Client) {
	// 带重试的消费
	fmt.Println("带重试的消费...")
	messages, err := client.Consumer().ConsumeWithRetry(ctx, fluvio.ConsumeOptions{
		Topic:       "advanced-topic-1",
		Group:       "advanced-group",
		MaxMessages: 5,
		AutoCommit:  true,
	})
	if err != nil {
		log.Printf("消费失败: %v", err)
		return
	}
	fmt.Printf("✓ 消费到 %d 条消息\n", len(messages))

	// 手动提交偏移量
	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		err = client.Consumer().CommitOffset(ctx, fluvio.CommitOffsetOptions{
			Topic:  "advanced-topic-1",
			Group:  "advanced-group",
			Offset: lastMessage.Offset + 1,
		})
		if err != nil {
			log.Printf("提交偏移量失败: %v", err)
		} else {
			fmt.Printf("✓ 偏移量提交成功: %d\n", lastMessage.Offset+1)
		}
	}
}

func demonstrateStreamConsumer(ctx context.Context, client *fluvio.Client) {
	// 创建带超时的上下文
	streamCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fmt.Println("开始流式消费 (10秒)...")
	stream, err := client.Consumer().ConsumeStream(streamCtx, fluvio.StreamConsumeOptions{
		Topic: "advanced-topic-2",
		Group: "stream-group",
	})
	if err != nil {
		log.Printf("创建流式消费失败: %v", err)
		return
	}

	messageCount := 0
	for {
		select {
		case msg, ok := <-stream:
			if !ok {
				fmt.Printf("✓ 流式消费结束，共收到 %d 条消息\n", messageCount)
				return
			}
			if msg.Error != nil {
				log.Printf("流式消费错误: %v", msg.Error)
				continue
			}
			messageCount++
			fmt.Printf("  流式消息 %d: [%s] %s\n", messageCount, msg.Message.Key, msg.Message.Value)
		case <-streamCtx.Done():
			fmt.Printf("✓ 流式消费超时，共收到 %d 条消息\n", messageCount)
			return
		}
	}
}

func demonstrateConcurrentOperations(ctx context.Context, client *fluvio.Client) {
	var wg sync.WaitGroup
	
	// 并发生产者
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				_, err := client.Producer().Produce(ctx, 
					fmt.Sprintf("并发消息 P%d-M%d", producerID, j+1),
					fluvio.ProduceOptions{
						Topic: "advanced-topic-3",
						Key:   fmt.Sprintf("producer-%d-msg-%d", producerID, j+1),
					})
				if err != nil {
					log.Printf("生产者 %d 消息 %d 失败: %v", producerID, j+1, err)
				}
			}
			fmt.Printf("✓ 生产者 %d 完成\n", producerID)
		}(i)
	}

	// 并发消费者
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
				Topic:       "advanced-topic-3",
				Group:       fmt.Sprintf("concurrent-group-%d", consumerID),
				MaxMessages: 10,
			})
			if err != nil {
				log.Printf("消费者 %d 失败: %v", consumerID, err)
				return
			}
			fmt.Printf("✓ 消费者 %d 收到 %d 条消息\n", consumerID, len(messages))
		}(i)
	}

	wg.Wait()
	fmt.Println("✓ 并发操作完成")
}

func demonstrateErrorHandling(ctx context.Context, client *fluvio.Client) {
	// 尝试操作不存在的主题
	fmt.Println("测试错误处理...")
	_, err := client.Producer().Produce(ctx, "测试消息", fluvio.ProduceOptions{
		Topic: "non-existent-topic",
	})
	if err != nil {
		fmt.Printf("✓ 预期错误被正确处理: %v\n", err)
	}

	// 带重试的操作
	fmt.Println("测试重试机制...")
	_, err = client.Producer().ProduceWithRetry(ctx, "重试测试消息", fluvio.ProduceOptions{
		Topic: "advanced-topic-1",
		Key:   "retry-test",
	})
	if err != nil {
		log.Printf("重试失败: %v", err)
	} else {
		fmt.Println("✓ 重试成功")
	}
}
