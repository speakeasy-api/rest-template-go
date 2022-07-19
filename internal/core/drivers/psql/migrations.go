package psql

import (
	"context"

	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/errors"
	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/logging"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// import file driver for migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	ErrDriverInit  = errors.Error("failed to initialize postgres driver")
	ErrMigrateInit = errors.Error("failed to initialize migration driver")
	ErrMigration   = errors.Error("failed to migrate database")
)

func (d *Driver) MigratePostgres(ctx context.Context, migrationsPath string) error {
	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		return ErrDriverInit.Wrap(err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return ErrMigrateInit.Wrap(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logging.From(ctx).Info("no migrations to run")
		} else {
			return ErrMigration.Wrap(err)
		}
	}

	logging.From(ctx).Info("migrations successfully run")

	return nil
}

func (d *Driver) RevertMigrations(ctx context.Context, migrationsPath string) error {
	driver, err := postgres.WithInstance(d.db.DB, &postgres.Config{})
	if err != nil {
		return ErrDriverInit.Wrap(err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return ErrMigrateInit.Wrap(err)
	}

	if err := m.Down(); err != nil {
		return ErrMigration.Wrap(err)
	}

	logging.From(ctx).Info("migrations successfully reverted")

	return nil
}
