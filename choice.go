package choice

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type chooser struct {
	terminal *terminal
	render   *render
	list     []string
}

func NewChooser(list []string) *chooser {
	return &chooser{
		terminal: newTerminal(),
		render:   newRender(),
		list:     list,
	}
}

func (c *chooser) init() {
	c.terminal.setup()
	c.render.winSize = c.terminal.getWinSize()
	c.filter()
	c.render.render()
}

func (c *chooser) filter() {
	var result []string
	for _, v := range c.list {
		if strings.Contains(v, c.render.buffer.text) {
			result = append(result, v)
		}
	}
	c.render.completion = newCompletion(c.render.limit(result))
}

func (c *chooser) readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	log.Println("[INFO] readBuffer start")
	for {
		select {
		case <-stopCh:
			log.Println("[INFO] stop readBuffer")
			return
		default:
			if b, err := c.terminal.read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *chooser) response(b []byte) (bool, string) {
	switch key := getKey(b); key {
	case ignore:
		log.Println("[INFO] Not defined.")
	case displayable:
		c.render.buffer.insert(string(b))
		c.filter()
	case enter:
		if 0 <= c.render.completion.target {
			return true, c.render.completion.suggestions[c.render.completion.target]
		}
		return true, ""
	case controlC:
		return true, ""
	case controlN:
		c.render.completion.next()
	case controlP:
		c.render.completion.previous()
	default:
		if function, ok := keyBindCmds[key]; ok {
			function(c.render.buffer)
			c.filter()
		}
	}

	return false, ""
}

func (c *chooser) Run() string {
	if l := os.Getenv("LOG_PATH"); l == "" {
		log.SetOutput(ioutil.Discard)
	} else if f, err := os.OpenFile(l, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666); err != nil {
		log.SetOutput(ioutil.Discard)
	} else {
		defer f.Close()
		log.SetOutput(f)
		log.Println("[INFO] Logging is enabled.")
	}

	c.init()
	defer c.terminal.restore()

	bufCh := make(chan []byte, 128)
	stopReadBufCh := make(chan struct{})
	go c.readBuffer(bufCh, stopReadBufCh)

	exitCh := make(chan int)
	winSizeCh := make(chan *winSize)
	stopHandleSignalCh := make(chan struct{})
	go c.handleSignals(exitCh, winSizeCh, stopHandleSignalCh)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, text := c.response(b); shouldExit {
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				c.render.clear()
				return text
			}
			c.render.render()

		case code := <-exitCh:
			os.Exit(code)

		case w := <-winSizeCh:
			c.render.winSize = w
			c.render.render()

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
