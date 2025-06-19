package errors

import (
	"fmt"
)

// ErrorCode 错误代码
type ErrorCode string

const (
	// 连接相关错误
	ErrConnection   ErrorCode = "CONNECTION_ERROR"
	ErrTimeout      ErrorCode = "TIMEOUT_ERROR"
	ErrNetworkError ErrorCode = "NETWORK_ERROR"

	// 认证相关错误
	ErrAuthentication ErrorCode = "AUTHENTICATION_ERROR"
	ErrAuthorization  ErrorCode = "AUTHORIZATION_ERROR"

	// 参数相关错误
	ErrInvalidArgument ErrorCode = "INVALID_ARGUMENT"
	ErrMissingArgument ErrorCode = "MISSING_ARGUMENT"

	// 资源相关错误
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrAlreadyExists ErrorCode = "ALREADY_EXISTS"
	ErrResourceLimit ErrorCode = "RESOURCE_LIMIT"

	// 操作相关错误
	ErrInternal    ErrorCode = "INTERNAL_ERROR"
	ErrUnavailable ErrorCode = "UNAVAILABLE"
	ErrCancelled   ErrorCode = "CANCELLED"
	ErrOperation   ErrorCode = "OPERATION_ERROR"

	// 业务相关错误
	ErrValidation    ErrorCode = "VALIDATION_ERROR"
	ErrBusinessLogic ErrorCode = "BUSINESS_LOGIC_ERROR"
)

// FluvioError Fluvio错误类型
type FluvioError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Details map[string]interface{}
}

// Error 实现error接口
func (e *FluvioError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *FluvioError) Unwrap() error {
	return e.Cause
}

// WithDetail 添加详细信息
func (e *FluvioError) WithDetail(key string, value interface{}) *FluvioError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// New 创建新的错误
func New(code ErrorCode, message string) *FluvioError {
	return &FluvioError{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
	}
}

// Wrap 包装错误
func Wrap(code ErrorCode, message string, cause error) *FluvioError {
	return &FluvioError{
		Code:    code,
		Message: message,
		Cause:   cause,
		Details: make(map[string]interface{}),
	}
}

// IsCode 检查错误代码
func IsCode(err error, code ErrorCode) bool {
	if fluvioErr, ok := err.(*FluvioError); ok {
		return fluvioErr.Code == code
	}
	return false
}

// GetCode 获取错误代码
func GetCode(err error) ErrorCode {
	if fluvioErr, ok := err.(*FluvioError); ok {
		return fluvioErr.Code
	}
	return ErrInternal
}

// IsRetryable 检查错误是否可重试
func IsRetryable(err error) bool {
	code := GetCode(err)
	switch code {
	case ErrTimeout, ErrNetworkError, ErrUnavailable, ErrInternal:
		return true
	default:
		return false
	}
}

// IsTemporary 检查错误是否是临时的
func IsTemporary(err error) bool {
	code := GetCode(err)
	switch code {
	case ErrTimeout, ErrNetworkError, ErrUnavailable:
		return true
	default:
		return false
	}
}

// IsPermanent 检查错误是否是永久的
func IsPermanent(err error) bool {
	code := GetCode(err)
	switch code {
	case ErrAuthentication, ErrAuthorization, ErrInvalidArgument, ErrNotFound, ErrAlreadyExists:
		return true
	default:
		return false
	}
}
