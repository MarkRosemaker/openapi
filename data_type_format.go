package openapi

import "slices"

// Format defines additional formats to provide fine detail for primitive data types.
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
	// FormatByte represents a byte.
	FormatByte Format = "byte"
	// FormatBinary represents a binary.
	FormatBinary Format = "binary"
	// FormatDate represents a date.
	FormatDate Format = "date"
	// FormatDateTime represents a date-time.
	FormatDateTime Format = "date-time"
	// FormatPassword represents a password. It's a hint to UIs to obscure input.
	FormatPassword Format = "password"
	// FormatDuration represents a duration.
	FormatDuration Format = "duration"
	// FormatUUID represents a UUID.
	FormatUUID Format = "uuid"
	// FormatEmail represents an email.
	FormatEmail Format = "email"
	// FormatURI represents a URI.
	FormatURI Format = "uri"
	// FormatZipCode represents a zip code.
	FormatZipCode Format = "zip-code"
)

var allFormats = []Format{
	FormatInt32, FormatInt64, FormatFloat, FormatDouble, FormatByte, FormatBinary, FormatDate, FormatDateTime, FormatPassword,
	FormatDuration, FormatUUID, FormatEmail, FormatURI, FormatZipCode,
}

// Validate validates the format.
func (f Format) Validate() error {
	if slices.Contains(allFormats, f) {
		return nil
	}

	return &ErrInvalid[Format]{
		Value: f,
		Enum:  allFormats,
	}
}
