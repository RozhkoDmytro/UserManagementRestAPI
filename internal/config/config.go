package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/constants"
)

type Config struct {
	AppPort     string `required:"true" split_words:"true"`
	PostgresURI string `required:"true" split_words:"true"`
	JwtKey      string `required:"true" split_words:"true"`
}

func NewConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == constants.EmptyString {
		return nil, apperrors.EnvConfigVarError.AppendMessage("please set CONFIG_PATH")
	}

	err := godotenv.Load(configPath)
	if err != nil {
		return nil, apperrors.EnvConfigLoadError.AppendMessage(err)
	}

	config := &Config{}
	err = envconfig.Process("", config)
	if err != nil {
		return nil, apperrors.EnvConfigParseError.AppendMessage(err)
	}

	return config, nil
}
