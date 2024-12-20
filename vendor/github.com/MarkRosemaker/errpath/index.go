package errpath

import (
	"fmt"
	"strings"
)

var _ pathWriter = (*ErrIndex)(nil)

// ErrIndex is an error that occurred in a slice. It contains the index of the element.
type ErrIndex struct {
	// The index of the slice where the error occurred.
	Index int
	// The underlying error.
	Err error
}

// Error returns the index and the error message.
// However, it makes sense to wrap `ErrIndex` in an `ErrField` so this method is rarely called.
func (e *ErrIndex) Error() string {
	b := &strings.Builder{}
	e.writePath(b)
	return writePath(b, e.Err)
}

// writePath appends the index to the path by writing it in brackets after the previous path
func (e *ErrIndex) writePath(b *strings.Builder) {
	fmt.Fprintf(b, "[%d]", e.Index)
}

// Unwrap returns the wrapped error.
func (e *ErrIndex) Unwrap() error { return e.Err }
