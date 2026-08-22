package asyncapi

import (
	"encoding/json/jsontext"
	"errors"

	"github.com/MarkRosemaker/errpath"
)

// ErrEmptyMessageExample is returned when a message example neither has headers nor a payload.
var ErrEmptyMessageExample = errors.New("must contain either headers and/or payload")

// MessageExample represents an example of a [Message] object
// and MUST contain either headers and/or payload fields.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#messageExampleObject
type MessageExample struct {
	// The value of this field MUST validate against the headers of the message.
	Headers jsontext.Value `json:"headers,omitempty" yaml:"headers,omitempty"`
	// The value of this field MUST validate against the payload of the message.
	Payload jsontext.Value `json:"payload,omitempty" yaml:"payload,omitempty"`
	// A machine-friendly name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// A short summary of what the example is about.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the message example for correctness.
func (ex *MessageExample) Validate() error {
	if len(ex.Headers) == 0 && len(ex.Payload) == 0 {
		return ErrEmptyMessageExample
	}

	return validateExtensions(ex.Extensions)
}

// MessageExamples is a list of examples of a message.
type MessageExamples []*MessageExample

// Validate validates each example.
func (exs MessageExamples) Validate() error {
	for i, ex := range exs {
		if err := ex.Validate(); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
