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
	ProfileImageBaseUrl string `envconfig:"PROFILE_IMAGE_BASE_URL" required:"true" default:"https://api.goyotashi.kmconner.net/images"`
	DBHost              string `envconfig:"DB_HOST" default:"goyotashi-mysql"`
	DBPort              uint16 `envconfig:"DB_PORT" default:"3306"`
	DBDatabaseName      string `envconfig:"DB_NAME" default:"goyotashi"`
	DBUser              string `envconfig:"DB_USER" default:"goyotashi"`
	DBPassword          string `envconfig:"DB_PASSWORD" default:"password"`
	HotpepperApiKey     string `envconfig:"HOTPEPPER_KEY" required:"true"`
}
