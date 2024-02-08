package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
)

type Player struct {
	Id          string
	Name        string
	Points      float32
	conn        *websocket.Conn
	game        *Game
	egress      chan *Message
	errorEgress chan []byte
}

func NewPlayer(name string, game *Game, conn *websocket.Conn) *Player {
	return &Player{
		Id:          uuid.NewString(),
		Name:        name,
		conn:        conn,
		game:        game,
		egress:      make(chan *Message),
		errorEgress: make(chan []byte),
	}
}

func (p *Player) ReadMessages() {
	defer p.game.RemovePlayer(p)
	for {
		_, payload, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				log.Errorf("error reading message: %v", err)
			}
			break
		}

		message, err := deserializePayload(payload)
		if err != nil {
			log.Errorf("Error deserializing the payload: %s", err)
			p.errorEgress <- []byte(err.Error())
			continue
		}

		log.Infof("Sending message [%s] to all other players.", payload)
		for _, v := range p.game.Players {
			if v != p {
				v.egress <- message
			}
		}
	}
}

func (p *Player) WriteMessage() {
	defer p.game.RemovePlayer(p)
	for {
		select {
		case message, ok := <-p.egress:
			if !ok {
				if err := p.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Infof("connection closed: %s\n", err)
				}
				return
			}

			payload, err := json.Marshal(message)

			if err != nil {
				log.Errorf("Error when marshaling message: %s, error message: %s", message, err)
				continue
			}

			if err := p.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Errorf("Failed to send message: %v", err)
			}

		case errorPayload, ok := <-p.errorEgress:
			if !ok {
				if err := p.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Infof("connection closed: %s\n", err)
				}
				return
			}

			if err := p.conn.WriteMessage(websocket.TextMessage, errorPayload); err != nil {
				log.Errorf("Failed to send error message: %v", err)
			}
		}
	}
}
