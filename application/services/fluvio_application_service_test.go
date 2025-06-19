package services

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/repositories"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
)

// MockMessageRepository 模拟消息仓储
type MockMessageRepository struct {
	produceError      error
	produceBatchError error
	consumeMessages   []*entities.Message
	consumeError      error
}

func (m *MockMessageRepository) Produce(ctx context.Context, message *entities.Message) error {
	if m.produceError != nil {
		return m.produceError
	}
	// 模拟服务器设置MessageID
	if message.MessageID == "" {
		message.MessageID = "server-generated-id"
	}
	message.Partition = 0
	message.Offset = 123
	return nil
}

func (m *MockMessageRepository) ProduceBatch(ctx context.Context, messages []*entities.Message) error {
	return m.produceBatchError
}

func (m *MockMessageRepository) Consume(ctx context.Context, topic string, partition int32, offset int64, maxMessages int) ([]*entities.Message, error) {
	if m.consumeError != nil {
		return nil, m.consumeError
	}
	return m.consumeMessages, nil
}

func (m *MockMessageRepository) ConsumeFiltered(ctx context.Context, topic string, filters []*valueobjects.FilterCondition, maxMessages int) ([]*entities.Message, error) {
	return m.consumeMessages, m.consumeError
}

func (m *MockMessageRepository) ConsumeStream(ctx context.Context, topic string, partition int32, offset int64) (<-chan *entities.Message, error) {
	return nil, nil
}

func (m *MockMessageRepository) GetOffset(ctx context.Context, topic string, partition int32, consumerGroup string) (int64, error) {
	return 0, nil
}

func (m *MockMessageRepository) CommitOffset(ctx context.Context, topic string, partition int32, consumerGroup string, offset int64) error {
	return nil
}

// MockTopicRepository 模拟主题仓储
type MockTopicRepository struct {
	createTopicResponse   *dtos.CreateTopicResponse
	deleteTopicResponse   *dtos.DeleteTopicResponse
	listTopicsResponse    *dtos.ListTopicsResponse
	describeTopicResponse *dtos.DescribeTopicResponse
	createTopicError      error
	deleteTopicError      error
	listTopicsError       error
	describeTopicError    error
}

func (m *MockTopicRepository) CreateTopic(ctx context.Context, req *dtos.CreateTopicRequest) (*dtos.CreateTopicResponse, error) {
	return m.createTopicResponse, m.createTopicError
}

func (m *MockTopicRepository) DeleteTopic(ctx context.Context, req *dtos.DeleteTopicRequest) (*dtos.DeleteTopicResponse, error) {
	return m.deleteTopicResponse, m.deleteTopicError
}

func (m *MockTopicRepository) ListTopics(ctx context.Context, req *dtos.ListTopicsRequest) (*dtos.ListTopicsResponse, error) {
	return m.listTopicsResponse, m.listTopicsError
}

func (m *MockTopicRepository) DescribeTopic(ctx context.Context, req *dtos.DescribeTopicRequest) (*dtos.DescribeTopicResponse, error) {
	return m.describeTopicResponse, m.describeTopicError
}

func (m *MockTopicRepository) Create(ctx context.Context, topic *entities.Topic) error {
	return nil
}

func (m *MockTopicRepository) Delete(ctx context.Context, name string) error {
	return nil
}

func (m *MockTopicRepository) Exists(ctx context.Context, name string) (bool, error) {
	return true, nil
}

func (m *MockTopicRepository) List(ctx context.Context) ([]*entities.Topic, error) {
	return nil, nil
}

func (m *MockTopicRepository) GetByName(ctx context.Context, name string) (*entities.Topic, error) {
	return nil, nil
}

func (m *MockTopicRepository) GetDetail(ctx context.Context, name string) (*entities.Topic, error) {
	return nil, nil
}

func (m *MockTopicRepository) GetStats(ctx context.Context, name string) (*repositories.TopicStats, error) {
	return nil, nil
}

func (m *MockTopicRepository) GetPartitionStats(ctx context.Context, name string, partition int32) (*repositories.PartitionStats, error) {
	return nil, nil
}

// MockAdminRepository 模拟管理仓储
type MockAdminRepository struct{}

func (m *MockAdminRepository) DescribeCluster(ctx context.Context, req *dtos.DescribeClusterRequest) (*dtos.DescribeClusterResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) ListBrokers(ctx context.Context, req *dtos.ListBrokersRequest) (*dtos.ListBrokersResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) ListConsumerGroups(ctx context.Context, req *dtos.ListConsumerGroupsRequest) (*dtos.ListConsumerGroupsResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) DescribeConsumerGroup(ctx context.Context, req *dtos.DescribeConsumerGroupRequest) (*dtos.DescribeConsumerGroupResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) ListSmartModules(ctx context.Context, req *dtos.ListSmartModulesRequest) (*dtos.ListSmartModulesResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) CreateSmartModule(ctx context.Context, req *dtos.CreateSmartModuleRequest) (*dtos.CreateSmartModuleResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) DeleteSmartModule(ctx context.Context, req *dtos.DeleteSmartModuleRequest) (*dtos.DeleteSmartModuleResponse, error) {
	return nil, nil
}

