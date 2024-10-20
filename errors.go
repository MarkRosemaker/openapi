package openapi

import (
	"errors"
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

// ErrField is an error that occurred in a field.
type ErrField struct {
	Field string
	Err   error
}

// Error returns the whole path to the field and the error message.
func (e *ErrField) Error() string {
	b := &strings.Builder{}
	b.WriteString(e.Field)

	err := e.Err
	for {
		switch et := err.(type) {
		case *ErrField:
			b.WriteString(".")
			b.WriteString(et.Field)
		case *ErrIndex:
			fmt.Fprintf(b, "[%d]", et.Index)
		case *ErrKey:
			fmt.Fprintf(b, "[%q]", et.Key)
		case *ErrRequired:
			b.WriteString(".")
			b.WriteString(et.Target)
			b.WriteString(" is required")
			return b.String()
		default:
			fmt.Fprintf(b, ": %v", err)
			return b.String()
		}

		err = errors.Unwrap(err)
	}
}

// Unwrap returns the wrapped error.
func (e *ErrField) Unwrap() error { return e.Err }

// ErrIndex is an error that occurred in a slice. It contains the index of the element.
type ErrIndex struct {
	Index int
	Err   error
}

// Error returns the index and the error message.
// However, most of the time, it is wrapped in an ErrField and the index is shown as part of a path.
func (e *ErrIndex) Error() string { return fmt.Sprintf("%d: %v", e.Index, e.Err) }

// Unwrap returns the wrapped error.
func (e *ErrIndex) Unwrap() error { return e.Err }

// ErrKey is an error that occurred in a map. It contains the key of the element.
type ErrKey struct {
	Key string
	Err error
}

// Error returns the key and the error message.
// However, most of the time, it is wrapped in an ErrField and the key is shown as part of a path.
func (e *ErrKey) Error() string { return fmt.Sprintf("%q: %v", e.Key, e.Err) }

// Unwrap returns the wrapped error.
func (e *ErrKey) Unwrap() error { return e.Err }
