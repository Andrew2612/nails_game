package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Move struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	GameID   uuid.UUID
	PlayerID uuid.UUID
	Position int
}
