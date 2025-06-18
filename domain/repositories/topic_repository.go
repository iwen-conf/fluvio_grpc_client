package repositories

import (
	"context"
	
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
)

// TopicRepository 主题仓储接口
type TopicRepository interface {
	// 主题管理
	Create(ctx context.Context, topic *entities.Topic) error
	Delete(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
	
	// 主题查询
	List(ctx context.Context) ([]*entities.Topic, error)
	GetByName(ctx context.Context, name string) (*entities.Topic, error)
	GetDetail(ctx context.Context, name string) (*entities.Topic, error)
	
	// 主题统计
	GetStats(ctx context.Context, name string) (*TopicStats, error)
	GetPartitionStats(ctx context.Context, name string, partition int32) (*PartitionStats, error)
}

// TopicStats 主题统计信息
type TopicStats struct {
	Topic              string
	TotalMessageCount  int64
	TotalSizeBytes     int64
	PartitionCount     int32
	PartitionStats     []*PartitionStats
}

// PartitionStats 分区统计信息
type PartitionStats struct {
	PartitionID      int32
	MessageCount     int64
	TotalSizeBytes   int64
	HighWatermark    int64
	LowWatermark     int64
}