package openapi

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MarkRosemaker/errpath"
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
// [Specification]: https://spec.openapis.org/oas/v3.2.0.html#schema-object
type Schema struct {
	// The name of the schema.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short description of the schema.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Specifies the data type of the property.
	Type DataType `json:"type,omitempty" yaml:"type,omitempty"`
	// Whether null is a valid value too, written as a "type" of [Type, "null"].
	Nullable bool `json:"-" yaml:"-"`
	// Further refines the data type.
	Format Format `json:"format,omitempty" yaml:"format,omitempty"`

	// AllOf validates the value against ALL of the given schemas.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	AllOf SchemaList `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	// OneOf validates the value against EXACTLY ONE of the given schemas.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	OneOf SchemaList `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	// AnyOf validates the value against AT LEAST ONE of the given schemas.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	AnyOf SchemaList `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// Not validates the value against the negation of the given schema — the value must NOT validate against it.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	Not *Schema `json:"not,omitzero" yaml:"not,omitempty"`

	// Integer / Number

	// The minimum value of the number.
	Min *float64 `json:"minimum,omitzero" yaml:"minimum,omitempty"`
	// The maximum value of the number.
	Max *float64 `json:"maximum,omitzero" yaml:"maximum,omitempty"`
	// A value the number must be greater than.
	ExclusiveMin *float64 `json:"exclusiveMinimum,omitzero" yaml:"exclusiveMinimum,omitempty"`
	// A value the number must be less than.
	ExclusiveMax *float64 `json:"exclusiveMaximum,omitzero" yaml:"exclusiveMaximum,omitempty"`
	// A number the value must be a multiple of.
	MultipleOf *float64 `json:"multipleOf,omitzero" yaml:"multipleOf,omitempty"`

	// String

	// The minimum length of the string, in characters.
	MinLength uint `json:"minLength,omitzero" yaml:"minLength,omitempty"`
	// The maximum length of the string, in characters.
	MaxLength *uint `json:"maxLength,omitzero" yaml:"maxLength,omitempty"`
	// An ECMA-262 regular expression the string must match, compiled with Go's regexp; see pattern.go.
	Pattern *regexp.Regexp `json:"pattern,omitzero" yaml:"pattern,omitempty"`
	// A list of possible values. Per JSON Schema 2020-12, enum may contain any JSON type.
	Enum []jsontext.Value `json:"enum,omitzero" yaml:"enum,omitempty"`
	// The one value allowed, of any JSON type.
	Const jsontext.Value `json:"const,omitzero" yaml:"const,omitempty"`

	// Array

	// The minimum number of items in the array.
	MinItems uint `json:"minItems,omitzero" yaml:"minItems,omitempty"`
	// The maximum number of items in the array.
	MaxItems *uint `json:"maxItems,omitzero" yaml:"maxItems,omitempty"`
	// Whether the items of the array must all be different.
	UniqueItems bool `json:"uniqueItems,omitzero" yaml:"uniqueItems,omitempty"`
	// PrefixItems validates the array positionally: the first element
	// against the first schema here, the second against the second, and so
	// on. Items still applies to any element beyond the ones listed here.
	// See JSON Schema 2020-12, "prefixItems".
	PrefixItems SchemaList `json:"prefixItems,omitempty" yaml:"prefixItems,omitempty"`
	// The items of the array; without it, any items are allowed (JSON Schema 2020-12).
	// The empty schema for `items` indicates a media type of `application/octet-stream`.
	Items *Schema `json:"items,omitzero" yaml:"items,omitempty"`

	// Object

	// For object types, defines the properties of the object
	Properties Schemas `json:"properties,omitempty" yaml:"properties,omitempty"`
	// Which properties are required.
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`
	// Applies to properties not listed in Properties: a schema for their values, or whether they are allowed at all.
	AdditionalProperties *AdditionalProperties `json:"additionalProperties,omitzero" yaml:"additionalProperties,omitempty"`
	// The maximum number of properties of the object.
	MaxProperties *uint `json:"maxProperties,omitzero" yaml:"maxProperties,omitempty"`
	// A schema every property name of the object must match.
	PropertyNames *Schema `json:"propertyNames,omitzero" yaml:"propertyNames,omitempty"`
	// Tells which of the composed schemas a payload is, by the value of one of its properties.
	Discriminator *Discriminator `json:"discriminator,omitzero" yaml:"discriminator,omitempty"`

	// special encoding for binary data
	ContentMediaType string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`
	ContentEncoding  string `json:"contentEncoding,omitempty"  yaml:"contentEncoding,omitempty"`

	// Specifies the default value of the property if no value is provided.
	Default jsontext.Value `json:"default,omitzero" yaml:"default,omitempty"`

	Example jsontext.Value `json:"example,omitzero" yaml:"example,omitzero"`
	// Example values the schema accepts, of any JSON type.
	Examples []jsontext.Value `json:"examples,omitempty" yaml:"examples,omitempty"`
	// Whether the schema should no longer be used.
	Deprecated bool `json:"deprecated,omitzero" yaml:"deprecated,omitempty"`
	// Whether the value is only ever sent in responses, such as an ID the server assigns, and left out of requests.
	ReadOnly bool `json:"readOnly,omitzero" yaml:"readOnly,omitempty"`
	// Whether the value is only ever sent in requests, such as a password, and left out of responses.
	WriteOnly bool `json:"writeOnly,omitzero" yaml:"writeOnly,omitempty"`

	// Another schema this one also applies, written last as "$ref"; see schema_json.go.
	Ref *SchemaRef `json:"-" yaml:"-"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",embed" yaml:"-"`

	// an index to the original location of this object
	idx int
}

