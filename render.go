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
	prefix     string
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
	fmt.Println(r.prefix + r.buffer.text)
}

func (r *render) shortenSuggestions() []string {
	var suggestions []string
	for _, suggestion := range r.completion.suggestions {
		runeSuggestion := []rune(suggestion)
		if len(runeSuggestion) <= int(r.winSize.col) {
			suggestions = append(suggestions, suggestion)
		} else {
			suggestions =
				append(suggestions, string(runeSuggestion[:r.winSize.col]))
		}
	}
	return suggestions
}

func (r *render) renderSuggestions() {
	if r.completion.target < 0 {
		return
	}
	suggestions := make([]string, r.completion.length())
	copy(suggestions, r.shortenSuggestions())
	suggestions[r.completion.target] =
		fmt.Sprintf(selectedSuggestionFormat, suggestions[r.completion.target])
	fmt.Print(strings.Join(suggestions, "\n"))
}

func (r *render) cursorColPosition() int {
	return r.buffer.cursorPosition + len(r.prefix) + 1
}

func (r *render) restoreCursorPosition() {
	fmt.Printf("\x1b[%dA\x1b[%dG", r.completion.length(), r.cursorColPosition())
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
		prefix: prefix,
	}
}
