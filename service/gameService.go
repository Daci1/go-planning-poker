package service

import (
	"errors"
	"go-planning-poker/model"
	"sync"
)

type GameService struct {
	games model.GameList
	sync.RWMutex
}

var lock = &sync.RWMutex{}
var singleInstance *GameService

func GetGameService() *GameService {
	lock.Lock()
	defer lock.Unlock()
	if singleInstance == nil {
		singleInstance = &GameService{
			games: make(model.GameList),
		}
		singleInstance.games["default"] = &model.Game{
			Players: make(model.PlayerList),
			GameId:  "default",
		}
	}
	return singleInstance
}

func (gs *GameService) CreateGame() string {
	gs.Lock()
	defer gs.Unlock()

	game := model.NewGame()
	gs.games[game.GameId] = game

	return game.GameId
}

func (gs *GameService) FindGame(gId string) (*model.Game, error) {
	gs.Lock()
	gs.Unlock()

	if game, ok := gs.games[gId]; ok {
		return game, nil
	}
	return nil, errors.New("game not found")
}

func (gs *GameService) DeleteGame(gId string) error {
	gs.Lock()
	defer gs.Unlock()

	if _, ok := gs.games[gId]; ok {
		gs.games[gId].RemoveAllPlayers()
		delete(gs.games, gId)
		return nil
	}

	return errors.New("game not found")
}
