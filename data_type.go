package openapi

import "slices"

// Data types in the OAS are based on the types supported by the JSON Schema Specification Draft 2020-12.
// Note that `integer` as a type is also supported and is defined as a JSON number without a fraction or exponent part.
// Models are defined using the Schema Object, which is a superset of JSON Schema Specification Draft 2020-12.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#data-types
type DataType string

const (
	TypeInteger DataType = "integer" // format: int32, int64
	TypeNumber  DataType = "number"  // format: float, double
	TypeString  DataType = "string"  // format: password
	TypeArray   DataType = "array"
	TypeBoolean DataType = "boolean"
	TypeObject  DataType = "object"
)

var allDataTypes = []DataType{
	TypeInteger, TypeNumber, TypeString, TypeArray, TypeBoolean, TypeObject,
}

func (d DataType) Validate() error {
	if slices.Contains(allDataTypes, d) {
		return nil
	}

	return &ErrInvalid[DataType]{
		Value: d,
		Enum:  allDataTypes,
	}
}
