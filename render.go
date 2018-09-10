package choice

import (
	"fmt"
	"strings"
)

const (
	prefix                   = ">>> "
	selectedSuggestionFormat = "\x1b[1;37m%s\x1b[m"
)

type render struct {
	prefix        string
	buffer        *buffer
	completion    *completion
	startingPoint int // Starting point of display.
	winSize       *winSize
}

func newRender() *render {
	return &render{
		buffer: newBuffer(),
		prefix: prefix,
	}
}

func (r *render) render() {
	r.clear()
	r.renderBuffer()
	numOfSuggestions := r.renderSuggestions()
	r.restoreCursorPosition(numOfSuggestions)
}

func (r *render) clear() {
	fmt.Print("\x1b[1G\x1b[J")
}

func (r *render) next() {
	r.completion.next()
	if r.endPoint() <= r.completion.target {
		r.startingPoint += 1
	}
}

func (r *render) previous() {
	r.completion.previous()
	if r.completion.target < r.startingPoint {
		r.startingPoint -= 1
	}
}

func (r *render) endPoint() int {
	return r.startingPoint + int(r.winSize.row) - 2
}

func (r *render) renderBuffer() {
	fmt.Println(r.prefix + r.buffer.text)
}

func (r *render) renderSuggestions() int {
	if r.completion.target < 0 {
		return 0
	}
	var suggestions []string
	for i := r.startingPoint; i < r.endPoint() && i < r.completion.length(); i++ {
		suggestions = append(
			suggestions, r.shortenSuggestion(r.completion.suggestions[i]))
	}
	suggestions[r.relativePositionOfTarget()] =
		fmt.Sprintf(selectedSuggestionFormat, suggestions[r.relativePositionOfTarget()])
	fmt.Print(strings.Join(suggestions, "\n"))

	return len(suggestions)
}

func (r *render) restoreCursorPosition(numOfSuggestions int) {
	fmt.Printf("\x1b[%dA\x1b[%dG", numOfSuggestions, r.cursorColPosition())
}

func (r *render) shortenSuggestion(suggestion string) string {
	runeSuggestion := []rune(suggestion)
	if len(runeSuggestion) <= int(r.winSize.col) {
		return suggestion
	}
	return string(runeSuggestion[:r.winSize.col])
}

func (r *render) relativePositionOfTarget() int {
	return r.completion.target - r.startingPoint
}

func (r *render) cursorColPosition() int {
	return r.buffer.cursorPosition + len(r.prefix) + 1
}
