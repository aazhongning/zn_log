package main

import (
	"time"
	"zn_log/internal"
	"zn_log/internal/async_file"
)

func test(){
	//fileLog := async_file.NewFileLogger("debug", "./", "test", 3*1024) // 向文件打印,非异步
	fileLog := async_file.NewFileLogger("debug", "./", "test", 3*1024) // 向文件打印,异步
	consoleLog := internal.NewConsoleLogger("debug")

	for {
		fileLog.Debug("Debug%v", "试试")
		fileLog.Info("Info")
		fileLog.Warning("Warning")
		fileLog.Error("Error")
		fileLog.Fatal("Fatal")

		consoleLog.Debug("Debug")
		time.Sleep(3*time.Second)
	}
}

func main() {
	test()
}