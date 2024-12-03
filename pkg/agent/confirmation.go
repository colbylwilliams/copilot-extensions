package agent

import (
	"encoding/json"
	"errors"
)

// Confirmation represents a confirmation that the agent has sent to the user.
type Confirmation struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Message      string `json:"message"`
	Confirmation any    `json:"confirmation"`
}

var (
	ConfirmationTypeAction = ConfirmationType{"action"}
)

type ConfirmationType struct {
	name string
}

func (c *ConfirmationType) UnmarshalJSON(data []byte) error {
	var action string
	if err := json.Unmarshal(data, &action); err != nil {
		return err
	}

	if action != ConfirmationTypeAction.name {
		return errors.New("invalid agent confirmation type")
	}

	c.name = action
	return nil
}

func (c ConfirmationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.name)
}

// ClientConfirmation represents a confirmation that the user has accepted or dismissed.
type ClientConfirmation struct {
	State        ClientConfirmationState `json:"state"`
	Confirmation any                     `json:"confirmation"`
}

var (
	ClientConfirmationStateAccepted  = ClientConfirmationState{"accepted"}
	ClientConfirmationStateDismissed = ClientConfirmationState{"dismissed"}
)

type ClientConfirmationState struct {
	name string
}

func (c *ClientConfirmationState) UnmarshalJSON(data []byte) error {
	var state string
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	if state != ClientConfirmationStateAccepted.name && state != ClientConfirmationStateDismissed.name {
		return errors.New("invalid client confirmation state")
	}

	c.name = state
	return nil
}

func (c ClientConfirmationState) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.name)
}
