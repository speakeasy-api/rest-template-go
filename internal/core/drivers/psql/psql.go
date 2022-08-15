package psql

import (
	"context"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // imports the postgres driver
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
)

const (
	// ErrConnect is returned when we cannot connect to the database.
	ErrConnect = errors.Error("failed to connect to postgres db")
	// ErrClose is returned when we cannot close the database.
	ErrClose = errors.Error("failed to close postgres db connection")
)

// Config represents the configuration for our postgres database.
type Config struct {
	DSN string `env:"POSTGRES_DSN" validate:"required"`
}

// Driver provides an implementation for connecting to a postgres database.
type Driver struct {
	cfg Config
	db  *sqlx.DB
}

// New instantiates a instance of the Driver.
func New(cfg Config) *Driver {
	return &Driver{
		cfg: cfg,
	}
}

// Connect connects to the database.
func (d *Driver) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", d.cfg.DSN)
	if err != nil {
		return ErrConnect.Wrap(err)
	}

	d.db = db

	return nil
}

// Close closes the database connection.
func (d *Driver) Close(ctx context.Context) error {
	if err := d.db.Close(); err != nil {
		return ErrClose.Wrap(err)
	}

	return nil
}

// GetDB returns the underlying database connection.
func (d *Driver) GetDB() *sqlx.DB {
	return d.db
}
