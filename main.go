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

	e := editor.New()
	e.SetRawMode()
	defer e.ResetRawMode()

	e.RefreshScreen()
	if err := e.DrawRows(); err != nil {
		panic(err)
	}

	keys := e.ReadKey(ctx)
	for k := range keys {
		switch {
		case k == '\x1b':
			a := <-keys
			b := <-keys
			if a == '[' {
				switch b {
				case 'A':
					e.Cy--
				case 'B':
					e.Cy++
				case 'C':
					e.Cx++
				case 'D':
					e.Cx--
				}
				e.MoveCursor()
			}
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
