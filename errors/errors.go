package errors

import (
	"fmt"
)

// ErrorCode 错误代码
type ErrorCode int

const (
	// ErrUnknown 未知错误
	ErrUnknown ErrorCode = iota
	// ErrConnection 连接错误
	ErrConnection
	// ErrTimeout 超时错误
	ErrTimeout
	// ErrInvalidConfig 无效配置
	ErrInvalidConfig
	// ErrTopicNotFound 主题不存在
	ErrTopicNotFound
	// ErrPermissionDenied 权限拒绝
	ErrPermissionDenied
	// ErrInvalidArgument 无效参数
	ErrInvalidArgument
	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable
	// ErrResourceExhausted 资源耗尽
	ErrResourceExhausted
	// ErrInternal 内部错误
	ErrInternal
)

// Error SDK错误类型
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
}

// Error 实现error接口
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code.String(), e.Message)
}

// Unwrap 返回底层错误
func (e *Error) Unwrap() error {
	return e.Cause
}

// String 返回错误代码的字符串表示
func (c ErrorCode) String() string {
	switch c {
	case ErrUnknown:
		return "UNKNOWN"
	case ErrConnection:
		return "CONNECTION"
	case ErrTimeout:
		return "TIMEOUT"
	case ErrInvalidConfig:
		return "INVALID_CONFIG"
	case ErrTopicNotFound:
		return "TOPIC_NOT_FOUND"
	case ErrPermissionDenied:
		return "PERMISSION_DENIED"
	case ErrInvalidArgument:
		return "INVALID_ARGUMENT"
	case ErrServiceUnavailable:
		return "SERVICE_UNAVAILABLE"
	case ErrResourceExhausted:
		return "RESOURCE_EXHAUSTED"
	case ErrInternal:
		return "INTERNAL"
	default:
		return "UNKNOWN"
	}
}

// New 创建一个新的SDK错误
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装一个错误
func Wrap(code ErrorCode, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsCode 检查错误是否为指定的错误代码
func IsCode(err error, code ErrorCode) bool {
	if sdkErr, ok := err.(*Error); ok {
		return sdkErr.Code == code
	}
	return false
}

// GetCode 获取错误代码
func GetCode(err error) ErrorCode {
	if sdkErr, ok := err.(*Error); ok {
		return sdkErr.Code
	}
	return ErrUnknown
}

// 预定义的常用错误
var (
	ErrConnectionFailed    = New(ErrConnection, "连接失败")
	ErrTimeoutExceeded     = New(ErrTimeout, "操作超时")
	ErrConfigInvalid       = New(ErrInvalidConfig, "配置无效")
	ErrTopicNotExists      = New(ErrTopicNotFound, "主题不存在")
	ErrAccessDenied        = New(ErrPermissionDenied, "访问被拒绝")
	ErrInvalidParam        = New(ErrInvalidArgument, "参数无效")
	ErrServiceDown         = New(ErrServiceUnavailable, "服务不可用")
	ErrResourceLimit       = New(ErrResourceExhausted, "资源限制")
	ErrInternalFailure     = New(ErrInternal, "内部错误")
)
