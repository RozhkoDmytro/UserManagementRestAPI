package database

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
	postgresURI := os.Getenv("POSTGRES_URI")
	if postgresURI == "" {
		return nil, errors.New("can't finde POSTGRES_URI")
	}

	// Parse the PostgreSQL URI
	parsedURL, err := url.Parse(postgresURI)
	if err != nil {
		return nil, err
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
