package client

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/errors"
	"github.com/iwen-conf/fluvio_grpc_client/logger"
	pb "github.com/iwen-conf/fluvio_grpc_client/proto/fluvio_service"
	"github.com/iwen-conf/fluvio_grpc_client/types"

	"google.golang.org/grpc"
)

// Consumer 消息消费者
type Consumer struct {
	client *Client
}

// NewConsumer 创建消息消费者
func NewConsumer(client *Client) *Consumer {
	return &Consumer{
		client: client,
	}
}

// Consume 消费消息
func (c *Consumer) Consume(ctx context.Context, opts types.ConsumeOptions) ([]*types.Message, error) {
	if c.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Group == "" {
		opts.Group = "default"
	}

	if opts.MaxMessages <= 0 {
		opts.MaxMessages = 1
	}

	req := &pb.ConsumeRequest{
		Topic:       opts.Topic,
		Group:       opts.Group,
		Offset:      opts.Offset,
		Partition:   opts.Partition,
		MaxMessages: opts.MaxMessages,
	}

	c.client.logger.Debug("消费消息",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset},
		logger.Field{Key: "max_messages", Value: opts.MaxMessages})

	var messages []*types.Message
	err := c.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.Consume(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "消费消息失败", err)
		}

		if resp.GetError() != "" {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		// 转换消息
		for _, msg := range resp.GetMessages() {
			message := &types.Message{
				Topic:     opts.Topic,
				Key:       msg.GetKey(),
				Value:     msg.GetMessage(),
				Headers:   msg.GetHeaders(),
				Offset:    msg.GetOffset(),
				Partition: msg.GetPartition(),
			}

			if msg.GetTimestamp() > 0 {
				message.Timestamp = time.Unix(msg.GetTimestamp(), 0)
			}

			messages = append(messages, message)
		}

		return nil
	})

	if err != nil {
		c.client.logger.Error("消费消息失败",
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "group", Value: opts.Group},
			logger.Field{Key: "error", Value: err})
		return nil, err
	}

	c.client.logger.Info("消息消费成功",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "count", Value: len(messages)})

	return messages, nil
}

// ConsumeStream 流式消费消息
func (c *Consumer) ConsumeStream(ctx context.Context, opts types.StreamConsumeOptions) (<-chan *StreamMessage, error) {
	if c.client.isClosed() {
		return nil, errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return nil, errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Group == "" {
		opts.Group = "default"
	}

	req := &pb.StreamConsumeRequest{
		Topic:     opts.Topic,
		Group:     opts.Group,
		Offset:    opts.Offset,
		Partition: opts.Partition,
	}

	c.client.logger.Debug("开始流式消费",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset})

	messageChan := make(chan *StreamMessage, 100) // 缓冲通道

	go func() {
		defer close(messageChan)

		err := c.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
			client := pb.NewFluvioServiceClient(conn)

			stream, err := client.StreamConsume(ctx, req)
			if err != nil {
				return errors.Wrap(errors.ErrInternal, "创建流式消费失败", err)
			}

			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				resp, err := stream.Recv()
				if err == io.EOF {
					c.client.logger.Info("流式消费结束",
						logger.Field{Key: "topic", Value: opts.Topic})
					return nil
				}
				if err != nil {
					return errors.Wrap(errors.ErrInternal, "接收流式消息失败", err)
				}

				// 转换消息
				message := &types.Message{
					Topic:     opts.Topic,
					Key:       resp.GetKey(),
					Value:     resp.GetMessage(),
					Headers:   resp.GetHeaders(),
					Offset:    resp.GetOffset(),
					Partition: resp.GetPartition(),
				}

				if resp.GetTimestamp() > 0 {
					message.Timestamp = time.Unix(resp.GetTimestamp(), 0)
				}

				streamMsg := &StreamMessage{
					Message: message,
					Error:   nil,
				}

				select {
				case messageChan <- streamMsg:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})

		if err != nil {
			c.client.logger.Error("流式消费失败",
				logger.Field{Key: "topic", Value: opts.Topic},
				logger.Field{Key: "error", Value: err})

			// 发送错误消息
			streamMsg := &StreamMessage{
				Message: nil,
				Error:   err,
			}
			select {
			case messageChan <- streamMsg:
			default:
			}
		}
	}()

	return messageChan, nil
}

