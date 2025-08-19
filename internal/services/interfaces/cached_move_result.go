package interfaces

import "nails_game/internal/models"

type CachedMoveResult struct {
	Game *models.Game
	ETag string
}
