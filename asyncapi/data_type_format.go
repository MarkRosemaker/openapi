package asyncapi

import "slices"

// Format defines additional formats to provide fine detail for primitive data types.
//
// The format property is an open string-valued property, and can have any value to support
// documentation needs, so an unknown format is not an error.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#dataTypeFormat
type Format string

const (
	// FormatInt32 represents a signed 32 bits integer.
	FormatInt32 Format = "int32"
	// FormatInt64 represents a signed 64 bits integer.
	FormatInt64 Format = "int64"
	// FormatFloat represents a float number.
	FormatFloat Format = "float"
	// FormatDouble represents a double number.
	FormatDouble Format = "double"
	// FormatByte represents base64 encoded characters.
	FormatByte Format = "byte"
	// FormatBinary represents any sequence of octets.
	FormatBinary Format = "binary"
	// FormatDate represents a date as defined by full-date in RFC3339.
	FormatDate Format = "date"
	// FormatDateTime represents a date-time as defined by date-time in RFC3339.
	FormatDateTime Format = "date-time"
	// FormatPassword is a hint to UIs that the input needs to be obscured.
	FormatPassword Format = "password"
)

// allFormats are the formats defined by the AsyncAPI Specification.
var allFormats = []Format{
	FormatInt32, FormatInt64,
	FormatFloat, FormatDouble,
	FormatByte, FormatBinary,
	FormatDate, FormatDateTime,
	FormatPassword,
}

// IsKnown reports whether the format is one of the formats defined by the AsyncAPI Specification.
//
// Formats such as "email" or "uuid" can be used even though they are not defined by the
// specification, so a format that is not known is not necessarily invalid.
func (f Format) IsKnown() bool { return slices.Contains(allFormats, f) }
