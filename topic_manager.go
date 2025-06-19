package fluvio

import (
	"context"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/application/services"
	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// TopicManager 主题管理器
type TopicManager struct {
	appService *services.FluvioApplicationService
	logger     logging.Logger
	connected  *bool
}

// CreateTopicOptions 创建主题选项
type CreateTopicOptions struct {
	Partitions        int32             `json:"partitions,omitempty"`
	ReplicationFactor int32             `json:"replication_factor,omitempty"`
	Config            map[string]string `json:"config,omitempty"`
}

// TopicInfo 主题信息
type TopicInfo struct {
	Name              string            `json:"name"`
	Partitions        int32             `json:"partitions"`
	ReplicationFactor int32             `json:"replication_factor"`
	Config            map[string]string `json:"config,omitempty"`
}

// Create 创建主题
func (t *TopicManager) Create(ctx context.Context, name string, opts *CreateTopicOptions) error {
	if !*t.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}

	if opts == nil {
		opts = &CreateTopicOptions{
			Partitions:        1,
			ReplicationFactor: 1,
		}
	}

	t.logger.Debug("Creating topic",
		logging.Field{Key: "name", Value: name},
		logging.Field{Key: "partitions", Value: opts.Partitions})

	req := &dtos.CreateTopicRequest{
		Name:              name,
		Partitions:        opts.Partitions,
		ReplicationFactor: opts.ReplicationFactor,
		Config:            opts.Config,
	}

	resp, err := t.appService.CreateTopic(ctx, req)
	if err != nil {
		t.logger.Error("Failed to create topic", logging.Field{Key: "error", Value: err})
		return err
	}

	if !resp.Success {
		return errors.New(errors.ErrOperation, resp.Error)
	}

	t.logger.Info("Topic created successfully", logging.Field{Key: "name", Value: name})
	return nil
}

// Delete 删除主题
func (t *TopicManager) Delete(ctx context.Context, name string) error {
	if !*t.connected {
		return errors.New(errors.ErrConnection, "client not connected")
	}

	t.logger.Debug("Deleting topic", logging.Field{Key: "name", Value: name})

	req := &dtos.DeleteTopicRequest{
		Name: name,
	}

	resp, err := t.appService.DeleteTopic(ctx, req)
	if err != nil {
		t.logger.Error("Failed to delete topic", logging.Field{Key: "error", Value: err})
		return err
	}

	if !resp.Success {
		return errors.New(errors.ErrOperation, resp.Error)
	}

	t.logger.Info("Topic deleted successfully", logging.Field{Key: "name", Value: name})
	return nil
}

// List 列出主题
func (t *TopicManager) List(ctx context.Context) ([]string, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	t.logger.Debug("Listing topics")

	req := &dtos.ListTopicsRequest{}

	resp, err := t.appService.ListTopics(ctx, req)
	if err != nil {
		t.logger.Error("Failed to list topics", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	t.logger.Info("Topics listed successfully", logging.Field{Key: "count", Value: len(resp.Topics)})
	return resp.Topics, nil
}

// Info 获取主题信息
func (t *TopicManager) Info(ctx context.Context, name string) (*TopicInfo, error) {
	if !*t.connected {
		return nil, errors.New(errors.ErrConnection, "client not connected")
	}

	t.logger.Debug("Getting topic info", logging.Field{Key: "name", Value: name})

	req := &dtos.DescribeTopicRequest{
		Name: name,
	}

	resp, err := t.appService.DescribeTopic(ctx, req)
	if err != nil {
		t.logger.Error("Failed to get topic info", logging.Field{Key: "error", Value: err})
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(errors.ErrOperation, resp.Error)
	}

	info := &TopicInfo{
		Name:       resp.Topic.Name,
		Partitions: resp.Topic.Partitions,
		Config:     resp.Topic.Config,
	}

	t.logger.Info("Topic info retrieved successfully", logging.Field{Key: "name", Value: name})
	return info, nil
}

// Exists 检查主题是否存在
func (t *TopicManager) Exists(ctx context.Context, name string) (bool, error) {
	if !*t.connected {
		return false, errors.New(errors.ErrConnection, "client not connected")
	}

	t.logger.Debug("Checking if topic exists", logging.Field{Key: "name", Value: name})

	topics, err := t.List(ctx)
	if err != nil {
		return false, err
	}

	for _, topic := range topics {
		if topic == name {
			t.logger.Debug("Topic exists", logging.Field{Key: "name", Value: name})
			return true, nil
		}
	}

	t.logger.Debug("Topic does not exist", logging.Field{Key: "name", Value: name})
	return false, nil
}

// CreateIfNotExists 如果主题不存在则创建
func (t *TopicManager) CreateIfNotExists(ctx context.Context, name string, opts *CreateTopicOptions) (bool, error) {
	exists, err := t.Exists(ctx, name)
	if err != nil {
		return false, err
	}

	if exists {
		t.logger.Debug("Topic already exists", logging.Field{Key: "name", Value: name})
		return false, nil
	}

	err = t.Create(ctx, name, opts)
	if err != nil {
		return false, err
	}

	return true, nil
}