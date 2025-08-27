package interfaces

import (
	"github.com/google/uuid"
	"nails_game/internal/models"
)

type PlayerRepository interface {
	Create(player *models.Player) error
	GetByID(id uuid.UUID) (*models.Player, error)
	GetWithGames(id uuid.UUID) (*models.Player, error)
	Update(player *models.Player) error
}
