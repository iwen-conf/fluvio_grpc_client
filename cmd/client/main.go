package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/internal/cli"
	"github.com/iwen-conf/fluvio_grpc_client/internal/client"
	"github.com/iwen-conf/fluvio_grpc_client/internal/config"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"

	"github.com/iwen-conf/colorprint/clr"
)

const (
	topicName = "test-topic"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load("internal/config/config.json")
	if err != nil {
		log.Fatal("加载配置失败: ", err)
	}

	// 2. 建立 gRPC 连接
	conn, err := client.Connect(&cfg.Server)
	if err != nil {
		log.Fatal("无法连接 gRPC 服务器: ", err)
	}
	defer conn.Close()

	// 3. 创建服务客户端
	fluvioClient := client.NewFluvioServiceClient(conn)

	// --- 启动时自动健康检查 ---
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	cli.PrintRequest("健康检查 ...")
	hcResp, err := fluvioClient.HealthCheck(ctx)
	if err != nil {
		cancel()
		log.Fatalf("%s", clr.FGColor(fmt.Sprintf("健康检查调用失败: %v", err), clr.Red))
	}
	cli.PrintResponseSuccess("健康检查: ok=%v, 消息=%s", hcResp.GetOk(), hcResp.GetMessage())
	cancel()

	// --- 获取主题列表 ---
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	cli.PrintRequest("获取主题列表 ...")
	ltResp, err := fluvioClient.ListTopics(ctx)
	if err != nil {
		cancel()
		log.Fatalf("%s", clr.FGColor(fmt.Sprintf("获取主题列表调用失败: %v", err), clr.Red))
	}
	cli.PrintResponseInfo("主题列表: %v", ltResp.GetTopics())
	cancel()

	// --- 自动创建主题（如不存在） ---
	topicExists := false
	for _, t := range ltResp.GetTopics() {
		if t == topicName {
			topicExists = true
			break
		}
	}
	if !topicExists {
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		cli.PrintRequest("创建主题: 名称=%s, 分区数=1", topicName)
		ctResp, err := fluvioClient.CreateTopic(ctx, &pb.CreateTopicRequest{Topic: topicName, Partitions: 1})
		if err != nil || !ctResp.GetSuccess() {
			cancel()
			errMsg := ""
			if err != nil {
				errMsg = err.Error()
			} else if ctResp != nil {
				errMsg = ctResp.GetError()
			}
			log.Fatalf("%s", clr.FGColor(fmt.Sprintf("创建主题调用失败: %s", errMsg), clr.Red))
		}
		cli.PrintResponseSuccess("创建主题成功, 错误信息=%s", ctResp.GetError())
		cancel()
	} else {
		cli.PrintInfo("主题 '%s' 已存在，跳过创建", topicName)
	}

	// --- 进入交互模式 ---
	cli.PrintWelcome()
	handler := cli.NewHandler(fluvioClient, topicName)
	reader := bufio.NewReader(os.Stdin)

	for {
		cli.PrintPrompt()
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}
		command := parts[0]
		args := parts[1:]

		exit := handler.HandleCommand(command, args)
		if exit {
			break
		}
	}

	cli.PrintExit()
}
