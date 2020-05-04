package choice

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type stringer string

func (s stringer) String() string {
	return string(s)
}

type chooser struct {
	terminal *terminal
	render   *render
	list     []fmt.Stringer
}

func NewChooser(list interface{}) (*chooser, error) {
	terminal, err := newTerminal()
	if err != nil {
		return &chooser{}, err
	}

	switch value := list.(type) {
	case []string:
		stringers := make([]fmt.Stringer, 0, len(value))
		for _, v := range value {
			stringers = append(stringers, stringer(v))
		}
		return &chooser{
			terminal: terminal,
			render:   newRender(),
			list:     stringers,
		}, nil
	case []fmt.Stringer:
		return &chooser{
			terminal: terminal,
			render:   newRender(),
			list:     value,
		}, nil
	}
	return &chooser{}, errors.New("expected []string or []fmt.Stringer")
}

func ToString(stringers []fmt.Stringer) []string {
	strings := make([]string, 0, len(stringers))
	for _, v := range stringers {
		strings = append(strings, string(v.String()))
	}
	return strings
}

func (c *chooser) init() error {
	if err := c.terminal.setup(); err != nil {
		return err
	}
	c.render.winSize = c.terminal.getWinSize()
	c.filter()
	c.render.renderSuggestions()
	return nil
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
	var suggestions []fmt.Stringer
	var indexes []int
	for i, v := range c.list {
		if c.contains(v.String()) {
			suggestions = append(suggestions, v)
			indexes = append(indexes, i)
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

func (c *chooser) response(b []byte) (bool, []fmt.Stringer) {
	switch key := getKey(b); key {
	case displayable:
		c.render.buffer.insert(string(b))
		c.filter()
	case enter:
		if len(c.render.register) == 0 {
			return true, []fmt.Stringer{c.render.completion.getSuggestion()}
		}
		var list []fmt.Stringer
		for _, v := range c.render.register {
			list = append(list, c.list[v])
		}
		return true, list
	case tab:
		c.render.updateRegister()
	case controlC:
		return true, make([]fmt.Stringer, 0)
	case question:
		c.render.renderKeyBindings()
		return false, make([]fmt.Stringer, 0)
	default:
		if keyBindingCmd, ok := keyBindingBufferCmds[key]; ok {
			keyBindingCmd.function(c.render.buffer)
			c.filter()
		} else if keyBindingCmd, ok := keyBindingRenderCmds[key]; ok {
			keyBindingCmd.function(c.render)
		}
	}
	c.render.renderSuggestions()

	return false, make([]fmt.Stringer, 0)
}

func (c *chooser) Run() []fmt.Stringer {
	if err := c.init(); err != nil {
		panic(err)
	}
	defer func() {
		if err := c.terminal.restore(); err != nil {
			panic(err)
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
				clearScreen()
				return texts
			}

		case code := <-exitCh:
			os.Exit(code)

		case w := <-winSizeCh:
			c.render.winSize = w
			c.render.renderSuggestions()

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
