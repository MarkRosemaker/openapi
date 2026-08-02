package asyncapi

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ServerVariable is an object representing a Server Variable for server URL template substitution.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#serverVariableObject
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options are from a limited set.
	Enum []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	// The default value to use for substitution, and to send, if an alternate value is not supplied.
	Default string `json:"default,omitempty" yaml:"default,omitempty"`
	// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// An array of examples of the server variable.
	Examples []string `json:"examples,omitempty" yaml:"examples,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the server variable for correctness.
func (v *ServerVariable) Validate() error {
	// either the array has entries or it is not defined
	if v.Enum != nil && len(v.Enum) == 0 {
		return errors.New("enum array must not be empty")
	}

	// if the enum is defined, the default value MUST exist in the enum's values
	if len(v.Enum) > 0 && v.Default != "" && !slices.Contains(v.Enum, v.Default) {
		return fmt.Errorf("default value %q must exist in the enum's values", v.Default)
	}

	v.Description = strings.TrimSpace(v.Description)

	return validateExtensions(v.Extensions)
}

func (l *loader) collectServerVariableRef(v *ServerVariableRef, ref ref) {
	collectRef(l, v, l.serverVariables, ref)
}

func (l *loader) resolveServerVariableRef(v *ServerVariableRef) error {
	return resolveRef(v, l.serverVariables, nil)
}
