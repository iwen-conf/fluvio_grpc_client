package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
	"github.com/iwen-conf/fluvio_grpc_client/types"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Producer 消息生产者
type Producer struct {
	client *Client
}

// NewProducer 创建消息生产者
func NewProducer(client *Client) *Producer {
	return &Producer{
		client: client,
	}
}

// Produce 发送单条消息
func (p *Producer) Produce(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	if p.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if message == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "消息内容不能为空")
	}

	// 构建请求
	req := &pb.ProduceRequest{
		Topic:   opts.Topic,
		Message: message,
		Key:     opts.Key,
		Headers: opts.Headers,
	}

	// 设置消息ID
	if opts.MessageID != "" {
		req.MessageId = opts.MessageID
	}

	// 设置时间戳
	if opts.Timestamp != nil {
		req.Timestamp = timestamppb.New(*opts.Timestamp)
	} else {
		req.Timestamp = timestamppb.New(time.Now())
	}

	p.client.logger.Debug("发送消息",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "key", Value: opts.Key},
		logger.Field{Key: "message_length", Value: len(message)})

	var result *types.ProduceResult
	err := p.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.Produce(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "发送消息失败", err)
		}

		result = &types.ProduceResult{
			MessageID: resp.GetMessageId(),
			Offset:    0, // proto中没有offset字段
			Success:   resp.GetSuccess(),
			Error:     resp.GetError(),
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		p.client.logger.Error("发送消息失败",
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "error", Value: err})
		return nil, err
	}

	p.client.logger.Info("消息发送成功",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "message_id", Value: result.MessageID},
		logger.Field{Key: "offset", Value: result.Offset})

	return result, nil
}

// ProduceBatch 批量发送消息
func (p *Producer) ProduceBatch(ctx context.Context, messages []types.Message) (*types.BatchProduceResult, error) {
	if p.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if len(messages) == 0 {
		return nil, errors.New(errors.ErrInvalidArgument, "消息列表不能为空")
	}

	// 构建批量请求
	var requests []*pb.ProduceRequest
	for _, msg := range messages {
		if msg.Topic == "" {
			return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
		}

		req := &pb.ProduceRequest{
			Topic:   msg.Topic,
			Message: msg.Value,
			Key:     msg.Key,
			Headers: msg.Headers,
		}

		// 设置消息ID
		if msg.MessageID != "" {
			req.MessageId = msg.MessageID
		}

		if !msg.Timestamp.IsZero() {
			req.Timestamp = timestamppb.New(msg.Timestamp)
		} else {
			req.Timestamp = timestamppb.New(time.Now())
		}

		requests = append(requests, req)
	}

	batchReq := &pb.BatchProduceRequest{
		Topic:    messages[0].Topic, // 使用第一个消息的主题
		Messages: requests,
	}

	p.client.logger.Debug("批量发送消息",
		logger.Field{Key: "count", Value: len(messages)})

	var result *types.BatchProduceResult
	err := p.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.BatchProduce(ctx, batchReq)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "批量发送消息失败", err)
		}

		// 转换结果
		var results []*types.ProduceResult
		var errorMsgs []string

		successList := resp.GetSuccess()
		errorList := resp.GetError()

		overallSuccess := true
		for i := 0; i < len(messages); i++ {
			success := i < len(successList) && successList[i]
			errorMsg := ""
			if i < len(errorList) {
				errorMsg = errorList[i]
			}

			results = append(results, &types.ProduceResult{
				MessageID: fmt.Sprintf("batch-%d", i),
				Offset:    0, // proto中没有offset字段
				Success:   success,
				Error:     errorMsg,
			})

			if !success {
				overallSuccess = false
				if errorMsg != "" {
					errorMsgs = append(errorMsgs, errorMsg)
				}
			}
		}

		result = &types.BatchProduceResult{
			Results:    results,
			TotalCount: len(results),
			Success:    overallSuccess,
			Errors:     errorMsgs,
		}

		if !overallSuccess {
			return errors.New(errors.ErrInternal, "批量发送部分失败")
		}

		return nil
	})

	if err != nil {
		p.client.logger.Error("批量发送消息失败",
			logger.Field{Key: "count", Value: len(messages)},
			logger.Field{Key: "error", Value: err})
		return nil, err
	}

	successCount := 0
	for _, r := range result.Results {
		if r.Success {
			successCount++
		}
	}

	p.client.logger.Info("批量消息发送完成",
		logger.Field{Key: "total", Value: len(messages)},
		logger.Field{Key: "success", Value: successCount},
		logger.Field{Key: "failed", Value: len(messages) - successCount})

	return result, nil
}

