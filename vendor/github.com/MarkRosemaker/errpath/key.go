package errpath

import (
	"fmt"
	"strings"
)

var _ pathWriter = (*ErrKey)(nil)

// ErrKey is an error that occurred in a map. It contains the key of the element.
type ErrKey struct {
	// The key of the map where the error occurred.
	Key string
	// The underlying error.
	Err error
}

// Error returns the key and the error message.
// However, it makes sense to wrap `ErrKey` in an `ErrField` so this method is rarely called.
func (e *ErrKey) Error() string {
	b := &strings.Builder{}
	e.writePath(b)
	return writePath(b, e.Err)
}

// writePath appends the key to the path by writing it in brackets and in quotes after the previous path
func (e *ErrKey) writePath(b *strings.Builder) {
	fmt.Fprintf(b, "[%q]", e.Key)
}

// Unwrap returns the wrapped error.
func (e *ErrKey) Unwrap() error { return e.Err }
