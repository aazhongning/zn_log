package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// FileLogger 往文件里面写日志相关代码
type FileLogger struct {
	level       LogLevel
	filePath    string // 日志文件保存路径
	fileName    string // 日志文件名
	maxFileSize int64  // 最大的文件大小
	fileObj     *os.File
	errFileObj  *os.File
}

// NewFileLogger ...
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	logLevel, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}

	fl := &FileLogger{
		level:       logLevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
	}

	err = fl.initFile() // 打开文件,获取文件对象
	if err != nil {
		panic(err)
	}
	return fl
}

func (f *FileLogger) initFile() error {
	// 创建记录正确的日志文件
	fullFileName := filepath.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file faild,err:%v", err)
		return err
	}
	// 创建错误的日志文件
	errFileObj, err := os.OpenFile(fullFileName+"_err.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open err log file faild,err:%v", err)
		return err
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}

func (f *FileLogger) enable(logLevel LogLevel) bool {
	return f.level <= logLevel
}

func (f *FileLogger) log(lv LogLevel, format string, args ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, args...)      // 合并输出
		funcName, fileName, lineNo := getInfo(3) // 三层调用
		now := time.Now().Format("2006-01-02 03:04:06")
		lvStr := getLogString(lv)
		if f.checkSize(f.fileObj) {
			newFile, err := f.splitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newFile
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s \n", now, lvStr, fileName, funcName, lineNo, msg)
		if lv >= ERROR {
			if f.checkSize(f.errFileObj) {
				newFile, err := f.splitFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = newFile
			}
			// 如果记录日志级别大于或等于ERROR，则再记录一份到LogErr的文件中
			fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s \n", now, lvStr, fileName, funcName, lineNo, msg)
		}
	}
}

// Debug ...
func (f *FileLogger) Debug(format string, args ...interface{}) {
	f.log(DEBUG, format, args...)
}

// Trace ...
func (f *FileLogger) Trace(format string, args ...interface{}) {
	f.log(TRACE, format, args...)
}

// Info ...
func (f *FileLogger) Info(format string, args ...interface{}) {
	f.log(INFO, format, args...)
}

// Warning ...
func (f *FileLogger) Warning(format string, args ...interface{}) {
	f.log(WARNING, format, args...)
}

// Error ...
func (f *FileLogger) Error(format string, args ...interface{}) {
	f.log(ERROR, format, args...)
}

// Fatal ...
func (f *FileLogger) Fatal(format string, args ...interface{}) {
	f.log(FATAL, format, args...)
}

// Close 关闭文件资源
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}

// 获取文件大小，判断是否要进行切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err%v\n", err)
		return false
	}
	// 如果当前文件的size大于设定的size，则返回true，否则返回false
	return fileInfo.Size() >= f.maxFileSize
}

func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	nowStr := time.Now().Format("20060102150405000")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return nil, err
	}
	logName := filepath.Join(f.filePath, fileInfo.Name())
	newlogName := fmt.Sprintf("%s.%s.bak", logName, nowStr)
	// 1. 关闭当前文件
	file.Close()
	// 2. 备份一个 rename
	os.Rename(logName, newlogName)
	// 3. 打开一个新的日志文件

	fileObj, err := os.OpenFile(logName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed, err:%v", err)
		return nil, err
	}
	// 4. 将打开的文件赋值给 fl.FileObj
	return fileObj, nil
}

