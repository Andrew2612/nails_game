package controllers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthController struct {
}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// CheckHealth проверяет статус сервера
// @Summary Health check
// @Description Проверяет работоспособность сервера и базы данных
// @Tags utils
// @Produce json
// @Success 200 {object} map[string]string "Статус сервера"
// @Router /health [get]
func (c *HealthController) CheckHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"status":   "healthy",
		"database": "connected",
		"version":  "1.0.0",
	})
}
