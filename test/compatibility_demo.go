package main

import (
	"context"
	"fmt"
	"log"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("🚀 测试向后兼容性...")
	fmt.Println("📡 连接地址: 101.43.173.154:50051")
	fmt.Println()

	// 测试1: 使用便捷函数
	fmt.Println("=== 测试1: 便捷函数 ===")
	testConvenienceFunctions()

	fmt.Println()

	// 测试2: 使用不同的客户端配置
	fmt.Println("=== 测试2: 不同配置 ===")
	testDifferentConfigurations()

	fmt.Println()

	// 测试3: 测试错误处理
	fmt.Println("=== 测试3: 错误处理 ===")
	testErrorHandling()

	fmt.Println()
	fmt.Println("🎉 向后兼容性测试完成！")
}

func testConvenienceFunctions() {
	// 测试QuickStart
	fmt.Println("⚡ 测试QuickStart...")
	client, err := fluvio.QuickStart("101.43.173.154", 50051)
	if err != nil {
		log.Printf("❌ QuickStart失败: %v", err)
		return
	}
	defer client.Close()
	fmt.Println("✅ QuickStart成功")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.HealthCheck(ctx)
	if err != nil {
		log.Printf("❌ 健康检查失败: %v", err)
	} else {
		fmt.Println("✅ QuickStart客户端健康检查成功")
	}

	// 测试SimpleProducer
	fmt.Println("📤 测试SimpleProducer...")
	producer, err := fluvio.SimpleProducer("101.43.173.154", 50051)
	if err != nil {
		log.Printf("❌ SimpleProducer创建失败: %v", err)
	} else {
		defer producer.Close()
		fmt.Println("✅ SimpleProducer创建成功")
		
		err = producer.HealthCheck(ctx)
		if err != nil {
			log.Printf("❌ SimpleProducer健康检查失败: %v", err)
		} else {
			fmt.Println("✅ SimpleProducer健康检查成功")
		}
	}

	// 测试SimpleConsumer
	fmt.Println("📥 测试SimpleConsumer...")
	consumer, err := fluvio.SimpleConsumer("101.43.173.154", 50051)
	if err != nil {
		log.Printf("❌ SimpleConsumer创建失败: %v", err)
	} else {
		defer consumer.Close()
		fmt.Println("✅ SimpleConsumer创建成功")
		
		err = consumer.HealthCheck(ctx)
		if err != nil {
			log.Printf("❌ SimpleConsumer健康检查失败: %v", err)
		} else {
			fmt.Println("✅ SimpleConsumer健康检查成功")
		}
	}

	// 测试HighThroughputClient
	fmt.Println("🚀 测试HighThroughputClient...")
	htClient, err := fluvio.HighThroughputClient("101.43.173.154", 50051)
	if err != nil {
		log.Printf("❌ HighThroughputClient创建失败: %v", err)
	} else {
		defer htClient.Close()
		fmt.Println("✅ HighThroughputClient创建成功")
		
		err = htClient.HealthCheck(ctx)
		if err != nil {
			log.Printf("❌ HighThroughputClient健康检查失败: %v", err)
		} else {
			fmt.Println("✅ HighThroughputClient健康检查成功")
		}
	}

	// 测试TestClient
	fmt.Println("🧪 测试TestClient...")
	testClient, err := fluvio.TestClient("101.43.173.154", 50051)
	if err != nil {
		log.Printf("❌ TestClient创建失败: %v", err)
	} else {
		defer testClient.Close()
		fmt.Println("✅ TestClient创建成功")
		
		err = testClient.HealthCheck(ctx)
		if err != nil {
			log.Printf("❌ TestClient健康检查失败: %v", err)
		} else {
			fmt.Println("✅ TestClient健康检查成功")
		}
	}
}

func testDifferentConfigurations() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试不同的配置组合
	configs := []struct {
		name string
		opts []fluvio.ClientOption
	}{
		{
			name: "基本配置",
			opts: []fluvio.ClientOption{
				fluvio.WithServer("101.43.173.154", 50051),
				fluvio.WithTimeout(3*time.Second, 5*time.Second),
			},
		},
		{
			name: "详细配置",
			opts: []fluvio.ClientOption{
				fluvio.WithServer("101.43.173.154", 50051),
				fluvio.WithTimeout(5*time.Second, 10*time.Second),
				fluvio.WithLogLevel(fluvio.LevelWarn),
				fluvio.WithMaxRetries(2),
				fluvio.WithPoolSize(3),
			},
		},
		{
			name: "高性能配置",
			opts: []fluvio.ClientOption{
				fluvio.WithServer("101.43.173.154", 50051),
				fluvio.WithTimeout(2*time.Second, 30*time.Second),
				fluvio.WithLogLevel(fluvio.LevelError),
				fluvio.WithMaxRetries(5),
				fluvio.WithPoolSize(10),
				fluvio.WithKeepAlive(30*time.Second),
			},
		},
	}

	for _, config := range configs {
		fmt.Printf("🔧 测试%s...\n", config.name)
		client, err := fluvio.New(config.opts...)
		if err != nil {
			log.Printf("❌ %s创建失败: %v", config.name, err)
			continue
		}
		defer client.Close()

		err = client.HealthCheck(ctx)
		if err != nil {
			log.Printf("❌ %s健康检查失败: %v", config.name, err)
		} else {
			fmt.Printf("✅ %s测试成功\n", config.name)
		}
	}
}

func testErrorHandling() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接到不存在的服务器
	fmt.Println("🔌 测试连接错误处理...")
	client, err := fluvio.New(
		fluvio.WithServer("192.168.1.999", 99999), // 不存在的地址
		fluvio.WithTimeout(1*time.Second, 2*time.Second),
		fluvio.WithMaxRetries(1),
	)
	if err != nil {
		fmt.Printf("✅ 预期的连接错误: %v\n", err)
	} else {
		defer client.Close()
		err = client.HealthCheck(ctx)
		if err != nil {
			fmt.Printf("✅ 预期的健康检查错误: %v\n", err)
		} else {
			fmt.Println("❌ 意外成功连接到不存在的服务器")
		}
	}

	// 测试Ping到不存在的服务器
	fmt.Println("🏓 测试Ping错误处理...")
	duration, err := fluvio.Ping(ctx, "192.168.1.999", 99999)
	if err != nil {
		fmt.Printf("✅ 预期的Ping错误: %v\n", err)
	} else {
		fmt.Printf("❌ 意外成功Ping到不存在的服务器，延迟: %v\n", duration)
	}

	// 测试超时处理
	fmt.Println("⏰ 测试超时处理...")
	shortCtx, shortCancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer shortCancel()
	
	client, err = fluvio.New(
		fluvio.WithServer("101.43.173.154", 50051),
		fluvio.WithTimeout(1*time.Second, 2*time.Second),
	)
	if err != nil {
		log.Printf("❌ 创建客户端失败: %v", err)
		return
	}
	defer client.Close()

	err = client.HealthCheck(shortCtx)
	if err != nil {
		fmt.Printf("✅ 预期的超时错误: %v\n", err)
	} else {
		fmt.Println("❌ 意外成功，应该超时")
	}
}
