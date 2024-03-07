package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go-planning-poker/model"
	"go-planning-poker/service"
	"net/http"
)

func checkOrigin(r *http.Request) bool {
	return true
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}
)

func GameWSHandler(c echo.Context) error {
	gameService := service.GetGameService()
	gameId := c.Param("game")
	game, err := gameService.FindGame(gameId)

	if err != nil {
		log.Errorf("Error when trying connecting to the game %s: [%s]", gameId, err)
		return c.String(http.StatusNotFound, fmt.Sprintf("Error when trying connecting to the game %s: [%s]", gameId, err))
	}

	log.Infof("New connection: %s", c.Request().RemoteAddr)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	player := model.NewPlayer("somePlayer", game, ws)
	go player.WriteMessage()
	game.AddPlayerToGame(player)
	player.ReadMessages()

	return nil
}
