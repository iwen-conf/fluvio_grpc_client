package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	// 创建客户端
	client, err := fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(5*time.Second, 10*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
	)
	if err != nil {
		log.Fatal("创建客户端失败:", err)
	}
	defer client.Close()

	ctx := context.Background()

	// 健康检查
	fmt.Println("执行健康检查...")
	err = client.HealthCheck(ctx)
	if err != nil {
		log.Fatal("健康检查失败:", err)
	}
	fmt.Println("✓ 健康检查成功")

	// 创建主题（如果不存在）
	topicName := "example-topic"
	fmt.Printf("确保主题 '%s' 存在...\n", topicName)
	_, err = client.Topic().CreateIfNotExists(ctx, fluvio.CreateTopicOptions{
		Name:       topicName,
		Partitions: 1,
	})
	if err != nil {
		log.Fatal("创建主题失败:", err)
	}
	fmt.Println("✓ 主题已就绪")

	// 生产消息
	fmt.Println("生产消息...")
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Hello from Fluvio SDK! Message #%d", i)
		result, err := client.Producer().Produce(ctx, message, fluvio.ProduceOptions{
			Topic: topicName,
			Key:   fmt.Sprintf("key-%d", i),
		})
		if err != nil {
			log.Printf("生产消息 %d 失败: %v", i, err)
			continue
		}
		fmt.Printf("✓ 消息 %d 发送成功: %s\n", i, result.MessageID)
	}

	// 等待一下确保消息已被处理
	time.Sleep(1 * time.Second)

	// 消费消息
	fmt.Println("消费消息...")
	messages, err := client.Consumer().Consume(ctx, fluvio.ConsumeOptions{
		Topic:       topicName,
		Group:       "example-group",
		MaxMessages: 10,
		Offset:      0, // 从头开始
	})
	if err != nil {
		log.Fatal("消费消息失败:", err)
	}

	fmt.Printf("✓ 收到 %d 条消息:\n", len(messages))
	for i, msg := range messages {
		fmt.Printf("  %d. [%s] %s (offset: %d)\n", i+1, msg.Key, msg.Value, msg.Offset)
	}

	// 列出主题
	fmt.Println("列出所有主题...")
	topicsResult, err := client.Topic().List(ctx)
	if err != nil {
		log.Fatal("列出主题失败:", err)
	}
	fmt.Printf("✓ 找到 %d 个主题: %v\n", len(topicsResult.Topics), topicsResult.Topics)

	// 获取集群信息
	fmt.Println("获取集群信息...")
	clusterResult, err := client.Admin().DescribeCluster(ctx)
	if err != nil {
		log.Printf("获取集群信息失败: %v", err)
	} else {
		fmt.Printf("✓ 集群状态: %s, 控制器ID: %d\n",
			clusterResult.Cluster.Status, clusterResult.Cluster.ControllerID)
	}

	fmt.Println("✓ 基本示例完成!")
}
