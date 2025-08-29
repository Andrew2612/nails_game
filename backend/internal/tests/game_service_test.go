package tests

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	"nails_game/internal/models"
	"nails_game/internal/models/enums"
	services "nails_game/internal/services/implemenatation"
	"nails_game/internal/tests/mocks"
)

func createTestGame() *models.Game {
	firstPlayerID := uuid.New()
	return &models.Game{
		ID:              uuid.New(),
		Line:            make([]enums.PositionState, 9),
		Status:          enums.InProgress,
		CurrentPlayerID: firstPlayerID,
		FirstPlayerID:   firstPlayerID,
		SecondPlayerID:  uuid.New(),
		MoveCount:       0,
	}
}

func TestGameService_MakeMove(t *testing.T) {
	game := createTestGame()
	move := models.Move{
		GameID:   game.ID,
		PlayerID: game.FirstPlayerID,
		Position: 0,
	}

	mockGameRepo := new(mocks.MockGameRepository)
	mockPlayerRepo := new(mocks.MockPlayerRepository)

	mockGameRepo.On("GetByID", game.ID).Return(game, nil)
	mockGameRepo.On("Update", mock.Anything).Return(nil)

	service := services.NewGameService(mockGameRepo, mockPlayerRepo, 3)
	result, err := service.MakeMove(move)

	require.NoError(t, err)
	assert.Equal(t, game.ID, result.Game.ID)
	assert.NotEmpty(t, result.ETag)

	mockGameRepo.AssertExpectations(t)
}

func TestGameService_MakeMove_ReturnsCachedResult(t *testing.T) {
	game := createTestGame()
	firstMove := models.Move{
		GameID:   game.ID,
		PlayerID: game.FirstPlayerID,
		Position: 1,
	}

	secondMove := models.Move{
		GameID:   game.ID,
		PlayerID: game.SecondPlayerID,
		Position: 2,
	}

	mockGameRepo := new(mocks.MockGameRepository)
	mockPlayerRepo := new(mocks.MockPlayerRepository)

	mockGameRepo.On("GetByID", game.ID).Return(game, nil)
	mockGameRepo.On("Update", mock.Anything).Return(nil)

	service := services.NewGameService(mockGameRepo, mockPlayerRepo, 3)
	firstResult, err := service.MakeMove(firstMove)
	require.NoError(t, err)

	secondResult, err := service.MakeMove(secondMove)
	require.NoError(t, err)

	assert.NotEqual(t, firstResult.ETag, secondResult.ETag)
	assert.Equal(t, firstResult.Game.ID, secondResult.Game.ID)

	mockGameRepo.AssertNumberOfCalls(t, "GetByID", 2)
}

func TestGameService_MakeMove_NotPlayersTurn(t *testing.T) {
	game := createTestGame()
	game.CurrentPlayerID = game.SecondPlayerID

	move := models.Move{
		GameID:   game.ID,
		PlayerID: game.FirstPlayerID,
		Position: 1,
	}

	mockGameRepo := new(mocks.MockGameRepository)
	mockPlayerRepo := new(mocks.MockPlayerRepository)

	mockGameRepo.On("GetByID", game.ID).Return(game, nil)

	service := services.NewGameService(mockGameRepo, mockPlayerRepo, 3)
	_, err := service.MakeMove(move)

	assert.Error(t, err)
	assert.Equal(t, errors.New("not this player's turn"), err)
}

func TestGameService_MakeMove_PositionOccupied(t *testing.T) {
	game := createTestGame()
	game.Line[1] = enums.FirstPlayer

	move := models.Move{
		GameID:   game.ID,
		PlayerID: game.FirstPlayerID,
		Position: 1,
	}

	mockGameRepo := new(mocks.MockGameRepository)
	mockPlayerRepo := new(mocks.MockPlayerRepository)

	mockGameRepo.On("GetByID", game.ID).Return(game, nil)

	service := services.NewGameService(mockGameRepo, mockPlayerRepo, 3)
	_, err := service.MakeMove(move)

	assert.Error(t, err)
	assert.Equal(t, errors.New("position is already taken"), err)
}

func TestGameService_MakeMove_PlayerNotAssigned(t *testing.T) {
	game := createTestGame()
	randomPlayerID := uuid.New()

	move := models.Move{
		GameID:   game.ID,
		PlayerID: randomPlayerID,
		Position: 1,
	}

	mockGameRepo := new(mocks.MockGameRepository)
	mockPlayerRepo := new(mocks.MockPlayerRepository)

	mockGameRepo.On("GetByID", game.ID).Return(game, nil)

	service := services.NewGameService(mockGameRepo, mockPlayerRepo, 3)
	_, err := service.MakeMove(move)

	assert.Error(t, err)
	assert.Equal(t, errors.New("player is not in this game"), err)
}

func TestGameService_MakeMove_GameEnded(t *testing.T) {
	game := createTestGame()
	game.Status = enums.FirstPlayerWon

	move := models.Move{
		GameID:   game.ID,
		PlayerID: game.FirstPlayerID,
		Position: 1,
	}

	mockGameRepo := new(mocks.MockGameRepository)
	mockPlayerRepo := new(mocks.MockPlayerRepository)

	mockGameRepo.On("GetByID", game.ID).Return(game, nil)

	service := services.NewGameService(mockGameRepo, mockPlayerRepo, 3)
	_, err := service.MakeMove(move)

	assert.Error(t, err)
	assert.Equal(t, errors.New("game has already ended"), err)
}
