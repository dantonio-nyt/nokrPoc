package config

import "github.com/kelseyhightower/envconfig"

type HermesServiceConfig struct {
	Env       string `envconfig:"ENV"`
	Port      string `envconfig:"PORT" default:"8080"`
}

func GetConfig() (*HermesServiceConfig, error) {
	var cfg HermesServiceConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}