package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

func input(ctx context.Context) chan byte {
	c := make(chan byte)
	go func() {
		<-ctx.Done()
		close(c)
	}()
	go func() {
		keyBuf := make([]byte, 1)
		for {
			_, err := os.Stdin.Read(keyBuf)
			if err != nil {
				panic(err)
			}
			c <- keyBuf[0]
		}
	}()
	return c
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	var t unix.Termios
	if err := termios.Tcgetattr(0, &t); err != nil {
		panic(err)
	}
	defer func(orig unix.Termios) {
		termios.Tcsetattr(0, unix.TCSAFLUSH, &orig)
	}(t)
	t.Lflag &^= unix.ECHO
	t.Lflag &^= unix.ICANON
	if err := termios.Tcsetattr(0, unix.TCSAFLUSH, &t); err != nil {
		panic(err)
	}
	for c := range input(ctx) {
		switch c {
		case 'q':
			cancel()
		default:
			fmt.Printf("%d (%c) ", c, c)
		}
	}
}
