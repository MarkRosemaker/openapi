package asyncapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// OperationAction describes whether the application sends messages to a channel or receives messages from it.
type OperationAction string

const (
	// OperationActionSend is used when it's expected that the application will send a message to the given channel.
	OperationActionSend OperationAction = "send"
	// OperationActionReceive is used when the application should expect receiving messages from the given channel.
	OperationActionReceive OperationAction = "receive"
)

var allOperationActions = []OperationAction{
	OperationActionSend,
	OperationActionReceive,
}

// Validate validates the operation action.
func (a OperationAction) Validate() error {
	if slices.Contains(allOperationActions, a) {
		return nil
	}

	return &errpath.ErrInvalid[OperationAction]{
		Value: a,
		Enum:  allOperationActions,
	}
}
