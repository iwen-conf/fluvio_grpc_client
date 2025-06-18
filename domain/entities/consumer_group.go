package entities

import (
	"time"
)

// ConsumerGroup 表示消费组实体
type ConsumerGroup struct {
	// 基本信息
	GroupID     string
	Description string
	
	// 状态信息
	State       string
	Protocol    string
	ProtocolType string
	
	// 成员信息
	Members []*ConsumerMember
	
	// 偏移量信息
	Offsets []*ConsumerOffset
	
	// 元数据
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ConsumerMember 消费组成员
type ConsumerMember struct {
	MemberID   string
	ClientID   string
	ClientHost string
	Assignment []*PartitionAssignment
}

// PartitionAssignment 分区分配
type PartitionAssignment struct {
	Topic     string
	Partition int32
}

// ConsumerOffset 消费偏移量
type ConsumerOffset struct {
	Topic     string
	Partition int32
	Offset    int64
	Metadata  string
	UpdatedAt time.Time
}

// NewConsumerGroup 创建新的消费组
func NewConsumerGroup(groupID string) *ConsumerGroup {
	now := time.Now()
	return &ConsumerGroup{
		GroupID:   groupID,
		Members:   make([]*ConsumerMember, 0),
		Offsets:   make([]*ConsumerOffset, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddMember 添加成员
func (cg *ConsumerGroup) AddMember(member *ConsumerMember) {
	cg.Members = append(cg.Members, member)
	cg.UpdatedAt = time.Now()
}

// RemoveMember 移除成员
func (cg *ConsumerGroup) RemoveMember(memberID string) {
	for i, member := range cg.Members {
		if member.MemberID == memberID {
			cg.Members = append(cg.Members[:i], cg.Members[i+1:]...)
			cg.UpdatedAt = time.Now()
			break
		}
	}
}

// UpdateOffset 更新偏移量
func (cg *ConsumerGroup) UpdateOffset(topic string, partition int32, offset int64) {
	for _, co := range cg.Offsets {
		if co.Topic == topic && co.Partition == partition {
			co.Offset = offset
			co.UpdatedAt = time.Now()
			return
		}
	}
	
	// 如果不存在，创建新的偏移量记录
	cg.Offsets = append(cg.Offsets, &ConsumerOffset{
		Topic:     topic,
		Partition: partition,
		Offset:    offset,
		UpdatedAt: time.Now(),
	})
	cg.UpdatedAt = time.Now()
}

// GetOffset 获取偏移量
func (cg *ConsumerGroup) GetOffset(topic string, partition int32) (int64, bool) {
	for _, co := range cg.Offsets {
		if co.Topic == topic && co.Partition == partition {
			return co.Offset, true
		}
	}
	return 0, false
}

// IsActive 检查消费组是否活跃
func (cg *ConsumerGroup) IsActive() bool {
	return len(cg.Members) > 0 && cg.State == "Stable"
}