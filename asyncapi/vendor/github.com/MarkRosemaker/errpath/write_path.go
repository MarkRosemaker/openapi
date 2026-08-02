package errpath

import (
	"errors"
	"fmt"
	"strings"
)

// A pathWriter allows us to continue a chain of errors.
// Instead of "info: name is required", we can have "info.name is required".
type pathWriter interface{ writePath(*strings.Builder) }

// a joined error is the result of calling `errors.Join()` or similar error joiner
type joinedErr interface{ Unwrap() []error }

func writePath(b *strings.Builder, err error) string {
	// check if the error is a pathWriter
	if w, ok := err.(pathWriter); ok {
		// write the element of the path concerning this error
		w.writePath(b)
		// continue writing the path of the underlying error
		return writePath(b, errors.Unwrap(err))
	}

	// check if we have a collection of errors
	if joined, ok := err.(joinedErr); ok {
		// multiple errors, print them as separate lines
		// but with the same prefix up until here
		prefix := b.String()
		for i, e := range joined.Unwrap() {
			if i != 0 {
				b.WriteByte('\n')
				b.WriteString(prefix)
			}

			writePath(b, e)
		}

		return b.String() // end of the path
	}

	// we end the path and write the error that we don't know yet
	if err != nil {
		fmt.Fprintf(b, ": %v", err)
	}

	return b.String() // end of the path
}
