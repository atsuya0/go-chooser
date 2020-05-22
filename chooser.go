package chooser

import (
	"os"
	"strings"
	"sync"
	"time"
)

type input interface {
	read() ([]byte, error)
	getWinSize() (*winSize, error)
	setup() error
	restore() error
}

type chooser struct {
	terminal input
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
		render:   newRender(os.Stderr, len(list)),
		list:     list,
	}, nil
}

func (c *chooser) init() error {
	if err := c.terminal.setup(); err != nil {
		return err
	}
	if winSize, err := c.terminal.getWinSize(); err != nil {
		return err
	} else {
		c.render.winSize = winSize
	}
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
	var suggestions []string
	var indexes []int
	for i, v := range c.list {
		if c.contains(v) {
			suggestions = append(suggestions, v)
			indexes = append(indexes, i)
		}
	}
	c.render.completion = newCompletion(suggestions, indexes)
}

func (c *chooser) readBuffer(bufCh chan []byte, stopCh chan struct{}, wg *sync.WaitGroup) {
	for {
		select {
		case <-stopCh:
			wg.Done()
			return
		default:
			if b, err := c.terminal.read(); err == nil && !(len(b) == 1 && b[0] == 0) {
				bufCh <- b
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (c *chooser) response(b []byte, isMultiple bool) (bool, []int, []string) {
	switch key := getKey(b); key {
	case displayable:
		c.render.buffer.insert(string(b))
		c.filter()
	case enter:
		if len(c.render.register) == 0 {
			return true,
				[]int{c.render.completion.getIndex()},
				[]string{c.render.completion.getSuggestion()}
		}
		var strings []string
		for _, index := range c.render.register {
			strings = append(strings, c.list[index])
		}
		return true, c.render.register, strings
	case controlC:
		return true, make([]int, 0), make([]string, 0)
	case question:
		c.render.renderKeyBindings(isMultiple)
		return false, make([]int, 0), make([]string, 0)
	default:
		if keyBindingCmd, ok := keyBindingBufferCmds[key]; ok {
			keyBindingCmd.function(c.render.buffer)
			c.filter()
		} else if keyBindingCmd, ok := keyBindingRenderCmds[key]; ok {
			keyBindingCmd.function(c.render, isMultiple)
		}
	}
	c.render.renderSuggestions()

	return false, make([]int, 0), make([]string, 0)
}

func (c *chooser) Run() ([]int, []string, error) {
	if err := c.init(); err != nil {
		return make([]int, 0), make([]string, 0), err
	}
	defer func() {
		if err := c.terminal.restore(); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	bufCh := make(chan []byte, 128)
	wg.Add(1)
	go c.readBuffer(bufCh, stopCh, &wg)

	exitCh := make(chan int)
	winSizeCh := make(chan *winSize)
	errCh := make(chan error)
	wg.Add(1)
	go c.handleSignals(exitCh, winSizeCh, errCh, stopCh, &wg)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, indexes, strings := c.response(b, true); shouldExit {
				c.render.clearScreen()
				close(stopCh)
				wg.Wait()
				return indexes, strings, nil
			}

		case code := <-exitCh:
			os.Exit(code)

		case w := <-winSizeCh:
			c.render.winSize = w
			c.render.renderSuggestions()

		case err := <-errCh:
			return make([]int, 0), make([]string, 0), err
		}
	}
}

func (c *chooser) SingleRun() (int, string, error) {
	if err := c.init(); err != nil {
		return -1, "", err
	}
	defer func() {
		if err := c.terminal.restore(); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup
	stopCh := make(chan struct{})

	bufCh := make(chan []byte, 128)
	wg.Add(1)
	go c.readBuffer(bufCh, stopCh, &wg)

	exitCh := make(chan int)
	winSizeCh := make(chan *winSize)
	errCh := make(chan error)
	wg.Add(1)
	go c.handleSignals(exitCh, winSizeCh, errCh, stopCh, &wg)

	for {
		select {
		case b := <-bufCh:
			if shouldExit, indexes, strings := c.response(b, false); shouldExit {
				c.render.clearScreen()
				close(stopCh)
				wg.Wait()
				if len(indexes) > 0 {
					return indexes[0], strings[0], nil
				}
				return -1, "", nil
			}

		case code := <-exitCh:
			os.Exit(code)

		case w := <-winSizeCh:
			c.render.winSize = w
			c.render.renderSuggestions()

		case err := <-errCh:
			return -1, "", err
		}
	}
}
