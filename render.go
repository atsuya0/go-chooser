package chooser

import (
	"fmt"
	"io"
	"strings"
)

const (
	promptHeight = 3
)

type prompt struct {
	text   string
	length int
}

func newPrompt(text string) prompt {
	return prompt{
		text:   fmt.Sprintf("\x1b[1;34m%s\x1b[m ", text),
		length: len(text) + 1,
	}
}

type symbol struct {
	normal   string
	cursor   string
	selected string
	length   int
}

func newSymbol() symbol {
	return symbol{normal: "  ", cursor: "> ", selected: "* ", length: 2}
}

type format struct {
	normal   string
	cursor   string
	selected string
}

func newFormat() format {
	return format{normal: "%s", cursor: "\x1b[1;37m%s\x1b[m", selected: "\x1b[1;34m%s\x1b[m"}
}

type render struct {
	buffer     *buffer
	completion *completion
	// The starting index of display suggestions.
	// If The number of lines suggestions is larger than the height of the screen.
	startingIndex   int
	winSize         *winSize
	heldCompletions []int
	out             io.Writer
	headerFormat    string
	symbol          symbol
	format          format
	prompt          prompt
}

func newRender(out io.Writer, len int) *render {
	return &render{
		buffer:          newBuffer(),
		startingIndex:   0,
		heldCompletions: make([]int, 0),
		out:             out,
		headerFormat:    fmt.Sprintf("\x1b[32m%%d/%d (%%d)\x1b[m", len),
		symbol:          newSymbol(),
		format:          newFormat(),
		prompt:          newPrompt(">>>"),
	}
}

func (r *render) clearScreen() {
	r.print(cursorForward(1) + clearCursorEnd())
}

func (r *render) restoreCursorPosition(row, col int) {
	r.print(cursorUp(row) + cursorForward(col))
}

func (r *render) print(i ...interface{}) {
	fmt.Fprint(r.out, i...)
}

func (r *render) render(lines []string) {
	r.clearScreen()
	lines = append([]string{r.inputField(), r.header()}, lines...)
	r.print(strings.Join(lines, "\n"))
	r.restoreCursorPosition(len(lines)-1, r.cursorColPosition())
}

func (r *render) renderSuggestions() {
	var suggestionsToDisplay []string
	for i := r.startingIndex; i < r.endingIndex() && i < r.completion.length(); i++ {
		suggestionsToDisplay = append(suggestionsToDisplay,
			r.assignSymbol(i)+fmt.Sprintf(r.assignFormat(i), r.shortenLine(r.completion.suggestions[i])))
	}
	r.render(suggestionsToDisplay)
}

func (r *render) renderKeyBindings(isMultiple bool) {
	var keyBindings []string
	for _, v := range bufferKeyBindingCmds {
		keyBindings = append(keyBindings,
			r.shortenLine(fmt.Sprintf("%s: %s", v.key, v.description)))
	}
	for _, v := range renderKeyBindingCmds {
		if isMultiple && v.key == "tab" {
			continue
		}
		keyBindings = append(keyBindings,
			r.shortenLine(fmt.Sprintf("%s: %s", v.key, v.description)))
	}
	for _, v := range keyBindingCmds {
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

func (r *render) inputField() string {
	return r.prompt.text + r.buffer.text
}

func (r *render) header() string {
	if len(r.heldCompletions) > 0 {
		return fmt.Sprintf(r.headerFormat, r.completion.length(), len(r.heldCompletions))
	}
	return fmt.Sprintf(r.headerFormat, r.completion.length(), 0)
}

func (r *render) shortenLine(line string) string {
	displayableWidth := int(r.winSize.col) - r.symbol.length
	runes := []rune(line)
	if len(runes) <= displayableWidth {
		return line
	}
	return string(runes[:displayableWidth:displayableWidth])
}

func (r *render) cursorColPosition() int {
	return r.buffer.cursorPosition + r.prompt.length + 1
}

// cursorSymbol is the highest priority.
func (r *render) assignSymbol(i int) string {
	if r.completion.target == i {
		return r.symbol.cursor
	}

	for _, v := range r.heldCompletions {
		if v == r.completion.indexes[i] {
			return r.symbol.selected
		}
	}

	return r.symbol.normal
}

// selectedFormat is the highest priority.
func (r *render) assignFormat(i int) string {
	for _, v := range r.heldCompletions {
		if v == r.completion.indexes[i] {
			return r.format.selected
		}
	}

	if r.completion.target == i {
		return r.format.cursor
	}

	return r.format.normal
}

func (r *render) holdCompletion() {
	index := r.completion.getIndex()
	if index < 0 {
		return
	}
	for i, v := range r.heldCompletions {
		if v == index {
			r.heldCompletions = append(r.heldCompletions[:i:i], r.heldCompletions[i+1:]...)
			return
		}
	}
	r.heldCompletions = append(r.heldCompletions, r.completion.indexes[r.completion.target])
	r.next()
}
