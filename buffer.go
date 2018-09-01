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

func (b *buffer) insert(char string) {
	runeText := b.isConvertRune()
	b.text = string(runeText[:b.cursorPosition]) + char + string(runeText[b.cursorPosition:])
	b.goRightChar()
}

func (b *buffer) del(idx int) {
	runeText := b.isConvertRune()
	b.text = string(runeText[:idx]) + string(runeText[idx+1:])
}

func (b *buffer) deleteChar() {
	if b.length() <= b.cursorPosition {
		return
	}
	b.del(b.cursorPosition)

	if b.length() <= b.cursorPosition {
		b.goLeftChar()
	}
}

func (b *buffer) deleteBeforeChar() {
	if b.cursorPosition-1 < 0 {
		return
	}
	b.del(b.cursorPosition - 1)
	b.goLeftChar()
}

func (b *buffer) goLineBeginning() {
	b.cursorPosition = 0
}

func (b *buffer) goLineEnd() {
	b.cursorPosition = b.length()
}

func (b *buffer) goRightChar() {
	if b.length() <= b.cursorPosition {
		return
	}
	b.cursorPosition += 1
}

func (b *buffer) goLeftChar() {
	if b.cursorPosition <= 0 {
		return
	}
	b.cursorPosition -= 1
}

func (b *buffer) backwardKillLine() {
	runeText := b.isConvertRune()
	b.text = string(runeText[b.cursorPosition:])
	b.cursorPosition = 0
}

func (b *buffer) init() {
	b.cursorPosition = 0
	b.text = ""
}
