package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/mizuochikeita/re/editor"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	e := editor.New()
	if err := e.SetRawMode(); err != nil {
		panic(err)
	}
	defer e.ResetRawMode()

	e.RefreshScreen()
	if err := e.DrawRows(); err != nil {
		panic(err)
	}
	keys := e.ReadKey(ctx)
	for k := range keys {
		e.HandleKey(k, cancel)
	}
}