func (m *MockAdminRepository) DescribeSmartModule(ctx context.Context, req *dtos.DescribeSmartModuleRequest) (*dtos.DescribeSmartModuleResponse, error) {
	return nil, nil
}

// TestFluvioApplicationService_ProduceMessage 测试生产消息
func TestFluvioApplicationService_ProduceMessage(t *testing.T) {
	tests := []struct {
		name            string
		request         *dtos.ProduceMessageRequest
		mockError       error
		expectedError   bool
		expectedSuccess bool
	}{
		{
			name: "成功生产消息",
			request: &dtos.ProduceMessageRequest{
				Topic:   "test-topic",
				Key:     "test-key",
				Value:   "test-value",
				Headers: map[string]string{"source": "test"},
			},
			mockError:       nil,
			expectedError:   false,
			expectedSuccess: true,
		},
		{
			name: "生产消息失败",
			request: &dtos.ProduceMessageRequest{
				Topic: "test-topic",
				Key:   "test-key",
				Value: "test-value",
			},
			mockError:       errors.New("production failed"),
			expectedError:   true,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟仓储
			mockMessageRepo := &MockMessageRepository{
				produceError: tt.mockError,
			}
			mockTopicRepo := &MockTopicRepository{}
			mockAdminRepo := &MockAdminRepository{}

			// 创建应用服务
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			service := NewFluvioApplicationService(mockMessageRepo, mockTopicRepo, mockAdminRepo, logger)

			// 执行测试
			resp, err := service.ProduceMessage(context.Background(), tt.request)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				if resp.Success != false {
					t.Errorf("期望Success为false，但得到 %v", resp.Success)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("期望Success %v，但得到 %v", tt.expectedSuccess, resp.Success)
				}
				if tt.expectedSuccess && resp.MessageID == "" {
					t.Errorf("期望MessageID不为空")
				}
			}
		})
	}
}

// TestFluvioApplicationService_ConsumeMessage 测试消费消息
func TestFluvioApplicationService_ConsumeMessage(t *testing.T) {
	tests := []struct {
		name            string
		request         *dtos.ConsumeMessageRequest
		mockMessages    []*entities.Message
		mockError       error
		expectedError   bool
		expectedSuccess bool
		expectedCount   int
	}{
		{
			name: "成功消费消息",
			request: &dtos.ConsumeMessageRequest{
				Topic:       "test-topic",
				Partition:   0,
				Offset:      0,
				MaxMessages: 10,
			},
			mockMessages: []*entities.Message{
				{
					ID:        "msg1",
					MessageID: "msg1",
					Topic:     "test-topic",
					Key:       "key1",
					Value:     []byte("value1"),
					Partition: 0,
					Offset:    0,
				},
				{
					ID:        "msg2",
					MessageID: "msg2",
					Topic:     "test-topic",
					Key:       "key2",
					Value:     []byte("value2"),
					Partition: 0,
					Offset:    1,
				},
			},
			mockError:       nil,
			expectedError:   false,
			expectedSuccess: true,
			expectedCount:   2,
		},
		{
			name: "消费消息失败",
			request: &dtos.ConsumeMessageRequest{
				Topic:       "test-topic",
				Partition:   0,
				Offset:      0,
				MaxMessages: 10,
			},
			mockMessages:    nil,
			mockError:       errors.New("consumption failed"),
			expectedError:   true,
			expectedSuccess: false,
			expectedCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟仓储
			mockMessageRepo := &MockMessageRepository{
				consumeMessages: tt.mockMessages,
				consumeError:    tt.mockError,
			}
			mockTopicRepo := &MockTopicRepository{}
			mockAdminRepo := &MockAdminRepository{}

			// 创建应用服务
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			service := NewFluvioApplicationService(mockMessageRepo, mockTopicRepo, mockAdminRepo, logger)

			// 执行测试
			resp, err := service.ConsumeMessage(context.Background(), tt.request)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
				if resp.Success != false {
					t.Errorf("期望Success为false，但得到 %v", resp.Success)
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("期望Success %v，但得到 %v", tt.expectedSuccess, resp.Success)
				}
				if len(resp.Messages) != tt.expectedCount {
					t.Errorf("期望消息数量 %d，但得到 %d", tt.expectedCount, len(resp.Messages))
				}
			}
		})
	}
}
