package main

import (
	"context"
	"faceittechtest/internal/app"
	"faceittechtest/internal/app/drivers/psql"
	"faceittechtest/internal/app/listeners/http"
	"faceittechtest/internal/app/logging"
	"faceittechtest/internal/config"
	"faceittechtest/internal/events"
	httptransport "faceittechtest/internal/transport/http"
	"faceittechtest/internal/users"
	"faceittechtest/internal/users/store"

	"go.uber.org/zap"
)

func main() {
	app.Start(appStart)
}

func appStart(ctx context.Context, a *app.App) ([]app.Listener, error) {
	// Load configuration from config/config.yaml which contains details such as DB connection params
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Connect to the postgres DB
	db := psql.New(cfg.PSQL)
	if err := db.Connect(ctx); err != nil {
		return nil, err
	}
	a.OnShutdown(func() {
		// Shutdown connection when server terminated
		logging.From(ctx).Info("shutting down db connection")
		if err := db.Close(ctx); err != nil {
			logging.From(ctx).Error("failed to close db connection", zap.Error(err))
		}
	})

	// Run our migrations which will update the DB or create it if it doesn't exist
	if err := db.MigratePostgres(ctx, "file://migrations"); err != nil {
		return nil, err
	}
	a.OnShutdown(func() {
		// Temp for development so database is cleared on shutdown
		if err := db.RevertMigrations(ctx, "file://migrations"); err != nil {
			logging.From(ctx).Error("failed to revert migrations", zap.Error(err))
		}
	})

	// Instantiate and connect all our classes
	us := store.New(db.GetDB())
	e := events.New()
	u := users.New(us, e)

	httpServer := httptransport.New(u, db.GetDB())

	// Create a HTTP server
	h, err := http.New(httpServer, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	// Start listening for HTTP requests
	return []app.Listener{
		h,
	}, nil
}
