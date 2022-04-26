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

func ReadKey(ctx context.Context) chan rune {
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
	t.Cc[unix.VMIN] = 0
	t.Cc[unix.VTIME] = 1
	termios.Tcsetattr(0, unix.TCIFLUSH, &t)
	return orig
}

func ResetRawMode(orig unix.Termios) {
	termios.Tcsetattr(0, unix.TCIFLUSH, &orig)
}

func HandleKey(k rune) error {
	fmt.Printf("%d (%c) ", k, k)
	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	defer ResetRawMode(SetRawMode())
	for k := range ReadKey(ctx) {
		switch {
		case unicode.IsControl(k):
			continue
		case k == 'q':
			cancel()
		default:
			HandleKey(k)
		}
	}
}
