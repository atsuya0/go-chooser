package choice

import (
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
	if err := c.terminal.setup(); err != nil {
		log.Fatalf("%+v\n", err)
	}
	c.render.winSize = c.terminal.getWinSize()
	c.filter()
	c.render.render()
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
	var suggestions []string
	var indexes []int
	for index, str := range c.list {
		if c.contains(str) {
			suggestions = append(suggestions, str)
			indexes = append(indexes, index)
		}
	}
	c.render.completion = newCompletion(suggestions, indexes)
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

func (c *chooser) response(b []byte) (bool, []string) {
	switch key := getKey(b); key {
	case displayable:
		c.render.buffer.insert(string(b))
		c.filter()
	case enter:
		if c.render.completion.target < 0 {
			return true, make([]string, 0)
		}
		if len(c.render.register) == 0 {
			return true, []string{c.render.completion.suggestions[c.render.completion.target]}
		}
		var texts []string
		for _, v := range c.render.register {
			texts = append(texts, c.list[v])
		}
		return true, texts
	case tab:
		c.render.updateRegister()
	case controlC:
		return true, make([]string, 0)
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

	return false, make([]string, 0)
}

func (c *chooser) Run() []string {
	c.init()
	defer func() {
		if err := c.terminal.restore(); err != nil {
			log.Fatalf("%+v\n", err)
		}
	}()

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
			if shouldExit, texts := c.response(b); shouldExit {
				stopReadBufCh <- struct{}{}
				stopHandleSignalCh <- struct{}{}
				clear()
				return texts
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
