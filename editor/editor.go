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
	Cols            int
	Rows            int
	Buffer          []string
}

func New() *Editor {
	var orig unix.Termios
	if err := termios.Tcgetattr(0, &orig); err != nil {
		panic(err)
	}
	return &Editor{
		Buffer: []string{
			"hello world",
		},
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
	fmt.Print("\x1b[2J")
}

func (e *Editor) MoveCursor() {
	fmt.Printf("\x1b[%d;%dH", e.Cy+1, e.Cx+1)
}

func (e *Editor) DrawRows() error {
	e.HideCursor()
	defer e.ShowCursor()
	fmt.Print("\x1b[1;1H")
	if err := e.UpdateWindowSize(); err != nil {
		return err
	}
	for i := 0; i < e.Rows; i++ {
		if i < len(e.Buffer) {
			fmt.Print(e.Buffer[i])
		} else {
			fmt.Print("~")
		}
		if i < e.Rows-1 {
			fmt.Print("\r\n")
		}
	}
	e.MoveCursor()
	return nil
}

func (e *Editor) HandleKey(k Key, cancel func()) error {
	switch {
	case k.IsControl():
		switch k.Value {
		case ToControl('Q'):
			e.RefreshScreen()
			cancel()
		case ToControl('A'):
			e.MoveBeginning()
		case ToControl('E'):
			e.MoveEnd()
		}
	case k.IsEscaped():
		switch k.EscapedSequence[0] {
		case 'A':
			e.MoveAbove()
		case 'B':
			e.MoveBelow()
		case 'C':
			e.MoveRight()
		case 'D':
			e.MoveLeft()
		}
	default:
		fmt.Printf("%d (%c) ", k.Value, k.Value)
	}
	return nil
}

func (e *Editor) MoveAbove() {
	if e.Cy <= 0 {
		return
	}
	e.Cy--
	e.MoveCursor()
}

func (e *Editor) MoveBelow() {
	if e.Cy >= e.Rows-1 {
		return
	}
	e.Cy++
	e.MoveCursor()
}

func (e *Editor) MoveRight() {
	if e.Cx >= e.Cols-1 {
		return
	}
	e.Cx++
	e.MoveCursor()
}

func (e *Editor) MoveLeft() {
	if e.Cx <= 0 {
		return
	}
	e.Cx--
	e.MoveCursor()
}

func (e *Editor) MoveBeginning() {
	e.Cx = 0
	e.MoveCursor()
}

func (e *Editor) MoveEnd() {
	e.Cx = e.Cols - 1
	e.MoveCursor()
}

func (e *Editor) HideCursor() {
	fmt.Print("\x1b[?25l")
}

func (e *Editor) ShowCursor() {
	fmt.Print("\x1b[?25h")
}

func (e *Editor) ReadRune(ctx context.Context) chan rune {
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

func (e *Editor) ReadKey(ctx context.Context) chan Key {
	ks := make(chan Key)
	go func() {
		defer close(ks)
		rs := e.ReadRune(ctx)
		for r := range rs {
			switch {
			case r == '\x1b':
				a := <-rs
				b := <-rs
				if a == '[' {
					switch b {
					case 'A', 'B', 'C', 'D':
						ks <- Key{
							EscapedSequence: []rune{b},
						}
					}
				}
			default:
				ks <- Key{
					Value: r,
				}
			}
		}
	}()
	return ks
}

func (e *Editor) UpdateWindowSize() error {
	w, err := unix.IoctlGetWinsize(1, unix.TIOCGWINSZ)
	if err != nil {
		return err
	}
	e.Cols = int(w.Col)
	e.Rows = int(w.Row)
	return nil
}
