package chooser

import (
	"os"
	"os/signal"
	"syscall"
)

func (c *chooser) handleSignals(exitCh chan int, winSizeCh chan *winSize, stopCh chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		select {
		case <-stopCh:
			return
		case signal := <-ch:
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
