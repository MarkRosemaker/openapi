package errpath

import "strings"

var _ pathWriter = (*ErrRequired)(nil)

// ErrRequired signals that a required value is missing.
type ErrRequired struct{}

// Error fulfills the error interface.
// Without a previous error path, it simply says "a value is required".
// Naturally, this is not very helpful, so it makes sense to wrap `ErrRequired` in another error such as `ErrField` so the user knows that the field was required.
func (e *ErrRequired) Error() string {
	b := &strings.Builder{}
	b.WriteString("a value")
	e.writePath(b)
	return b.String()
}

// writePath appends `" is required"` to the path so that the user knows the given path was required
func (e *ErrRequired) writePath(b *strings.Builder) {
	b.WriteString(" is required")
}
