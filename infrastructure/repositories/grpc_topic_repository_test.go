package repositories

import (
	"context"
	"os"
	"testing"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
)

// MockTopicGRPCClient 模拟主题管理的gRPC客户端
type MockTopicGRPCClient struct {
	MockGRPCClient
	createTopicResponse   *pb.CreateTopicReply
	deleteTopicResponse   *pb.DeleteTopicReply
	listTopicsResponse    *pb.ListTopicsReply
	describeTopicResponse *pb.DescribeTopicReply
	createTopicError      error
	deleteTopicError      error
	listTopicsError       error
	describeTopicError    error
}

func (m *MockTopicGRPCClient) CreateTopic(ctx context.Context, req *pb.CreateTopicRequest) (*pb.CreateTopicReply, error) {
	return m.createTopicResponse, m.createTopicError
}

func (m *MockTopicGRPCClient) DeleteTopic(ctx context.Context, req *pb.DeleteTopicRequest) (*pb.DeleteTopicReply, error) {
	return m.deleteTopicResponse, m.deleteTopicError
}

func (m *MockTopicGRPCClient) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsReply, error) {
	return m.listTopicsResponse, m.listTopicsError
}

func (m *MockTopicGRPCClient) DescribeTopic(ctx context.Context, req *pb.DescribeTopicRequest) (*pb.DescribeTopicReply, error) {
	return m.describeTopicResponse, m.describeTopicError
}

func (m *MockTopicGRPCClient) Close() error {
	return nil
}

// TestGRPCTopicRepository_CreateTopic 测试创建主题
func TestGRPCTopicRepository_CreateTopic(t *testing.T) {
	tests := []struct {
		name            string
		request         *dtos.CreateTopicRequest
		mockResponse    *pb.CreateTopicReply
		mockError       error
		expectedError   bool
		expectedSuccess bool
	}{
		{
			name: "成功创建主题",
			request: &dtos.CreateTopicRequest{
				Name:              "test-topic",
				Partitions:        3,
				ReplicationFactor: 1,
			},
			mockResponse: &pb.CreateTopicReply{
				Success: true,
			},
			mockError:       nil,
			expectedError:   false,
			expectedSuccess: true,
		},
		{
			name: "创建主题失败",
			request: &dtos.CreateTopicRequest{
				Name:       "test-topic",
				Partitions: 3,
			},
			mockResponse: &pb.CreateTopicReply{
				Success: false,
				Error:   "topic already exists",
			},
			mockError:       nil,
			expectedError:   false,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockTopicGRPCClient{
				createTopicResponse: tt.mockResponse,
				createTopicError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCTopicRepository(mockClient, logger).(*GRPCTopicRepository)

			// 执行测试
			resp, err := repo.CreateTopic(context.Background(), tt.request)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("期望Success %v，但得到 %v", tt.expectedSuccess, resp.Success)
				}
			}
		})
	}
}

// TestGRPCTopicRepository_DeleteTopic 测试删除主题
func TestGRPCTopicRepository_DeleteTopic(t *testing.T) {
	tests := []struct {
		name            string
		request         *dtos.DeleteTopicRequest
		mockResponse    *pb.DeleteTopicReply
		mockError       error
		expectedError   bool
		expectedSuccess bool
	}{
		{
			name: "成功删除主题",
			request: &dtos.DeleteTopicRequest{
				Name: "test-topic",
			},
			mockResponse: &pb.DeleteTopicReply{
				Success: true,
			},
			mockError:       nil,
			expectedError:   false,
			expectedSuccess: true,
		},
		{
			name: "删除不存在的主题",
			request: &dtos.DeleteTopicRequest{
				Name: "non-existent-topic",
			},
			mockResponse: &pb.DeleteTopicReply{
				Success: false,
				Error:   "topic not found",
			},
			mockError:       nil,
			expectedError:   false,
			expectedSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockTopicGRPCClient{
				deleteTopicResponse: tt.mockResponse,
				deleteTopicError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCTopicRepository(mockClient, logger).(*GRPCTopicRepository)

			// 执行测试
			resp, err := repo.DeleteTopic(context.Background(), tt.request)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if resp.Success != tt.expectedSuccess {
					t.Errorf("期望Success %v，但得到 %v", tt.expectedSuccess, resp.Success)
				}
			}
		})
	}
}

