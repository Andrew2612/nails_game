package dtos

import "github.com/google/uuid"

// CreateGameRequest represents request for creating a game
// @Description Запрос на создание игры
type CreateGameRequest struct {
	LineSize       int       `json:"line_size"`
	FirstPlayerID  uuid.UUID `json:"firstPlayerId"`
	SecondPlayerID uuid.UUID `json:"secondPlayerId"`
}
