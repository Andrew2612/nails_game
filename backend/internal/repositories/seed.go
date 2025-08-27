package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nails_game/internal/models"
	"os"
)

type PlayerSeed struct {
	ID           string `json:"Id"`
	Name         string `json:"Name"`
	Email        string `json:"Email"`
	PasswordHash string `json:"PasswordHash"`
}

func SeedDatabase(db *gorm.DB, seedFilePath string) error {
	var count int64
	if err := db.Model(&models.Player{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	file, err := os.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	var players []PlayerSeed
	if err := json.Unmarshal(file, &players); err != nil {
		return fmt.Errorf("failed to unmarshal seed data: %w", err)
	}

	for _, p := range players {
		playerID, err := uuid.Parse(p.ID)
		if err != nil {
			return fmt.Errorf("invalid player ID in seed data: %w", err)
		}

		player := &models.Player{
			ID:           playerID,
			Name:         p.Name,
			Email:        p.Email,
			PasswordHash: p.PasswordHash,
		}

		if err := db.Create(player).Error; err != nil {
			return fmt.Errorf("failed to seed player: %w", err)
		}
	}

	return nil
}
