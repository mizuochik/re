package editor

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"syscall"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

type Editor struct {
	OriginalTermios *unix.Termios
}

func (e *Editor) SetRawMode() {
	var orig unix.Termios
	if err := termios.Tcgetattr(0, &orig); err != nil {
		panic(err)
	}
	t := orig
	t.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	t.Oflag &^= syscall.OPOST
	t.Cflag |= syscall.CS8
	t.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	t.Cc[unix.VMIN] = 0
	t.Cc[unix.VTIME] = 1
	termios.Tcsetattr(0, unix.TCIFLUSH, &t)
	e.OriginalTermios = &orig
}

func (e *Editor) ResetRawMode() {
	termios.Tcsetattr(0, unix.TCIFLUSH, e.OriginalTermios)
}

func (e *Editor) RefreshScreen() {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")
}

func (e *Editor) DrawRows() error {
	_, col, err := e.WindowSize()
	if err != nil {
		return err
	}
	for i := 0; i < col; i++ {
		fmt.Print("~\r\n")
	}
	return nil
}

func (e *Editor) HandleKey(k rune) error {
	fmt.Printf("%d (%c) ", k, k)
	return nil
}

func (e *Editor) ReadKey(ctx context.Context) chan rune {
	c := make(chan rune)
	go func() {
		<-ctx.Done()
		close(c)
	}()
	go func() {
		rd := bufio.NewReader(os.Stdin)
		for {
			r, _, err := rd.ReadRune()
			if err != nil {
				if err.Error() == "EOF" {
					continue // timeout
				}
				panic(err)
			}
			c <- r
		}
	}()
	return c
}

func (e *Editor) WindowSize() (int, int, error) {
	w, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ)
	if err != nil {
		return 0, 0, err
	}
	return int(w.Row), int(w.Col), nil
}
