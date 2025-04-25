package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "fluvio_grpc_client/proto/fluvio_service"

	"github.com/iwen-conf/colorprint/clr"
	"google.golang.org/grpc"
)

// FluvioServiceClient 封装了与 FluvioService gRPC 服务的交互
type FluvioServiceClient struct {
	client pb.FluvioServiceClient
}

// NewFluvioServiceClient 创建一个新的 FluvioServiceClient
func NewFluvioServiceClient(conn *grpc.ClientConn) *FluvioServiceClient {
	return &FluvioServiceClient{
		client: pb.NewFluvioServiceClient(conn),
	}
}

// HealthCheck 调用 FluvioService 的 HealthCheck 方法
func (c *FluvioServiceClient) HealthCheck(ctx context.Context) (*pb.HealthCheckReply, error) {
	// 使用黄色打印调用信息
	log.Println(clr.FGColor("正在调用 FluvioService.HealthCheck...", clr.Yellow))

	// 如果未提供 context，则创建一个带超时的默认 context
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) // 5秒超时
		defer cancel()
	}

	reply, err := c.client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		// 使用红色打印失败信息
		errMsg := fmt.Sprintf("FluvioService.HealthCheck 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}

	if reply.Ok {
		// 使用绿色打印成功信息
		successMsg := fmt.Sprintf("FluvioService.HealthCheck 成功: %s", reply.Message)
		log.Println(clr.FGColor(successMsg, clr.Green))
	} else {
		// 使用红色打印失败信息
		failMsg := fmt.Sprintf("FluvioService.HealthCheck 失败: %s", reply.Message)
		log.Println(clr.FGColor(failMsg, clr.Red))
	}
	return reply, nil
}

