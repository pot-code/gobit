package util

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ExitHandler func(ctx context.Context)

type LivenessProbe func(ctx context.Context) error

type LifecycleManager struct {
	eh []ExitHandler
	lp []LivenessProbe
}

func NewLifecycleManager() *LifecycleManager {
	return &LifecycleManager{
		eh: []ExitHandler{},
	}
}

func (em *LifecycleManager) AddLivenessProbe(p LivenessProbe) {
	em.lp = append(em.lp, p)
}

// Probe check if all services are alive, return error if one is dead
func (em *LifecycleManager) Probe(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, p := range em.lp {
		if err := p(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (em *LifecycleManager) OnExit(handler ExitHandler) {
	em.eh = append(em.eh, handler)
}

func (em *LifecycleManager) Exit(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		cancel()
		os.Exit(1)
	}()

	for _, h := range em.eh {
		h(ctx)
	}
}

func (em *LifecycleManager) WaitExitSignal(timeout time.Duration) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	em.Exit(timeout)
}
