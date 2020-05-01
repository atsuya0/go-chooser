package choice

import (
	"syscall"
	"unsafe"

	"github.com/pkg/term/termios"
	"golang.org/x/xerrors"
)

const maxReadBytes = 1024

type terminal struct {
	fd              int
	originalTermios syscall.Termios
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
	n, err := syscall.Read(t.fd, buf)
	if err != nil {
		return []byte{}, err
	}
	return buf[:n], nil
}

func (t *terminal) setRawMode() error {
	org := t.originalTermios
	org.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	org.Cc[syscall.VTIME] = 0
	org.Cc[syscall.VMIN] = 1

	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &org)
}

func (t *terminal) resetMode() error {
	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &t.originalTermios)
}

func (t *terminal) getWinSize() *winSize {
	ws := &ioctlWinsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(t.fd),
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
	if err := syscall.SetNonblock(t.fd, true); err != nil {
		return xerrors.Errorf("Cannot set non blocking mode: %w", err)
	}
	if err := t.setRawMode(); err != nil {
		return xerrors.Errorf("Cannot set raw mode: %w", err)
	}
	return nil
}

func (t *terminal) restore() error {
	if err := syscall.SetNonblock(t.fd, false); err != nil {
		return xerrors.Errorf("Cannot set blocking mode: %w", err)
	}
	if err := t.resetMode(); err != nil {
		return xerrors.Errorf("Cannot reset from raw mode: %w", err)
	}
	return nil
}

func newTerminal() (*terminal, error) {
	fd, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		return &terminal{}, err
	}

	var org syscall.Termios
	if err := termios.Tcgetattr(uintptr(fd), &org); err != nil {
		return &terminal{}, err
	}

	return &terminal{
		fd:              fd,
		originalTermios: org,
	}, nil
}
