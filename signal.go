package chooser

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type stopCh struct {
	wg *sync.WaitGroup
	ch chan struct{}
}

func (c *stopCh) close() {
	close(c.ch)
	c.wg.Wait()
}

type winSizeCh struct {
	winSize chan *winSize
	err     chan error
}

func (c *chooser) handleSignals(exitCh chan int, winSizeCh winSizeCh, stopCh stopCh) {
	defer stopCh.wg.Done()

	ch := make(chan os.Signal, 1)
	defer signal.Stop(ch)
	defer close(ch)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGWINCH,
	)

	for {
		select {
		case <-stopCh.ch:
			return
		case signal := <-ch:
			switch signal {
			case syscall.SIGINT, syscall.SIGQUIT:
				exitCh <- 0
			case syscall.SIGTERM:
				exitCh <- 1
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
