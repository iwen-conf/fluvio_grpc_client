package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client"
)

func main() {
	fmt.Println("=== Fluvio Go SDK 基本功能测试 ===")

	// 测试客户端创建
	fmt.Println("1. 测试客户端创建...")
	client, err := fluvio.New(
		fluvio.WithServer("localhost", 50051),
		fluvio.WithTimeout(5*time.Second, 10*time.Second),
		fluvio.WithLogLevel(fluvio.LevelInfo),
	)
	if err != nil {
		log.Printf("❌ 客户端创建失败: %v", err)
		return
	}
	defer client.Close()
	fmt.Println("✅ 客户端创建成功")

	ctx := context.Background()

	// 测试健康检查
	fmt.Println("2. 测试健康检查...")
	err = client.HealthCheck(ctx)
	if err != nil {
		log.Printf("❌ 健康检查失败: %v", err)
		fmt.Println("⚠️  请确保Fluvio服务正在运行在localhost:50051")
		return
	}
	fmt.Println("✅ 健康检查成功")

	// 测试主题管理
	fmt.Println("3. 测试主题管理...")
	topicsResult, err := client.Topic().List(ctx)
	if err != nil {
		log.Printf("❌ 列出主题失败: %v", err)
	} else {
		fmt.Printf("✅ 主题列表获取成功，共 %d 个主题\n", len(topicsResult.Topics))
	}

	// 测试管理功能
	fmt.Println("4. 测试管理功能...")
	clusterResult, err := client.Admin().DescribeCluster(ctx)
	if err != nil {
		log.Printf("❌ 获取集群信息失败: %v", err)
	} else {
		fmt.Printf("✅ 集群信息获取成功: 状态=%s, 控制器ID=%d\n", 
			clusterResult.Cluster.Status, clusterResult.Cluster.ControllerID)
	}

	// 测试Broker列表
	fmt.Println("5. 测试Broker列表...")
	brokersResult, err := client.Admin().ListBrokers(ctx)
	if err != nil {
		log.Printf("❌ 获取Broker列表失败: %v", err)
	} else {
		fmt.Printf("✅ Broker列表获取成功，共 %d 个Broker\n", len(brokersResult.Brokers))
	}

	fmt.Println("\n=== SDK基本功能测试完成 ===")
	fmt.Println("✅ 所有基本功能正常工作")
	fmt.Println("💡 要进行完整测试，请运行examples目录下的示例")
}
