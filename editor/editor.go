package editor

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

type Editor struct {
	OriginalTermios unix.Termios
	Cx              int
	Cy              int
	Cols            int
	Rows            int
	Buffer          []string
	Vscroll         int
	Screen          *Screen
}

func New() *Editor {
	return &Editor{
		Screen: &Screen{},
	}
}

func (e *Editor) SetRawMode() error {
	if err := termios.Tcgetattr(0, &e.OriginalTermios); err != nil {
		return err
	}
	t := e.OriginalTermios
	t.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	t.Oflag &^= syscall.OPOST
	t.Cflag |= syscall.CS8
	t.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	t.Cc[unix.VMIN] = 0
	t.Cc[unix.VTIME] = 1
	termios.Tcsetattr(0, unix.TCIFLUSH, &t)
	return nil
}

func (e *Editor) ResetRawMode() {
	termios.Tcsetattr(0, unix.TCIFLUSH, &e.OriginalTermios)
}

func (e *Editor) ClearScreen() {
	e.HideCursor()
	defer e.ShowCursor()
	fmt.Print("\x1b[2J")
}

func (e *Editor) RefreshCursor() {
	x, y := e.Screen.CursorPosition()
	fmt.Printf("\x1b[%d;%dH", y+2, x+1) // for status bar
}

func (e *Editor) MoveCursorRelative(x, y int) {
	e.Screen.MoveCursorHorizontally(x)
	e.Screen.MoveCursorVertically(y)
	e.RefreshCursor()
}

func (e *Editor) MoveCursorAbsolute(x, y int) {
	// TBD
	e.RefreshCursor()
}

func (e *Editor) DrawStatusBar() {
	left := "re/main.go"
	right := "Saved"
	padding := e.Cols - len(left) - len(right)
	fmt.Print("\x1b[2K")
	fmt.Print("\x1b[37;40m") // white on black
	fmt.Print(left)
	fmt.Print(strings.Repeat(" ", padding))
	fmt.Print(right)
	fmt.Print("\x1b[0m") // reset color
	fmt.Print("\r\n")
}

func (e *Editor) RefreshScreen() error {
	e.HideCursor()
	defer e.ShowCursor()
	fmt.Print("\x1b[1;1H")
	if err := e.UpdateWindowSize(); err != nil {
		return err
	}
	e.DrawStatusBar()
	rows := e.Screen.View()
	for i := 0; i < e.Rows; i++ {
		fmt.Print("\x1b[2K")
		if i < len(rows) {
			fmt.Print(rows[i].Body)
		} else {
			fmt.Print("~")
		}
		if i < e.Rows-1 {
			fmt.Print("\r\n")
		}
	}
	e.RefreshCursor()
	return nil
}

func (e *Editor) OpenFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	r := bufio.NewReader(f)
	var buf []string
	var line []byte
outer:
	for {
		line = line[:0]
		for {
			bs, isPrefix, err := r.ReadLine()
			if errors.Is(err, io.EOF) {
				break outer
			}
			if err != nil {
				panic(err)
			}
			line = append(line, bs...)
			if !isPrefix {
				break
			}
		}
		buf = append(buf, string(line))
	}
	e.Buffer = buf
	e.UpdateWindowSize()
	e.Screen.Update(buf)
	return nil
}

func (e *Editor) HandleKey(k Key, cancel func()) error {
	switch {
	case k.IsControl():
		switch k.Value {
		case ToControl('P'):
			e.MoveAbove()
		case ToControl('N'):
			e.MoveBelow()
		case ToControl('F'):
			e.MoveRight()
		case ToControl('B'):
			e.MoveLeft()
		case ToControl('A'):
			e.MoveBeginning()
		case ToControl('E'):
			e.MoveEnd()
		case ToControl('U'):
			e.MoveCursorRelative(0, -e.Rows/2)
		case ToControl('D'):
			e.MoveCursorRelative(0, e.Rows/2)
		case ToControl('Q'):
			e.ClearScreen()
			cancel()
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
	e.MoveCursorRelative(0, -1)
}

func (e *Editor) MoveBelow() {
	e.MoveCursorRelative(0, 1)
}

func (e *Editor) MoveRight() {
	e.MoveCursorRelative(1, 0)
}

func (e *Editor) MoveLeft() {
	e.MoveCursorRelative(-1, 0)
}

func (e *Editor) MoveBeginning() {
	e.MoveCursorAbsolute(0, -1)
}

func (e *Editor) MoveEnd() {
	e.MoveCursorAbsolute(e.Cols-1, -1)
}

func (e *Editor) Scroll(rows int) {
	e.Vscroll += rows
	maxVscroll := len(e.Buffer) - e.Rows*3/4
	if e.Vscroll > maxVscroll {
		e.Vscroll = maxVscroll
	}
	minVscroll := 0
	if e.Vscroll < minVscroll {
		e.Vscroll = minVscroll
	}
	e.RefreshScreen()
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
	e.Rows = int(w.Row) - 1 // for status bar
	e.Screen.Width = int(w.Col)
	e.Screen.Height = int(w.Row) - 1 // for status bar
	return nil
}

func (e *Editor) Debugf(format string, a ...interface{}) {
	f, err := os.OpenFile("/tmp/re.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0700)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintf(f, format, a...)
	fmt.Fprintln(f)
}
