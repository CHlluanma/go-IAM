package server

import (
	"os"
	"syscall"
)

// 关闭信号： os.Interrupt  syscall.SIGTERM
var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
