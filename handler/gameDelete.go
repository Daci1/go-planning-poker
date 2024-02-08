package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go-planning-poker/service"
	"net/http"
)

func DeleteGameHandler(c echo.Context) error {
	gameToDelete := c.Param("game")
	gameService := service.GetGameService()

	err := gameService.DeleteGame(gameToDelete)
	if err != nil {
		log.Errorf("Error occurred when deleting game: [%s]", err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error occurred when deleting game: [%s]", err))
	}

	return c.String(http.StatusOK, "")
}
