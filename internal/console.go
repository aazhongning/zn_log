package internal

import (
	"fmt"
	"time"
)

// 往终端写日志

// ConsoleLogger ...
type ConsoleLogger struct {
	level LogLevel
}

// NewConsoleLogger 构造函数 ...
func NewConsoleLogger(levelStr string) ConsoleLogger {
	level, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{
		level: level,
	}
}

func (c ConsoleLogger) enable(logLevel LogLevel) bool {
	return c.level <= logLevel
}

func (c ConsoleLogger) log(lv LogLevel, format string, args ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, args...)      // 合并输出
		funcName, fileName, lineNo := getInfo(3) // 三层调用
		now := time.Now().Format("2006-01-02 03:04:06")
		lvStr := getLogString(lv)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s \n", now, lvStr, fileName, funcName, lineNo, msg)
	}
}

// Debug ...
func (c ConsoleLogger) Debug(format string, args ...interface{}) {
	c.log(DEBUG, format, args...)
}

// TRACE ...
func (c ConsoleLogger) Trace(format string, args ...interface{}) {
	c.log(TRACE, format, args...)
}

// Info ...
func (c ConsoleLogger) Info(format string, args ...interface{}) {
	c.log(INFO, format, args...)
}

// Warning ...
func (c ConsoleLogger) Warning(format string, args ...interface{}) {
	c.log(WARNING, format, args...)
}

// Error ...
func (c ConsoleLogger) Error(format string, args ...interface{}) {
	c.log(ERROR, format, args...)

}

// Fatal ...
func (c ConsoleLogger) Fatal(format string, args ...interface{}) {
	c.log(FATAL, format, args...)
}


