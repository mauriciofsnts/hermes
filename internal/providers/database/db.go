package database

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mauriciofsnts/hermes/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupConnection() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		panic("Failed to load env file")
	}

	host := config.Hermes.PG.Host
	port := config.Hermes.PG.Port
	user := config.Hermes.PG.User
	dbname := config.Hermes.PG.DBName
	password := config.Hermes.PG.Password

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, strconv.Itoa(port))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic("Failed to connect to database")
	}

	return db
}

func CloseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()

	if err != nil {
		panic("Failed to close connection")
	}

	dbSQL.Close()
}
