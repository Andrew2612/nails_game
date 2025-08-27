package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string
	Email        string `gorm:"unique"`
	PasswordHash string
	Games        []*Game `gorm:"many2many:player_games;"`
}
