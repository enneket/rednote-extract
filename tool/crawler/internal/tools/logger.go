package tools

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger 日志接口
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// ConsoleLogger 控制台日志实现
type ConsoleLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewLogger 创建日志记录器实例
func NewLogger() Logger {
	return &ConsoleLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info 记录信息日志
func (l *ConsoleLogger) Info(format string, args ...interface{}) {
	l.infoLogger.Printf(getLogPrefix()+format, args...)
}

// Error 记录错误日志
func (l *ConsoleLogger) Error(format string, args ...interface{}) {
	l.errorLogger.Printf(getLogPrefix()+format, args...)
}

// getLogPrefix 获取日志前缀
func getLogPrefix() string {
	return fmt.Sprintf("[%s] ", time.Now().Format("2006-01-02 15:04:05"))
}
