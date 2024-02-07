package model

import (
	"github.com/google/uuid"
	"sync"
)

type PlayerList map[uuid.UUID]*Player

type Game struct {
	Players PlayerList
	sync.RWMutex
}

func (g *Game) AddPlayer(p *Player) {
	g.Lock()
	defer g.Unlock()

	g.Players[p.Id] = p
}

func NewGame() *Game {
	return &Game{
		Players: make(PlayerList),
	}
}
