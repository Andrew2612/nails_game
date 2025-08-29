package repositories

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"nails_game/internal/models"
	"nails_game/internal/models/dtos"
)

func InitDB(logger *logrus.Logger) (*gorm.DB, *dtos.Config, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Player{},
		&models.Game{},
		&models.Move{},
	); err != nil {
		return nil, nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	if err := SeedDatabase(db, "seed_players.json"); err != nil {
		logger.WithError(err).Warn("Database seeding failed - continuing without seed data")
	}

	return db, cfg, nil
}

func loadConfig() (*dtos.Config, error) {
	lineSize, err := strconv.Atoi(os.Getenv("LINE_SIZE"))
	if err != nil {
		return nil, fmt.Errorf("invalid LINE_SIZE value: %v", err)
	}

	return &dtos.Config{
		DatabaseConfig: dtos.DatabaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		GameSettings: dtos.GameSettings{
			LineSize: lineSize,
		},
	}, nil
}