// ProduceAsync 异步发送消息
func (p *Producer) ProduceAsync(ctx context.Context, message string, opts types.ProduceOptions) <-chan *AsyncProduceResult {
	resultChan := make(chan *AsyncProduceResult, 1)

	go func() {
		defer close(resultChan)

		result, err := p.Produce(ctx, message, opts)
		resultChan <- &AsyncProduceResult{
			Result: result,
			Error:  err,
		}
	}()

	return resultChan
}

// AsyncProduceResult 异步生产结果
type AsyncProduceResult struct {
	Result *types.ProduceResult
	Error  error
}

// ProduceWithRetry 带重试的消息发送
func (p *Producer) ProduceWithRetry(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	var result *types.ProduceResult
	err := p.client.withRetry(ctx, func(retryCtx context.Context) error {
		var err error
		result, err = p.Produce(retryCtx, message, opts)
		return err
	})
	return result, err
}

// BatchProduceWithRetry 带重试的批量消息发送
func (p *Producer) BatchProduceWithRetry(ctx context.Context, messages []types.Message) (*types.BatchProduceResult, error) {
	var result *types.BatchProduceResult
	err := p.client.withRetry(ctx, func(retryCtx context.Context) error {
		var err error
		result, err = p.ProduceBatch(retryCtx, messages)
		return err
	})
	return result, err
}

// ProducerConfig 生产者配置
type ProducerConfig struct {
	DefaultTopic    string            `json:"default_topic"`
	DefaultHeaders  map[string]string `json:"default_headers"`
	BatchSize       int               `json:"batch_size"`
	BatchTimeout    time.Duration     `json:"batch_timeout"`
	EnableBatching  bool              `json:"enable_batching"`
	EnableAsync     bool              `json:"enable_async"`
	AsyncBufferSize int               `json:"async_buffer_size"`
}

// DefaultProducerConfig 返回默认生产者配置
func DefaultProducerConfig() *ProducerConfig {
	return &ProducerConfig{
		DefaultHeaders:  make(map[string]string),
		BatchSize:       100,
		BatchTimeout:    1 * time.Second,
		EnableBatching:  false,
		EnableAsync:     false,
		AsyncBufferSize: 1000,
	}
}

// AdvancedProducer 高级生产者
type AdvancedProducer struct {
	producer *Producer
	config   *ProducerConfig

	// 批处理相关
	batchMutex  sync.Mutex
	batchBuffer []types.Message
	batchTimer  *time.Timer

	// 异步处理相关
	asyncChan   chan *asyncProduceRequest
	asyncWg     sync.WaitGroup
	asyncCtx    context.Context
	asyncCancel context.CancelFunc

	// 状态
	closed bool
	mu     sync.RWMutex
}

// asyncProduceRequest 异步生产请求
type asyncProduceRequest struct {
	message    string
	opts       types.ProduceOptions
	resultChan chan *AsyncProduceResult
}

