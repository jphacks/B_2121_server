package config

import (
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/xerrors"
)

func LoadServerConfig() (*ServerConfig, error) {
	var conf ServerConfig
	if err := envconfig.Process("", &conf); err != nil {
		return nil, xerrors.Errorf("failed to load configuration: %w", err)
	}
	return &conf, nil
}

type ServerConfig struct {
	ProfileImageBaseUrl string `envconfig:"PROFILE_IMAGE_BASE_URL" required:"true" default:"http://localhost:8080/images"`
}
