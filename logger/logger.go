package logger

import (
	"fmt"
	"log"
	"os"
)

// Level 日志级别
type Level int

const (
	// LevelDebug 调试级别
	LevelDebug Level = iota
	// LevelInfo 信息级别
	LevelInfo
	// LevelWarn 警告级别
	LevelWarn
	// LevelError 错误级别
	LevelError
)

// String 返回日志级别的字符串表示
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	SetLevel(level Level)
	GetLevel() Level
}

// DefaultLogger 默认日志实现
type DefaultLogger struct {
	level  Level
	logger *log.Logger
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger(level Level) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Debug 输出调试日志
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	if l.level <= LevelDebug {
		l.output(LevelDebug, msg, fields...)
	}
}

// Info 输出信息日志
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	if l.level <= LevelInfo {
		l.output(LevelInfo, msg, fields...)
	}
}

// Warn 输出警告日志
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	if l.level <= LevelWarn {
		l.output(LevelWarn, msg, fields...)
	}
}

// Error 输出错误日志
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	if l.level <= LevelError {
		l.output(LevelError, msg, fields...)
	}
}

// SetLevel 设置日志级别
func (l *DefaultLogger) SetLevel(level Level) {
	l.level = level
}

// GetLevel 获取日志级别
func (l *DefaultLogger) GetLevel() Level {
	return l.level
}

// output 输出日志
func (l *DefaultLogger) output(level Level, msg string, fields ...Field) {
	var fieldsStr string
	if len(fields) > 0 {
		fieldsStr = " ["
		for i, field := range fields {
			if i > 0 {
				fieldsStr += ", "
			}
			fieldsStr += fmt.Sprintf("%s=%v", field.Key, field.Value)
		}
		fieldsStr += "]"
	}
	
	l.logger.Printf("[%s] %s%s", level.String(), msg, fieldsStr)
}

// NoopLogger 空日志实现
type NoopLogger struct{}

// NewNoopLogger 创建空日志器
func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

// Debug 空实现
func (l *NoopLogger) Debug(msg string, fields ...Field) {}

// Info 空实现
func (l *NoopLogger) Info(msg string, fields ...Field) {}

// Warn 空实现
func (l *NoopLogger) Warn(msg string, fields ...Field) {}

// Error 空实现
func (l *NoopLogger) Error(msg string, fields ...Field) {}

// SetLevel 空实现
func (l *NoopLogger) SetLevel(level Level) {}

// GetLevel 空实现
func (l *NoopLogger) GetLevel() Level {
	return LevelError
}

// 全局默认日志器
var defaultLogger Logger = NewDefaultLogger(LevelInfo)

// SetDefault 设置默认日志器
func SetDefault(logger Logger) {
	defaultLogger = logger
}

// GetDefault 获取默认日志器
func GetDefault() Logger {
	return defaultLogger
}

// Debug 使用默认日志器输出调试日志
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

// Info 使用默认日志器输出信息日志
func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

// Warn 使用默认日志器输出警告日志
func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

// Error 使用默认日志器输出错误日志
func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}
