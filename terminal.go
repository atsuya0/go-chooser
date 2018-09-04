package choice

import (
	"log"
	"syscall"
	"unsafe"

	"github.com/pkg/term/termios"
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
	termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &org)

	return nil
}

func (t *terminal) resetMode() error {
	return termios.Tcsetattr(uintptr(t.fd), termios.TCSANOW, &t.originalTermios)
}

func (t *terminal) getWinSize() *winSize {
	ws := &ioctlWinsize{}
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(t.fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return &winSize{
		row: ws.Row,
		col: ws.Col,
	}
}

func (t *terminal) setup() {
	if err := syscall.SetNonblock(t.fd, true); err != nil {
		log.Fatalln("Cannot set non blocking mode.")
	}
	if err := t.setRawMode(); err != nil {
		log.Fatalln("Cannot set raw mode.")
	}
}

func (t *terminal) restore() error {
	if err := syscall.SetNonblock(t.fd, false); err != nil {
		log.Fatalln("Cannot set blocking mode.")
		return err
	}
	if err := t.resetMode(); err != nil {
		log.Fatalln("Cannot reset from raw mode.")
		return err
	}
	return nil
}

func newTerminal() *terminal {
	fd, err := syscall.Open("/dev/tty", syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatalln(err)
	}

	var org syscall.Termios
	if err := termios.Tcgetattr(uintptr(fd), &org); err != nil {
		log.Fatalln(err)
	}

	return &terminal{
		fd:              fd,
		originalTermios: org,
	}
}
