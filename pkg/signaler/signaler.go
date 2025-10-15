package signaler

import (
	"fmt"
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

func waitForInterrupt() <-chan os.Signal {
	return s
}

func WaitForInterrupt(cancel <-chan struct{}) error {
	interrupt := waitForInterrupt()
	select {
	case s := <-interrupt:
		return fmt.Errorf("received signal: %s", s)
	case <-cancel:
		return fmt.Errorf("Captured %v, shutdown requested.\n", interrupt)
	}
}