// StreamMessage 流式消息
type StreamMessage struct {
	Message *types.Message
	Error   error
}

// CommitOffset 提交偏移量
func (c *Consumer) CommitOffset(ctx context.Context, opts types.CommitOffsetOptions) error {
	if c.client.isClosed() {
		return errors.New(errors.ErrConnection, "客户端已关闭")
	}

	if opts.Topic == "" {
		return errors.New(errors.ErrInvalidArgument, "主题名称不能为空")
	}

	if opts.Group == "" {
		opts.Group = "default"
	}

	req := &pb.CommitOffsetRequest{
		Topic:     opts.Topic,
		Group:     opts.Group,
		Offset:    opts.Offset,
		Partition: opts.Partition,
	}

	c.client.logger.Debug("提交偏移量",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset})

	err := c.client.withConnection(ctx, func(conn *grpc.ClientConn) error {
		client := pb.NewFluvioServiceClient(conn)

		resp, err := client.CommitOffset(ctx, req)
		if err != nil {
			return errors.Wrap(errors.ErrInternal, "提交偏移量失败", err)
		}

		if !resp.GetSuccess() {
			return errors.New(errors.ErrInternal, resp.GetError())
		}

		return nil
	})

	if err != nil {
		c.client.logger.Error("提交偏移量失败",
			logger.Field{Key: "topic", Value: opts.Topic},
			logger.Field{Key: "group", Value: opts.Group},
			logger.Field{Key: "offset", Value: opts.Offset},
			logger.Field{Key: "error", Value: err})
		return err
	}

	c.client.logger.Info("偏移量提交成功",
		logger.Field{Key: "topic", Value: opts.Topic},
		logger.Field{Key: "group", Value: opts.Group},
		logger.Field{Key: "offset", Value: opts.Offset})

	return nil
}

// ConsumeWithRetry 带重试的消息消费
func (c *Consumer) ConsumeWithRetry(ctx context.Context, opts types.ConsumeOptions) ([]*types.Message, error) {
	var messages []*types.Message
	err := c.client.withRetry(ctx, func(retryCtx context.Context) error {
		var err error
		messages, err = c.Consume(retryCtx, opts)
		return err
	})
	return messages, err
}

// CommitOffsetWithRetry 带重试的偏移量提交
func (c *Consumer) CommitOffsetWithRetry(ctx context.Context, opts types.CommitOffsetOptions) error {
	return c.client.withRetry(ctx, func(retryCtx context.Context) error {
		return c.CommitOffset(retryCtx, opts)
	})
}

// ConsumerConfig 消费者配置
type ConsumerConfig struct {
	DefaultTopic       string        `json:"default_topic"`
	DefaultGroup       string        `json:"default_group"`
	AutoCommit         bool          `json:"auto_commit"`
	AutoCommitInterval time.Duration `json:"auto_commit_interval"`
	MaxMessages        int32         `json:"max_messages"`
	PollTimeout        time.Duration `json:"poll_timeout"`
	EnableMetrics      bool          `json:"enable_metrics"`
}

// DefaultConsumerConfig 返回默认消费者配置
func DefaultConsumerConfig() *ConsumerConfig {
	return &ConsumerConfig{
		DefaultGroup:       "default",
		AutoCommit:         true,
		AutoCommitInterval: 5 * time.Second,
		MaxMessages:        10,
		PollTimeout:        30 * time.Second,
		EnableMetrics:      false,
	}
}

// AdvancedConsumer 高级消费者
type AdvancedConsumer struct {
	consumer *Consumer
	config   *ConsumerConfig

	// 自动提交相关
	autoCommitTicker *time.Ticker
	lastCommitOffset map[string]int64 // topic -> offset
	commitMutex      sync.RWMutex

	// 状态管理
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	closed bool
	mu     sync.RWMutex

	// 指标统计
	metrics ConsumerMetrics
}

