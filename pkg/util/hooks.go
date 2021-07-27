package util

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func OnExit(timeout time.Duration, cb func(ctx context.Context)) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		cancel()
		os.Exit(1)
	}()
	cb(ctx)
}
