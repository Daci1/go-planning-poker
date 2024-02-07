package handler

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go-planning-poker/model"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var game = model.NewGame()

func GameWSHandler(c echo.Context) error {
	log.Infof("New connection: %s", c.Request().RemoteAddr)
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	player := model.NewPlayer("somePlayer", game, ws)
	game.AddPlayer(player)
	go player.WriteMessage()
	player.ReadMessages()

	return nil
}