// SchemaRef is a schema's "$ref": the location of another schema, and that schema once the document is loaded.
type SchemaRef struct {
	// The location of the schema, e.g. "#/components/schemas/Pet".
	Identifier string
	// The schema Identifier points to, set when the document is loaded.
	Value *Schema
}

// Replace makes s a copy of v while keeping s's place in the ordered map that holds it.
// Assigning *s = *v instead would also copy v's place, reordering the properties or components around s.
func (s *Schema) Replace(v *Schema) {
	idx := s.idx
	*s = *v
	s.idx = idx
}

func getIndexSchema(s *Schema) int              { return s.idx }
func setIndexSchema(s *Schema, idx int) *Schema { s.idx = idx; return s }

// Validate checks the schema for correctness.
func (s *Schema) Validate() error { return s.validate(false) }

// validate checks the schema; extended says another component schema builds on it through allOf.
func (s *Schema) validate(extended bool) error {
	s.Description = strings.TrimSpace(s.Description)

	if s.Ref != nil && s.Ref.Value == nil {
		return &errpath.ErrField{Field: "$ref", Err: fmt.Errorf("%q was not resolved", s.Ref.Identifier)}
	}

	// type is optional (JSON Schema 2020-12): without one, every type's keywords apply, each to instances of its type.
	if s.Type != "" {
		if err := s.Type.Validate(); err != nil {
			return &errpath.ErrField{Field: "type", Err: err}
		}
	}

	if s.Format != "" {
		if err := s.Format.Validate(); err != nil {
			return &errpath.ErrField{Field: "format", Err: err}
		}
	}

	// validate if format is valid for type
	switch s.Format {
	case "": // no format
	case FormatInt32, FormatInt64, FormatUint, FormatUint32, FormatUint64:
		if !s.allows(TypeInteger) {
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for integer type, got %s", s.typeOrNone()),
			}}
		}
	case FormatFloat, FormatDouble:
		if !s.allows(TypeNumber) {
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for number type, got %s", s.typeOrNone()),
			}}
		}
	case FormatEmail, FormatPassword,
		FormatUUID, FormatURI, FormatURIRef, FormatZipCode,
		FormatIPv4, FormatIPv6:
		if !s.allows(TypeString) {
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for string type, got %s", s.typeOrNone()),
			}}
		}
	case FormatDuration, FormatDate, FormatDateTime:
		switch s.Type {
		case "", TypeInteger, TypeString:
		default:
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for integer or string type, got %s", s.typeOrNone()),
			}}
		}
	case FormatByte, FormatBinary:
		switch s.Type {
		case "", TypeString:
		default:
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for string type, got %s", s.typeOrNone()),
			}}
		}
	default:
		return fmt.Errorf("unimplemented format: %s", s.Format)
	}

	// String

	if !s.allows(TypeString) {
		for _, kw := range []struct {
			field string
			set   bool
		}{
			{"minLength", s.MinLength != 0},
			{"maxLength", s.MaxLength != nil},
			{"pattern", s.Pattern != nil},
			{"contentMediaType", s.ContentMediaType != ""},
			{"contentEncoding", s.ContentEncoding != ""},
		} {
			if kw.set {
				return &errpath.ErrField{Field: kw.field, Err: &errpath.ErrInvalid[string]{
					Message: fmt.Sprintf("only valid for string type, got %s", s.typeOrNone()),
				}}
			}
		}
	}

	if s.MaxLength != nil && s.MinLength > *s.MaxLength {
		return &errpath.ErrField{Field: "minLength", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinLength,
			Message: fmt.Sprintf("minLength is greater than maxLength (%d > %d)", s.MinLength, *s.MaxLength),
		}}
	}

	for i, v := range s.AllOf {
		if err := v.Validate(); err != nil {
			return &errpath.ErrField{
				Field: "allOf",
				Err:   &errpath.ErrIndex{Index: i, Err: err},
			}
		}
	}

	for i, v := range s.OneOf {
		if err := v.Validate(); err != nil {
			return &errpath.ErrField{
				Field: "oneOf",
				Err:   &errpath.ErrIndex{Index: i, Err: err},
			}
		}
	}

	for i, v := range s.AnyOf {
		if err := v.Validate(); err != nil {
			return &errpath.ErrField{
				Field: "anyOf",
				Err:   &errpath.ErrIndex{Index: i, Err: err},
			}
		}
	}

	if s.Not != nil {
		if err := s.Not.Validate(); err != nil {
			return &errpath.ErrField{Field: "not", Err: err}
		}
	}

	// Integer / Number

	// validate min and max
	if s.Type == TypeInteger {
		if s.Min != nil && *s.Min != float64(int(*s.Min)) {
			return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.Min,
				Message: "not an integer",
			}}
		}

		if s.Max != nil && *s.Max != float64(int(*s.Max)) {
			return &errpath.ErrField{Field: "maximum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.Max,
				Message: "not an integer",
			}}
		}

		if s.ExclusiveMin != nil && *s.ExclusiveMin != float64(int(*s.ExclusiveMin)) {
			return &errpath.ErrField{Field: "exclusiveMinimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.ExclusiveMin,
				Message: "not an integer",
			}}
		}

		if s.ExclusiveMax != nil && *s.ExclusiveMax != float64(int(*s.ExclusiveMax)) {
			return &errpath.ErrField{Field: "exclusiveMaximum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.ExclusiveMax,
				Message: "not an integer",
			}}
		}

		if s.MultipleOf != nil && *s.MultipleOf != float64(int(*s.MultipleOf)) {
			return &errpath.ErrField{Field: "multipleOf", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.MultipleOf,
				Message: "not an integer",
			}}
		}
	}

	if s.allows(TypeNumber) || s.allows(TypeInteger) {
		if s.Min != nil && s.Max != nil && *s.Min > *s.Max {
			return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.Min,
				Message: fmt.Sprintf("minimum is greater than maximum (%v > %v)", *s.Min, *s.Max),
			}}
		}

		if s.Min != nil && s.ExclusiveMax != nil && *s.Min >= *s.ExclusiveMax {
			return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.Min,
				Message: fmt.Sprintf("minimum is not less than exclusiveMaximum (%v >= %v)", *s.Min, *s.ExclusiveMax),
			}}
		}

		if s.ExclusiveMin != nil && s.Max != nil && *s.ExclusiveMin >= *s.Max {
			return &errpath.ErrField{Field: "exclusiveMinimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.ExclusiveMin,
				Message: fmt.Sprintf("exclusiveMinimum is not less than maximum (%v >= %v)", *s.ExclusiveMin, *s.Max),
			}}
		}

		if s.ExclusiveMin != nil && s.ExclusiveMax != nil && *s.ExclusiveMin >= *s.ExclusiveMax {
			return &errpath.ErrField{Field: "exclusiveMinimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.ExclusiveMin,
				Message: fmt.Sprintf("exclusiveMinimum is not less than exclusiveMaximum (%v >= %v)", *s.ExclusiveMin, *s.ExclusiveMax),
			}}
		}
	} else if s.Min != nil {
		return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.Min,
			Message: fmt.Sprintf("only valid for number type, got %s", s.typeOrNone()),
		}}
	} else if s.Max != nil {
		return &errpath.ErrField{Field: "maximum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.Max,
			Message: fmt.Sprintf("only valid for number type, got %s", s.typeOrNone()),
		}}
	} else if s.ExclusiveMin != nil {
		return &errpath.ErrField{Field: "exclusiveMinimum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.ExclusiveMin,
			Message: fmt.Sprintf("only valid for number type, got %s", s.typeOrNone()),
		}}
	} else if s.ExclusiveMax != nil {
		return &errpath.ErrField{Field: "exclusiveMaximum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.ExclusiveMax,
			Message: fmt.Sprintf("only valid for number type, got %s", s.typeOrNone()),
		}}
	} else if s.MultipleOf != nil {
		return &errpath.ErrField{Field: "multipleOf", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.MultipleOf,
			Message: fmt.Sprintf("only valid for number type, got %s", s.typeOrNone()),
		}}
	}

	if s.MultipleOf != nil && *s.MultipleOf <= 0 {
		return &errpath.ErrField{Field: "multipleOf", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.MultipleOf,
			Message: "must be greater than 0",
		}}
	}

	// String / Enum

	// Per JSON Schema 2020-12, enum and const can hold any JSON type; validate each value's kind matches the schema type.
	if s.Type != "" {
		for i, ev := range s.Enum {
			if !s.allowsKindOf(ev) {
				return &errpath.ErrField{Field: "enum", Err: &errpath.ErrIndex{Index: i, Err: &errpath.ErrInvalid[any]{
					Value:   jsonDisplayValue(ev),
					Message: fmt.Sprintf("must be a %s value", s.Type),
				}}}
			}
		}

		if s.Const != nil && !s.allowsKindOf(s.Const) {
			return &errpath.ErrField{Field: "const", Err: &errpath.ErrInvalid[any]{
				Value:   jsonDisplayValue(s.Const),
				Message: fmt.Sprintf("must be a %s value", s.Type),
			}}
		}

		if s.Example != nil && !s.allowsKindOf(s.Example) {
			return &errpath.ErrField{Field: "example", Err: &errpath.ErrInvalid[any]{
				Value:   jsonDisplayValue(s.Example),
				Message: fmt.Sprintf("must be a %s value", s.Type),
			}}
		}

		for i, ev := range s.Examples {
			if !s.allowsKindOf(ev) {
				return &errpath.ErrField{Field: "examples", Err: &errpath.ErrIndex{Index: i, Err: &errpath.ErrInvalid[any]{
					Value:   jsonDisplayValue(ev),
					Message: fmt.Sprintf("must be a %s value", s.Type),
				}}}
			}
		}
	}

	// Array

	// validate min and max items
	if s.allows(TypeArray) {
		if s.MaxItems != nil && s.MinItems > *s.MaxItems {
			return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
				Value:   s.MinItems,
				Message: fmt.Sprintf("minItems is greater than maxItems (%d > %d)", s.MinItems, *s.MaxItems),
			}}
		}

		for i, v := range s.PrefixItems {
			if err := v.Validate(); err != nil {
				return &errpath.ErrField{
					Field: "prefixItems",
					Err:   &errpath.ErrIndex{Index: i, Err: err},
				}
			}
		}

		// empty schema for items indicates a media type of application/octet-stream.
		if s.Items != nil && !s.Items.isEmpty() {
			if err := s.Items.Validate(); err != nil {
				return &errpath.ErrField{Field: "items", Err: err}
			}
		}
	} else if s.MinItems != 0 {
		return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinItems,
			Message: fmt.Sprintf("only valid for array type, got %s", s.typeOrNone()),
		}}
	} else if s.MaxItems != nil {
		return &errpath.ErrField{Field: "maxItems", Err: &errpath.ErrInvalid[uint]{
			Value:   *s.MaxItems,
			Message: fmt.Sprintf("only valid for array type, got %s", s.typeOrNone()),
		}}
	} else if s.UniqueItems {
		return &errpath.ErrField{Field: "uniqueItems", Err: &errpath.ErrInvalid[bool]{
			Value:   true,
			Message: fmt.Sprintf("only valid for array type, got %s", s.typeOrNone()),
		}}
	} else if len(s.PrefixItems) != 0 {
		return &errpath.ErrField{Field: "prefixItems", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for array type, got %s", s.typeOrNone()),
		}}
	} else if s.Items != nil {
		return &errpath.ErrField{Field: "items", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for array type, got %s", s.typeOrNone()),
		}}
	}

	// Object

	if s.allows(TypeObject) {
		if err := s.Properties.Validate(); err != nil {
			return &errpath.ErrField{Field: "properties", Err: err}
		}

		for i, r := range s.Required {
			// without a type, required is a constraint on whatever object properties hold
			if _, ok := s.Properties[r]; ok || s.Type == "" {
				continue
			}

			return &errpath.ErrField{
				Field: "required",
				Err: &errpath.ErrIndex{Index: i, Err: &errpath.ErrInvalid[string]{
					Value:   r,
					Message: "property does not exist",
				}},
			}
		}

		if s.AdditionalProperties != nil {
			if err := s.AdditionalProperties.Validate(); err != nil {
				return &errpath.ErrField{Field: "additionalProperties", Err: err}
			}
		}

		if s.PropertyNames != nil {
			if err := s.PropertyNames.Validate(); err != nil {
				return &errpath.ErrField{Field: "propertyNames", Err: err}
			}
		}

		if s.MaxProperties != nil && uint(len(s.Required)) > *s.MaxProperties {
			return &errpath.ErrField{Field: "maxProperties", Err: &errpath.ErrInvalid[uint]{
				Value:   *s.MaxProperties,
				Message: fmt.Sprintf("fewer than the %d required properties", len(s.Required)),
			}}
		}
	} else if s.Properties != nil {
		return &errpath.ErrField{Field: "properties", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for object type, got %s", s.typeOrNone()),
		}}
	} else if s.Required != nil {
		return &errpath.ErrField{Field: "required", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for object type, got %s", s.typeOrNone()),
		}}
	} else if s.AdditionalProperties != nil {
		return &errpath.ErrField{Field: "additionalProperties", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for object type, got %s", s.typeOrNone()),
		}}
	} else if s.MaxProperties != nil {
		return &errpath.ErrField{Field: "maxProperties", Err: &errpath.ErrInvalid[uint]{
			Value:   *s.MaxProperties,
			Message: fmt.Sprintf("only valid for object type, got %s", s.typeOrNone()),
		}}
	} else if s.PropertyNames != nil {
		return &errpath.ErrField{Field: "propertyNames", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for object type, got %s", s.typeOrNone()),
		}}
	}

	if s.Discriminator != nil {
		// a parent schema may carry the discriminator for the schemas extending it through allOf
		if len(s.OneOf) == 0 && len(s.AnyOf) == 0 && len(s.AllOf) == 0 && !extended {
			return &errpath.ErrField{Field: "discriminator", Err: &errpath.ErrInvalid[string]{
				Message: "only valid with oneOf, anyOf or allOf, or on a component schema another extends through allOf",
			}}
		}

		if err := s.Discriminator.Validate(); err != nil {
			return &errpath.ErrField{Field: "discriminator", Err: err}
		}
	}

	// validate default
	if len(s.Default) > 0 {
		defaultTypeErr := func() error {
			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[any]{
				Value:   jsonDisplayValue(s.Default),
				Message: fmt.Sprintf("does not match schema type, got %s", s.typeOrNone()),
			}}
		}

		switch s.Type {
		case TypeString:
			if s.Default.Kind() != jsontext.KindString {
				return defaultTypeErr()
			}
		case TypeNumber:
			if s.Default.Kind() != jsontext.KindNumber {
				return defaultTypeErr()
			}
		case TypeInteger:
			if s.Default.Kind() != jsontext.KindNumber || !isJSONInteger(s.Default) {
				return defaultTypeErr()
			}
		case TypeBoolean:
			if s.Default.Kind() != jsontext.KindTrue && s.Default.Kind() != jsontext.KindFalse {
				return defaultTypeErr()
			}
		case TypeArray:
			if s.Default.Kind() != jsontext.KindBeginArray {
				return defaultTypeErr()
			}
		case TypeObject:
			if s.Default.Kind() != jsontext.KindBeginObject {
				return defaultTypeErr()
			}
		case TypeNull:
			if s.Default.Kind() != jsontext.KindNull {
				return defaultTypeErr()
			}
		}

		if len(s.Enum) > 0 {
			found := false
			for _, ev := range s.Enum {
				if bytes.Equal(ev, s.Default) {
					found = true
					break
				}
			}

			if !found {
				parts := make([]string, len(s.Enum))
				for i, ev := range s.Enum {
					parts[i] = ev.String()
				}

				return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[any]{
					Value:   jsonDisplayValue(s.Default),
					Message: fmt.Sprintf("is not one of the enums ([%s])", strings.Join(parts, " ")),
				}}
			}
		}
	}

	return validateExtensions(s.Extensions)
}

