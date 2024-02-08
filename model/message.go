package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
	"math"
	"strconv"
)

type Action string
type VotedPoints float32

const (
	OnVote          Action = "OnVote"
	OnRevealResults Action = "OnRevealResults"
)

type Message struct {
	Sender      uuid.UUID   `json:"sender"`
	Action      Action      `json:"action"`
	VotedPoints VotedPoints `json:"votedPoints" validate:"votedPointsValidator"`
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
	validate.RegisterValidation("votedPointsValidator", m.votedPointsValidator)

	return validate.Struct(m)
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

func (v *VotedPoints) UnmarshalText(b []byte) error {
	strVotedPoint := string(b)
	votedPoint, err := strconv.ParseFloat(strVotedPoint, 32)
	if err != nil {
		log.Errorf("Error parsing votedPoints [%s]: [%s]", strVotedPoint, err)
		return errors.New(fmt.Sprintf("Error parsing votedPoints [%s]: [%s]", strVotedPoint, err))
	}

	*v = VotedPoints(math.Floor(votedPoint*100) / 100)
	return nil
}

func (m *Message) votedPointsValidator(fl validator.FieldLevel) bool {
	return m.Action == OnVote && m.VotedPoints >= 0
}
