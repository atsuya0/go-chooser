package selector

import (
	"fmt"
	"strings"
)

const (
	formatSelectedSuggestion = "\x1b[37m\x1b[1m%s\x1b[39m\x1b[0m"
)

type render struct {
	buffer     *buffer
	completion *completion
	winSize    *winSize
}

func (r *render) limit(suggestions []string) []string {
	if len(suggestions) <= int(r.winSize.row-2) {
		return suggestions
	} else {
		return suggestions[:r.winSize.row-2]
	}
}

func (r *render) clear() {
	fmt.Print("\x1b[1G\x1b[J")
}

func (r *render) renderBuffer() {
	fmt.Println(r.buffer.text)
}

func (r *render) renderSuggestions() {
	if r.completion.target < 0 {
		return
	}
	suggestions := make([]string, r.completion.length())
	copy(suggestions, r.completion.suggestions)
	suggestions[r.completion.target] =
		fmt.Sprintf(formatSelectedSuggestion, suggestions[r.completion.target])
	fmt.Print(strings.Join(suggestions, "\n"))
}

func (r *render) restoreCursorPosition() {
	fmt.Printf("\x1b[%dA\x1b[%dG", r.completion.length(), r.buffer.cursorPosition+1)
}

func (r *render) render() {
	r.clear()
	r.renderBuffer()
	r.renderSuggestions()
	r.restoreCursorPosition()
}

func newRender() *render {
	return &render{
		buffer: newBuffer(),
	}
}
