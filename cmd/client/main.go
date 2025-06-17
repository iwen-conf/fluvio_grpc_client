package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	fluvio "github.com/iwen-conf/fluvio_grpc_client"
	"github.com/iwen-conf/fluvio_grpc_client/internal/cli"

	"github.com/iwen-conf/colorprint/clr"
)

const (
	topicName = "test-topic"
)

var (
	configFile = flag.String("config", "internal/config/config.json", "配置文件路径")
	host       = flag.String("host", "localhost", "Fluvio服务器地址")
	port       = flag.Int("port", 50051, "Fluvio服务器端口")
	logLevel   = flag.String("log-level", "info", "日志级别 (debug, info, warn, error)")
)

// parseLogLevel 解析日志级别
func parseLogLevel(level string) fluvio.Level {
	switch strings.ToLower(level) {
	case "debug":
		return fluvio.LevelDebug
	case "info":
		return fluvio.LevelInfo
	case "warn":
		return fluvio.LevelWarn
	case "error":
		return fluvio.LevelError
	default:
		return fluvio.LevelInfo
	}
}

func main() {
	flag.Parse()

	// 1. 创建SDK客户端
	var client *fluvio.Client
	var err error

	// 优先使用命令行参数
	if *host != "localhost" || *port != 50051 {
		client, err = fluvio.ConnectWithAddress(*host, *port,
			fluvio.WithTimeout(5*time.Second, 10*time.Second),
			fluvio.WithLogLevel(parseLogLevel(*logLevel)),
			fluvio.WithMaxRetries(3),
		)
	} else {
		// 尝试从配置文件加载
		if _, err := os.Stat(*configFile); err == nil {
			cfg, err := fluvio.LoadConfigFromFile(*configFile)
			if err != nil {
				log.Printf("加载配置文件失败，使用默认配置: %v", err)
				client, err = fluvio.QuickStart("localhost", 50051)
			} else {
				client, err = fluvio.NewWithConfig(cfg)
			}
		} else {
			// 使用默认配置
			client, err = fluvio.QuickStart("localhost", 50051)
		}
	}

	if err != nil {
		log.Fatal("无法连接 Fluvio 服务器: ", err)
	}
	defer client.Close()

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
