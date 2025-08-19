package implemenatation

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"

	"nails_game/internal/models"
	"nails_game/internal/models/enums"
	repositories "nails_game/internal/repositories/interfaces"
	services "nails_game/internal/services/interfaces"
)

type gameService struct {
	gameRepo   repositories.GameRepository
	playerRepo repositories.PlayerRepository
	lineSize   int
	cache      map[string]services.CachedMoveResult
	cacheMutex sync.RWMutex
	moveLocks  map[string]*sync.Mutex
	locksMutex sync.Mutex
}

func NewGameService(
	gameRepo repositories.GameRepository,
	playerRepo repositories.PlayerRepository,
	lineSize int,
) services.GameService {
	return &gameService{
		gameRepo:   gameRepo,
		playerRepo: playerRepo,
		lineSize:   lineSize,
		cache:      make(map[string]services.CachedMoveResult),
		moveLocks:  make(map[string]*sync.Mutex),
	}
}

func (s *gameService) CreateGame(lineSize int, firstPlayerID, secondPlayerID uuid.UUID) (*models.Game, error) {
	if _, err := s.playerRepo.GetByID(firstPlayerID); err != nil {
		return nil, fmt.Errorf("first player not found: %w", err)
	}
	if _, err := s.playerRepo.GetByID(secondPlayerID); err != nil {
		return nil, fmt.Errorf("second player not found: %w", err)
	}

	game := &models.Game{
		ID:              uuid.New(),
		Line:            make([]enums.PositionState, lineSize),
		Status:          enums.InProgress,
		CurrentPlayerID: firstPlayerID,
		FirstPlayerID:   firstPlayerID,
		SecondPlayerID:  secondPlayerID,
		MoveCount:       0,
	}

	if err := s.gameRepo.Create(game); err != nil {
		return nil, fmt.Errorf("failed to create game: %w", err)
	}

	return game, nil
}

func (s *gameService) MakeMove(move models.Move) (*services.CachedMoveResult, error) {
	cacheKey := s.generateCacheKey(move)

	s.cacheMutex.RLock()
	cached, exists := s.cache[cacheKey]
	s.cacheMutex.RUnlock()

	if exists {
		return &cached, nil
	}

	lock := s.getMoveLock(cacheKey)
	lock.Lock()
	defer lock.Unlock()

	s.cacheMutex.RLock()
	cached, exists = s.cache[cacheKey]
	s.cacheMutex.RUnlock()

	if exists {
		return &cached, nil
	}

	game, err := s.gameRepo.GetByID(move.GameID)
	if err != nil {
		return nil, fmt.Errorf("game not found: %w", err)
	}

	if move.Position < 0 || move.Position >= len(game.Line) {
		return nil, errors.New("invalid move position")
	}

	if game.Status != enums.InProgress {
		return nil, errors.New("game has already ended")
	}

	if game.FirstPlayerID != move.PlayerID && game.SecondPlayerID != move.PlayerID {
		return nil, errors.New("player is not in this game")
	}

	if game.CurrentPlayerID != move.PlayerID {
		return nil, errors.New("not this player's turn")
	}

	if game.Line[move.Position] != enums.Empty {
		return nil, errors.New("position is already taken")
	}

	if move.PlayerID == game.FirstPlayerID {
		game.Line[move.Position] = enums.FirstPlayer
	} else {
		game.Line[move.Position] = enums.SecondPlayer
	}

	game.MoveCount++
	game.CurrentPlayerID = s.getNextPlayerID(game)
	game.Status = s.checkGameStatus(game)

	if err := s.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	result := services.CachedMoveResult{
		Game: game,
		ETag: uuid.New().String(),
	}

	s.cacheMutex.Lock()
	s.cache[cacheKey] = result
	s.cacheMutex.Unlock()

	go func() {
		time.Sleep(3 * time.Minute)
		s.cacheMutex.Lock()
		delete(s.cache, cacheKey)
		s.cacheMutex.Unlock()
	}()

	return &result, nil
}

func (s *gameService) GetGame(gameID uuid.UUID) (*models.Game, error) {
	return s.gameRepo.GetByID(gameID)
}

func (s *gameService) getNextPlayerID(game *models.Game) uuid.UUID {
	if game.CurrentPlayerID == game.FirstPlayerID {
		return game.SecondPlayerID
	}
	return game.FirstPlayerID
}

func (s *gameService) checkGameStatus(game *models.Game) enums.GameStatus {
	allPosTaken := true
	for _, pos := range game.Line {
		if pos == enums.Empty {
			allPosTaken = false
			break
		}
	}

	if !allPosTaken {
		return enums.InProgress
	}

	firstPlayerSum := s.getNailsSum(game, enums.FirstPlayer)
	secondPlayerSum := s.getNailsSum(game, enums.SecondPlayer)

	if firstPlayerSum > secondPlayerSum {
		return enums.FirstPlayerWon
	}
	return enums.SecondPlayerWon
}

func (s *gameService) getNailsSum(game *models.Game, state enums.PositionState) int {
	var nails []int
	for i, pos := range game.Line {
		if pos == state {
			nails = append(nails, i)
		}
	}

	n := len(nails)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return nails[0]
	}
	if n == 2 {
		return nails[1] - nails[0]
	}

	twoBack := nails[1] - nails[0]
	oneBack := nails[2] - nails[1] + twoBack

	for i := 3; i < n; i++ {
		oldOneBack := oneBack
		oneBack = nails[i] - nails[i-1] + min(oneBack, twoBack)
		twoBack = oldOneBack
	}

	return oneBack
}

func (s *gameService) generateCacheKey(move models.Move) string {
	return fmt.Sprintf("move:%s:%s:%d", move.GameID, move.PlayerID, move.Position)
}

func (s *gameService) getMoveLock(key string) *sync.Mutex {
	s.locksMutex.Lock()
	defer s.locksMutex.Unlock()

	if lock, exists := s.moveLocks[key]; exists {
		return lock
	}

	lock := &sync.Mutex{}
	s.moveLocks[key] = lock
	return lock
}
