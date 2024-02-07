package model

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type Action string

const (
	OnVote          Action = "OnVote"
	OnRevealResults Action = "OnRevealResults"
)

type Message struct {
	Sender uuid.UUID `json:"sender"`
	Action Action    `json:"action"`
}

func deserializePayload(payload []byte) (*Message, error) {
	message := Message{}
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	return &message, nil
}

func (a *Action) UnmarshalText(b []byte) error {
	action := string(b)
	switch action {
	case string(OnVote):
	case string(OnRevealResults):
		break
	default:
		return &Error{
			Name:        InvalidActionError,
			Description: fmt.Sprintf("Invalid action send from client: %s", action),
		}
	}
	*a = Action(action)
	return nil
}
