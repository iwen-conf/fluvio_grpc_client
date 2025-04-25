package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/internal/client"
	"github.com/iwen-conf/fluvio_grpc_client/internal/config"

	"github.com/iwen-conf/colorprint/clr"
)

func TestGrpcHealthCheck(t *testing.T) {
	cfg, err := config.Load("../internal/config/config.json")
	if err != nil {
		t.Fatal(clr.FGColor(fmt.Sprintf("加载配置失败: %v", err), clr.Red))
	}

	conn, err := client.Connect(&cfg.Server)
	if err != nil {
		t.Fatal(clr.FGColor(fmt.Sprintf("无法连接 gRPC 服务器: %v", err), clr.Red))
	}
	defer conn.Close()

	grpcClient := client.NewFluvioServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reply, err := grpcClient.HealthCheck(ctx)
	if err != nil {
		t.Error(clr.FGColor(fmt.Sprintf("健康检查调用失败: %v", err), clr.Red))
		return
	}
	if reply.GetOk() {
		t.Log(clr.FGColor(fmt.Sprintf("gRPC 健康检查通过: %s", reply.GetMessage()), clr.Green))
	} else {
		t.Error(clr.FGColor(fmt.Sprintf("gRPC 健康检查未通过: %s", reply.GetMessage()), clr.Red))
	}
}
