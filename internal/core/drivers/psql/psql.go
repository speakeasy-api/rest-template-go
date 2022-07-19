package psql

import (
	"context"
	"fmt"

	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	ErrConnect = errors.Error("failed to connect to postgres db")
	ErrClose   = errors.Error("failed to close postgres db connection")
)

type Config struct {
	Username string `yaml:"username"` // TODO get these from secrets in the future
	Password string `yaml:"password"` // TODO get these from secrets in the future
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type Driver struct {
	cfg Config
	db  *sqlx.DB
}

func New(cfg Config) *Driver {
	return &Driver{
		cfg: cfg,
	}
}

func (d *Driver) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", d.cfg.Username, d.cfg.Password, d.cfg.Host, d.cfg.Port, d.cfg.Database))
	if err != nil {
		return ErrConnect.Wrap(err)
	}

	d.db = db

	return nil
}

func (d *Driver) Close(ctx context.Context) error {
	if err := d.db.Close(); err != nil {
		return ErrClose.Wrap(err)
	}

	return nil
}

func (d *Driver) GetDB() *sqlx.DB {
	return d.db
}
