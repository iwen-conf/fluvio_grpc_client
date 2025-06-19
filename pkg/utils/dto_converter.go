package utils

import (
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/application/dtos"
	"github.com/iwen-conf/fluvio_grpc_client/domain/entities"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DTOConverter DTO转换器
type DTOConverter struct{}

// NewDTOConverter 创建DTO转换器
func NewDTOConverter() *DTOConverter {
	return &DTOConverter{}
}

// MessageEntityToDTO 将消息实体转换为DTO
func (c *DTOConverter) MessageEntityToDTO(message *entities.Message) *dtos.MessageDTO {
	if message == nil {
		return nil
	}

	return &dtos.MessageDTO{
		ID:        message.ID,
		MessageID: message.MessageID,
		Topic:     message.Topic,
		Key:       message.Key,
		Value:     string(message.Value),
		Headers:   message.Headers,
		Partition: message.Partition,
		Offset:    message.Offset,
		Timestamp: message.Timestamp,
	}
}

// MessageEntitiesToDTOs 批量转换消息实体为DTO
func (c *DTOConverter) MessageEntitiesToDTOs(messages []*entities.Message) []*dtos.MessageDTO {
	if messages == nil {
		return nil
	}

	dtos := make([]*dtos.MessageDTO, len(messages))
	for i, message := range messages {
		dtos[i] = c.MessageEntityToDTO(message)
	}
	return dtos
}

// MessageDTOToEntity 将消息DTO转换为实体
func (c *DTOConverter) MessageDTOToEntity(dto *dtos.MessageDTO) *entities.Message {
	if dto == nil {
		return nil
	}

	message := entities.NewMessage(dto.Key, dto.Value)
	message.ID = dto.ID
	message.MessageID = dto.MessageID
	message.Topic = dto.Topic
	message.Headers = dto.Headers
	message.Partition = dto.Partition
	message.Offset = dto.Offset
	message.Timestamp = dto.Timestamp

	return message
}

// ProduceRequestToEntity 将生产请求DTO转换为实体
func (c *DTOConverter) ProduceRequestToEntity(req *dtos.ProduceMessageRequest) *entities.Message {
	if req == nil {
		return nil
	}

	message := entities.NewMessage(req.Key, req.Value)
	message.Topic = req.Topic

	if req.MessageID != "" {
		message.WithMessageID(req.MessageID)
	}

	if req.Headers != nil {
		message.WithHeaders(req.Headers)
	}

	return message
}

// MessageEntityToProtoRequest 将消息实体转换为protobuf请求
func (c *DTOConverter) MessageEntityToProtoRequest(message *entities.Message) *pb.ProduceRequest {
	if message == nil {
		return nil
	}

	req := &pb.ProduceRequest{
		Topic:     message.Topic,
		Message:   string(message.Value),
		Key:       message.Key,
		Headers:   message.Headers,
		MessageId: message.MessageID,
	}

	// 设置时间戳
	if !message.Timestamp.IsZero() {
		req.Timestamp = timestamppb.New(message.Timestamp)
	}

	return req
}

// MessageEntitiesToProtoRequests 批量转换消息实体为protobuf请求
func (c *DTOConverter) MessageEntitiesToProtoRequests(messages []*entities.Message) []*pb.ProduceRequest {
	if messages == nil {
		return nil
	}

	requests := make([]*pb.ProduceRequest, len(messages))
	for i, message := range messages {
		requests[i] = c.MessageEntityToProtoRequest(message)
	}
	return requests
}

// ProtoMessageToEntity 将protobuf消息转换为实体
func (c *DTOConverter) ProtoMessageToEntity(protoMsg *pb.ConsumedMessage) *entities.Message {
	if protoMsg == nil {
		return nil
	}

	message := entities.NewMessage(protoMsg.GetKey(), protoMsg.GetMessage())
	message.MessageID = protoMsg.GetMessageId()
	message.Headers = protoMsg.GetHeaders()
	message.Partition = protoMsg.GetPartition()
	message.Offset = protoMsg.GetOffset()
	message.Timestamp = time.Unix(protoMsg.GetTimestamp(), 0)

	return message
}

// ProtoMessagesToEntities 批量转换protobuf消息为实体
func (c *DTOConverter) ProtoMessagesToEntities(protoMsgs []*pb.ConsumedMessage) []*entities.Message {
	if protoMsgs == nil {
		return nil
	}

	entities := make([]*entities.Message, len(protoMsgs))
	for i, protoMsg := range protoMsgs {
		entities[i] = c.ProtoMessageToEntity(protoMsg)
	}
	return entities
}

// ConsumedMessageToEntity 将ConsumedMessage转换为实体
func (c *DTOConverter) ConsumedMessageToEntity(consumedMsg *pb.ConsumedMessage) *entities.Message {
	if consumedMsg == nil {
		return nil
	}

	message := entities.NewMessage(consumedMsg.GetKey(), consumedMsg.GetMessage())
	message.MessageID = consumedMsg.GetMessageId()
	message.Headers = consumedMsg.GetHeaders()
	message.Partition = consumedMsg.GetPartition()
	message.Offset = consumedMsg.GetOffset()
	message.Timestamp = time.Unix(consumedMsg.GetTimestamp(), 0)

	return message
}

// ConsumedMessagesToEntities 批量转换ConsumedMessage为实体
func (c *DTOConverter) ConsumedMessagesToEntities(consumedMsgs []*pb.ConsumedMessage) []*entities.Message {
	if consumedMsgs == nil {
		return nil
	}

	entities := make([]*entities.Message, len(consumedMsgs))
	for i, consumedMsg := range consumedMsgs {
		entities[i] = c.ConsumedMessageToEntity(consumedMsg)
	}
	return entities
}

// TopicEntityToDTO 将主题实体转换为DTO
func (c *DTOConverter) TopicEntityToDTO(topic *entities.Topic) *dtos.TopicDTO {
	if topic == nil {
		return nil
	}

	return &dtos.TopicDTO{
		Name:              topic.Name,
		Partitions:        topic.Partitions,
		ReplicationFactor: topic.ReplicationFactor,
		Config:            topic.Config,
	}
}

// TopicEntitiesToDTOs 批量转换主题实体为DTO
func (c *DTOConverter) TopicEntitiesToDTOs(topics []*entities.Topic) []*dtos.TopicDTO {
	if topics == nil {
		return nil
	}

	dtos := make([]*dtos.TopicDTO, len(topics))
	for i, topic := range topics {
		dtos[i] = c.TopicEntityToDTO(topic)
	}
	return dtos
}

// CreateTopicRequestToEntity 将创建主题请求转换为实体
func (c *DTOConverter) CreateTopicRequestToEntity(req *dtos.CreateTopicRequest) *entities.Topic {
	if req == nil {
		return nil
	}

	return &entities.Topic{
		Name:              req.Name,
		Partitions:        req.Partitions,
		ReplicationFactor: req.ReplicationFactor,
		Config:            req.Config,
	}
}

// SmartModuleDTOToProtoSpec 将SmartModule DTO转换为protobuf规格
func (c *DTOConverter) SmartModuleDTOToProtoSpec(dto *dtos.SmartModuleDTO) *pb.SmartModuleSpec {
	if dto == nil {
		return nil
	}

	return &pb.SmartModuleSpec{
		Name:        dto.Name,
		Version:     dto.Version,
		Description: dto.Description,
		InputKind:   pb.SmartModuleInput_SMART_MODULE_INPUT_STREAM,   // 默认值
		OutputKind:  pb.SmartModuleOutput_SMART_MODULE_OUTPUT_STREAM, // 默认值
	}
}

// ProtoSpecToSmartModuleDTO 将protobuf规格转换为SmartModule DTO
func (c *DTOConverter) ProtoSpecToSmartModuleDTO(spec *pb.SmartModuleSpec) *dtos.SmartModuleDTO {
	if spec == nil {
		return nil
	}

	return &dtos.SmartModuleDTO{
		Name:        spec.GetName(),
		Version:     spec.GetVersion(),
		Description: spec.GetDescription(),
	}
}

// ProtoSpecsToSmartModuleDTOs 批量转换protobuf规格为SmartModule DTO
func (c *DTOConverter) ProtoSpecsToSmartModuleDTOs(specs []*pb.SmartModuleSpec) []*dtos.SmartModuleDTO {
	if specs == nil {
		return nil
	}

	dtos := make([]*dtos.SmartModuleDTO, len(specs))
	for i, spec := range specs {
		dtos[i] = c.ProtoSpecToSmartModuleDTO(spec)
	}
	return dtos
}

// BuildProduceResponse 构建生产响应
func (c *DTOConverter) BuildProduceResponse(message *entities.Message, success bool, errorMsg string) *dtos.ProduceMessageResponse {
	response := &dtos.ProduceMessageResponse{
		Success: success,
		Error:   errorMsg,
	}

	if message != nil {
		response.MessageID = message.MessageID
		response.Topic = message.Topic
		response.Partition = message.Partition
		response.Offset = message.Offset
	}

	return response
}

// BuildConsumeResponse 构建消费响应
func (c *DTOConverter) BuildConsumeResponse(messages []*entities.Message, success bool, errorMsg string) *dtos.ConsumeMessageResponse {
	response := &dtos.ConsumeMessageResponse{
		Success: success,
		Error:   errorMsg,
	}

	if messages != nil {
		response.Messages = c.MessageEntitiesToDTOs(messages)
		response.Count = len(messages)
	}

	return response
}

// TimeToProtoTimestamp 将时间转换为protobuf时间戳
func (c *DTOConverter) TimeToProtoTimestamp(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	return timestamppb.New(t)
}

// ProtoTimestampToTime 将protobuf时间戳转换为时间
func (c *DTOConverter) ProtoTimestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}
