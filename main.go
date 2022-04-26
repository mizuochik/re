package main

import (
	"context"
	"os/signal"
	"syscall"
	"unicode"

	"github.com/mizuochikeita/re/editor"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	e := &editor.Editor{}

	e.SetRawMode()
	defer e.ResetRawMode()

	e.RefreshScreen()
	if err := e.DrawRows(); err != nil {
		panic(err)
	}

	for k := range e.ReadKey(ctx) {
		switch {
		case unicode.IsControl(k):
			continue
		case k == 'q':
			e.RefreshScreen()
			cancel()
		default:
			e.HandleKey(k)
		}
	}
}
