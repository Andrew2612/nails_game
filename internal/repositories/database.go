package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nails_game/internal/models"
)

func InitDB(host, user, password, dbname, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Автомиграции
	if err := db.AutoMigrate(
		&models.Player{},
		&models.Game{},
		&models.Move{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	return db, nil
}
