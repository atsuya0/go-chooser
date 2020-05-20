package chooser

import (
	"bufio"
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
)

type testTerminal struct {
	buf *bytes.Buffer
	row uint16
	col uint16
}

func (t *testTerminal) read() ([]byte, error) {
	buf := make([]byte, 1024)
	n, err := t.buf.Read(buf)
	if err != nil {
		return []byte{}, err
	}
	t.buf.Reset()
	return buf[:n], nil
}

func (t *testTerminal) getWinSize() (*winSize, error) {
	return &winSize{
		row: t.row,
		col: t.col,
	}, nil
}

func (t *testTerminal) setup() error   { return nil }
func (t *testTerminal) restore() error { return nil }

func newTestChooser(input *bytes.Buffer, output *bytes.Buffer, list []string, size ...uint16) *chooser {
	var row, col uint16
	if len(size) > 1 {
		row, col = size[0], size[1]
	} else {
		row, col = 100, 100
	}
	return &chooser{
		terminal: &testTerminal{
			buf: input,
			row: row,
			col: col,
		},
		render: newRender(output, len(list)),
		list:   list,
	}
}

func compiledScreenEscapeSequence() *regexp.Regexp {
	return regexp.MustCompile("\x1b\\[[0-9]*[A-Z]")
}

func compiledFormatEscapeSequence() *regexp.Regexp {
	return regexp.MustCompile("\x1b\\[[0-9]*;?[0-9]*m")
}

type ioBuf struct {
	i *bytes.Buffer
	o *bytes.Buffer
}

func (b *ioBuf) readLines() ([]string, error) {
	time.Sleep(time.Millisecond * 100)
	var lines []string
	scanner := bufio.NewScanner(b.o)
	for scanner.Scan() {
		line := compiledFormatEscapeSequence().ReplaceAllString(scanner.Text(), "")
		if line == "" {
			continue
		}
		line = compiledScreenEscapeSequence().ReplaceAllString(line, "")
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	b.o.Reset()
	if len(lines) <= promptHeight-1 {
		return make([]string, 0), nil
	}
	return lines[promptHeight-1:], nil
}

func (b *ioBuf) write(bytes []byte) ([]string, error) {
	if _, err := b.i.Write(bytes); err != nil {
		return make([]string, 0), err
	}
	return b.readLines()
}

func (b *ioBuf) writeString(s string) ([]string, error) {
	if _, err := b.i.WriteString(s); err != nil {
		return make([]string, 0), err
	}
	return b.readLines()
}

func setupTestChooser(size ...uint16) (ioBuf, []string, *chooser) {
	io := ioBuf{
		i: new(bytes.Buffer),
		o: new(bytes.Buffer),
	}
	list := []string{"a1", "a2", "a3", "a4", "b1", "b2", "b3", "b4"}
	if len(size) > 1 {
		return io, list, newTestChooser(io.i, io.o, list, size...)
	}
	return io, list, newTestChooser(io.i, io.o, list)
}

func TestChooserInputString(t *testing.T) {
	io, list, chooser := setupTestChooser()
	go chooser.Run()

	if lines, err := io.readLines(); err != nil {
		t.Error(err)
	} else if len(lines) != len(list) {
		t.Errorf("result %d, expected %d", len(lines), len(list))
	}

	lines, err := io.writeString("a")
	if err != nil {
		t.Error(err)
	}
	if len(lines) != 4 {
		t.Errorf("result %d, expected %d", len(lines), 4)
	}

	lines, err = io.writeString("1")
	if err != nil {
		t.Error(err)
	}
	if len(lines) != 1 {
		t.Errorf("result %d, expected %d", len(lines), 1)
	}
}

func TestChooserInputBytes(t *testing.T) {
	io, list, chooser := setupTestChooser()
	go func() {
		results, _ := chooser.Run()
		if results[0] != list[1] {
			t.Errorf("result %s, expected %s", results[0], list[1])
		}
	}()

	if lines, err := io.readLines(); err != nil {
		t.Error(err)
	} else if len(lines) != len(list) {
		t.Errorf("result %d, expected %d", len(lines), len(list))
	}

	// C-n
	if _, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	}
	// Enter
	if _, err := io.write([]byte{0xa}); err != nil {
		t.Error(err)
	}
}

func TestChooserMultipleSelection(t *testing.T) {
	io, list, chooser := setupTestChooser()
	go func() {
		expectedValues := []string{list[0], list[2]}
		results, _ := chooser.Run()
		for i, v := range results {
			if v != expectedValues[i] {
				t.Errorf("result %s, expected %s", v, expectedValues[i])
			}
		}
	}()

	if lines, err := io.readLines(); err != nil {
		t.Error(err)
	} else if len(lines) != len(list) {
		t.Errorf("result %d, expected %d", len(lines), len(list))
	}

	// tab
	if _, err := io.write([]byte{0x9}); err != nil {
		t.Error(err)
	}
	// C-n
	if _, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	}
	// C-n
	if _, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	}
	// tab
	if _, err := io.write([]byte{0x9}); err != nil {
		t.Error(err)
	}
	// C-n
	if lines, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	} else {
		if bool, err := regexp.MatchString(`^*`, lines[0]); err != nil {
			t.Log(err)
			t.Fail()
		} else if !bool {
			t.Errorf("'%s' does not start with *.", lines[0])
		}
		if bool, err := regexp.MatchString(`^*`, lines[2]); err != nil {
			t.Log(err)
			t.Fail()
		} else if !bool {
			t.Errorf("'%s' does not start with *.", lines[2])
		}
		if bool, err := regexp.MatchString(`^>`, lines[3]); err != nil {
			t.Log(err)
			t.Fail()
		} else if !bool {
			t.Errorf("'%s' does not start with >.", lines[3])
		}
	}
	// Enter
	if _, err := io.write([]byte{0xa}); err != nil {
		t.Error(err)
	}
}

func TestChooserScroll(t *testing.T) {
	numOfSuggestionsToDisplay := 2
	io, list, chooser := setupTestChooser(promptHeight+uint16(numOfSuggestionsToDisplay), 100)
	go chooser.Run()

	if lines, err := io.readLines(); err != nil {
		t.Error(err)
	} else if len(lines) != numOfSuggestionsToDisplay {
		t.Errorf("result %d, expected %d", len(lines), numOfSuggestionsToDisplay)
	}
	// C-n
	if _, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	}
	// C-n
	if _, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	}
	// C-n
	if lines, err := io.write([]byte{0xe}); err != nil {
		t.Error(err)
	} else {
		for i, line := range lines {
			if !strings.Contains(line, list[i+numOfSuggestionsToDisplay]) {
				t.Errorf("The %s is not included in the %s.", list[i+numOfSuggestionsToDisplay], line)
			}
		}
	}
}
