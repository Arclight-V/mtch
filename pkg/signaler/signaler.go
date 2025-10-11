package signaler

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	s = make(chan os.Signal, 1)
)

func init() {
	signs := []os.Signal{
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGABRT,
	}
	signal.Notify(s, signs...)
}

func WaitForInterrupt() <-chan os.Signal {
	return s
}
