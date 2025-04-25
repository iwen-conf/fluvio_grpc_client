package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"

	"github.com/iwen-conf/colorprint/clr"
	"google.golang.org/grpc"
)

// FluvioAdminServiceClient 封装了与 FluvioAdminService gRPC 服务的交互
type FluvioAdminServiceClient struct {
	client pb.FluvioAdminServiceClient
}

// NewFluvioAdminServiceClient 创建一个新的 FluvioAdminServiceClient
func NewFluvioAdminServiceClient(conn *grpc.ClientConn) *FluvioAdminServiceClient {
	if conn == nil {
		log.Println(clr.FGColor("尝试使用 nil 连接创建 FluvioAdminServiceClient", clr.Red))
		return nil // 或者返回错误
	}
	return &FluvioAdminServiceClient{
		client: pb.NewFluvioAdminServiceClient(conn),
	}
}

// Helper function for context with timeout
func getCtxWithTimeout(parentCtx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	return context.WithTimeout(parentCtx, duration)
}

// DescribeCluster 调用 FluvioAdminService 的 DescribeCluster 方法
func (c *FluvioAdminServiceClient) DescribeCluster(ctx context.Context) (*pb.DescribeClusterReply, error) {
	log.Println(clr.FGColor("正在调用 FluvioAdminService.DescribeCluster...", clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.DescribeCluster(localCtx, &pb.DescribeClusterRequest{})
	if err != nil {
		errMsg := fmt.Sprintf("FluvioAdminService.DescribeCluster 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	successMsg := fmt.Sprintf("FluvioAdminService.DescribeCluster 成功. Status: %s, ControllerID: %d", reply.Status, reply.ControllerId)
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// ListBrokers 调用 FluvioAdminService 的 ListBrokers 方法
func (c *FluvioAdminServiceClient) ListBrokers(ctx context.Context) (*pb.ListBrokersReply, error) {
	log.Println(clr.FGColor("正在调用 FluvioAdminService.ListBrokers...", clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.ListBrokers(localCtx, &pb.ListBrokersRequest{})
	if err != nil {
		errMsg := fmt.Sprintf("FluvioAdminService.ListBrokers 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	successMsg := fmt.Sprintf("FluvioAdminService.ListBrokers 成功，获取到 %d 个 Broker 信息", len(reply.Brokers))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// GetMetrics 调用 FluvioAdminService 的 GetMetrics 方法
func (c *FluvioAdminServiceClient) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.GetMetricsReply, error) {
	log.Println(clr.FGColor("正在调用 FluvioAdminService.GetMetrics...", clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 10*time.Second) // Metrics 可能需要更长时间
	defer cancel()

	if req == nil {
		req = &pb.GetMetricsRequest{}
	}

	reply, err := c.client.GetMetrics(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioAdminService.GetMetrics 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	successMsg := fmt.Sprintf("FluvioAdminService.GetMetrics 成功，获取到 %d 个指标", len(reply.Metrics))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}
