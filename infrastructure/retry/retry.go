package retry

import (
	"context"
	"fmt"
	"time"

	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts int           // 最大重试次数
	BaseDelay   time.Duration // 基础延迟时间
	MaxDelay    time.Duration // 最大延迟时间
	Multiplier  float64       // 延迟倍数
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}
}

// RetryableFunc 可重试的函数类型
type RetryableFunc func() error

// IsRetryableError 判断错误是否可重试
type IsRetryableError func(error) bool

// DefaultIsRetryableError 默认的可重试错误判断
func DefaultIsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// 检查常见的可重试错误
	errStr := err.Error()
	
	// 网络相关错误
	if contains(errStr, "connection refused") ||
		contains(errStr, "connection reset") ||
		contains(errStr, "timeout") ||
		contains(errStr, "temporary failure") ||
		contains(errStr, "service unavailable") ||
		contains(errStr, "deadline exceeded") {
		return true
	}
	
	return false
}

// contains 检查字符串是否包含子字符串（忽略大小写）
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		(s == substr || 
		 (len(s) > len(substr) && 
		  (s[:len(substr)] == substr || 
		   s[len(s)-len(substr):] == substr ||
		   containsInner(s, substr))))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Retry 执行重试逻辑
func Retry(ctx context.Context, config *RetryConfig, isRetryable IsRetryableError, fn RetryableFunc, logger logging.Logger) error {
	if config == nil {
		config = DefaultRetryConfig()
	}
	
	if isRetryable == nil {
		isRetryable = DefaultIsRetryableError
	}
	
	var lastErr error
	delay := config.BaseDelay
	
	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		// 执行函数
		err := fn()
		if err == nil {
			if attempt > 1 {
				logger.Info("重试成功", 
					logging.Field{Key: "attempt", Value: attempt},
					logging.Field{Key: "total_attempts", Value: config.MaxAttempts})
			}
			return nil
		}
		
		lastErr = err
		
		// 检查是否可重试
		if !isRetryable(err) {
			logger.Debug("错误不可重试", 
				logging.Field{Key: "error", Value: err},
				logging.Field{Key: "attempt", Value: attempt})
			return err
		}
		
		// 如果是最后一次尝试，直接返回错误
		if attempt == config.MaxAttempts {
			logger.Error("重试次数已用完", 
				logging.Field{Key: "error", Value: err},
				logging.Field{Key: "attempts", Value: attempt})
			return fmt.Errorf("max retry attempts (%d) exceeded: %w", config.MaxAttempts, err)
		}
		
		// 记录重试日志
		logger.Warn("操作失败，准备重试", 
			logging.Field{Key: "error", Value: err},
			logging.Field{Key: "attempt", Value: attempt},
			logging.Field{Key: "delay", Value: delay})
		
		// 等待延迟时间
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		
		// 计算下次延迟时间（指数退避）
		delay = time.Duration(float64(delay) * config.Multiplier)
		if delay > config.MaxDelay {
			delay = config.MaxDelay
		}
	}
	
	return lastErr
}