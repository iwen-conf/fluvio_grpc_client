package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String 返回级别字符串
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
	case LevelFatal:
		return "FATAL"
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
	Fatal(msg string, fields ...Field)
	SetLevel(level Level)
	GetLevel() Level
	WithFields(fields ...Field) Logger
}

// StandardLogger 标准日志实现
type StandardLogger struct {
	logger *log.Logger
	level  Level
	fields []Field
}

// NewStandardLogger 创建标准日志器
func NewStandardLogger(output io.Writer, level Level) *StandardLogger {
	return &StandardLogger{
		logger: log.New(output, "", log.LstdFlags),
		level:  level,
		fields: make([]Field, 0),
	}
}

// NewDefaultLogger 创建默认日志器
func NewDefaultLogger() *StandardLogger {
	return NewStandardLogger(os.Stdout, LevelInfo)
}

// Debug 调试日志
func (l *StandardLogger) Debug(msg string, fields ...Field) {
	if l.level <= LevelDebug {
		l.log(LevelDebug, msg, fields...)
	}
}

// Info 信息日志
func (l *StandardLogger) Info(msg string, fields ...Field) {
	if l.level <= LevelInfo {
		l.log(LevelInfo, msg, fields...)
	}
}

// Warn 警告日志
func (l *StandardLogger) Warn(msg string, fields ...Field) {
	if l.level <= LevelWarn {
		l.log(LevelWarn, msg, fields...)
	}
}

// Error 错误日志
func (l *StandardLogger) Error(msg string, fields ...Field) {
	if l.level <= LevelError {
		l.log(LevelError, msg, fields...)
	}
}

// Fatal 致命错误日志
func (l *StandardLogger) Fatal(msg string, fields ...Field) {
	l.log(LevelFatal, msg, fields...)
	os.Exit(1)
}

// SetLevel 设置日志级别
func (l *StandardLogger) SetLevel(level Level) {
	l.level = level
}

// GetLevel 获取日志级别
func (l *StandardLogger) GetLevel() Level {
	return l.level
}

// WithFields 添加字段
func (l *StandardLogger) WithFields(fields ...Field) Logger {
	newLogger := &StandardLogger{
		logger: l.logger,
		level:  l.level,
		fields: make([]Field, len(l.fields)+len(fields)),
	}

	copy(newLogger.fields, l.fields)
	copy(newLogger.fields[len(l.fields):], fields)

	return newLogger
}

// log 记录日志
func (l *StandardLogger) log(level Level, msg string, fields ...Field) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 合并字段
	allFields := make([]Field, len(l.fields)+len(fields))
	copy(allFields, l.fields)
	copy(allFields[len(l.fields):], fields)

	// 构建日志消息
	logMsg := fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), msg)

	// 添加字段
	if len(allFields) > 0 {
		logMsg += " |"
		for _, field := range allFields {
			logMsg += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	l.logger.Println(logMsg)
}

// ParseLevel 解析日志级别
func ParseLevel(levelStr string) (Level, error) {
	switch levelStr {
	case "debug", "DEBUG":
		return LevelDebug, nil
	case "info", "INFO":
		return LevelInfo, nil
	case "warn", "WARN":
		return LevelWarn, nil
	case "error", "ERROR":
		return LevelError, nil
	case "fatal", "FATAL":
		return LevelFatal, nil
	default:
		return LevelInfo, fmt.Errorf("unknown log level: %s", levelStr)
	}
}
