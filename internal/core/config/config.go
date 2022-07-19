package config

import (
	"io/ioutil"

	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/drivers/psql"
	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/errors"
	"github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/listeners/http"

	"gopkg.in/yaml.v2"
)

const (
	ErrRead      = errors.Error("failed to read file")
	ErrUnmarshal = errors.Error("failed to unmarshal file")
)

type AppConfig struct {
	HTTP http.Config `yaml:"http"`
	PSQL psql.Config `yaml:"psql"`
}

func Load(cfg interface{}) error {
	data, err := ioutil.ReadFile("config/config.yaml") // TODO support different environments
	if err != nil {
		return ErrRead.Wrap(err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return ErrUnmarshal.Wrap(err)
	}

	return nil
}
