package openapi

import (
	"fmt"
	"regexp"
	"strings"
)

// The Schema Object allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays. This object is a superset of the JSON Schema Specification Draft 2020-12.
//
// For more information about the properties, see JSON Schema Core and JSON Schema Validation.
//
// Unless stated otherwise, the property definitions follow those of JSON Schema and do not add any additional semantics.
// Where JSON Schema indicates that behavior is defined by the application (e.g. for annotations), OAS also defers the definition of semantics to the application consuming the OpenAPI document.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#schema-object
type Schema struct {
	// The name of the schema.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short description of the schema.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Specifies the data type of the property.
	Type DataType `json:"type,omitempty" yaml:"type,omitempty"`
	// Further refines the data type.
	Format Format `json:"format,omitempty" yaml:"format,omitempty"`

	// For object types, defines the properties of the object
	Properties Schemas `json:"properties,omitzero" yaml:"properties,omitzero"`
	// Which properties are required.
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`

	// String

	// The pattern is used to validate the string.
	Pattern *regexp.Regexp `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Number
	Min *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Max *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`

	// Array

	// The minimum number of items in the array.
	MinItems uint `json:"minItems,omitzero" yaml:"minItems,omitzero"`
	// The maximum number of items in the array.
	MaxItems *uint `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	// The items of the array. When the type is array, this property is REQUIRED.
	// The empty schema for `items` indicates a media type of `application/octet-stream`.
	Items *SchemaRef `json:"items,omitzero" yaml:"items,omitzero"`

	// Object

	AdditionalProperties *SchemaRef `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

	// special encoding for binary data
	ContentMediaType string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`
	ContentEncoding  string `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`

	Example any `json:"example,omitempty" yaml:"example,omitempty"`
}

func (s *Schema) Validate() error {
	s.Description = strings.TrimSpace(s.Description)

	if s.Type == "" {
		return &ErrField{Field: "type", Err: &ErrRequired{}}
	}

	if err := s.Type.Validate(); err != nil {
		return &ErrField{Field: "type", Err: err}
	}

	switch s.Type {
	case TypeArray:
		if s.Items == nil {
			return &ErrField{Field: "items", Err: &ErrRequired{}}
		}
	}

	if s.Format != "" {
		if err := s.Format.Validate(); err != nil {
			return &ErrField{Field: "format", Err: err}
		}
	}

	// validate if format is valid for type
	switch s.Format {
	case "": // no format
	case FormatInt32, FormatInt64:
		if s.Type != TypeInteger {
			return &ErrField{Field: "format", Err: &ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for integer type, got %s", s.Type),
			}}
		}
	case FormatFloat, FormatDouble:
		if s.Type != TypeNumber {
			return &ErrField{Field: "format", Err: &ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
			}}
		}
	case FormatPassword:
		if s.Type != TypeString {
			return &ErrField{Field: "format", Err: &ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for string type, got %s", s.Type),
			}}
		}
	}

	// validate min and max
	if s.Min != nil || s.Max != nil {
		switch s.Type {
		case TypeInteger:
			if s.Min != nil {
				if *s.Min != float64(int(*s.Min)) {
					return &ErrField{Field: "minimum", Err: &ErrInvalid[float64]{
						Value:   *s.Min,
						Message: "not an integer",
					}}
				}
			}

			if s.Max != nil {
				if *s.Max != float64(int(*s.Max)) {
					return &ErrField{Field: "maximum", Err: &ErrInvalid[float64]{
						Value:   *s.Max,
						Message: "not an integer",
					}}
				}
			}

			fallthrough
		case TypeNumber:
			if s.Min != nil && s.Max != nil && *s.Min > *s.Max {
				return &ErrField{Field: "minimum", Err: &ErrInvalid[float64]{
					Value:   *s.Min,
					Message: fmt.Sprintf("minimum is greater than maximum (%v > %v)", *s.Min, *s.Max),
				}}
			}
		default:
			if s.Min != nil {
				return &ErrField{Field: "minimum", Err: &ErrInvalid[float64]{
					Value:   *s.Min,
					Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
				}}
			}

			return &ErrField{Field: "maximum", Err: &ErrInvalid[float64]{
				Value:   *s.Max,
				Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
			}}
		}
	}

	// validate min and max items
	if s.MinItems != 0 || s.MaxItems != nil {
		// must be array
		if s.Type != TypeArray {
			if s.MinItems != 0 {
				return &ErrField{Field: "minItems", Err: &ErrInvalid[uint]{
					Value:   s.MinItems,
					Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
				}}
			}

			return &ErrField{Field: "maxItems", Err: &ErrInvalid[uint]{
				Value:   *s.MaxItems,
				Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
			}}
		}

		if s.MaxItems != nil && s.MinItems > *s.MaxItems {
			return &ErrField{Field: "minItems", Err: &ErrInvalid[uint]{
				Value:   s.MinItems,
				Message: fmt.Sprintf("minItems is greater than maxItems (%d > %d)", s.MinItems, *s.MaxItems),
			}}
		}
	}

	if s.Items != nil {
		if s.Type != TypeArray {
			return &ErrField{Field: "items", Err: &ErrInvalid[string]{
				Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
			}}
		}

		if err := s.Items.Validate(); err != nil {
			return &ErrField{Field: "items", Err: err}
		}
	}

	return nil
}
