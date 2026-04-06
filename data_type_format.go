package openapi

// Format defines additional formats to provide fine detail for primitive data types.
// The format property is an open string-valued property and can have any value.
// Formats such as "email", "uuid", and so on, MAY be used even though they are not
// defined by this specification. Tools that do not recognize a specific format MUST
// treat it like a plain schema object and MUST NOT generate an error.
// See: https://spec.openapis.org/oas/v3.2.0.html#data-types
type Format string

const (
	// FormatInt32 represents a signed 32 bits integer.
	FormatInt32 Format = "int32"
	// FormatInt64 represents a signed 64 bits integer.
	FormatInt64 Format = "int64"
	// FormatUint represents an unsigned integer.
	FormatUint Format = "uint"
	// FormatUint32 represents an unsigned 32 bits integer.
	FormatUint32 Format = "uint32"
	// FormatUint64 represents an unsigned 64 bits integer.
	FormatUint64 Format = "uint64"
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
	// FormatIPv4 represents an IPv4 address.
	FormatIPv4 Format = "ipv4"
	// FormatIPv6 represents an IPv6 address.
	FormatIPv6 Format = "ipv6"
)
