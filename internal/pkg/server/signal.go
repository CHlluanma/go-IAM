package server

import (
	"os"
	"os/signal"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownHandler      chan os.Signal
)

// SetUpSignalHandler 注册SIGTERM和SIGINT。返回一个停止通道
// 在其中一个信号时关闭。如果捕捉到第二个信号，程序就停止运行
// 以退出码1结束。
func SetUpSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler)
	shutdownHandler = make(chan os.Signal, 2)

	stop := make(chan struct{})

	signal.Notify(shutdownHandler, shutdownSignals...)

	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1)
	}()

	return stop
}

// RequestShutdown 模拟接收到的作为关机信号的事件
// 返回是否通知了处理程序
func RequestShutdown() bool {
	if shutdownHandler == nil {
		select {
		case shutdownHandler <- shutdownSignals[0]:
			return true
		default:
		}
	}
	return false
}
