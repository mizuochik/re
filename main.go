package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unicode"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

func input(ctx context.Context) chan rune {
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
				panic(err)
			}
			c <- r
		}
	}()
	return c
}

func SetRawMode() unix.Termios {
	var orig unix.Termios
	if err := termios.Tcgetattr(0, &orig); err != nil {
		panic(err)
	}
	t := orig
	t.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	t.Oflag &^= syscall.OPOST
	t.Cflag |= syscall.CS8
	t.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	termios.Tcsetattr(0, unix.TCIFLUSH, &t)
	return orig
}

func ResetRawMode(orig unix.Termios) {
	termios.Tcsetattr(0, unix.TCIFLUSH, &orig)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	orig := SetRawMode()
	defer ResetRawMode(orig)
	for c := range input(ctx) {
		switch {
		case unicode.IsControl(c):
			continue
		case c == 'q':
			cancel()
		default:
			fmt.Printf("%d (%c) ", c, c)
		}
	}
}
