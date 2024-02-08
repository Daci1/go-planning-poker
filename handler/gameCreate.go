package handler

import (
	"github.com/labstack/echo/v4"
	"go-planning-poker/service"
	"net/http"
)

type createGameResponse struct {
	GameId string `json:"gameId"`
}

func CreateGame(c echo.Context) error {
	gameService := service.GetGameService()
	gameId := gameService.CreateGame()

	return c.JSON(http.StatusOK, &createGameResponse{GameId: gameId})
}
