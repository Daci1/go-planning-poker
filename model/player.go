package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"go-planning-poker/util"
)

type Player struct {
	Id          uuid.UUID
	Name        string
	VotedPoints float32
	conn        *websocket.Conn
	Game        *Game
	egress      chan *Message
	errorEgress chan []byte
}

func NewPlayer(name string, game *Game, conn *websocket.Conn) *Player {
	return &Player{
		Id:          uuid.New(),
		Name:        name,
		conn:        conn,
		Game:        game,
		egress:      make(chan *Message),
		errorEgress: make(chan []byte),
	}
}

func (p *Player) ReadMessages() {
	defer p.Game.RemovePlayer(p)
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

		switch message.Action {
		case OnJoin:
			p.egress <- message
			for _, v := range p.Game.Players {
				if v != p {
					v.egress <- message
				}
			}
			break
		case OnVote:
			if val, ok := message.Payload["votedPoints"].(float64); ok {
				truncatedFloat := util.TruncateFloat32(float32(val), 2)
				p.VotedPoints = truncatedFloat
				onVoteMessage := &Message{
					Action: OnVote,
					Payload: map[string]interface{}{
						"player": p.JSON(false),
					},
				}

				log.Infof("Sending message [%s] to all other players.", onVoteMessage)

				for _, v := range p.Game.Players {
					if v != p {
						v.egress <- onVoteMessage
					}
				}
			} else {
				p.errorEgress <- []byte("invalid voted points value, must be float64")
			}
			break
		case OnRevealResults:
			log.Infof("Sending message [%s] to all players.", payload)
			message.Payload = map[string]interface{}{
				"players": p.Game.SerializePlayers(true),
			}
			for _, v := range p.Game.Players {
				v.egress <- message
			}
			break
		}
	}
}

func (p *Player) WriteMessage() {
	defer p.Game.RemovePlayer(p)
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

func (p *Player) JSON(shouldIncludeVotedPoints bool) map[string]interface{} {
	serializedPlayer := map[string]interface{}{
		"id":   p.Id.String(),
		"name": p.Name,
	}

	if shouldIncludeVotedPoints {
		serializedPlayer["votedPoints"] = p.VotedPoints
	}

	return serializedPlayer
}
