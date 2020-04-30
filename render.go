package choice

import (
	"fmt"
	"strings"
)

const (
	prompt         = ">>> "
	normalSymbol   = "  "
	cursorSymbol   = "> "
	selectedSymbol = "* "
	normalFormat   = "%s"
	cursorFormat   = "\x1b[1;37m%s\x1b[m"
	selectedFormat = "\x1b[1;34m%s\x1b[m"
	promptHeight   = 2
)

type render struct {
	buffer     *buffer
	completion *completion
	// The starting index of display suggestions.
	// If The number of lines suggestions is larger than the height of the screen.
	startingIndex int
	winSize       *winSize
	register      []int
}

func newRender() *render {
	return &render{
		buffer:        newBuffer(),
		startingIndex: 0,
		register:      make([]int, 0),
	}
}

func (r *render) render(lines []string) {
	clearScreen()
	r.renderBuffer()
	fmt.Print(strings.Join(lines, "\n"))
	r.restoreCursorPosition(len(lines))
}

func (r *render) renderSuggestions() {
	var suggestionsToDisplay []string
	for i := r.startingIndex; i < r.endingIndex() && i < r.completion.length(); i++ {
		suggestionsToDisplay = append(suggestionsToDisplay,
			r.shortenLine(r.assignSymbol(i)+fmt.Sprintf(r.assignFormat(i), r.completion.suggestions[i])))
	}
	r.render(suggestionsToDisplay)
}

func (r *render) renderKeyBindings() {
	var keyBindings []string
	for _, v := range keyBindingBufferCmds {
		keyBindings = append(keyBindings,
			r.shortenLine(fmt.Sprintf("%s: %s", v.key, v.description)))
	}
	for _, v := range keyBindingRenderCmds {
		keyBindings = append(keyBindings,
			r.shortenLine(fmt.Sprintf("%s: %s", v.key, v.description)))
	}
	r.render(keyBindings)
}

func (r *render) next() {
	r.completion.next()
	if r.endingIndex() <= r.completion.target {
		r.startingIndex += 1
	}
}

func (r *render) previous() {
	r.completion.previous()
	if r.completion.target < r.startingIndex {
		r.startingIndex -= 1
	}
}

// The ending index of display suggestions.
func (r *render) endingIndex() int {
	return r.startingIndex + int(r.winSize.row) - promptHeight
}

func (r *render) renderBuffer() {
	fmt.Println(prompt + r.buffer.text)
}

func (r *render) restoreCursorPosition(numOfSuggestions int) {
	fmt.Print(cursorUp(numOfSuggestions), setColCursor(r.cursorColPosition()))
}

func (r *render) shortenLine(line string) string {
	displayableWidth := int(r.winSize.col)
	runes := []rune(line)
	if len(runes) <= displayableWidth {
		return line
	}
	return string(runes[:displayableWidth:displayableWidth])
}

func (r *render) relativePositionOfTarget() int {
	return r.completion.target - r.startingIndex
}

func (r *render) cursorColPosition() int {
	return r.buffer.cursorPosition + len(prompt) + 1
}

// cursorSymbol is the highest priority.
func (r *render) assignSymbol(i int) string {
	if r.relativePositionOfTarget() == i {
		return cursorSymbol
	}

	for _, v := range r.register {
		if v == r.completion.indexes[i] {
			return selectedSymbol
		}
	}

	return normalSymbol
}

// selectedFormat is the highest priority.
func (r *render) assignFormat(i int) string {
	for _, v := range r.register {
		if v == r.completion.indexes[i] {
			return selectedFormat
		}
	}

	if r.relativePositionOfTarget() == i {
		return cursorFormat
	}

	return normalFormat
}

func (r *render) updateRegister() {
	index := r.completion.getIndex()
	if index < 0 {
		return
	}
	for i, v := range r.register {
		if v == index {
			r.register = append(r.register[:i:i], r.register[i+1:]...)
			return
		}
	}
	r.register = append(r.register, r.completion.indexes[r.completion.target])
}
