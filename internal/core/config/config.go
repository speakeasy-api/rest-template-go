package config

import (
	"os"

	"github.com/speakeasy-api/rest-template-go/internal/core/drivers/psql"
	"github.com/speakeasy-api/rest-template-go/internal/core/errors"
	"github.com/speakeasy-api/rest-template-go/internal/core/listeners/http"
	"gopkg.in/yaml.v2"
)

const (
	// ErrRead is returned when we cannot read the config file.
	ErrRead = errors.Error("failed to read file")
	// ErrUnmarshal is returned when we cannot unmarshal the config file.
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

// AppConfig represents the configuration of our application.
type AppConfig struct {
	HTTP http.Config `yaml:"http"`
	PSQL psql.Config `yaml:"psql"`
}

// Load loads the configuration from a yaml file on disk.
func Load(cfg interface{}) error {
	data, err := os.ReadFile("config/config.yaml") // TODO support different environments
	if err != nil {
		return ErrRead.Wrap(err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return ErrUnmarshal.Wrap(err)
	}

	return nil
}
