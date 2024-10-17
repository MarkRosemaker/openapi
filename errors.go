package openapi

import (
	"fmt"
	"strings"
)

// ErrRequired signals that a required field is missing.
type ErrRequired struct{ Target string }

// Error returns the name of the missing field.
func (e *ErrRequired) Error() string {
	return fmt.Sprintf("%s is required", e.Target)
}

// ErrInvalid signals that a field has an invalid value.
type ErrInvalid struct {
	Target string
	Value  string
	// Enum    []string
	Message string
}

// Error returns helpful information about the invalid field and how to fix it.
func (e *ErrInvalid) Error() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%s (%q) is invalid", e.Target, e.Value)
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
		return b.String()
	}

	// if len(e.Enum) > 0 {
	// 	b.WriteString(", must be one of: ")
	// 	b.WriteString(strings.Join(e.Enum, ", "))
	// }

	return b.String()
}
