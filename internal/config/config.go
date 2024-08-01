package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/constants"
)

type Config struct {
	AppPort  string `required:"true" split_words:"true"`
	Postgres *PostgresConfig
	Baseauth *BaseauthConfig
}

type PostgresConfig struct {
	Host     string
	Username string
	Password string
	Dbname   string
	Timezone string
}
type BaseauthConfig struct {
	Username string
	Password string
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