// NewAdvancedProducer 创建高级生产者
func NewAdvancedProducer(client *Client, config *ProducerConfig) *AdvancedProducer {
	if config == nil {
		config = DefaultProducerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	ap := &AdvancedProducer{
		producer:    NewProducer(client),
		config:      config,
		batchBuffer: make([]types.Message, 0, config.BatchSize),
		asyncChan:   make(chan *asyncProduceRequest, config.AsyncBufferSize),
		asyncCtx:    ctx,
		asyncCancel: cancel,
	}

	// 启动异步处理器
	if config.EnableAsync {
		ap.startAsyncProcessor()
	}

	// 启动批处理定时器
	if config.EnableBatching {
		ap.startBatchTimer()
	}

	return ap
}

// Produce 发送消息（支持批处理和异步）
func (ap *AdvancedProducer) Produce(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	if ap.isClosed() {
		return nil, errors.New(errors.ErrConnection, "生产者已关闭")
	}

	// 设置默认值
	if opts.Topic == "" && ap.config.DefaultTopic != "" {
		opts.Topic = ap.config.DefaultTopic
	}

	if opts.Headers == nil {
		opts.Headers = make(map[string]string)
	}

	// 合并默认headers
	for k, v := range ap.config.DefaultHeaders {
		if _, exists := opts.Headers[k]; !exists {
			opts.Headers[k] = v
		}
	}

	// 如果启用了批处理，添加到批处理缓冲区
	if ap.config.EnableBatching {
		return ap.addToBatch(ctx, message, opts)
	}

	// 如果启用了异步处理
	if ap.config.EnableAsync {
		return ap.produceAsync(ctx, message, opts)
	}

	// 同步发送
	return ap.producer.Produce(ctx, message, opts)
}

// addToBatch 添加消息到批处理缓冲区
func (ap *AdvancedProducer) addToBatch(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	ap.batchMutex.Lock()
	defer ap.batchMutex.Unlock()

	msg := types.Message{
		Topic:   opts.Topic,
		Key:     opts.Key,
		Value:   message,
		Headers: opts.Headers,
	}

	if opts.Timestamp != nil {
		msg.Timestamp = *opts.Timestamp
	} else {
		msg.Timestamp = time.Now()
	}

	ap.batchBuffer = append(ap.batchBuffer, msg)

	// 如果达到批处理大小，立即发送
	if len(ap.batchBuffer) >= ap.config.BatchSize {
		return ap.flushBatch(ctx)
	}

	// 返回一个占位结果
	return &types.ProduceResult{
		Success:   true,
		MessageID: "batched",
	}, nil
}

// flushBatch 刷新批处理缓冲区
func (ap *AdvancedProducer) flushBatch(ctx context.Context) (*types.ProduceResult, error) {
	if len(ap.batchBuffer) == 0 {
		return nil, nil
	}

	messages := make([]types.Message, len(ap.batchBuffer))
	copy(messages, ap.batchBuffer)
	ap.batchBuffer = ap.batchBuffer[:0] // 清空缓冲区

	// 发送批处理消息
	batchResult, err := ap.producer.ProduceBatch(ctx, messages)
	if err != nil {
		return nil, err
	}

	// 返回最后一个消息的结果
	if len(batchResult.Results) > 0 {
		return batchResult.Results[len(batchResult.Results)-1], nil
	}

	return &types.ProduceResult{
		Success: batchResult.Success,
	}, nil
}

// produceAsync 异步发送消息
func (ap *AdvancedProducer) produceAsync(ctx context.Context, message string, opts types.ProduceOptions) (*types.ProduceResult, error) {
	resultChan := make(chan *AsyncProduceResult, 1)

	req := &asyncProduceRequest{
		message:    message,
		opts:       opts,
		resultChan: resultChan,
	}

	select {
	case ap.asyncChan <- req:
		// 请求已发送，等待结果
		select {
		case result := <-resultChan:
			return result.Result, result.Error
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// startAsyncProcessor 启动异步处理器
func (ap *AdvancedProducer) startAsyncProcessor() {
	ap.asyncWg.Add(1)
	go func() {
		defer ap.asyncWg.Done()

		for {
			select {
			case req := <-ap.asyncChan:
				result, err := ap.producer.Produce(ap.asyncCtx, req.message, req.opts)
				req.resultChan <- &AsyncProduceResult{
					Result: result,
					Error:  err,
				}
				close(req.resultChan)

			case <-ap.asyncCtx.Done():
				return
			}
		}
	}()
}

// startBatchTimer 启动批处理定时器
func (ap *AdvancedProducer) startBatchTimer() {
	ap.batchTimer = time.AfterFunc(ap.config.BatchTimeout, func() {
		ap.batchMutex.Lock()
		defer ap.batchMutex.Unlock()

		if len(ap.batchBuffer) > 0 {
			// 创建一个后台上下文来刷新批处理
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			ap.flushBatch(ctx)
		}

		// 重新设置定时器
		if !ap.isClosed() {
			ap.batchTimer.Reset(ap.config.BatchTimeout)
		}
	})
}

// Flush 强制刷新所有待处理的消息
func (ap *AdvancedProducer) Flush(ctx context.Context) error {
	if ap.config.EnableBatching {
		ap.batchMutex.Lock()
		defer ap.batchMutex.Unlock()

		if len(ap.batchBuffer) > 0 {
			_, err := ap.flushBatch(ctx)
			return err
		}
	}

	return nil
}

// Close 关闭高级生产者
func (ap *AdvancedProducer) Close() error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if ap.closed {
		return nil
	}

	ap.closed = true

	// 停止异步处理
	if ap.asyncCancel != nil {
		ap.asyncCancel()
	}

	// 停止批处理定时器
	if ap.batchTimer != nil {
		ap.batchTimer.Stop()
	}

	// 刷新剩余的批处理消息
	if ap.config.EnableBatching {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		ap.Flush(ctx)
	}

	// 等待异步处理完成
	ap.asyncWg.Wait()

	return nil
}

// isClosed 检查是否已关闭
func (ap *AdvancedProducer) isClosed() bool {
	ap.mu.RLock()
	defer ap.mu.RUnlock()
	return ap.closed
}

// GetStats 获取生产者统计信息
func (ap *AdvancedProducer) GetStats() ProducerStats {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	ap.batchMutex.Lock()
	batchBufferSize := len(ap.batchBuffer)
	ap.batchMutex.Unlock()

	return ProducerStats{
		BatchBufferSize: batchBufferSize,
		AsyncQueueSize:  len(ap.asyncChan),
		Closed:          ap.closed,
	}
}

// ProducerStats 生产者统计信息
type ProducerStats struct {
	BatchBufferSize int  `json:"batch_buffer_size"`
	AsyncQueueSize  int  `json:"async_queue_size"`
	Closed          bool `json:"closed"`
}
