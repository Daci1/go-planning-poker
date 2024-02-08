package model

import (
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"sync"
)

type GameList map[string]*Game
type PlayerList map[string]*Player

type Game struct {
	GameId  string
	Players PlayerList
	sync.RWMutex
}

func (g *Game) AddPlayer(p *Player) {
	g.Lock()
	defer g.Unlock()

	g.Players[p.Id] = p
}

func (g *Game) RemovePlayer(p *Player) {
	g.Lock()
	defer g.Unlock()

	if _, ok := g.Players[p.Id]; ok {
		log.Infof("Removing player %s from game %s", p.Id, g.GameId)
		p.conn.Close()
		delete(g.Players, p.Id)
	}
}

func (g *Game) RemoveAllPlayers() {
	g.Lock()
	defer g.Unlock()

	log.Infof("Removing all players from game %s", g.GameId)

	for _, p := range g.Players {
		p.conn.Close()
		delete(g.Players, p.Id)
	}
}

func NewGame() *Game {
	return &Game{
		GameId:  uuid.NewString(),
		Players: make(PlayerList),
	}
}
