package utils

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RetryableFunc 可重试的函数类型
type RetryableFunc func() error

// RetryableContextFunc 带上下文的可重试函数类型
type RetryableContextFunc func(ctx context.Context) error

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries       int           `json:"max_retries"`
	InitialBackoff   time.Duration `json:"initial_backoff"`
	MaxBackoff       time.Duration `json:"max_backoff"`
	BackoffMultiple  float64       `json:"backoff_multiple"`
	EnableJitter     bool          `json:"enable_jitter"`
	JitterFactor     float64       `json:"jitter_factor"`
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:      3,
		InitialBackoff:  1 * time.Second,
		MaxBackoff:      30 * time.Second,
		BackoffMultiple: 2.0,
		EnableJitter:    true,
		JitterFactor:    0.1,
	}
}

// Retryer 重试器
type Retryer struct {
	config *RetryConfig
	logger logging.Logger
}

// NewRetryer 创建重试器
func NewRetryer(config *RetryConfig, logger logging.Logger) *Retryer {
	if config == nil {
		config = DefaultRetryConfig()
	}
	return &Retryer{
		config: config,
		logger: logger,
	}
}

// Retry 执行重试
func (r *Retryer) Retry(fn RetryableFunc) error {
	return r.RetryWithContext(context.Background(), func(ctx context.Context) error {
		return fn()
	})
}

// RetryWithContext 带上下文执行重试
func (r *Retryer) RetryWithContext(ctx context.Context, fn RetryableContextFunc) error {
	var lastErr error
	backoff := r.config.InitialBackoff

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			r.logger.Debug("重试操作",
				logging.Field{Key: "attempt", Value: attempt},
				logging.Field{Key: "backoff", Value: backoff},
				logging.Field{Key: "last_error", Value: lastErr})

			// 等待退避时间
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}

			// 计算下次退避时间
			backoff = r.calculateNextBackoff(backoff)
		}

		// 执行函数
		err := fn(ctx)
		if err == nil {
			if attempt > 0 {
				r.logger.Info("重试成功", logging.Field{Key: "attempts", Value: attempt + 1})
			}
			return nil
		}

		lastErr = err

		// 检查是否应该重试
		if !r.shouldRetry(err) {
			r.logger.Debug("错误不可重试", logging.Field{Key: "error", Value: err})
			break
		}

		// 检查上下文是否已取消
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	r.logger.Error("重试失败",
		logging.Field{Key: "max_attempts", Value: r.config.MaxRetries + 1},
		logging.Field{Key: "final_error", Value: lastErr})

	return lastErr
}

// calculateNextBackoff 计算下次退避时间
func (r *Retryer) calculateNextBackoff(current time.Duration) time.Duration {
	next := time.Duration(float64(current) * r.config.BackoffMultiple)
	if next > r.config.MaxBackoff {
		next = r.config.MaxBackoff
	}

	// 添加抖动
	if r.config.EnableJitter && r.config.JitterFactor > 0 {
		jitterAmount := float64(next) * r.config.JitterFactor
		jitterTime := time.Duration(jitterAmount * (2*rand.Float64() - 1))
		next += jitterTime

		if next < 0 {
			next = current / 2
		}
	}

	return next
}

// shouldRetry 判断错误是否应该重试
func (r *Retryer) shouldRetry(err error) bool {
	// 检查SDK错误
	if fluvioErr, ok := err.(*errors.FluvioError); ok {
		return errors.IsRetryable(fluvioErr)
	}

	// 检查gRPC状态错误
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Aborted:
			return true
		case codes.InvalidArgument, codes.NotFound, codes.PermissionDenied, codes.Unauthenticated:
			return false
		default:
			return true
		}
	}

	// 默认重试
	return true
}

// BackoffStrategy 退避策略接口
type BackoffStrategy interface {
	Next() time.Duration
	Reset()
}

// ExponentialBackoff 指数退避策略
type ExponentialBackoff struct {
	Initial    time.Duration
	Max        time.Duration
	Multiplier float64
	current    time.Duration
}

// NewExponentialBackoff 创建指数退避策略
func NewExponentialBackoff(initial, max time.Duration, multiplier float64) *ExponentialBackoff {
	return &ExponentialBackoff{
		Initial:    initial,
		Max:        max,
		Multiplier: multiplier,
		current:    initial,
	}
}

// Next 获取下一个退避时间
func (eb *ExponentialBackoff) Next() time.Duration {
	current := eb.current
	eb.current = time.Duration(float64(eb.current) * eb.Multiplier)
	if eb.current > eb.Max {
		eb.current = eb.Max
	}
	return current
}

// Reset 重置退避时间
func (eb *ExponentialBackoff) Reset() {
	eb.current = eb.Initial
}

// LinearBackoff 线性退避策略
type LinearBackoff struct {
	Initial   time.Duration
	Max       time.Duration
	Increment time.Duration
	current   time.Duration
}

// NewLinearBackoff 创建线性退避策略
func NewLinearBackoff(initial, max, increment time.Duration) *LinearBackoff {
	return &LinearBackoff{
		Initial:   initial,
		Max:       max,
		Increment: increment,
		current:   initial,
	}
}

// Next 获取下一个退避时间
func (lb *LinearBackoff) Next() time.Duration {
	current := lb.current
	lb.current += lb.Increment
	if lb.current > lb.Max {
		lb.current = lb.Max
	}
	return current
}

// Reset 重置退避时间
func (lb *LinearBackoff) Reset() {
	lb.current = lb.Initial
}

// JitteredBackoff 带抖动的退避策略
type JitteredBackoff struct {
	base   *ExponentialBackoff
	jitter float64 // 抖动因子 (0.0 - 1.0)
}

// NewJitteredBackoff 创建带抖动的退避策略
func NewJitteredBackoff(initial, max time.Duration, multiplier, jitter float64) *JitteredBackoff {
	return &JitteredBackoff{
		base:   NewExponentialBackoff(initial, max, multiplier),
		jitter: jitter,
	}
}

// Next 获取下一个退避时间（带抖动）
func (jb *JitteredBackoff) Next() time.Duration {
	baseTime := jb.base.Next()
	if jb.jitter <= 0 {
		return baseTime
	}

	// 添加随机抖动
	jitterAmount := float64(baseTime) * jb.jitter
	jitterTime := time.Duration(jitterAmount * (2*rand.Float64() - 1))

	result := baseTime + jitterTime
	if result < 0 {
		result = baseTime / 2
	}

	return result
}

// Reset 重置退避时间
func (jb *JitteredBackoff) Reset() {
	jb.base.Reset()
}

// RetryWithBackoff 使用自定义退避策略重试
func RetryWithBackoff(ctx context.Context, fn RetryableContextFunc, strategy BackoffStrategy, maxRetries int, logger logging.Logger) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			backoff := strategy.Next()
			
			if logger != nil {
				logger.Debug("重试操作",
					logging.Field{Key: "attempt", Value: attempt},
					logging.Field{Key: "backoff", Value: backoff})
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		}

		err := fn(ctx)
		if err == nil {
			if attempt > 0 && logger != nil {
				logger.Info("重试成功", logging.Field{Key: "attempts", Value: attempt + 1})
			}
			return nil
		}

		lastErr = err

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if logger != nil {
		logger.Error("重试失败",
			logging.Field{Key: "max_attempts", Value: maxRetries + 1},
			logging.Field{Key: "final_error", Value: lastErr})
	}

	return lastErr
}