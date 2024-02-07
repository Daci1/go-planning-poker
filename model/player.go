package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	Id          uuid.UUID
	Name        string
	conn        *websocket.Conn
	game        *Game
	egress      chan *Message
	errorEgress chan []byte
	// maybe error channel idk
}

func NewPlayer(name string, game *Game, conn *websocket.Conn) *Player {
	return &Player{
		Id:          uuid.New(),
		Name:        name,
		conn:        conn,
		game:        game,
		egress:      make(chan *Message),
		errorEgress: make(chan []byte),
	}
}

func (p *Player) ReadMessages() {
	for {
		messageType, payload, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				fmt.Printf("error reading message: %v", err)
			}
			break
		}

		message, err := deserializePayload(payload)
		if err != nil {
			fmt.Printf("Error deserializing the payload: %s", err)
			p.errorEgress <- []byte(err.Error())
			continue
		}

		for _, v := range p.game.Players {
			if v != p {
				v.egress <- message
			}
		}
		fmt.Println(messageType)
		fmt.Println(string(payload))
	}
}

func (p *Player) WriteMessage() {
	for {
		select {
		case message, ok := <-p.egress:
			if !ok {
				if err := p.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed: %s\n", err)
				}
				return
			}

			payload, err := json.Marshal(message)

			if err != nil {
				log.Printf("Error when marshaling message: %s, error message: %s", message, err)
				continue
			}

			if err := p.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Printf("Failed to send message: %v", err)
			}

		case errorPayload, ok := <-p.errorEgress:
			if !ok {
				if err := p.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Printf("connection closed: %s\n", err)
				}
				return
			}

			if err := p.conn.WriteMessage(websocket.TextMessage, errorPayload); err != nil {
				log.Printf("Failed to send error message: %v", err)
			}
		}
	}
}
