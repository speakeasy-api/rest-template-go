package config

import (
	"context"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"github.com/speakeasy-api/rest-template-go/internal/core/config"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/core/logging"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	// ErrInvalidEnvironment is returned when the SPEAKEASY_ENVIRONMENT environment variable is not set.
	ErrInvalidEnvironment = errors.Error("SPEAKEASY_ENVIRONMENT is not set")
	// ErrValidation is returned when the configuration is invalid.
	ErrValidation = errors.Error("invalid configuration")
	// ErrEnvVars is returned when the environment variables are invalid.
	ErrEnvVars = errors.Error("failed parsing env vars")
	// ErrRead is returned when the configuration file cannot be read.
	ErrRead = errors.Error("failed to read file")
	// ErrUnmarshal is returned when the configuration file cannot be unmarshalled.
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

var (
	baseConfigPath = "config/config.yaml"
	envConfigPath  = "config/config-%s.yaml"
)

// Config represents the configuration of our application.
type Config struct {
	config.AppConfig `yaml:",inline"`
}

// Load loads the configuration from the config/config.yaml file.
func Load(ctx context.Context) (*Config, error) {
	cfg := &Config{}

	if err := loadFromFiles(ctx, cfg); err != nil {
		return nil, err
	}

	if err := env.Parse(cfg); err != nil {
		return nil, ErrEnvVars.Wrap(err)
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return nil, ErrValidation.Wrap(err)
	}

	return cfg, nil
}

func loadFromFiles(ctx context.Context, cfg any) error {
	environ := os.Getenv("SPEAKEASY_ENVIRONMENT")
	if environ == "" {
		return ErrInvalidEnvironment
	}

	if err := loadYaml(ctx, baseConfigPath, cfg); err != nil {
		return err
	}

	p := fmt.Sprintf(envConfigPath, environ)

	if _, err := os.Stat(p); !errors.Is(err, os.ErrNotExist) {
		if err := loadYaml(ctx, p, cfg); err != nil {
			return err
		}
	}

	return nil
}

func loadYaml(ctx context.Context, filename string, cfg any) error {
	logging.From(ctx).Info("Loading configuration", zap.String("path", filename))

	data, err := os.ReadFile(filename)
	if err != nil {
		return ErrRead.Wrap(err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return ErrUnmarshal.Wrap(err)
	}

	return nil
}
