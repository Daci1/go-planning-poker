package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Action string

const (
	OnInitialJoin   Action = "OnInitialJoin"
	OnJoin          Action = "OnJoin"
	OnVote          Action = "OnVote"
	OnRevealResults Action = "OnRevealResults"
)

type StringOrNumber string
type Payload map[string]interface{}

type Message struct {
	Action  Action  `json:"action"`
	Payload Payload `json:"payload,omitempty"`
}

func deserializePayload(payload []byte) (*Message, error) {
	message := Message{}
	if err := json.Unmarshal(payload, &message); err != nil {
		return nil, err
	}

	if err := validateMessage(&message); err != nil {
		return nil, err
	}

	return &message, nil
}

func validateMessage(m *Message) error {
	validate := validator.New()

	return validate.Struct(m)
}

func (a *Action) UnmarshalText(b []byte) error {
	action := string(b)
	switch action {
	case string(OnInitialJoin):
	case string(OnJoin):
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
