package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"nails_game/internal/models"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) GetByID(id uuid.UUID) (*models.Game, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Game), args.Error(1)
}

func (m *MockGameRepository) Update(game *models.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func (m *MockGameRepository) Create(game *models.Game) error {
	args := m.Called(game)
	return args.Error(0)
}
