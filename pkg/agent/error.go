package agent

import (
	"encoding/json"
	"errors"
)

// Error represents an error that occurred during the agent request.
type Error struct {
	Type       ErrorType `json:"type"`
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	Identifier string    `json:"identifier"`
}

func (a *Error) Error() string {
	return a.Message
}

type ErrorType struct {
	name string
}

func (a ErrorType) String() string {
	return a.name
}

var (
	ErrorTypeReference = ErrorType{"reference"}
	ErrorTypeFunction  = ErrorType{"function"}
	ErrorTypeAgent     = ErrorType{"agent"}
)

func (a *ErrorType) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	if name != ErrorTypeAgent.name && name != ErrorTypeFunction.name && name != ErrorTypeReference.name {
		return errors.New("invalid agent error type, got " + name)
	}

	a.name = name
	return nil
}

func (a ErrorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.name)
}
