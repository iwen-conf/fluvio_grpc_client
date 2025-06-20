package repositories

import (
	"context"
	
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	"github.com/iwen-conf/fluvio_grpc_client/domain/valueobjects"
)

// MessageRepository 消息仓储接口
type MessageRepository interface {
	// 生产消息
	Produce(ctx context.Context, message *entities.Message) error
	ProduceBatch(ctx context.Context, messages []*entities.Message) error
	
	// 消费消息
	Consume(ctx context.Context, topic string, partition int32, offset int64, maxMessages int, group string) ([]*entities.Message, error)
	ConsumeFiltered(ctx context.Context, topic string, filters []*valueobjects.FilterCondition, maxMessages int) ([]*entities.Message, error)
	
	// 流式消费
	ConsumeStream(ctx context.Context, topic string, partition int32, offset int64, group string) (<-chan *entities.Message, error)
	
	// 偏移量管理
	GetOffset(ctx context.Context, topic string, partition int32, consumerGroup string) (int64, error)
	CommitOffset(ctx context.Context, topic string, partition int32, consumerGroup string, offset int64) error
}