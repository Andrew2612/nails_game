package dtos

import "github.com/google/uuid"

// MoveRequest represents a move request
// @Description Запрос на выполнение хода
type MoveRequest struct {
	PlayerID uuid.UUID `json:"playerId"`
	Position int       `json:"position"`
}
