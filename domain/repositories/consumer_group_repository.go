package repositories

import (
	"context"
	
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
)

// ConsumerGroupRepository 消费组仓储接口
type ConsumerGroupRepository interface {
	// 消费组管理
	Create(ctx context.Context, group *entities.ConsumerGroup) error
	Delete(ctx context.Context, groupID string) error
	
	// 消费组查询
	List(ctx context.Context) ([]*entities.ConsumerGroup, error)
	GetByID(ctx context.Context, groupID string) (*entities.ConsumerGroup, error)
	
	// 成员管理
	AddMember(ctx context.Context, groupID string, member *entities.ConsumerMember) error
	RemoveMember(ctx context.Context, groupID string, memberID string) error
	
	// 偏移量管理
	UpdateOffset(ctx context.Context, groupID string, topic string, partition int32, offset int64) error
	GetOffsets(ctx context.Context, groupID string) ([]*entities.ConsumerOffset, error)
}