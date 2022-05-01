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
	Cx              int
	Cy              int
}

func New() *Editor {
	var orig unix.Termios
	if err := termios.Tcgetattr(0, &orig); err != nil {
		panic(err)
	}
	return &Editor{
		OriginalTermios: &orig,
	}
}

func (e *Editor) SetRawMode() {
	t := *e.OriginalTermios
	t.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	t.Oflag &^= syscall.OPOST
	t.Cflag |= syscall.CS8
	t.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	t.Cc[unix.VMIN] = 0
	t.Cc[unix.VTIME] = 1
	termios.Tcsetattr(0, unix.TCIFLUSH, &t)
}

func (e *Editor) ResetRawMode() {
	termios.Tcsetattr(0, unix.TCIFLUSH, e.OriginalTermios)
}

func (e *Editor) RefreshScreen() {
	e.HideCursor()
	defer e.ShowCursor()
	e.MoveCursor()
	fmt.Print("\x1b[2J")
}

func (e *Editor) MoveCursor() {
	fmt.Printf("\x1b[%d;%dH", e.Cy+1, e.Cx+1)
}

func (e *Editor) DrawRows() error {
	e.HideCursor()
	defer e.ShowCursor()

	_, col, err := e.WindowSize()
	if err != nil {
		return err
	}
	for i := 0; i < col; i++ {
		fmt.Print("~")
		if i < col-1 {
			fmt.Print("\r\n")
		}
	}
	e.MoveCursor()
	return nil
}

func (e *Editor) HandleKey(k rune) error {
	fmt.Printf("%d (%c) ", k, k)
	return nil
}

func (e *Editor) HideCursor() {
	fmt.Print("\x1b[?25l")
}

func (e *Editor) ShowCursor() {
	fmt.Print("\x1b[?25h")
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
