package psql

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // import file driver for migrate
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
)

const (
	// ErrDriverInit is returned when we cannot initialize the driver.
	ErrDriverInit = errors.Error("failed to initialize postgres driver")
	// ErrMigrateInit is returned when we cannot initialize the migrate driver.
	ErrMigrateInit = errors.Error("failed to initialize migration driver")
	// ErrMigration is returned when we cannot run a migration.
	ErrMigration = errors.Error("failed to migrate database")
)

// MigratePostgres migrates the database to the latest version.
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

// RevertMigrations reverts the database to the previous version.
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
