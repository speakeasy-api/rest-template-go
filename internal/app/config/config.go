package config

import (
	"faceittechtest/internal/app"
	"faceittechtest/internal/app/drivers/psql"
	"faceittechtest/internal/app/listeners/http"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	ErrRead      = app.Error("failed to read file")
	ErrUnmarshal = app.Error("failed to unmarshal file")
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
