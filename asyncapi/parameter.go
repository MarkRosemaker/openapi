package asyncapi

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Parameter describes a parameter included in a channel address.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#parameterObject
type Parameter struct {
	// An enumeration of string values to be used if the substitution options are from a limited set.
	Enum []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	// The default value to use for substitution, and to send, if an alternate value is not supplied.
	Default string `json:"default,omitempty" yaml:"default,omitempty"`
	// An optional description for the parameter. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// An array of examples of the parameter value.
	Examples []string `json:"examples,omitempty" yaml:"examples,omitempty"`
	// A runtime expression that specifies the location of the parameter value.
	Location RuntimeExpression `json:"location,omitempty" yaml:"location,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the parameter for correctness.
func (p *Parameter) Validate() error {
	// either the array has entries or it is not defined
	if p.Enum != nil && len(p.Enum) == 0 {
		return errors.New("enum array must not be empty")
	}

	// if the enum is defined, the default value MUST exist in the enum's values
	if len(p.Enum) > 0 && p.Default != "" && !slices.Contains(p.Enum, p.Default) {
		return fmt.Errorf("default value %q must exist in the enum's values", p.Default)
	}

	p.Description = strings.TrimSpace(p.Description)

	if p.Location != "" {
		if err := p.Location.Validate(); err != nil {
			return &errpath.ErrField{Field: "location", Err: err}
		}
	}

	return validateExtensions(p.Extensions)
}

func (l *loader) collectParameterRef(p *ParameterRef, ref ref) {
	collectRef(l, p, l.parameters, ref)
}

func (l *loader) resolveParameterRef(p *ParameterRef) error {
	return resolveRef(p, l.parameters, nil)
}
