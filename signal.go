package chooser

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type winSizeCh struct {
	winSize chan *winSize
	err     chan error
}

func (c *chooser) handleSignals(exitCh chan int, winSizeCh winSizeCh, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

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
				if winSize, err := c.terminal.getWinSize(); err != nil {
					winSizeCh.err <- err
				} else {
					winSizeCh.winSize <- winSize
				}
			}
		}
	}
}