// allows reports whether the schema's keywords for type t apply: it is of type t, or of no type at all.
func (s *Schema) allows(t DataType) bool { return s.Type == "" || s.Type == t }

// typeOrNone names the schema's type for error messages.
func (s *Schema) typeOrNone() string {
	if s.Type == "" {
		return "no type"
	}

	return string(s.Type)
}

// allowsKindOf reports whether v's kind is one the schema's type allows: its
// Type's, or null when the schema is nullable.
func (s *Schema) allowsKindOf(v jsontext.Value) bool {
	return s.Nullable && v.Kind() == jsontext.KindNull || enumKindMatchesType(v, s.Type)
}

// enumKindMatchesType reports whether a JSON value's kind is compatible with the given DataType.
// For TypeInteger it additionally requires the number to be a whole number.
func enumKindMatchesType(v jsontext.Value, t DataType) bool {
	switch t {
	case TypeString:
		return v.Kind() == jsontext.KindString
	case TypeNumber:
		return v.Kind() == jsontext.KindNumber
	case TypeInteger:
		return v.Kind() == jsontext.KindNumber && isJSONInteger(v)
	case TypeBoolean:
		return v.Kind() == jsontext.KindTrue || v.Kind() == jsontext.KindFalse
	case TypeNull:
		return v.Kind() == jsontext.KindNull
	case TypeArray:
		return v.Kind() == jsontext.KindBeginArray
	case TypeObject:
		return v.Kind() == jsontext.KindBeginObject
	default:
		return true
	}
}

