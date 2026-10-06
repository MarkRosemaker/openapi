package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// DataType is a type of the data a schema describes. Data types in the OAS are based on the types supported by the JSON Schema Specification Draft 2020-12.
// Note that `integer` as a type is also supported and is defined as a JSON number without a fraction or exponent part.
// Models are defined using the Schema Object, which is a superset of JSON Schema Specification Draft 2020-12.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.2.0.html#data-types
type DataType string

// The data types a schema can declare.
const (
	TypeInteger DataType = "integer" // format: int32, int64
	TypeNumber  DataType = "number"  // format: float, double
	TypeString  DataType = "string"  // format: password
	TypeArray   DataType = "array"
	TypeBoolean DataType = "boolean"
	TypeObject  DataType = "object"
	// TypeNull represents a null value.
	// See: https://spec.openapis.org/oas/v3.2.0.html#data-types
	TypeNull DataType = "null"
)

var allDataTypes = []DataType{
	TypeInteger, TypeNumber, TypeString, TypeArray, TypeBoolean, TypeObject, TypeNull,
}

// Validate returns an error if the DataType breaks the specification.
func (d DataType) Validate() error {
	if slices.Contains(allDataTypes, d) {
		return nil
	}

	return &errpath.ErrInvalid[DataType]{
		Value: d,
		Enum:  allDataTypes,
	}
}
