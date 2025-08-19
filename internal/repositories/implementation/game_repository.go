package implementation

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nails_game/internal/models"
	"nails_game/internal/repositories/interfaces"
)

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) interfaces.GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) Create(game *models.Game) error {
	return r.db.Create(game).Error
}

func (r *gameRepository) GetByID(id uuid.UUID) (*models.Game, error) {
	var game models.Game
	if err := r.db.First(&game, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("game not found")
		}
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) Update(game *models.Game) error {
	return r.db.Save(game).Error
}
