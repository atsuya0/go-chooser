package selector

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func (s *selector) handleSignals(exitCh chan int, winSizeCh chan *winSize, stopCh chan struct{}) {
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
			log.Println("[INFO] stop handleSignals")
			return
		case signal := <-sigCh:
			switch signal {
			case syscall.SIGINT:
				log.Println("[SIGNAL] Catch SIGINT")
				exitCh <- 0

			case syscall.SIGTERM:
				log.Println("[SIGNAL] Catch SIGTERM")
				exitCh <- 1

			case syscall.SIGQUIT:
				log.Println("[SIGNAL] Catch SIGQUIT")
				exitCh <- 0

			case syscall.SIGWINCH:
				log.Println("[SIGNAL] Catch SIGWINCH")
				winSizeCh <- s.terminal.getWinSize()
			}
		}
	}
}
