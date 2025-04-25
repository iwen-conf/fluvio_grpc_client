package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fluvio_grpc_client/internal/client"
	pb "fluvio_grpc_client/proto/fluvio_service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	client    *client.FluvioServiceClient
	topicName string
}

func NewHandler(client *client.FluvioServiceClient, topicName string) *Handler {
	return &Handler{
		client:    client,
		topicName: topicName,
	}
}

func (h *Handler) HandleCommand(command string, args []string) (exit bool) {
	switch command {
	case "exit", "quit":
		return true
	case "help":
		PrintHelp()
	case "produce":
		h.handleProduce(args)
	case "batch_produce":
		h.handleBatchProduce(args)
	case "consume":
		h.handleConsume(args)
	case "health":
		h.handleHealthCheck()
	case "topics":
		h.handleListTopics()
	case "delete_topic":
		h.handleDeleteTopic(args)
	case "create_topic":
		h.handleCreateTopic(args)
	default:
		PrintError("未知命令，输入 help 查看可用命令。")
	}
	return false
}

func (h *Handler) handleProduce(args []string) {
	if len(args) < 1 {
		PrintError("用法: produce <消息内容>")
		return
	}
	msg := strings.Join(args, " ")
	headers := map[string]string{"from": "cli"}
	timestampProto := timestamppb.New(time.Now())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	PrintRequest("生产消息: 主题=%s, 键=%s, 内容=%s, headers=%v, timestamp=%s", h.topicName, "key1", msg, headers, timestampProto.AsTime())
	produceReq := &pb.ProduceRequest{
		Topic:     h.topicName,
		Message:   msg,
		Key:       "key1",
		Headers:   headers,
		Timestamp: timestampProto,
	}
	prResp, err := h.client.Produce(ctx, produceReq)
	if err != nil {
		PrintError("生产消息调用失败: %v", err)
		return
	}
	if prResp.GetSuccess() {
		PrintResponseSuccess("生产消息成功")
	} else {
		PrintResponseError("生产消息失败: %s", prResp.GetError())
	}
}

func (h *Handler) handleBatchProduce(args []string) {
	if len(args) < 1 {
		PrintError("用法: batch_produce <消息1,消息2,...>")
		return
	}
	msgs := strings.Join(args, " ")
	msgArr := strings.Split(msgs, ",")
	var reqs []*pb.ProduceRequest
	headers := map[string]string{"from": "cli-batch"}
	timestampProto := timestamppb.New(time.Now())
	for i, m := range msgArr {
		m = strings.TrimSpace(m)
		if m == "" {
			continue
		}
		reqs = append(reqs, &pb.ProduceRequest{
			Topic:     h.topicName,
			Message:   m,
			Key:       fmt.Sprintf("key%d", i+1),
			Headers:   headers,
			Timestamp: timestampProto,
		})
	}
	if len(reqs) == 0 {
		PrintError("没有有效的消息可供批量生产")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	PrintRequest("批量生产消息: 主题=%s, 条数=%d, headers=%v, timestamp=%s", h.topicName, len(reqs), headers, timestampProto.AsTime())
	batchReq := &pb.BatchProduceRequest{
		Topic:    h.topicName,
		Messages: reqs,
	}
	batchResp, err := h.client.BatchProduce(ctx, batchReq)
	if err != nil {
		PrintError("批量生产消息调用失败: %v", err)
		return
	}
	PrintResponseSuccess("批量生产结果: success=%v, error=%v", batchResp.GetSuccess(), batchResp.GetError())
}

func (h *Handler) handleConsume(args []string) {
	max := 1
	offset := int64(0)
	group := "default"
	if len(args) > 0 {
		if v, err := strconv.Atoi(args[0]); err == nil && v > 0 {
			max = v
		}
	}
	if len(args) > 1 {
		if v, err := strconv.ParseInt(args[1], 10, 64); err == nil && v >= 0 {
			offset = v
		}
	}
	if len(args) > 2 {
		group = args[2]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	PrintRequest("消费消息: 主题=%s, 最大消息数=%d, 起始偏移=%d, 消费组=%s", h.topicName, max, offset, group)
	consumeReq := &pb.ConsumeRequest{
		Topic:       h.topicName,
		MaxMessages: int32(max),
		Offset:      offset,
		Group:       group,
	}
	crResp, err := h.client.Consume(ctx, consumeReq)
	if err != nil {
		PrintError("消费消息调用失败: %v", err)
		return
	}
	if crResp.GetError() == "" {
		if len(crResp.GetMessages()) == 0 {
			PrintResponseInfo("未消费到新消息")
		} else {
			PrintResponseInfo("消费到的消息:")
			for _, msg := range crResp.GetMessages() {
				PrintResponseInfo("偏移=%d, Key=%s, 内容=%s", msg.GetOffset(), msg.GetKey(), msg.GetMessage())
			}
			PrintResponseInfo("建议下次消费偏移: %d", crResp.GetNextOffset())
		}
	} else {
		PrintResponseError("消费消息失败: %s", crResp.GetError())
	}
}

func (h *Handler) handleHealthCheck() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	PrintRequest("健康检查 ...")
	hcResp, err := h.client.HealthCheck(ctx)
	if err != nil {
		PrintError("健康检查调用失败: %v", err)
		return
	}
	PrintResponseSuccess("健康检查: ok=%v, 消息=%s", hcResp.GetOk(), hcResp.GetMessage())
}

func (h *Handler) handleListTopics() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	PrintRequest("获取主题列表 ...")
	ltResp, err := h.client.ListTopics(ctx)
	if err != nil {
		PrintError("获取主题列表调用失败: %v", err)
		return
	}
	PrintResponseInfo("主题列表: %v", ltResp.GetTopics())
}

func (h *Handler) handleDeleteTopic(args []string) {
	if len(args) < 1 {
		PrintError("用法: delete_topic <主题名>")
		return
	}
	topic := args[0]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	PrintRequest("删除主题: %s", topic)
	dtResp, err := h.client.DeleteTopic(ctx, &pb.DeleteTopicRequest{Topic: topic})
	if err != nil {
		PrintError("删除主题调用失败: %v", err)
		return
	}
	if dtResp.GetSuccess() {
		PrintResponseSuccess("删除主题成功")
	} else {
		PrintResponseError("删除主题失败: %s", dtResp.GetError())
	}
}

func (h *Handler) handleCreateTopic(args []string) {
	if len(args) < 1 {
		PrintError("用法: create_topic <主题名> [分区数]")
		return
	}
	topic := args[0]
	partitions := int32(1)
	if len(args) > 1 {
		if v, err := strconv.Atoi(args[1]); err == nil && v > 0 {
			partitions = int32(v)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	PrintRequest("创建主题: 名称=%s, 分区数=%d", topic, partitions)
	ctResp, err := h.client.CreateTopic(ctx, &pb.CreateTopicRequest{Topic: topic, Partitions: partitions})
	if err != nil {
		PrintError("创建主题调用失败: %v", err)
		return
	}
	if ctResp.GetSuccess() {
		PrintResponseSuccess("创建主题成功, 错误信息=%s", ctResp.GetError())
	} else {
		PrintResponseError("创建主题失败: %s", ctResp.GetError())
	}
}
