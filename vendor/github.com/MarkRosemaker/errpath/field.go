package errpath

import "strings"

var _ pathWriter = (*ErrField)(nil)

// ErrField is an error that occurred in a field.
type ErrField struct {
	// The name of the field where the error occurred.
	Field string
	// The underlying error.
	Err error
}

// Error returns the whole path to the field and the error message.
func (e *ErrField) Error() string {
	b := &strings.Builder{}
	b.WriteString(e.Field)
	return writePath(b, e.Err)
}

// writePath appends the field to the path by separating it with a dot from the previous path
func (e *ErrField) writePath(b *strings.Builder) {
	b.WriteString(".")
	b.WriteString(e.Field)
}

// Unwrap returns the wrapped error.
func (e *ErrField) Unwrap() error { return e.Err }
