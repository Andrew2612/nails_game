package dtos

import (
	"github.com/google/uuid"
	"nails_game/internal/models/enums"
)

// GameStateResponse represents game state
// @Description Состояние игры
type GameStateResponse struct {
	GameID          uuid.UUID             `json:"gameId"`
	Status          string                `json:"status"`
	CurrentPlayerID uuid.UUID             `json:"currentPlayerId"`
	Line            []enums.PositionState `json:"line"`
	MoveCount       int                   `json:"moveCount"`
}
