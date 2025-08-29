package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "nails_game/docs"

	"nails_game/internal/controllers"
	"nails_game/internal/helpers"
	database "nails_game/internal/repositories"
	repositories "nails_game/internal/repositories/implementation"
	services "nails_game/internal/services/implemenatation"
)

// @title Nails Game API
// @version 1.0
// @description API для игры в гвоздики
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := helpers.InitLogger()

	if err := godotenv.Load(".env"); err != nil {
		logger.Warnf("Could not load .env file: %v", err)
	}

	db, cfg, err := database.InitDB(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize database")
	}

	logger.Info("Application configuration loaded")

	gameRepo := repositories.NewGameRepository(db)
	playerRepo := repositories.NewPlayerRepository(db)

	gameService := services.NewGameService(gameRepo, playerRepo, cfg.LineSize)

	gameController := controllers.NewGameController(gameService)
	healthController := controllers.NewHealthController()

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/api/game", gameController.CreateGame)
	e.POST("/api/game/:gameId/move", gameController.MakeMove)
	e.GET("/api/game/:gameId", gameController.GetGame)

	e.GET("/health", healthController.CheckHealth)

	logger.WithField("port", cfg.Port).Info("Starting server")
	if err := e.Start(":" + cfg.Port); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}
