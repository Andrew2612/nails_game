package dtos

import "github.com/google/uuid"

// CreateGameResponse represents response for created game
// @Description Ответ с созданной игрой
type CreateGameResponse struct {
	GameID         uuid.UUID `json:"gameId"`
	LineSize       int       `json:"lineSize"`
	FirstPlayerID  uuid.UUID `json:"firstPlayerId"`
	SecondPlayerID uuid.UUID `json:"secondPlayerId"`
	Status         string    `json:"status"`
}
