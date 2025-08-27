package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nails_game/internal/models/enums"
)

type Game struct {
	gorm.Model
	ID              uuid.UUID             `gorm:"type:uuid;primaryKey"`
	Line            []enums.PositionState `gorm:"type:integer[]"`
	Status          enums.GameStatus
	CurrentPlayerID uuid.UUID
	FirstPlayerID   uuid.UUID
	SecondPlayerID  uuid.UUID
	MoveCount       int
}
