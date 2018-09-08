package choice

import (
	"os"
	"os/signal"
	"syscall"
)

func (c *chooser) handleSignals(exitCh chan int, winSizeCh chan *winSize, stopCh chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		select {
		case <-stopCh:
			return
		case signal := <-sigCh:
			switch signal {
			case syscall.SIGINT:
				exitCh <- 0
			case syscall.SIGTERM:
				exitCh <- 1
			case syscall.SIGQUIT:
				exitCh <- 0
			case syscall.SIGWINCH:
				winSizeCh <- c.terminal.getWinSize()
			}
		}
	}
}
