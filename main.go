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

	if err := e.OpenFile("README.md"); err != nil {
		panic(err)
	}
	if err := e.RefreshScreen(); err != nil {
		panic(err)
	}
	for k := range e.ReadKey(ctx) {
		e.HandleKey(k, cancel)
	}
}
