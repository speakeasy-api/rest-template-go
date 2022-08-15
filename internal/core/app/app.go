package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"github.com/speakeasy-api/rest-template-go/internal/core/tracing"
	"go.uber.org/zap"
)

// Listener represents a type that can listen for incoming connections.
type Listener interface {
	Listen(context.Context) error
}

// OnShutdownFunc is a function that is called when the app is shutdown.
type OnShutdownFunc func()

// App represents the application run by this service.
type App struct {
	Name          string
	shutdownFuncs []OnShutdownFunc
}

// OnStart is a function that is called when the app is started.
type OnStart func(context.Context, *App) ([]Listener, error)

// Start starts the application.
func Start(onStart OnStart) {
	ctx := context.Background()

	a := &App{
		Name: "test-app", // TODO determine how to configure this
	}
	a.OnShutdown(func() {
		if err := logging.Sync(ctx); err != nil {
			log.Printf("failed to sync logging: %v", err)
		}
	})

	logging.From(ctx).Info("app starting...")

	if err := tracing.EnableTracing(ctx, a.Name, a); err != nil {
		logging.From(ctx).Fatal("failed to enable tracing", zap.Error(err))
	}

	listeners, err := onStart(ctx, a)
	if err != nil {
		logging.From(ctx).Fatal("failed to start app", zap.Error(err))
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		shutdown(ctx, a)
		os.Exit(1)
	}()

	var wg sync.WaitGroup

	for _, listener := range listeners {
		wg.Add(1)

		go func(l Listener) {
			defer wg.Done()

			err := l.Listen(ctx)
			if err != nil {
				logging.From(ctx).Error("listener failed", zap.Error(err))
			}
		}(listener)
	}

	wg.Wait()

	shutdown(ctx, a)
}

// OnShutdown registers a function that is called when the app is shutdown.
func (a *App) OnShutdown(onShutdown func()) {
	a.shutdownFuncs = append([]OnShutdownFunc{onShutdown}, a.shutdownFuncs...)
}

func shutdown(ctx context.Context, a *App) {
	for _, shutdownFunc := range a.shutdownFuncs {
		shutdownFunc()
	}

	logging.From(ctx).Info("app shutdown")
}
