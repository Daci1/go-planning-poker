package model

import (
	"encoding/json"
)

const (
	InvalidActionError = "InvalidActionError"
)

type Error struct {
	Name        string `json:"Name"`
	Description string `json:"description"`
}

func (e *Error) Error() string {
	// marshaledError, _ := json.Marshal(fmt.Sprintf("%s: %s", e.Name, e.Description))
	marshaledError, _ := json.Marshal(e)
	return string(marshaledError)
}
