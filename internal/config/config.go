package config

import "github.com/speakeasy-api/speakeasy-example-rest-service-go/internal/core/config"

// Config represents the configuration of our application
type Config struct {
	config.AppConfig `yaml:",inline"`
}

// Load loads the configuration from the config/config.yaml file
func Load() (*Config, error) {
	cfg := &Config{}

	err := config.Load(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
