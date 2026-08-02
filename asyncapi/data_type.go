package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// DataType is the type of a schema, based on the types supported by the JSON Schema Specification Draft 07.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#dataTypeFormat
type DataType string

const (
	// TypeInteger is a JSON number without a fraction or exponent part. Format: int32, int64
	TypeInteger DataType = "integer"
	// TypeNumber is a JSON number. Format: float, double
	TypeNumber DataType = "number"
	// TypeString is a JSON string. Format: byte, binary, date, date-time, password
	TypeString DataType = "string"
	// TypeArray is a JSON array.
	TypeArray DataType = "array"
	// TypeBoolean is a JSON boolean.
	TypeBoolean DataType = "boolean"
	// TypeObject is a JSON object.
	TypeObject DataType = "object"
	// TypeNull is the JSON null value.
	TypeNull DataType = "null"
)

var allDataTypes = []DataType{
	TypeInteger, TypeNumber, TypeString, TypeArray, TypeBoolean, TypeObject, TypeNull,
}

// Validate validates the data type.
func (d DataType) Validate() error {
	if slices.Contains(allDataTypes, d) {
		return nil
	}

	return &errpath.ErrInvalid[DataType]{
		Value: d,
		Enum:  allDataTypes,
	}
}

// DataTypes is the value of the `type` keyword of a schema.
//
// JSON Schema allows a single type as well as a list of types,
// e.g. `"type": "string"` and `"type": ["string", "null"]` are both valid.
// A single type is marshalled back as a single type, not as a list.
type DataTypes []DataType

// Validate validates each data type.
func (ds DataTypes) Validate() error {
	for i, d := range ds {
		if err := d.Validate(); err != nil {
			if len(ds) == 1 { // don't report an index if there is only one type
				return err
			}

			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

// Contains reports whether the schema is of the given type.
func (ds DataTypes) Contains(d DataType) bool { return slices.Contains(ds, d) }

// String returns the types as a comma-separated list.
func (ds DataTypes) String() string {
	switch len(ds) {
	case 0:
		return ""
	case 1:
		return string(ds[0])
	}

	s := string(ds[0])
	for _, d := range ds[1:] {
		s += ", " + string(d)
	}

	return s
}

var _ json.UnmarshalerFrom = (*DataTypes)(nil)

// UnmarshalJSONFrom unmarshals either a single type or a list of types.
func (ds *DataTypes) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == '[' {
		var types []DataType
		if err := json.UnmarshalDecode(dec, &types); err != nil {
			return err
		}

		*ds = types

		return nil
	}

	var d DataType
	if err := json.UnmarshalDecode(dec, &d); err != nil {
		return err
	}

	*ds = DataTypes{d}

	return nil
}

var _ json.MarshalerTo = (*DataTypes)(nil)

// MarshalJSONTo marshals a single type as a string and multiple types as a list.
func (ds *DataTypes) MarshalJSONTo(enc *jsontext.Encoder) error {
	if len(*ds) == 1 {
		return json.MarshalEncode(enc, (*ds)[0])
	}

	return json.MarshalEncode(enc, []DataType(*ds))
}
