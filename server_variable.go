package openapi

import (
	"errors"
	"fmt"
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// An object representing a Server Variable for server URL template substitution.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#server-variable-object
type ServerVariable struct {
	// An enumeration of string values to be used if the substitution options are from a limited set. The array MUST NOT be empty.
	Enum []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	// REQUIRED. The default value to use for substitution, which SHALL be sent if an alternate value is _not_ supplied. Note this behavior is different than the Schema Object's treatment of default values, because in those cases parameter values are optional. If the enum field is defined, the value MUST exist in the enum's values.
	Default string `json:"default" yaml:"default"`
	// An optional description for the server variable. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`

	// an index to the original location of this object
	idx int
}

func (s *ServerVariable) Validate() error {
	// either the array has entries or it is not defined
	if s.Enum != nil && len(s.Enum) == 0 {
		return errors.New("enum array must not be empty")
	}

	if s.Default == "" {
		return &errpath.ErrField{Field: "default", Err: &errpath.ErrRequired{}}
	}

	// If the enum is defined, the default value MUST exist in the enum's values.
	if len(s.Enum) > 0 {
		if !slices.Contains(s.Enum, s.Default) {
			return fmt.Errorf("default value %q must exist in the enum's values", s.Default)
		}
	}

	return nil
}

func getIndexServerVariable(v *ServerVariable) int                    { return v.idx }
func setIndexServerVariable(v *ServerVariable, i int) *ServerVariable { v.idx = i; return v }
