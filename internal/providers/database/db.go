package database

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/mauriciofsnts/hermes/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupConnection opens a connection to the configured PostgreSQL database.
func SetupConnection() (*gorm.DB, error) {
	cfg := config.Hermes.PG
	if cfg == nil {
		return nil, fmt.Errorf("database configuration is missing")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, strconv.Itoa(cfg.Port),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// CloseConnection closes the underlying database connection.
func CloseConnection(db *gorm.DB) error {
	if db == nil {
		return nil
	}

	dbSQL, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying database connection: %w", err)
	}

	if err := dbSQL.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
