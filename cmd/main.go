package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "nails_game/docs"

	"nails_game/internal/controllers"
	"nails_game/internal/models/dtos"
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
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := loadConfig()

	db, err := database.InitDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
	)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := database.SeedDatabase(db, "seed_players.json"); err != nil {
		log.Printf("Warning: Failed to seed database: %v", err)
	}

	gameRepo := repositories.NewGameRepository(db)
	playerRepo := repositories.NewPlayerRepository(db)

	gameService := services.NewGameService(gameRepo, playerRepo, cfg.GameSettings.LineSize)

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

	port := getEnv("PORT", "8080")
	e.Logger.Fatal(e.Start(":" + port))
}

func loadConfig() dtos.Config {
	var cfg dtos.Config

	lineSize, err := strconv.Atoi(getEnv("LINE_SIZE", "4"))
	if err != nil {
		log.Fatalf("Invalid LINE_SIZE value: %v", err)
	}
	cfg.GameSettings.LineSize = lineSize

	cfg.Database.Host = getEnv("POSTGRES_HOST", "localhost")
	cfg.Database.User = getEnv("POSTGRES_USER", "postgres")
	cfg.Database.Password = getEnv("POSTGRES_PASSWORD", "secret")
	cfg.Database.DBName = getEnv("POSTGRES_DB", "nails_db")
	cfg.Database.Port = getEnv("POSTGRES_PORT", "5432")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
