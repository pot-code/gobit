package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ExitHandler func(ctx context.Context)

type OnlineProbe func(ctx context.Context) error

type ExitManager struct {
	handlers []ExitHandler
	probes   []OnlineProbe
}

func NewExitManager() *ExitManager {
	return &ExitManager{
		handlers: []ExitHandler{},
	}
}

func (em *ExitManager) AddLivenessProbe(p OnlineProbe) {
	em.probes = append(em.probes, p)
}

// Probe check if all services are alive, return error if one is dead
func (em *ExitManager) Probe(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, p := range em.probes {
		if err := p(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (em *ExitManager) OnExit(handler ExitHandler) {
	em.handlers = append(em.handlers, handler)
}

func (em *ExitManager) Exit(timeout time.Duration) {
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

func (em *ExitManager) WaitExitSignal(timeout time.Duration) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch
}
