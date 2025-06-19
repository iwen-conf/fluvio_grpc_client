package utils

import (
	"fmt"

	"github.com/iwen-conf/fluvio_grpc_client/infrastructure/logging"
	"github.com/iwen-conf/fluvio_grpc_client/pkg/errors"
)

// GRPCResponseHandler gRPC响应处理器
type GRPCResponseHandler struct {
	logger logging.Logger
}

// NewGRPCResponseHandler 创建gRPC响应处理器
func NewGRPCResponseHandler(logger logging.Logger) *GRPCResponseHandler {
	return &GRPCResponseHandler{
		logger: logger,
	}
}

// HandleError 处理gRPC错误
func (h *GRPCResponseHandler) HandleError(err error, operation string, context map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// 构建日志字段
	fields := []logging.Field{
		{Key: "operation", Value: operation},
		{Key: "error", Value: err},
	}
	
	// 添加上下文字段
	for key, value := range context {
		fields = append(fields, logging.Field{Key: key, Value: value})
	}

	h.logger.Error(fmt.Sprintf("%s失败", operation), fields...)
	
	// 包装为Fluvio错误
	return errors.Wrap(errors.ErrOperation, fmt.Sprintf("%s failed", operation), err)
}

// HandleSuccessResponse 处理成功响应
func (h *GRPCResponseHandler) HandleSuccessResponse(operation string, context map[string]interface{}) {
	// 构建日志字段
	fields := []logging.Field{
		{Key: "operation", Value: operation},
	}
	
	// 添加上下文字段
	for key, value := range context {
		fields = append(fields, logging.Field{Key: key, Value: value})
	}

	h.logger.Info(fmt.Sprintf("%s成功", operation), fields...)
}

// ValidateResponse 验证响应并处理错误
func (h *GRPCResponseHandler) ValidateResponse(success bool, errorMsg, operation string, context map[string]interface{}) error {
	if success {
		h.HandleSuccessResponse(operation, context)
		return nil
	}

	if errorMsg == "" {
		errorMsg = "unknown error"
	}

	// 构建日志字段
	fields := []logging.Field{
		{Key: "operation", Value: operation},
		{Key: "error", Value: errorMsg},
	}
	
	// 添加上下文字段
	for key, value := range context {
		fields = append(fields, logging.Field{Key: key, Value: value})
	}

	h.logger.Error(fmt.Sprintf("%s被服务器拒绝", operation), fields...)
	
	return errors.New(errors.ErrOperation, fmt.Sprintf("%s failed: %s", operation, errorMsg))
}

// LogDebugOperation 记录调试操作
func (h *GRPCResponseHandler) LogDebugOperation(operation string, context map[string]interface{}) {
	// 构建日志字段
	fields := []logging.Field{
		{Key: "operation", Value: operation},
	}
	
	// 添加上下文字段
	for key, value := range context {
		fields = append(fields, logging.Field{Key: key, Value: value})
	}

	h.logger.Debug(operation, fields...)
}

// ResponseValidator 响应验证器接口
type ResponseValidator interface {
	GetSuccess() bool
	GetError() string
}

// ValidateGRPCResponse 通用gRPC响应验证
func (h *GRPCResponseHandler) ValidateGRPCResponse(resp ResponseValidator, operation string, context map[string]interface{}) error {
	return h.ValidateResponse(resp.GetSuccess(), resp.GetError(), operation, context)
}

// ContextBuilder 上下文构建器
type ContextBuilder struct {
	context map[string]interface{}
}

// NewContextBuilder 创建上下文构建器
func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		context: make(map[string]interface{}),
	}
}

// Add 添加上下文字段
func (cb *ContextBuilder) Add(key string, value interface{}) *ContextBuilder {
	cb.context[key] = value
	return cb
}

// AddIf 条件添加上下文字段
func (cb *ContextBuilder) AddIf(condition bool, key string, value interface{}) *ContextBuilder {
	if condition {
		cb.context[key] = value
	}
	return cb
}

// Build 构建上下文
func (cb *ContextBuilder) Build() map[string]interface{} {
	return cb.context
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int
	FailureCount int
	Errors       []error
}

// NewBatchOperationResult 创建批量操作结果
func NewBatchOperationResult() *BatchOperationResult {
	return &BatchOperationResult{
		Errors: make([]error, 0),
	}
}

// AddSuccess 添加成功计数
func (r *BatchOperationResult) AddSuccess() {
	r.SuccessCount++
}

// AddFailure 添加失败计数和错误
func (r *BatchOperationResult) AddFailure(err error) {
	r.FailureCount++
	if err != nil {
		r.Errors = append(r.Errors, err)
	}
}

// HasFailures 检查是否有失败
func (r *BatchOperationResult) HasFailures() bool {
	return r.FailureCount > 0
}

// GetSummaryError 获取汇总错误
func (r *BatchOperationResult) GetSummaryError() error {
	if !r.HasFailures() {
		return nil
	}
	
	return errors.New(errors.ErrOperation, 
		fmt.Sprintf("batch operation partially failed: %d success, %d failure", 
			r.SuccessCount, r.FailureCount))
}

// LogSummary 记录汇总日志
func (r *BatchOperationResult) LogSummary(handler *GRPCResponseHandler, operation string, context map[string]interface{}) {
	summaryContext := make(map[string]interface{})
	for k, v := range context {
		summaryContext[k] = v
	}
	summaryContext["success_count"] = r.SuccessCount
	summaryContext["failure_count"] = r.FailureCount
	
	if r.HasFailures() {
		handler.logger.Warn(fmt.Sprintf("%s部分失败", operation), 
			handler.buildLogFields("batch_summary", summaryContext)...)
	} else {
		handler.HandleSuccessResponse(operation, summaryContext)
	}
}

// buildLogFields 构建日志字段（内部方法）
func (h *GRPCResponseHandler) buildLogFields(operation string, context map[string]interface{}) []logging.Field {
	fields := []logging.Field{
		{Key: "operation", Value: operation},
	}
	
	for key, value := range context {
		fields = append(fields, logging.Field{Key: key, Value: value})
	}
	
	return fields
}
