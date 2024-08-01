package database

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		return SetupDatabaseFromConfig(cfg)
	}

	// Parse the PostgreSQL URI
	parsedURL, err := url.Parse(postgresURI)
	if err != nil {
		return SetupDatabaseFromConfig(cfg)
	}

	// Extract components from URI
	user := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()
	hostPort := strings.Split(parsedURL.Host, ":")
	host := hostPort[0]
	port := "5432" // Default port if not specified

	if len(hostPort) > 1 {
		port = hostPort[1]
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, parsedURL.Path[1:], port,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func SetupDatabaseFromConfig(cfg *config.Config) (*gorm.DB, error) {
	if cfg.Postgres == nil {
		return nil, &apperrors.NilPostgresConfigError
	}

	postgresCfg := cfg.Postgres
	splitHost := strings.Split(postgresCfg.Host, ":")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		splitHost[0], postgresCfg.Username, postgresCfg.Password, postgresCfg.Dbname, splitHost[1], postgresCfg.Timezone,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}
