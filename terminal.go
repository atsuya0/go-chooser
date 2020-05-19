package chooser

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/pkg/term/termios"
)

const maxReadBytes = 1024

type terminal struct {
	org syscall.Termios
}

type winSize struct {
	row uint16
	col uint16
}

type ioctlWinsize struct {
	Row uint16
	Col uint16
	X   uint16
	Y   uint16
}

func (t *terminal) read() ([]byte, error) {
	buf := make([]byte, maxReadBytes)
	n, err := syscall.Read(syscall.Stdin, buf)
	if err != nil {
		return []byte{}, err
	}
	return buf[:n], nil
}

func (t *terminal) setRawMode() error {
	org := t.org
	org.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	org.Cc[syscall.VTIME] = 0
	org.Cc[syscall.VMIN] = 1

	return termios.Tcsetattr(uintptr(syscall.Stdin), termios.TCSANOW, &org)
}

func (t *terminal) resetMode() error {
	return termios.Tcsetattr(uintptr(syscall.Stdin), termios.TCSANOW, &t.org)
}

func (t *terminal) getWinSize() *winSize {
	ws := &ioctlWinsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if errno != 0 {
		panic(errno)
	}

	return &winSize{
		row: ws.Row,
		col: ws.Col,
	}
}

func (t *terminal) setup() error {
	if err := syscall.SetNonblock(syscall.Stdin, true); err != nil {
		return fmt.Errorf("Cannot set non blocking mode: %w", err)
	}
	if err := t.setRawMode(); err != nil {
		return fmt.Errorf("Cannot set raw mode: %w", err)
	}
	return nil
}

func (t *terminal) restore() error {
	if err := syscall.SetNonblock(syscall.Stdin, false); err != nil {
		return fmt.Errorf("Cannot set blocking mode: %w", err)
	}
	if err := t.resetMode(); err != nil {
		return fmt.Errorf("Cannot reset from raw mode: %w", err)
	}
	return nil
}

func newTerminal() (*terminal, error) {
	var org syscall.Termios
	if err := termios.Tcgetattr(uintptr(syscall.Stdin), &org); err != nil {
		return &terminal{}, err
	}

	return &terminal{org: org}, nil
}
