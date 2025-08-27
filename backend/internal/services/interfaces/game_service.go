package interfaces

import (
	"github.com/google/uuid"

	"nails_game/internal/models"
)

type GameService interface {
	CreateGame(lineSize int, firstPlayerID, secondPlayerID uuid.UUID) (*models.Game, error)
	MakeMove(move models.Move) (*CachedMoveResult, error)
	GetGame(gameID uuid.UUID) (*models.Game, error)
}