// isJSONInteger reports whether a JSON number value represents a whole number.
func isJSONInteger(v jsontext.Value) bool {
	f, err := strconv.ParseFloat(string(v), 64)
	return err == nil && f == float64(int64(f))
}

// jsonDisplayValue converts a jsontext.Value to a typed Go value suitable for
// errpath.ErrInvalid display formatting.
func jsonDisplayValue(v jsontext.Value) any {
	switch v.Kind() {
	case jsontext.KindString:
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			return s
		}
	case jsontext.KindNumber:
		var f float64
		if err := json.Unmarshal(v, &f); err == nil {
			return f
		}
	case jsontext.KindTrue:
		return true
	case jsontext.KindFalse:
		return false
	default:
	}

	return string(v)
}

func (l *loader) collectSchema(s *Schema, ref ref) {
	l.schemas[ref.String()] = s // collect this schema
}

func (l *loader) resolveSchema(s *Schema) error {
	if s.Ref != nil {
		target, ok := l.schemas[s.Ref.Identifier]
		if !ok {
			return fmt.Errorf("couldn't resolve %q", s.Ref.Identifier)
		}

		s.Ref.Value = target
	}

	if err := l.resolveSchemaList(s.AllOf); err != nil {
		return &errpath.ErrField{Field: "allOf", Err: err}
	}

	if err := l.resolveSchemaList(s.OneOf); err != nil {
		return &errpath.ErrField{Field: "oneOf", Err: err}
	}

	if err := l.resolveSchemaList(s.AnyOf); err != nil {
		return &errpath.ErrField{Field: "anyOf", Err: err}
	}

	if s.Not != nil {
		if err := l.resolveSchema(s.Not); err != nil {
			return &errpath.ErrField{Field: "not", Err: err}
		}
	}

	if err := l.resolveSchemaList(s.PrefixItems); err != nil {
		return &errpath.ErrField{Field: "prefixItems", Err: err}
	}

	if s.Items != nil {
		if err := l.resolveSchema(s.Items); err != nil {
			return &errpath.ErrField{Field: "items", Err: err}
		}
	}

	if err := l.resolveSchemas(s.Properties); err != nil {
		return &errpath.ErrField{Field: "properties", Err: err}
	}

	if s.AdditionalProperties != nil && s.AdditionalProperties.Schema != nil {
		if err := l.resolveSchema(s.AdditionalProperties.Schema); err != nil {
			return &errpath.ErrField{Field: "additionalProperties", Err: err}
		}
	}

	if s.PropertyNames != nil {
		if err := l.resolveSchema(s.PropertyNames); err != nil {
			return &errpath.ErrField{Field: "propertyNames", Err: err}
		}
	}

	if s.Discriminator != nil {
		if err := l.resolveDiscriminator(s.Discriminator); err != nil {
			return &errpath.ErrField{Field: "discriminator", Err: err}
		}
	}

	return nil
}

