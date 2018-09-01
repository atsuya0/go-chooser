package selector

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type selector struct {
	terminal *terminal
	render   *render
	list     []string
}

func NewSelector(list []string) *selector {
	return &selector{
		terminal: newTerminal(),
		render:   newRender(),
		list:     list,
	}
}

func (s *selector) init() {
	s.terminal.setup()
	s.render.winSize = s.terminal.getWinSize()
	s.filter()
	s.render.render()
}

func (s *selector) filter() {
	var result []string
	for _, v := range s.list {
		if strings.Contains(v, s.render.buffer.text) {
			result = append(result, v)
		}
	}
	s.render.completion = newCompletion(s.render.limit(result))
}

func (s *selector) readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	log.Println("[INFO] readBuffer start")
	for {
		select {
		case <-stopCh:
			log.Println("[INFO] stop readBuffer")
			return
		default:
			if b, err := s.terminal.read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *selector) response(b []byte) (bool, string) {
	switch key := getKey(b); key {
	case ignore:
		log.Println("[INFO] Not defined.")
	case displayable:
		s.render.buffer.insert(string(b))
		s.filter()
	case enter:
		if 0 <= s.render.completion.target {
			return true, s.render.completion.suggestions[s.render.completion.target]
		}
		return true, ""
	case controlC:
		return true, ""
	case controlN:
		s.render.completion.next()
	case controlP:
		s.render.completion.previous()
	default:
		if function, ok := keyBindCmds[key]; ok {
			function(s.render.buffer)
			s.filter()
		}
	}

	return false, ""
}

func (s *selector) Run() string {
	if l := os.Getenv("LOG_PATH"); l == "" {
		log.SetOutput(ioutil.Discard)
	} else if f, err := os.OpenFile(l, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("[INFO] Logging is enabled.")
	}

	s.init()
	defer s.terminal.restore()

	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go s.readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *winSize)
	stopHandleSignalCh := make(chan struct{})
	go s.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, text := s.response(b); shouldExit {
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				s.render.clear()
				return text
			}
			s.render.render()

		case code := <-exitCh:
			os.Exit(code)

		case w := <-winSizeCh:
			s.render.winSize = w
			s.render.render()

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
