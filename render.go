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
	prompt     string
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
		prompt:        prompt,
		buffer:        newBuffer(),
		startingIndex: 0,
		register:      make([]int, 0),
	}
}

func (r *render) render() {
	clear()
	r.renderBuffer()
	numOfSuggestions := r.renderSuggestions()
	r.restoreCursorPosition(numOfSuggestions)
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
	fmt.Println(r.prompt + r.buffer.text)
}

func (r *render) renderSuggestions() int {
	if r.completion.target < 0 {
		return 0
	}
	var suggestionsToDisplay []string
	for i := r.startingIndex; i < r.endingIndex() && i < r.completion.length(); i++ {
		suggestionsToDisplay = append(suggestionsToDisplay,
			r.assignSymbol(i)+fmt.Sprintf(r.assignFormat(i), r.shortenSuggestion(r.completion.suggestions[i])))
	}
	fmt.Print(strings.Join(suggestionsToDisplay, "\n"))

	return len(suggestionsToDisplay)
}

func (r *render) restoreCursorPosition(numOfSuggestions int) {
	fmt.Print(cursorUp(numOfSuggestions), setColCursor(r.cursorColPosition()))
}

func (r *render) shortenSuggestion(suggestion string) string {
	displayableWidth := int(r.winSize.col) - len(normalSymbol)
	runeSuggestion := []rune(suggestion)
	if len(runeSuggestion) <= displayableWidth {
		return suggestion
	}
	return string(runeSuggestion[:displayableWidth:displayableWidth])
}

func (r *render) relativePositionOfTarget() int {
	return r.completion.target - r.startingIndex
}

func (r *render) cursorColPosition() int {
	return r.buffer.cursorPosition + len(r.prompt) + 1
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