// derefType is the schema's type, or for a reference without one, the type of the schema it points to.
func (s *Schema) derefType() DataType {
	// a document built in code has not been checked for cycles, so stop at the first schema seen twice
	seen := map[*Schema]bool{}
	for s.Type == "" && s.Ref != nil && s.Ref.Value != nil && !seen[s] {
		seen[s] = true
		s = s.Ref.Value
	}

	return s.Type
}

func (s *Schema) isEmpty() bool {
	return s == nil ||
		(s.Ref == nil && s.Type == "" && !s.Nullable && s.Format == "" &&
			len(s.AllOf) == 0 && len(s.OneOf) == 0 && len(s.AnyOf) == 0 && s.Not == nil &&
			s.Min == nil && s.Max == nil && s.ExclusiveMin == nil && s.ExclusiveMax == nil && s.MultipleOf == nil &&
			s.MinLength == 0 && s.MaxLength == nil && s.Pattern == nil &&
			s.MinItems == 0 && s.MaxItems == nil && !s.UniqueItems && len(s.PrefixItems) == 0 && s.Items == nil &&
			s.Properties == nil && s.Required == nil &&
			s.AdditionalProperties == nil && s.MaxProperties == nil && s.PropertyNames == nil && s.Discriminator == nil &&
			len(s.Examples) == 0 && !s.Deprecated && !s.ReadOnly && !s.WriteOnly &&
			s.ContentMediaType == "" && s.ContentEncoding == "" &&
			s.Const == nil &&
			s.Example == nil)
}
