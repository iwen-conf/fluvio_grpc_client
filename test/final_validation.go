package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("🎯 Fluvio Go SDK v2.0 最终验证测试")
	fmt.Println("📡 连接地址: 101.43.173.154:50051")
	fmt.Println()

	// 测试所有配置选项
	fmt.Println("=== 配置选项测试 ===")
	testConfigurationOptions()

	// 测试核心功能
	fmt.Println("\n=== 核心功能测试 ===")
	testCoreFunctionality()

	// 测试错误处理
	fmt.Println("\n=== 错误处理测试 ===")
	testErrorHandling()

	fmt.Println("\n🎉 最终验证测试完成！")
}

func testConfigurationOptions() {
	fmt.Println("🔧 测试各种配置选项...")

	configs := []struct {
		name string
		opts []fluvio.ClientOption
	}{
		{
			name: "基本配置",
			opts: []fluvio.ClientOption{
				fluvio.WithAddress("101.43.173.154", 50051),
			},
		},
		{
			name: "完整配置",
			opts: []fluvio.ClientOption{
				fluvio.WithAddress("101.43.173.154", 50051),
				fluvio.WithTimeout(30 * time.Second),
				fluvio.WithRetry(3, time.Second),
				fluvio.WithLogLevel(fluvio.LogLevelInfo),
				fluvio.WithConnectionPool(5, 5*time.Minute),
				fluvio.WithKeepAlive(30 * time.Second),
			},
		},
		{
			name: "不安全连接",
			opts: []fluvio.ClientOption{
				fluvio.WithAddress("101.43.173.154", 50051),
				fluvio.WithInsecure(),
			},
		},
	}

	for _, config := range configs {
		fmt.Printf("   测试%s...", config.name)
		client, err := fluvio.NewClient(config.opts...)
		if err != nil {
			fmt.Printf(" ❌ 失败: %v\n", err)
			continue
		}
		client.Close()
		fmt.Println(" ✅ 成功")
	}
}

func testCoreFunctionality() {
	fmt.Println("🚀 测试核心功能...")

	client, err := fluvio.NewClient(
		fluvio.WithAddress("101.43.173.154", 50051),
		fluvio.WithTimeout(30*time.Second),
		fluvio.WithLogLevel(fluvio.LogLevelWarn), // 减少日志输出
	)
	if err != nil {
		log.Fatalf("❌ 创建客户端失败: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 连接测试
	fmt.Print("   连接测试...")
	if err := client.Connect(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
		return
	}
	fmt.Println(" ✅ 成功")

	// 健康检查测试
	fmt.Print("   健康检查...")
	if err := client.HealthCheck(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Println(" ✅ 成功")
	}

	// Ping测试
	fmt.Print("   Ping测试...")
	if duration, err := client.Ping(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (%v)\n", duration)
	}

	// 主题管理测试
	fmt.Print("   主题管理...")
	if topics, err := client.Topics().List(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (%d个主题)\n", len(topics))
	}

	// 消息生产测试
	fmt.Print("   消息生产...")
	if result, err := client.Producer().SendString(ctx, "test-topic", "key", "value"); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (ID: %s)\n", result.MessageID)
	}

	// 消息消费测试
	fmt.Print("   消息消费...")
	if messages, err := client.Consumer().Receive(ctx, "test-topic", &fluvio.ReceiveOptions{
		Group:       "test-group",
		MaxMessages: 1,
	}); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (%d条消息)\n", len(messages))
	}

	// 管理功能测试
	fmt.Print("   集群管理...")
	if clusterInfo, err := client.Admin().ClusterInfo(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (状态: %s)\n", clusterInfo.Status)
	}

	fmt.Print("   Broker管理...")
	if brokers, err := client.Admin().Brokers(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (%d个Broker)\n", len(brokers))
	}

	fmt.Print("   消费者组管理...")
	if groups, err := client.Admin().ConsumerGroups(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (%d个组)\n", len(groups))
	}

	fmt.Print("   SmartModule管理...")
	if modules, err := client.Admin().SmartModules().List(ctx); err != nil {
		fmt.Printf(" ❌ 失败: %v\n", err)
	} else {
		fmt.Printf(" ✅ 成功 (%d个模块)\n", len(modules))
	}
}

func testErrorHandling() {
	fmt.Println("⚠️ 测试错误处理...")

	// 测试无效地址
	fmt.Print("   无效地址...")
	client, err := fluvio.NewClient(
		fluvio.WithAddress("invalid-host", 99999),
		fluvio.WithTimeout(1*time.Second),
	)
	if err != nil {
		fmt.Printf(" ✅ 预期错误: %v\n", err)
	} else {
		client.Close()
		fmt.Println(" ❌ 应该失败但成功了")
	}

	// 测试超时
	fmt.Print("   连接超时...")
	client, err = fluvio.NewClient(
		fluvio.WithAddress("192.168.1.999", 50051),
		fluvio.WithTimeout(100*time.Millisecond),
	)
	if err != nil {
		fmt.Printf(" ✅ 预期错误: %v\n", err)
	} else {
		defer client.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		
		if err := client.Connect(ctx); err != nil {
			fmt.Printf(" ✅ 预期连接错误: %v\n", err)
		} else {
			fmt.Println(" ❌ 应该超时但成功了")
		}
	}

	// 测试未连接状态下的操作
	fmt.Print("   未连接操作...")
	client, err = fluvio.NewClient(
		fluvio.WithAddress("101.43.173.154", 50051),
	)
	if err != nil {
		fmt.Printf(" ❌ 创建客户端失败: %v\n", err)
		return
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.HealthCheck(ctx); err != nil {
		fmt.Printf(" ✅ 预期错误: %v\n", err)
	} else {
		fmt.Println(" ❌ 应该失败但成功了")
	}
}