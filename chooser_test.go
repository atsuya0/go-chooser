package chooser

import (
	"bufio"
	"bytes"
	"regexp"
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

func (t *testTerminal) getWinSize() *winSize {
	return &winSize{
		row: t.row,
		col: t.col,
	}
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
		render: newRender(output),
		list:   list,
	}
}

func compiledEscapeSequence() (*regexp.Regexp, error) {
	return regexp.Compile("\x1b\\[[0-9]*[A-Z]")
}

type ioBuf struct {
	i *bytes.Buffer
	o *bytes.Buffer
}

func (b *ioBuf) readLines() ([]string, error) {
	csi, err := compiledEscapeSequence()
	if err != nil {
		return make([]string, 0), err
	}
	time.Sleep(time.Millisecond * 100)
	var lines []string
	scanner := bufio.NewScanner(b.o)
	for scanner.Scan() {
		if line := csi.ReplaceAllString(scanner.Text(), ""); line != "" {
			lines = append(lines, line)
		}
	}
	b.o.Reset()
	if len(lines) == 0 {
		return make([]string, 0), nil
	}
	return lines[1:], nil
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

func setupChooser() (ioBuf, []string, *chooser) {
	io := ioBuf{
		i: new(bytes.Buffer),
		o: new(bytes.Buffer),
	}
	list := []string{"a1", "a2", "a3", "a4", "b1", "b2", "b3", "b4"}
	chooser := newTestChooser(io.i, io.o, list)
	return io, list, chooser
}

func TestChooserInputString(t *testing.T) {
	io, list, chooser := setupChooser()
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
	io, list, chooser := setupChooser()
	go func() {
		results := chooser.Run()
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
