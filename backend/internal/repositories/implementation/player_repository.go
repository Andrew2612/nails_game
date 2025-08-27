package implementation

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nails_game/internal/models"
	"nails_game/internal/repositories/interfaces"
)

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) interfaces.PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) GetByID(id uuid.UUID) (*models.Player, error) {
	var player models.Player
	if err := r.db.First(&player, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("player not found")
		}
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) GetWithGames(id uuid.UUID) (*models.Player, error) {
	var player models.Player
	if err := r.db.Preload("Games").First(&player, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("player not found")
		}
		return nil, err
	}
	return &player, nil
}

func (r *playerRepository) Update(player *models.Player) error {
	return r.db.Save(player).Error
}
