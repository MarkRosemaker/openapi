package openapi

import (
	"errors"
	"fmt"
	"strings"
)

// A chainContinuer allows us to continue a chain of errors.
// Instead of "info: name is required", we can have "info.name is required".
type chainContinuer interface{ continueChain(*strings.Builder) }

var (
	_ chainContinuer = (*ErrRequired)(nil)
	_ chainContinuer = (*ErrInvalid[string])(nil)
	_ chainContinuer = (*ErrField)(nil)
	_ chainContinuer = (*ErrIndex)(nil)
	_ chainContinuer = (*ErrKey)(nil)
)

// ErrRequired signals that a required value is missing.
type ErrRequired struct{}

// Error fulfills the error interface.
// Without an error chain, it simply says "a value is required".
func (e *ErrRequired) Error() string {
	b := &strings.Builder{}
	b.WriteString("a value")
	e.continueChain(b)
	return b.String()
}

func (e *ErrRequired) continueChain(b *strings.Builder) {
	b.WriteString(" is required")
}

// ErrInvalid signals that a value is invalid.
type ErrInvalid[T comparable] struct {
	// The value that is invalid.
	Value T
	// An optional list of valid values.
	Enum []T
	// An optional message that explains the error.
	Message string
}

// Error returns helpful information about the invalid field and how to fix it.
func (e *ErrInvalid[_]) Error() string {
	b := &strings.Builder{}
	b.WriteString("a value")
	e.continueChain(b)
	return b.String()
}

func (e *ErrInvalid[T]) continueChain(b *strings.Builder) {
	switch val := any(e.Value).(type) {
	case bool:
		fmt.Fprintf(b, " (%t)", val)
	default:
		var zero T
		if e.Value != zero {
			fmt.Fprintf(b, " (%#v)", e.Value)
		}
	}

	b.WriteString(" is invalid")
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
		return
	}

	if len(e.Enum) > 0 {
		b.WriteString(", must be one of: ")
		enums := make([]string, len(e.Enum))
		for i, v := range e.Enum {
			enums[i] = fmt.Sprintf("%#v", v)
		}

		b.WriteString(strings.Join(enums, ", "))
	}
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
	return writeOutChain(b, e.Err)
}

func (e *ErrField) continueChain(b *strings.Builder) {
	b.WriteString(".")
	b.WriteString(e.Field)
}

func writeOutChain(b *strings.Builder, err error) string {
	for {
		// check if the error is a chainContinuer
		if cont, ok := err.(chainContinuer); ok {
			cont.continueChain(b)
			err = errors.Unwrap(err)
			continue
		}

		if err != nil {
			fmt.Fprintf(b, ": %v", err)
		}

		return b.String() // end of the chain
	}
}

// Unwrap returns the wrapped error.
func (e *ErrField) Unwrap() error { return e.Err }

// ErrIndex is an error that occurred in a slice. It contains the index of the element.
type ErrIndex struct {
	Index int
	Err   error
}

func (e *ErrIndex) continueChain(b *strings.Builder) {
	fmt.Fprintf(b, "[%d]", e.Index)
}

// Error returns the index and the error message.
// However, most of the time, it is wrapped in an ErrField and the index is shown as part of a path.
func (e *ErrIndex) Error() string {
	b := &strings.Builder{}
	e.continueChain(b)
	return writeOutChain(b, e.Err)
}

// Unwrap returns the wrapped error.
func (e *ErrIndex) Unwrap() error { return e.Err }

// ErrKey is an error that occurred in a map. It contains the key of the element.
type ErrKey struct {
	Key string
	Err error
}

func (e *ErrKey) continueChain(b *strings.Builder) {
	fmt.Fprintf(b, "[%q]", e.Key)
}

// Error returns the key and the error message.
// However, most of the time, it is wrapped in an ErrField and the key is shown as part of a path.
func (e *ErrKey) Error() string {
	b := &strings.Builder{}
	e.continueChain(b)
	return writeOutChain(b, e.Err)
}

// Unwrap returns the wrapped error.
func (e *ErrKey) Unwrap() error { return e.Err }
