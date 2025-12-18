package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
}

func (cfg Config) ConnectionPostgres() (*Postgres, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PsqlDB.User,
		cfg.PsqlDB.Password,
		cfg.PsqlDB.Host,
		cfg.PsqlDB.Port,
		cfg.PsqlDB.DBName)

	db, err := gorm.Open(postgres.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-1] Failed to connect to database " + cfg.PsqlDB.Host)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionPostgres-3] Failed to get database connection")
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.PsqlDB.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.PsqlDB.DBMaxIdle)

	return &Postgres{DB: db}, nil
}

// Tambahan khusus untuk Goose (butuh *sql.DB, bukan *gorm.DB)
func (cfg Config) ConnectionSqlDB() (*sql.DB, error) {
	dbConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PsqlDB.User,
		cfg.PsqlDB.Password,
		cfg.PsqlDB.Host,
		cfg.PsqlDB.Port,
		cfg.PsqlDB.DBName)

	sqlDB, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Error().Err(err).Msg("[ConnectionSqlDB-1] Failed to connect to database " + cfg.PsqlDB.Host)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.PsqlDB.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.PsqlDB.DBMaxIdle)

	return sqlDB, nil
}
