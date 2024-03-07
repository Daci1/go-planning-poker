package model

import (
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"sync"
)

type GameList map[string]*Game
type PlayerList map[uuid.UUID]*Player

type Game struct {
	GameId  string
	Players PlayerList
	sync.RWMutex
}

func (g *Game) AddPlayerToGame(p *Player) {
	g.Lock()
	defer g.Unlock()

	g.Players[p.Id] = p
	g.notifyOtherPlayers(p)
	g.initCurrentPlayerGame(p)
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

func (g *Game) notifyOtherPlayers(newPlayer *Player) {
	newPlayerJoinMessage := &Message{
		Action: OnJoin,
		Payload: map[string]interface{}{
			"player": map[string]string{
				"id":   newPlayer.Id.String(),
				"name": newPlayer.Name,
			},
		},
	}

	for _, player := range g.Players {
		if player != newPlayer {

			player.egress <- newPlayerJoinMessage
		}
	}
}

func (g *Game) initCurrentPlayerGame(p *Player) {
	currentGameData := &Message{
		Action: OnInitialJoin,
		Payload: map[string]interface{}{
			"game": map[string]interface{}{
				"id":      g.GameId,
				"players": serializePlayersExceptCurrent(p, false),
			},
		},
	}

	p.egress <- currentGameData
}

func serializePlayersExceptCurrent(p *Player, shouldIncludeVotedPoints bool) []map[string]interface{} {
	players := p.Game.Players
	serializedPlayers := make([]map[string]interface{}, 0)

	for _, player := range players {
		if player != p {
			serializedPlayer := player.JSON(shouldIncludeVotedPoints)
			serializedPlayers = append(serializedPlayers, serializedPlayer)
		}
	}

	return serializedPlayers
}

func (g *Game) SerializePlayers(shouldIncludeVotedPoints bool) []map[string]interface{} {
	players := g.Players
	serializedPlayers := make([]map[string]interface{}, 0)

	for _, player := range players {
		serializedPlayer := player.JSON(shouldIncludeVotedPoints)
		serializedPlayers = append(serializedPlayers, serializedPlayer)
	}

	return serializedPlayers
}
