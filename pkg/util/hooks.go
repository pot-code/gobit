package util

import (
	"context"
	"os"
	"os/signal"
	"time"
)

type ExitHandler func(ctx context.Context)

type ExitManager struct {
	handlers []ExitHandler
}

func NewExitManager() *ExitManager {
	return &ExitManager{
		handlers: []ExitHandler{},
	}
}

func (em ExitManager) Register(handler ExitHandler) {
	em.handlers = append(em.handlers, handler)
}

func (em ExitManager) Wait(timeout time.Duration) {
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

	for _, h := range em.handlers {
		h(ctx)
	}
}