// TestGRPCTopicRepository_ListTopics 测试列出主题
func TestGRPCTopicRepository_ListTopics(t *testing.T) {
	tests := []struct {
		name          string
		mockResponse  *pb.ListTopicsReply
		mockError     error
		expectedError bool
		expectedCount int
	}{
		{
			name: "成功列出主题",
			mockResponse: &pb.ListTopicsReply{
				Topics: []string{"topic1", "topic2", "topic3"},
			},
			mockError:     nil,
			expectedError: false,
			expectedCount: 3,
		},
		{
			name: "空主题列表",
			mockResponse: &pb.ListTopicsReply{
				Topics: []string{},
			},
			mockError:     nil,
			expectedError: false,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockTopicGRPCClient{
				listTopicsResponse: tt.mockResponse,
				listTopicsError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCTopicRepository(mockClient, logger).(*GRPCTopicRepository)

			// 执行测试
			resp, err := repo.ListTopics(context.Background(), &dtos.ListTopicsRequest{})

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if len(resp.Topics) != tt.expectedCount {
					t.Errorf("期望主题数量 %d，但得到 %d", tt.expectedCount, len(resp.Topics))
				}
			}
		})
	}
}

// TestGRPCTopicRepository_Exists 测试主题存在检查
func TestGRPCTopicRepository_Exists(t *testing.T) {
	tests := []struct {
		name           string
		topicName      string
		mockResponse   *pb.DescribeTopicReply
		mockError      error
		expectedExists bool
		expectedError  bool
	}{
		{
			name:      "主题存在",
			topicName: "existing-topic",
			mockResponse: &pb.DescribeTopicReply{
				Topic: "existing-topic",
				Error: "",
			},
			mockError:      nil,
			expectedExists: true,
			expectedError:  false,
		},
		{
			name:      "主题不存在",
			topicName: "non-existent-topic",
			mockResponse: &pb.DescribeTopicReply{
				Error: "topic not found",
			},
			mockError:      nil,
			expectedExists: false,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockTopicGRPCClient{
				describeTopicResponse: tt.mockResponse,
				describeTopicError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCTopicRepository(mockClient, logger).(*GRPCTopicRepository)

			// 执行测试
			exists, err := repo.Exists(context.Background(), tt.topicName)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if exists != tt.expectedExists {
					t.Errorf("期望存在性 %v，但得到 %v", tt.expectedExists, exists)
				}
			}
		})
	}
}

// TestGRPCTopicRepository_GetByName 测试根据名称获取主题
func TestGRPCTopicRepository_GetByName(t *testing.T) {
	tests := []struct {
		name          string
		topicName     string
		mockResponse  *pb.DescribeTopicReply
		mockError     error
		expectedTopic *entities.Topic
		expectedError bool
	}{
		{
			name:      "成功获取主题",
			topicName: "test-topic",
			mockResponse: &pb.DescribeTopicReply{
				Topic: "test-topic",
				Partitions: []*pb.PartitionInfo{
					{PartitionId: 0},
					{PartitionId: 1},
					{PartitionId: 2},
				},
				Config: map[string]string{
					"retention.ms": "86400000",
				},
			},
			mockError: nil,
			expectedTopic: &entities.Topic{
				Name:       "test-topic",
				Partitions: 3,
				Config: map[string]string{
					"retention.ms": "86400000",
				},
			},
			expectedError: false,
		},
		{
			name:      "主题不存在",
			topicName: "non-existent-topic",
			mockResponse: &pb.DescribeTopicReply{
				Error: "topic not found",
			},
			mockError:     nil,
			expectedTopic: nil,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟客户端
			mockClient := &MockTopicGRPCClient{
				describeTopicResponse: tt.mockResponse,
				describeTopicError:    tt.mockError,
			}

			// 创建仓储
			logger := logging.NewStandardLogger(os.Stdout, logging.LevelDebug)
			repo := NewGRPCTopicRepository(mockClient, logger).(*GRPCTopicRepository)

			// 执行测试
			topic, err := repo.GetByName(context.Background(), tt.topicName)

			// 验证结果
			if tt.expectedError {
				if err == nil {
					t.Errorf("期望错误但没有返回错误")
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误但返回了错误: %v", err)
				}
				if tt.expectedTopic == nil {
					if topic != nil {
						t.Errorf("期望nil主题，但得到 %v", topic)
					}
				} else {
					if topic == nil {
						t.Errorf("期望主题但得到nil")
					} else {
						if topic.Name != tt.expectedTopic.Name {
							t.Errorf("期望主题名称 %s，但得到 %s", tt.expectedTopic.Name, topic.Name)
						}
						if topic.Partitions != tt.expectedTopic.Partitions {
							t.Errorf("期望分区数 %d，但得到 %d", tt.expectedTopic.Partitions, topic.Partitions)
						}
					}
				}
			}
		})
	}
}
