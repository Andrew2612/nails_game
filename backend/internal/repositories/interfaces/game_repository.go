package interfaces

import (
	"github.com/google/uuid"
	"nails_game/internal/models"
)

type GameRepository interface {
	Create(game *models.Game) error
	GetByID(id uuid.UUID) (*models.Game, error)
	Update(game *models.Game) error
}
