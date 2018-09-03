package selector

type buffer struct {
	text           string
	cursorPosition int
}

func newBuffer() *buffer {
	return &buffer{text: "", cursorPosition: 0}
}

func (b *buffer) length() int {
	return len(b.isConvertRune())
}

func (b *buffer) isConvertRune() []rune {
	return []rune(b.text)
}

func (b *buffer) setBuffer(text string) {
	b.text = text
	b.cursorPosition = len(text)
}

func (b *buffer) insert(char string) {
	runeText := b.isConvertRune()
	b.text = string(runeText[:b.cursorPosition]) + char + string(runeText[b.cursorPosition:])
	b.forwardChar()
}

func (b *buffer) deleteChar(idx int) {
	runeText := b.isConvertRune()
	b.text = string(runeText[:idx]) + string(runeText[idx+1:])
}

func (b *buffer) deleteCharOnCursor() {
	if b.cursorPositionIsAtTheEnd() {
		return
	}
	b.deleteChar(b.cursorPosition)
}

func (b *buffer) backwardDeleteChar() {
	if b.cursorPositionIsAtTheBeginning() {
		return
	}
	b.deleteChar(b.cursorPosition - 1)
	b.backwardChar()
}

func (b *buffer) killLine() {
	runeText := b.isConvertRune()
	b.text = string(runeText[:b.cursorPosition])
}

func (b *buffer) backwardKillLine() {
	runeText := b.isConvertRune()
	b.text = string(runeText[b.cursorPosition:])
	b.cursorPosition = 0
}

func (b *buffer) backwardKillWord() {
	spacePosition := b.lastIndex(" ", b.cursorPosition-1)
	if spacePosition < 0 {
		b.backwardKillLine()
		return
	}
	runeText := b.isConvertRune()
	b.text = string(runeText[:spacePosition+1]) + string(runeText[b.cursorPosition:])
	b.cursorPosition = spacePosition + 1
}

func (b *buffer) beginningOfLine() {
	b.cursorPosition = 0
}

func (b *buffer) endOfLine() {
	b.cursorPosition = b.length()
}

func (b *buffer) forwardChar() {
	if b.cursorPositionIsAtTheEnd() {
		return
	}
	b.cursorPosition += 1
}

func (b *buffer) backwardChar() {
	if b.cursorPositionIsAtTheBeginning() {
		return
	}
	b.cursorPosition -= 1
}

func (b *buffer) indexOfTextFromCursor(substr string) int {
	if b.cursorPositionIsAtTheEnd() {
		return -1
	}
	return b.index(substr, b.cursorPosition)
}

func (b *buffer) index(substr string, beginningPoint int) int {
	if b.length() <= beginningPoint {
		return -1
	}
	runeText := b.isConvertRune()
	for i := beginningPoint; i < b.length(); i++ {
		if string(runeText[i]) == substr {
			return i
		}
	}
	return -1
}

func (b *buffer) lastIndexOfTextUpToCursor(substr string) int {
	if b.cursorPositionIsAtTheBeginning() {
		return -1
	}
	return b.lastIndex(substr, b.cursorPosition)
}

func (b *buffer) lastIndex(substr string, endPoint int) int {
	if endPoint <= 0 {
		return -1
	}
	runeText := b.isConvertRune()
	for i := endPoint - 1; 0 <= i; i-- {
		if string(runeText[i]) == substr {
			return i
		}
	}
	return -1
}

func (b *buffer) cursorPositionIsAtTheBeginning() bool {
	if b.cursorPosition <= 0 {
		return true
	}
	return false
}

func (b *buffer) cursorPositionIsAtTheEnd() bool {
	if b.length() <= b.cursorPosition {
		return true
	}
	return false
}

func (b *buffer) init() {
	b.cursorPosition = 0
	b.text = ""
}
