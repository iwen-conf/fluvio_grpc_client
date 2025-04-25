package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"fluvio_grpc_client/internal/client"
	"fluvio_grpc_client/internal/config"
	pb "fluvio_grpc_client/proto/fluvio_service"

	"github.com/iwen-conf/colorprint/clr"
)

const (
	testTopic = "test-topic"
)

func getGrpcClient(t *testing.T) *client.FluvioServiceClient {
	cfg, err := config.Load("../internal/config/config.json")
	if err != nil {
		t.Fatal(clr.FGColor(fmt.Sprintf("加载配置失败: %v", err), clr.Red))
	}
	conn, err := client.Connect(&cfg.Server)
	if err != nil {
		t.Fatal(clr.FGColor(fmt.Sprintf("无法连接 gRPC 服务器: %v", err), clr.Red))
	}
	t.Cleanup(func() { conn.Close() })
	return client.NewFluvioServiceClient(conn)
}

func TestHealthCheck(t *testing.T) {
	cli := getGrpcClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reply, err := cli.HealthCheck(ctx)
	if err != nil || !reply.GetOk() {
		t.Error(clr.FGColor(fmt.Sprintf("健康检查失败: %v, %s", err, reply.GetMessage()), clr.Red))
	} else {
		t.Log(clr.FGColor(fmt.Sprintf("健康检查通过: %s", reply.GetMessage()), clr.Green))
	}
}

func TestListTopics(t *testing.T) {
	cli := getGrpcClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reply, err := cli.ListTopics(ctx)
	if err != nil {
		t.Error(clr.FGColor(fmt.Sprintf("获取主题列表失败: %v", err), clr.Red))
	} else {
		t.Log(clr.FGColor(fmt.Sprintf("主题列表: %v", reply.GetTopics()), clr.Green))
	}
}

func TestCreateAndDeleteTopic(t *testing.T) {
	cli := getGrpcClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 创建
	_, _ = cli.DeleteTopic(ctx, &pb.DeleteTopicRequest{Topic: testTopic}) // 保证干净
	reply, err := cli.CreateTopic(ctx, &pb.CreateTopicRequest{Topic: testTopic, Partitions: 1, ReplicationFactor: 1})
	if err != nil || !reply.GetSuccess() {
		t.Error(clr.FGColor(fmt.Sprintf("创建主题失败: %v, %s", err, reply.GetError()), clr.Red))
	} else {
		t.Log(clr.FGColor("创建主题成功", clr.Green))
	}
	// 删除
	reply2, err2 := cli.DeleteTopic(ctx, &pb.DeleteTopicRequest{Topic: testTopic})
	if err2 != nil || !reply2.GetSuccess() {
		t.Error(clr.FGColor(fmt.Sprintf("删除主题失败: %v, %s", err2, reply2.GetError()), clr.Red))
	} else {
		t.Log(clr.FGColor("删除主题成功", clr.Green))
	}
}

func TestProduceAndConsume(t *testing.T) {
	cli := getGrpcClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 先确保主题存在
	_, _ = cli.CreateTopic(ctx, &pb.CreateTopicRequest{Topic: testTopic, Partitions: 1, ReplicationFactor: 1})
	// 生产
	msg := fmt.Sprintf("test-msg-%d", time.Now().UnixNano())
	prodReply, err := cli.Produce(ctx, &pb.ProduceRequest{
		Topic:   testTopic,
		Message: msg,
		Key:     "k1",
	})
	if err != nil || !prodReply.GetSuccess() {
		t.Error(clr.FGColor(fmt.Sprintf("生产消息失败: %v, %s", err, prodReply.GetError()), clr.Red))
	} else {
		t.Log(clr.FGColor("生产消息成功", clr.Green))
	}
	// 消费
	conReply, err := cli.Consume(ctx, &pb.ConsumeRequest{
		Topic:       testTopic,
		MaxMessages: 1,
		Offset:      0,
		Group:       "test-group",
	})
	if err != nil {
		t.Error(clr.FGColor(fmt.Sprintf("消费消息失败: %v", err), clr.Red))
	} else if len(conReply.GetMessages()) == 0 {
		t.Error(clr.FGColor("未消费到消息", clr.Red))
	} else {
		t.Log(clr.FGColor(fmt.Sprintf("消费到消息: %s", conReply.GetMessages()[0].GetMessage()), clr.Green))
	}
}

func TestBatchProduce(t *testing.T) {
	cli := getGrpcClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	msgs := []*pb.ProduceRequest{
		{Topic: testTopic, Message: "batch-1", Key: "k1"},
		{Topic: testTopic, Message: "batch-2", Key: "k2"},
	}
	reply, err := cli.BatchProduce(ctx, &pb.BatchProduceRequest{
		Topic:    testTopic,
		Messages: msgs,
	})
	if err != nil {
		t.Error(clr.FGColor(fmt.Sprintf("批量生产失败: %v", err), clr.Red))
	} else {
		t.Log(clr.FGColor(fmt.Sprintf("批量生产结果: %v", reply.GetSuccess()), clr.Green))
	}
}

func TestCommitOffset(t *testing.T) {
	cli := getGrpcClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 先生产一条消息
	_, _ = cli.Produce(ctx, &pb.ProduceRequest{
		Topic:   testTopic,
		Message: "offset-test",
		Key:     "k1",
	})
	// 消费一条消息
	conReply, _ := cli.Consume(ctx, &pb.ConsumeRequest{
		Topic:       testTopic,
		MaxMessages: 1,
		Offset:      0,
		Group:       "test-group",
	})
	if len(conReply.GetMessages()) == 0 {
		t.Skip("无消息可提交 offset")
	}
	offset := conReply.GetMessages()[0].GetOffset()
	commitReply, err := cli.CommitOffset(ctx, &pb.CommitOffsetRequest{
		Topic:  testTopic,
		Group:  "test-group",
		Offset: offset,
	})
	if err != nil || !commitReply.GetSuccess() {
		t.Error(clr.FGColor(fmt.Sprintf("提交 offset 失败: %v, %s", err, commitReply.GetError()), clr.Red))
	} else {
		t.Log(clr.FGColor("提交 offset 成功", clr.Green))
	}
}
