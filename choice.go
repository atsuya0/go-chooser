package choice

import (
	"os"
	"strings"
	"time"
)

type chooser struct {
	terminal *terminal
	render   *render
	list     []string
}

func NewChooser(list []string) (*chooser, error) {
	terminal, err := newTerminal()
	if err != nil {
		return &chooser{}, err
	}

	return &chooser{
		terminal: terminal,
		render:   newRender(),
		list:     list,
	}, nil
}

func (c *chooser) init() {
	c.terminal.setup()
	c.render.winSize = c.terminal.getWinSize()
	c.filter()
	c.render.render()
}

func (c *chooser) SetPrefix(prefix string) {
	c.render.prefix = prefix + " "
}

// It contains all whitespace-separated character strings.
func (c *chooser) contains(str string) bool {
	for _, substr := range strings.Fields(c.render.buffer.text) {
		if !(strings.Contains(str, substr)) {
			return false
		}
	}
	return true
}

// Filter the complement target.
func (c *chooser) filter() {
	var result []string
	for _, str := range c.list {
		if c.contains(str) {
			result = append(result, str)
		}
	}
	c.render.completion = newCompletion(result)
}

func (c *chooser) readBuffer(bufCh chan []byte, stopCh chan struct{}) {
	for {
		select {
		case <-stopCh:
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
		return false, ""
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
		c.render.next()
	case controlP:
		c.render.previous()
	default:
		if function, ok := keyBindCmds[key]; ok {
			function(c.render.buffer)
			c.filter()
		}
	}

	return false, ""
}

func (c *chooser) Run() string {
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
				clear()
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
