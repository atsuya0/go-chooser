package chooser

import "testing"

const (
	errFormatForText   = "[Text] output: %s, request: %s"
	errFormatForCursor = "[Cursor position] output: %d, request: %d"
	errFormatForIndex  = "[Index] output: %d, request: %d"
)

func TestBufferSetBuffer(t *testing.T) {
	patterns := []struct {
		text           string
		cursorPosition int
	}{
		{text: "01234", cursorPosition: len("01234")},
		{text: "0", cursorPosition: len("0")},
	}

	for _, pattern := range patterns {
		buffer := &buffer{}
		buffer.setBuffer(pattern.text)

		if buffer.text != pattern.text {
			t.Errorf(errFormatForText, buffer.text, pattern.text)
		}
		if buffer.cursorPosition != pattern.cursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.cursorPosition)
		}
	}
}

func TestBufferInsert(t *testing.T) {
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestText           string
		requestCursorPosition int
	}{
		{text: "01234", cursorPosition: 0,
			requestText: "^01234", requestCursorPosition: 1},
		{text: "01234", cursorPosition: len("01234"),
			requestText: "01234^", requestCursorPosition: len("01234") + 1},
		{text: "01234", cursorPosition: 2,
			requestText: "01^234", requestCursorPosition: 3},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.insert("^")
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferDeleteChar(t *testing.T) {
	patterns := []struct {
		text        string
		requestText string
		index       int
	}{
		{text: "01234", requestText: "1234", index: 0},
		{text: "01234", requestText: "0123", index: len("01234") - 1},
		{text: "01234", requestText: "0134", index: 2},
	}

	for _, pattern := range patterns {
		buffer := &buffer{text: pattern.text}
		buffer.deleteChar(pattern.index)
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
	}
}

func TestBufferDeleteCharOnCursor(t *testing.T) {
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestText           string
		requestCursorPosition int
	}{
		{text: "01234", cursorPosition: 0,
			requestText: "1234", requestCursorPosition: 0},
		{text: "01234", cursorPosition: len("01234"),
			requestText: "01234", requestCursorPosition: len("01234")},
		{text: "01234", cursorPosition: 2,
			requestText: "0134", requestCursorPosition: 2},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.deleteCharOnCursor()
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferBackwardDeleteChar(t *testing.T) {
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestText           string
		requestCursorPosition int
	}{
		{text: "01234", cursorPosition: 0,
			requestText: "01234", requestCursorPosition: 0},
		{text: "01234", cursorPosition: len("01234"),
			requestText: "0123", requestCursorPosition: len("0123")},
		{text: "01234", cursorPosition: 2,
			requestText: "0234", requestCursorPosition: 1},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.backwardDeleteChar()
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferKillLine(t *testing.T) {
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestText           string
		requestCursorPosition int
	}{
		{text: "01234", cursorPosition: 0,
			requestText: "", requestCursorPosition: 0},
		{text: "01234", cursorPosition: len("01234"),
			requestText: "01234", requestCursorPosition: len("01234")},
		{text: "01234", cursorPosition: 2,
			requestText: "01", requestCursorPosition: 2},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.killLine()
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferBackwardKillLine(t *testing.T) {
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestText           string
		requestCursorPosition int
	}{
		{text: "01234", cursorPosition: 0,
			requestText: "01234", requestCursorPosition: 0},
		{text: "01234", cursorPosition: len("01234"),
			requestText: "", requestCursorPosition: 0},
		{text: "01234", cursorPosition: 2,
			requestText: "234", requestCursorPosition: 0},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.backwardKillLine()
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferBackwardKillWord(t *testing.T) {
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestText           string
		requestCursorPosition int
	}{
		{text: "012 34", cursorPosition: 0,
			requestText: "012 34", requestCursorPosition: 0},
		{text: "012 34", cursorPosition: len("012 34"),
			requestText: "012 ", requestCursorPosition: len("012 ")},
		{text: "012 34", cursorPosition: 2,
			requestText: "2 34", requestCursorPosition: 0},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.backwardKillWord()
		if buffer.text != pattern.requestText {
			t.Errorf(errFormatForText, buffer.text, pattern.requestText)
		}
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferForwardChar(t *testing.T) {
	const text = "01234"
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestCursorPosition int
	}{
		{text: text, cursorPosition: 0, requestCursorPosition: 1},
		{text: text, cursorPosition: len(text), requestCursorPosition: len(text)},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.forwardChar()
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferBackwardChar(t *testing.T) {
	const text = "01234"
	patterns := []struct {
		text                  string
		cursorPosition        int
		requestCursorPosition int
	}{
		{text: text, cursorPosition: 0, requestCursorPosition: 0},
		{text: text, cursorPosition: len(text),
			requestCursorPosition: len(text) - 1},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		buffer.backwardChar()
		if buffer.cursorPosition != pattern.requestCursorPosition {
			t.Errorf(errFormatForCursor,
				buffer.cursorPosition, pattern.requestCursorPosition)
		}
	}
}

func TestBufferIndex(t *testing.T) {
	const text = "0-2-4"
	patterns := []struct {
		text           string
		beginningPoint int
		requestIndex   int
	}{
		{text: text, beginningPoint: 0, requestIndex: 1},
		{text: text, beginningPoint: 2, requestIndex: 3},
	}

	for _, pattern := range patterns {
		buffer := &buffer{text: pattern.text}
		index := buffer.index("-", pattern.beginningPoint)
		if index != pattern.requestIndex {
			t.Errorf(errFormatForIndex, index, pattern.requestIndex)
		}
	}
}

func TestBufferIndexOfTextFromCursor(t *testing.T) {
	const text = "0-2-4"
	patterns := []struct {
		text           string
		cursorPosition int
		requestIndex   int
	}{
		{text: text, cursorPosition: 0, requestIndex: 1},
		{text: text, cursorPosition: 2, requestIndex: 3},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		index := buffer.indexOfTextFromCursor("-")
		if index != pattern.requestIndex {
			t.Errorf(errFormatForIndex, index, pattern.requestIndex)
		}
	}
}

func TestBufferLastIndex(t *testing.T) {
	const text = "0-2-4"
	patterns := []struct {
		text         string
		endPoint     int
		requestIndex int
	}{
		{text: text, endPoint: len(text), requestIndex: 3},
		{text: text, endPoint: 2, requestIndex: 1},
	}

	for _, pattern := range patterns {
		buffer := &buffer{text: pattern.text}
		index := buffer.lastIndex("-", pattern.endPoint)
		if index != pattern.requestIndex {
			t.Errorf(errFormatForIndex, index, pattern.requestIndex)
		}
	}
}

func TestBufferLastIndexOfTextFromCursor(t *testing.T) {
	const text = "0-2-4"
	patterns := []struct {
		text           string
		cursorPosition int
		requestIndex   int
	}{
		{text: text, cursorPosition: len(text), requestIndex: 3},
		{text: text, cursorPosition: 2, requestIndex: 1},
	}

	for _, pattern := range patterns {
		buffer :=
			&buffer{text: pattern.text, cursorPosition: pattern.cursorPosition}
		index := buffer.lastIndexOfTextUpToCursor("-")
		if index != pattern.requestIndex {
			t.Errorf(errFormatForIndex, index, pattern.requestIndex)
		}
	}
}