// ListTopics 调用 FluvioService 的 ListTopics 方法
func (c *FluvioServiceClient) ListTopics(ctx context.Context) (*pb.ListTopicsReply, error) {
	// 使用黄色打印调用信息
	log.Println(clr.FGColor("正在调用 FluvioService.ListTopics...", clr.Yellow))
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second) // 5秒超时
		defer cancel()
	}

	reply, err := c.client.ListTopics(ctx, &pb.ListTopicsRequest{})
	if err != nil {
		// 使用红色打印失败信息
		errMsg := fmt.Sprintf("FluvioService.ListTopics 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	// 使用绿色打印成功信息
	successMsg := fmt.Sprintf("FluvioService.ListTopics 成功，获取到 %d 个主题", len(reply.Topics))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// Produce 调用 FluvioService 的 Produce 方法
func (c *FluvioServiceClient) Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.Produce (Topic: %s, Key: %s)...", req.Topic, req.Key), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.Produce(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.Produce 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		successMsg := fmt.Sprintf("FluvioService.Produce 成功. MessageID: %s", reply.MessageId)
		log.Println(clr.FGColor(successMsg, clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.Produce 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		// 即使 gRPC 调用成功，业务上也可能失败，返回一个错误
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// BatchProduce 调用 FluvioService 的 BatchProduce 方法
func (c *FluvioServiceClient) BatchProduce(ctx context.Context, req *pb.BatchProduceRequest) (*pb.BatchProduceReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.BatchProduce (Topic: %s, Count: %d)...", req.Topic, len(req.Messages)), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 15*time.Second) // Batch 可能需要更长时间
	defer cancel()

	reply, err := c.client.BatchProduce(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.BatchProduce 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	// 这里可以添加更详细的日志，报告成功/失败的数量
	log.Println(clr.FGColor("FluvioService.BatchProduce 调用完成.", clr.Green))
	return reply, nil
}

// Consume 调用 FluvioService 的 Consume 方法
func (c *FluvioServiceClient) Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.Consume (Topic: %s, Group: %s, Offset: %d, Partition: %d)...", req.Topic, req.Group, req.Offset, req.Partition), clr.Yellow))
	// Consume 可能需要等待，给更长的超时或允许外部控制
	localCtx, cancel := getCtxWithTimeout(ctx, 30*time.Second)
	defer cancel()

	reply, err := c.client.Consume(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.Consume 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Error != "" {
		failMsg := fmt.Sprintf("FluvioService.Consume 遇到错误: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	successMsg := fmt.Sprintf("FluvioService.Consume 成功，获取到 %d 条消息. Next Offset: %d", len(reply.Messages), reply.NextOffset)
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// StreamConsume 调用 FluvioService 的 StreamConsume 方法
// 注意：这个方法返回一个 stream 客户端，调用者需要负责接收消息
func (c *FluvioServiceClient) StreamConsume(ctx context.Context, req *pb.StreamConsumeRequest) (pb.FluvioService_StreamConsumeClient, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.StreamConsume (Topic: %s, Group: %s, Offset: %d, Partition: %d)...", req.Topic, req.Group, req.Offset, req.Partition), clr.Yellow))
	// 对于流式调用，通常使用传入的 context，不设置内部超时
	stream, err := c.client.StreamConsume(ctx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.StreamConsume 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	log.Println(clr.FGColor("FluvioService.StreamConsume 流已建立.", clr.Green))
	return stream, nil
}

// CommitOffset 调用 FluvioService 的 CommitOffset 方法
func (c *FluvioServiceClient) CommitOffset(ctx context.Context, req *pb.CommitOffsetRequest) (*pb.CommitOffsetReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.CommitOffset (Topic: %s, Group: %s, Offset: %d, Partition: %d)...", req.Topic, req.Group, req.Offset, req.Partition), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.CommitOffset(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.CommitOffset 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		log.Println(clr.FGColor("FluvioService.CommitOffset 成功.", clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.CommitOffset 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// --- 主题管理 --- //

// CreateTopic 调用 FluvioService 的 CreateTopic 方法
func (c *FluvioServiceClient) CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.CreateTopic (Topic: %s, Partitions: %d)...", req.Topic, req.Partitions), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 10*time.Second)
	defer cancel()

	reply, err := c.client.CreateTopic(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.CreateTopic 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		log.Println(clr.FGColor("FluvioService.CreateTopic 成功.", clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.CreateTopic 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// DeleteTopic 调用 FluvioService 的 DeleteTopic 方法
func (c *FluvioServiceClient) DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.DeleteTopic (Topic: %s)...", req.Topic), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 10*time.Second)
	defer cancel()

	reply, err := c.client.DeleteTopic(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.DeleteTopic 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		log.Println(clr.FGColor("FluvioService.DeleteTopic 成功.", clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.DeleteTopic 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// DescribeTopic 调用 FluvioService 的 DescribeTopic 方法
func (c *FluvioServiceClient) DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.DescribeTopic (Topic: %s)...", req.Topic), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.DescribeTopic(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.DescribeTopic 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Error != "" {
		failMsg := fmt.Sprintf("FluvioService.DescribeTopic 遇到错误: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	successMsg := fmt.Sprintf("FluvioService.DescribeTopic 成功. Topic: %s, Partitions: %d", reply.Topic, len(reply.Partitions))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// --- 消费者组管理 --- //

// ListConsumerGroups 调用 FluvioService 的 ListConsumerGroups 方法
func (c *FluvioServiceClient) ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error) {
	log.Println(clr.FGColor("正在调用 FluvioService.ListConsumerGroups...", clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.ListConsumerGroups(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.ListConsumerGroups 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Error != "" {
		failMsg := fmt.Sprintf("FluvioService.ListConsumerGroups 遇到错误: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	successMsg := fmt.Sprintf("FluvioService.ListConsumerGroups 成功，获取到 %d 个消费组", len(reply.Groups))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// DescribeConsumerGroup 调用 FluvioService 的 DescribeConsumerGroup 方法
func (c *FluvioServiceClient) DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.DescribeConsumerGroup (Group: %s)...", req.GroupId), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.DescribeConsumerGroup(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.DescribeConsumerGroup 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Error != "" {
		failMsg := fmt.Sprintf("FluvioService.DescribeConsumerGroup 遇到错误: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	successMsg := fmt.Sprintf("FluvioService.DescribeConsumerGroup 成功. Group: %s, Offset Count: %d", reply.GroupId, len(reply.Offsets))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// --- SmartModule 管理 --- //

// CreateSmartModule 调用 FluvioService 的 CreateSmartModule 方法
func (c *FluvioServiceClient) CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error) {
	name := "<nil>"
	if req.Spec != nil {
		name = req.Spec.Name
	}
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.CreateSmartModule (Name: %s)...", name), clr.Yellow))
	// Wasm upload might take longer
	localCtx, cancel := getCtxWithTimeout(ctx, 30*time.Second)
	defer cancel()

	reply, err := c.client.CreateSmartModule(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.CreateSmartModule 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		log.Println(clr.FGColor("FluvioService.CreateSmartModule 成功.", clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.CreateSmartModule 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// DeleteSmartModule 调用 FluvioService 的 DeleteSmartModule 方法
func (c *FluvioServiceClient) DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.DeleteSmartModule (Name: %s)...", req.Name), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 10*time.Second)
	defer cancel()

	reply, err := c.client.DeleteSmartModule(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.DeleteSmartModule 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		log.Println(clr.FGColor("FluvioService.DeleteSmartModule 成功.", clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.DeleteSmartModule 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// ListSmartModules 调用 FluvioService 的 ListSmartModules 方法
func (c *FluvioServiceClient) ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error) {
	log.Println(clr.FGColor("正在调用 FluvioService.ListSmartModules...", clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.ListSmartModules(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.ListSmartModules 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Error != "" {
		failMsg := fmt.Sprintf("FluvioService.ListSmartModules 遇到错误: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	successMsg := fmt.Sprintf("FluvioService.ListSmartModules 成功，获取到 %d 个 SmartModule", len(reply.Modules))
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// DescribeSmartModule 调用 FluvioService 的 DescribeSmartModule 方法
func (c *FluvioServiceClient) DescribeSmartModule(ctx context.Context, req *pb.DescribeSmartModuleRequest) (*pb.DescribeSmartModuleReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.DescribeSmartModule (Name: %s)...", req.Name), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 5*time.Second)
	defer cancel()

	reply, err := c.client.DescribeSmartModule(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.DescribeSmartModule 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Error != "" {
		failMsg := fmt.Sprintf("FluvioService.DescribeSmartModule 遇到错误: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	name := "<nil>"
	if reply.Spec != nil {
		name = reply.Spec.Name
	}
	successMsg := fmt.Sprintf("FluvioService.DescribeSmartModule 成功. Name: %s", name)
	log.Println(clr.FGColor(successMsg, clr.Green))
	return reply, nil
}

// UpdateSmartModule 调用 FluvioService 的 UpdateSmartModule 方法
func (c *FluvioServiceClient) UpdateSmartModule(ctx context.Context, req *pb.UpdateSmartModuleRequest) (*pb.UpdateSmartModuleReply, error) {
	log.Println(clr.FGColor(fmt.Sprintf("正在调用 FluvioService.UpdateSmartModule (Name: %s)...", req.Name), clr.Yellow))
	localCtx, cancel := getCtxWithTimeout(ctx, 30*time.Second) // Wasm upload might take longer
	defer cancel()

	reply, err := c.client.UpdateSmartModule(localCtx, req)
	if err != nil {
		errMsg := fmt.Sprintf("FluvioService.UpdateSmartModule 调用失败: %v", err)
		log.Println(clr.FGColor(errMsg, clr.Red))
		return nil, err
	}
	if reply.Success {
		log.Println(clr.FGColor("FluvioService.UpdateSmartModule 成功.", clr.Green))
	} else {
		failMsg := fmt.Sprintf("FluvioService.UpdateSmartModule 失败: %s", reply.Error)
		log.Println(clr.FGColor(failMsg, clr.Red))
		return reply, fmt.Errorf("%s", reply.Error)
	}
	return reply, nil
}

// Helper function (already defined in fluvio_admin.go, could be moved to a shared place)
// func getCtxWithTimeout(parentCtx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
// 	if parentCtx == nil {
// 		parentCtx = context.Background()
// 	}
// 	return context.WithTimeout(parentCtx, duration)
// }