// ConsumerMetrics 消费者指标
type ConsumerMetrics struct {
	MessagesConsumed int64 `json:"messages_consumed"`
	BytesConsumed    int64 `json:"bytes_consumed"`
	CommitsSucceeded int64 `json:"commits_succeeded"`
	CommitsFailed    int64 `json:"commits_failed"`
	LastConsumeTime  int64 `json:"last_consume_time"`
	LastCommitTime   int64 `json:"last_commit_time"`
}

// NewAdvancedConsumer 创建高级消费者
func NewAdvancedConsumer(client *Client, config *ConsumerConfig) *AdvancedConsumer {
	if config == nil {
		config = DefaultConsumerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	ac := &AdvancedConsumer{
		consumer:         NewConsumer(client),
		config:           config,
		lastCommitOffset: make(map[string]int64),
		ctx:              ctx,
		cancel:           cancel,
	}

	// 启动自动提交
	if config.AutoCommit {
		ac.startAutoCommit()
	}

	return ac
}

// Consume 消费消息（支持自动提交）
func (ac *AdvancedConsumer) Consume(ctx context.Context, opts types.ConsumeOptions) ([]*types.Message, error) {
	if ac.isClosed() {
		return nil, errors.New(errors.ErrConnection, "消费者已关闭")
	}

	// 设置默认值
	if opts.Topic == "" && ac.config.DefaultTopic != "" {
		opts.Topic = ac.config.DefaultTopic
	}

	if opts.Group == "" {
		opts.Group = ac.config.DefaultGroup
	}

	if opts.MaxMessages <= 0 {
		opts.MaxMessages = ac.config.MaxMessages
	}

	// 消费消息
	messages, err := ac.consumer.Consume(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 更新指标
	if ac.config.EnableMetrics {
		ac.updateMetrics(messages)
	}

	// 更新最后消费的偏移量
	if len(messages) > 0 {
		lastMessage := messages[len(messages)-1]
		ac.updateLastOffset(opts.Topic, lastMessage.Offset)
	}

	return messages, nil
}

// ConsumeStream 流式消费消息（支持自动提交）
func (ac *AdvancedConsumer) ConsumeStream(ctx context.Context, opts types.StreamConsumeOptions) (<-chan *StreamMessage, error) {
	if ac.isClosed() {
		return nil, errors.New(errors.ErrConnection, "消费者已关闭")
	}

	// 设置默认值
	if opts.Topic == "" && ac.config.DefaultTopic != "" {
		opts.Topic = ac.config.DefaultTopic
	}

	if opts.Group == "" {
		opts.Group = ac.config.DefaultGroup
	}

	// 获取原始流
	originalStream, err := ac.consumer.ConsumeStream(ctx, opts)
	if err != nil {
		return nil, err
	}

	// 创建包装的流
	wrappedStream := make(chan *StreamMessage, 100)

	ac.wg.Add(1)
	go func() {
		defer ac.wg.Done()
		defer close(wrappedStream)

		for {
			select {
			case msg, ok := <-originalStream:
				if !ok {
					return
				}

				// 转发消息
				select {
				case wrappedStream <- msg:
					// 更新指标和偏移量
					if msg.Message != nil {
						if ac.config.EnableMetrics {
							ac.updateMetrics([]*types.Message{msg.Message})
						}
						ac.updateLastOffset(opts.Topic, msg.Message.Offset)
					}
				case <-ctx.Done():
					return
				case <-ac.ctx.Done():
					return
				}

			case <-ctx.Done():
				return
			case <-ac.ctx.Done():
				return
			}
		}
	}()

	return wrappedStream, nil
}

// Poll 轮询消费消息
func (ac *AdvancedConsumer) Poll(ctx context.Context, opts types.ConsumeOptions) ([]*types.Message, error) {
	// 创建带超时的上下文
	pollCtx, cancel := context.WithTimeout(ctx, ac.config.PollTimeout)
	defer cancel()

	return ac.Consume(pollCtx, opts)
}

// CommitSync 同步提交偏移量
func (ac *AdvancedConsumer) CommitSync(ctx context.Context, topic, group string) error {
	ac.commitMutex.RLock()
	offset, exists := ac.lastCommitOffset[topic]
	ac.commitMutex.RUnlock()

	if !exists {
		return nil // 没有需要提交的偏移量
	}

	opts := types.CommitOffsetOptions{
		Topic:  topic,
		Group:  group,
		Offset: offset,
	}

	err := ac.consumer.CommitOffset(ctx, opts)
	if err != nil {
		if ac.config.EnableMetrics {
			ac.metrics.CommitsFailed++
		}
		return err
	}

	if ac.config.EnableMetrics {
		ac.metrics.CommitsSucceeded++
		ac.metrics.LastCommitTime = time.Now().Unix()
	}

	return nil
}

// startAutoCommit 启动自动提交
func (ac *AdvancedConsumer) startAutoCommit() {
	ac.autoCommitTicker = time.NewTicker(ac.config.AutoCommitInterval)

	ac.wg.Add(1)
	go func() {
		defer ac.wg.Done()

		for {
			select {
			case <-ac.autoCommitTicker.C:
				ac.commitMutex.RLock()
				topics := make([]string, 0, len(ac.lastCommitOffset))
				for topic := range ac.lastCommitOffset {
					topics = append(topics, topic)
				}
				ac.commitMutex.RUnlock()

				// 为每个主题提交偏移量
				for _, topic := range topics {
					ctx, cancel := context.WithTimeout(ac.ctx, 10*time.Second)
					ac.CommitSync(ctx, topic, ac.config.DefaultGroup)
					cancel()
				}

			case <-ac.ctx.Done():
				return
			}
		}
	}()
}

// updateLastOffset 更新最后消费的偏移量
func (ac *AdvancedConsumer) updateLastOffset(topic string, offset int64) {
	ac.commitMutex.Lock()
	defer ac.commitMutex.Unlock()

	if currentOffset, exists := ac.lastCommitOffset[topic]; !exists || offset > currentOffset {
		ac.lastCommitOffset[topic] = offset
	}
}

// updateMetrics 更新指标
func (ac *AdvancedConsumer) updateMetrics(messages []*types.Message) {
	ac.metrics.MessagesConsumed += int64(len(messages))
	ac.metrics.LastConsumeTime = time.Now().Unix()

	for _, msg := range messages {
		ac.metrics.BytesConsumed += int64(len(msg.Value))
	}
}

// GetMetrics 获取消费者指标
func (ac *AdvancedConsumer) GetMetrics() ConsumerMetrics {
	return ac.metrics
}

// Close 关闭高级消费者
func (ac *AdvancedConsumer) Close() error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if ac.closed {
		return nil
	}

	ac.closed = true

	// 停止自动提交
	if ac.autoCommitTicker != nil {
		ac.autoCommitTicker.Stop()
	}

	// 最后一次同步提交
	if ac.config.AutoCommit {
		ac.commitMutex.RLock()
		for topic := range ac.lastCommitOffset {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			ac.CommitSync(ctx, topic, ac.config.DefaultGroup)
			cancel()
		}
		ac.commitMutex.RUnlock()
	}

	// 取消上下文
	if ac.cancel != nil {
		ac.cancel()
	}

	// 等待所有goroutine完成
	ac.wg.Wait()

	return nil
}

// isClosed 检查是否已关闭
func (ac *AdvancedConsumer) isClosed() bool {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.closed
}

// GetStats 获取消费者统计信息
func (ac *AdvancedConsumer) GetStats() ConsumerStats {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	ac.commitMutex.RLock()
	pendingCommits := len(ac.lastCommitOffset)
	ac.commitMutex.RUnlock()

	return ConsumerStats{
		PendingCommits: pendingCommits,
		Metrics:        ac.GetMetrics(),
		Closed:         ac.closed,
	}
}

// ConsumerStats 消费者统计信息
type ConsumerStats struct {
	PendingCommits int             `json:"pending_commits"`
	Metrics        ConsumerMetrics `json:"metrics"`
	Closed         bool            `json:"closed"`
}
