package config

import "github.com/kelseyhightower/envconfig"

type HermesServiceConfig struct {
	Env   string `envconfig:"ENV"`
	Port  string `envconfig:"PORT" default:"8080"`
	Helix HelixConfig
}

type HelixConfig struct {
	HelixCredentials string `envconfig:"HELIX_CREDENTIALS"`
	HelixHost        string `envconfig:"HELIX_HOST" required:"true"`
	HelixTemplateID  int64  `envconfig:"HELIX_TEMPLATE_ID" required:"true"`
	HelixTrackingTag string `envconfig:"HELIX_TRACKING_TAG" required:"true"`
}

func GetConfig() (*HermesServiceConfig, error) {
	var cfg HermesServiceConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
