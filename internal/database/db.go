package database

import (
	"fmt"
	"strings"

	"gitlab.com/jkozhemiaka/web-layout/internal/apperrors"
	"gitlab.com/jkozhemiaka/web-layout/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(cfg *config.Config) (*gorm.DB, error) {
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
