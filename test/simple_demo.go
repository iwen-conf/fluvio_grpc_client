package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("🚀 开始简单测试 Fluvio Go SDK...")
	fmt.Println("📡 连接地址: 101.43.173.154:50051")

	// 创建客户端
	fmt.Println("📝 创建客户端...")
	client, err := fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(5*time.Second, 10*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
	)
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
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
	} else {
		fmt.Println("✅ 健康检查成功")
	}

	// 测试Ping
	fmt.Println("🏓 测试Ping...")
	duration, err := client.Ping(ctx)
	if err != nil {
		log.Printf("❌ Ping失败: %v", err)
	} else {
		fmt.Printf("✅ Ping成功，延迟: %v\n", duration)
	}

	// 测试主题列表
	fmt.Println("📋 获取主题列表...")
	topics, err := client.Topic().List(ctx)
	if err != nil {
		log.Printf("❌ 获取主题列表失败: %v", err)
	} else {
		fmt.Printf("✅ 获取主题列表成功，共 %d 个主题\n", len(topics.Topics))
		for i, topic := range topics.Topics {
			if i < 3 { // 只显示前3个
				fmt.Printf("   - %s\n", topic)
			}
		}
		if len(topics.Topics) > 3 {
			fmt.Printf("   ... 还有 %d 个主题\n", len(topics.Topics)-3)
		}
	}

	// 测试管理功能
	fmt.Println("🔧 测试管理功能...")
	brokers, err := client.Admin().ListBrokers(ctx)
	if err != nil {
		log.Printf("❌ 获取Broker列表失败: %v", err)
	} else {
		fmt.Printf("✅ 获取Broker列表成功，共 %d 个Broker\n", len(brokers.Brokers))
		for i, broker := range brokers.Brokers {
			if i < 2 { // 只显示前2个
				fmt.Printf("   - Broker %d: %s (%s)\n", broker.ID, broker.Addr, broker.Status)
			}
		}
	}

	fmt.Println("🎉 测试完成！")
}
