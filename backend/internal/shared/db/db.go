package db

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func Connect(cfg Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection with retries
	var lastErr error
	for i := 0; i < 5; i++ {
		if err := db.Ping(); err == nil {
			log.Info().Msg("Successfully connected to database")
			return db, nil
		} else {
			lastErr = err
			log.Warn().Err(err).Msgf("Failed to ping database (attempt %d/5), retrying in 2s...", i+1)
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf("failed to ping database after 5 attempts: %w", lastErr)
}
