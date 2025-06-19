package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// MockGRPCClient 模拟gRPC客户端
type MockGRPCClient struct {
	produceResponse      *pb.ProduceReply
	batchProduceResponse *pb.BatchProduceReply
	consumeResponse      *pb.ConsumeReply
	produceError         error
	batchProduceError    error
	consumeError         error
}

func (m *MockGRPCClient) Produce(ctx context.Context, req *pb.ProduceRequest) (*pb.ProduceReply, error) {
	return m.produceResponse, m.produceError
}

func (m *MockGRPCClient) BatchProduce(ctx context.Context, req *pb.BatchProduceRequest) (*pb.BatchProduceReply, error) {
	return m.batchProduceResponse, m.batchProduceError
}

func (m *MockGRPCClient) Consume(ctx context.Context, req *pb.ConsumeRequest) (*pb.ConsumeReply, error) {
	return m.consumeResponse, m.consumeError
}

// 实现其他必需的接口方法（简化版）
func (m *MockGRPCClient) StreamConsume(ctx context.Context, req *pb.StreamConsumeRequest) (pb.FluvioService_StreamConsumeClient, error) {
	return nil, nil
}

func (m *MockGRPCClient) CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) ListConsumerGroups(ctx context.Context, req *pb.ListConsumerGroupsRequest) (*pb.ListConsumerGroupsReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) DescribeConsumerGroup(ctx context.Context, req *pb.DescribeConsumerGroupRequest) (*pb.DescribeConsumerGroupReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) CommitOffset(ctx context.Context, req *pb.CommitOffsetRequest) (*pb.CommitOffsetReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) ListSmartModules(ctx context.Context, req *pb.ListSmartModulesRequest) (*pb.ListSmartModulesReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) CreateSmartModule(ctx context.Context, req *pb.CreateSmartModuleRequest) (*pb.CreateSmartModuleReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) DeleteSmartModule(ctx context.Context, req *pb.DeleteSmartModuleRequest) (*pb.DeleteSmartModuleReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) DescribeSmartModule(ctx context.Context, req *pb.DescribeSmartModuleRequest) (*pb.DescribeSmartModuleReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckReply, error) {
	return nil, nil
}

func (m *MockGRPCClient) Connect() error {
	return nil
}

func (m *MockGRPCClient) Disconnect() error {
	return nil
}

func (m *MockGRPCClient) IsConnected() bool {
	return true
}

func (m *MockGRPCClient) Close() error {
	return nil
}

// TestGRPCMessageRepository_Produce 测试消息生产
func TestGRPCMessageRepository_Produce(t *testing.T) {
	tests := []struct {
		name          string
		message       *entities.Message
		mockResponse  *pb.ProduceReply
		mockError     error
		expectedError bool
		expectedMsgID string
	}{
		{
			name: "成功生产消息",
			message: &entities.Message{
				Topic:     "test-topic",
				Key:       "test-key",
				Value:     []byte("test-value"),
				MessageID: "test-msg-id",
				Timestamp: time.Now(),
			},
			mockResponse: &pb.ProduceReply{
				MessageId: "server-msg-id",
				Success:   true,
			},
			mockError:     nil,
			expectedError: false,
			expectedMsgID: "server-msg-id",
		},
		{
			name: "生产消息失败",
			message: &entities.Message{
				Topic: "test-topic",
				Key:   "test-key",
				Value: []byte("test-value"),
			},
			mockResponse: &pb.ProduceReply{
				Success: false,
				Error:   "production failed",
			},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockGRPCClient{
				produceResponse: tt.mockResponse,
				produceError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCMessageRepository(mockClient, logger).(*GRPCMessageRepository)

			// 执行测试
			err := repo.Produce(context.Background(), tt.message)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if tt.expectedMsgID != "" && tt.message.MessageID != tt.expectedMsgID {
					t.Errorf("期望MessageID %s，但得到 %s", tt.expectedMsgID, tt.message.MessageID)
				}
			}
		})
	}
}

// TestGRPCMessageRepository_ProduceBatch 测试批量生产消息
func TestGRPCMessageRepository_ProduceBatch(t *testing.T) {
	tests := []struct {
		name          string
		messages      []*entities.Message
		mockResponse  *pb.BatchProduceReply
		mockError     error
		expectedError bool
	}{
		{
			name: "成功批量生产消息",
			messages: []*entities.Message{
				{Topic: "test-topic", Key: "key1", Value: []byte("value1")},
				{Topic: "test-topic", Key: "key2", Value: []byte("value2")},
			},
			mockResponse: &pb.BatchProduceReply{
				Success: []bool{true, true},
				Error:   []string{"", ""},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "主题不一致应该失败",
			messages: []*entities.Message{
				{Topic: "topic1", Key: "key1", Value: []byte("value1")},
				{Topic: "topic2", Key: "key2", Value: []byte("value2")},
			},
			expectedError: true,
		},
		{
			name: "部分消息失败",
			messages: []*entities.Message{
				{Topic: "test-topic", Key: "key1", Value: []byte("value1")},
				{Topic: "test-topic", Key: "key2", Value: []byte("value2")},
			},
			mockResponse: &pb.BatchProduceReply{
				Success: []bool{true, false},
				Error:   []string{"", "failed to produce"},
			},
			mockError:     nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockGRPCClient{
				batchProduceResponse: tt.mockResponse,
				batchProduceError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCMessageRepository(mockClient, logger).(*GRPCMessageRepository)

			// 执行测试
			err := repo.ProduceBatch(context.Background(), tt.messages)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
			}
		})
	}
}

// TestGRPCMessageRepository_FilterMatching 测试过滤匹配逻辑
func TestGRPCMessageRepository_FilterMatching(t *testing.T) {
	logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
	repo := &GRPCMessageRepository{logger: logger}

	message := &entities.Message{
		Key:   "test-key",
		Value: []byte("test-value"),
		Headers: map[string]string{
			"source": "test-service",
			"type":   "event",
		},
	}

	tests := []struct {
		name     string
		filter   *valueobjects.FilterCondition
		expected bool
	}{
		{
			name: "键值相等匹配",
			filter: &valueobjects.FilterCondition{
				Type:     valueobjects.FilterTypeKey,
				Operator: valueobjects.FilterOperatorEq,
				Value:    "test-key",
			},
			expected: true,
		},
		{
			name: "值包含匹配",
			filter: &valueobjects.FilterCondition{
				Type:     valueobjects.FilterTypeValue,
				Operator: valueobjects.FilterOperatorContains,
				Value:    "test",
			},
			expected: true,
		},
		{
			name: "头部字段匹配",
			filter: &valueobjects.FilterCondition{
				Type:     valueobjects.FilterTypeHeader,
				Field:    "source",
				Operator: valueobjects.FilterOperatorEq,
				Value:    "test-service",
			},
			expected: true,
		},
		{
			name: "通配符匹配",
			filter: &valueobjects.FilterCondition{
				Type:     valueobjects.FilterTypeKey,
				Operator: valueobjects.FilterOperatorRegex,
				Value:    "test-*",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repo.matchesFilter(message, tt.filter)
			if result != tt.expected {
				t.Errorf("期望 %v，但得到 %v", tt.expected, result)
			}
		})
	}
}
