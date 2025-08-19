package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"nails_game/internal/models"
	"nails_game/internal/models/dtos"
	"nails_game/internal/services/errors"
	services "nails_game/internal/services/interfaces"
)

type GameController struct {
	gameService services.GameService
}

func NewGameController(gameService services.GameService) *GameController {
	return &GameController{gameService: gameService}
}

// CreateGame создает новую игру
// @Summary Создать новую игру
// @Description Создает новую игру между двумя игроками
// @Tags games
// @Accept json
// @Produce json
// @Param request body dtos.CreateGameRequest true "Данные для создания игры"
// @Success 201 {object} dtos.CreateGameResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/game [post]
func (c *GameController) CreateGame(ctx echo.Context) error {
	var req dtos.CreateGameRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	game, err := c.gameService.CreateGame(req.LineSize, req.FirstPlayerID, req.SecondPlayerID)
	if err != nil {
		return handleServiceError(err)
	}

	resp := dtos.CreateGameResponse{
		GameID:         game.ID,
		LineSize:       len(game.Line),
		FirstPlayerID:  req.FirstPlayerID,
		SecondPlayerID: req.SecondPlayerID,
		Status:         game.Status.String(),
	}

	return ctx.JSON(http.StatusCreated, resp)
}

// MakeMove выполняет ход в игре
// @Summary Сделать ход
// @Description Выполняет ход в указанной игре
// @Tags games
// @Accept json
// @Produce json
// @Param gameId path string true "ID игры"
// @Param request body dtos.MoveRequest true "Данные хода"
// @Success 200 {object} dtos.GameStateResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/game/{gameId}/move [post]
func (c *GameController) MakeMove(ctx echo.Context) error {
	gameID, err := uuid.Parse(ctx.Param("gameId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid game ID")
	}

	var req dtos.MoveRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	move := models.Move{
		GameID:   gameID,
		PlayerID: req.PlayerID,
		Position: req.Position,
	}

	game, err := c.gameService.MakeMove(move)
	if err != nil {
		return handleServiceError(err)
	}

	resp := mapGameStateToResponse(game.Game)
	return ctx.JSON(http.StatusOK, resp)
}

// GetGame возвращает состояние игры
// @Summary Получить состояние игры
// @Description Возвращает текущее состояние указанной игры
// @Tags games
// @Produce json
// @Param gameId path string true "ID игры"
// @Success 200 {object} dtos.GameStateResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/game/{gameId} [get]
func (c *GameController) GetGame(ctx echo.Context) error {
	gameID, err := uuid.Parse(ctx.Param("gameId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid game ID")
	}

	game, err := c.gameService.GetGame(gameID)
	if err != nil {
		return handleServiceError(err)
	}

	resp := mapGameStateToResponse(game)
	return ctx.JSON(http.StatusOK, resp)
}

func mapGameStateToResponse(game *models.Game) dtos.GameStateResponse {
	return dtos.GameStateResponse{
		GameID:          game.ID,
		Status:          game.Status.String(),
		CurrentPlayerID: game.CurrentPlayerID,
		Line:            game.Line,
		MoveCount:       game.MoveCount,
	}
}

func handleServiceError(err error) *echo.HTTPError {
	switch err.(type) {
	case *errors.NotFoundError:
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	case *errors.UnauthorizedError:
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	case *errors.InvalidOperationError:
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}
