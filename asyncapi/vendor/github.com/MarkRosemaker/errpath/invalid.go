package errpath

import (
	"fmt"
	"reflect"
	"strings"
)

var _ pathWriter = (*ErrInvalid[string])(nil)

// ErrInvalid signals that a value is invalid.
type ErrInvalid[T any] struct {
	// The value that is invalid.
	Value T
	// An optional list of valid values.
	Enum []T
	// An optional message that explains the error.
	Message string
}

// Error returns helpful information about the invalid field and how to fix it.
// Without a previous error path, it simply calls it "a value".
// It makes sense to wrap `ErrInvalid` in another error such as `ErrField` so the user knows that the field is invalid.
func (e *ErrInvalid[_]) Error() string {
	b := &strings.Builder{}
	b.WriteString("a value")
	e.writePath(b)
	return b.String()
}

// writePath appends information to the path so that the user knows why the value is invalid
func (e *ErrInvalid[T]) writePath(b *strings.Builder) {
	if s := stringify(e.Value); s != "" {
		fmt.Fprintf(b, " (%s)", s)
	}

	b.WriteString(" is invalid")
	if e.Message != "" {
		b.WriteString(": ")
		b.WriteString(e.Message)
	}

	if len(e.Enum) == 0 {
		return
	}

	b.WriteString(", must be one of: ")
	enums := make([]string, len(e.Enum))
	for i, v := range e.Enum {
		enums[i] = stringify(v)
	}

	b.WriteString(strings.Join(enums, ", "))
}

func stringify(val any) string {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128:
		// write "false" and 0 even though it is the zero value
		return fmt.Sprint(v.Interface())
	default:
		// do not stringify zero values
		if v.IsZero() {
			return ""
		}

		return fmt.Sprintf("%#v", v.Interface())
	}
}
