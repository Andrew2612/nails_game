package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"nails_game/internal/models"
)

type MockPlayerRepository struct {
	mock.Mock
}

func (m *MockPlayerRepository) Create(player *models.Player) error {
	return m.Called(player).Error(0)
}

func (m *MockPlayerRepository) GetByID(id uuid.UUID) (*models.Player, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) GetWithGames(id uuid.UUID) (*models.Player, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Player), args.Error(1)
}

func (m *MockPlayerRepository) Update(player *models.Player) error {
	return m.Called(player).Error(0)
}
