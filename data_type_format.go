package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

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
	// FormatDuration represents a duration.
	FormatDuration Format = "duration"
	// FormatEmail represents an email.
	FormatEmail Format = "email"
	// FormatPassword represents a password. It's a hint to UIs to obscure input.
	FormatPassword Format = "password"
	// FormatUUID represents a UUID.
	FormatUUID Format = "uuid"
	// FormatURI represents a URI.
	FormatURI Format = "uri"
	// FormatURIRef represents a URI reference.
	FormatURIRef Format = "uriref"
	// FormatZipCode represents a zip code.
	FormatZipCode Format = "zip-code"
)

var allFormats = []Format{
	FormatInt32, FormatInt64, FormatFloat, FormatDouble,
	FormatByte, FormatBinary,
	FormatDate, FormatDateTime, FormatDuration,
	FormatEmail, FormatPassword,
	FormatUUID,
	FormatURI, FormatURIRef, FormatZipCode,
}

// Validate validates the format.
func (f Format) Validate() error {
	if slices.Contains(allFormats, f) {
		return nil
	}

	return &errpath.ErrInvalid[Format]{
		Value: f,
		Enum:  allFormats,
	}
}
