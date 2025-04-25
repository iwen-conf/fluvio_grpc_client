package client

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/internal/config"

	"github.com/iwen-conf/colorprint/clr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultConnectTimeout = 5 * time.Second

// Connect 根据提供的配置连接到 gRPC 服务器
// 它返回一个 gRPC 客户端连接或错误
func Connect(cfg *config.ServerConfig) (*grpc.ClientConn, error) {
	if cfg == nil {
		// 使用红色打印错误信息
		errMsg := "服务器配置不能为空"
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, fmt.Errorf("%s", errMsg)
	}
	serverAddr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	// 使用黄色打印连接信息
	log.Println(clr.FGColor(fmt.Sprintf("正在连接到 gRPC 服务器: %s", serverAddr), clr.Yellow))

	conn, err := grpc.NewClient(serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("建立配置服务失败 (%s): %w", serverAddr, err)
	}
	conn.Connect()
	// 读取 connectTimeout 配置
	connectTimeout := defaultConnectTimeout
	if cfg.ConnectTimeout > 0 {
		connectTimeout = time.Duration(cfg.ConnectTimeout) * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			break
		}
		if !conn.WaitForStateChange(ctx, state) {
			return nil, fmt.Errorf("等待连接状态变化超时或被取消")
		}
		if conn.GetState() == connectivity.TransientFailure || conn.GetState() == connectivity.Shutdown {
			return nil, fmt.Errorf("连接失败，当前状态: %v", conn.GetState())
		}
	}
	// 使用绿色打印成功信息
	log.Println(clr.FGColor("已成功连接到 gRPC 服务器", clr.Green))
	return conn, nil
}
