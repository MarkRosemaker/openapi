package openapi

import (
	"bytes"
	"encoding/json/jsontext"
	"fmt"
	"regexp"
	"slices"
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
	// Further refines the data type.
	Format Format `json:"format,omitempty" yaml:"format,omitempty"`

	// AllOf validates the value against ALL of the given schemas.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	AllOf SchemaRefList `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	// OneOf validates the value against EXACTLY ONE of the given schemas.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	OneOf SchemaRefList `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	// AnyOf validates the value against AT LEAST ONE of the given schemas.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	AnyOf SchemaRefList `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// Not validates the value against the negation of the given schema — the value must NOT validate against it.
	// See: https://spec.openapis.org/oas/v3.2.0.html#schema-object
	Not *SchemaRef `json:"not,omitempty" yaml:"not,omitempty"`

	// Integer / Number

	// The minimum value of the number.
	Min *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// The maximum value of the number.
	Max *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`

	// String

	// The pattern is used to validate the string.
	// This string SHOULD be a valid regular expression, according to the Ecma-262 Edition 5.1 regular expression dialect.
	// NOTE: We simply use text unmarshalling for this field. This guarantees that the regular expression is valid or we can't unmarshal.
	Pattern *regexp.Regexp `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	// A list of possible values. Per JSON Schema 2020-12, enum may contain any JSON type.
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`

	// Array

	// The minimum number of items in the array.
	MinItems uint `json:"minItems,omitzero" yaml:"minItems,omitempty"`
	// The maximum number of items in the array.
	MaxItems *uint `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	// The items of the array. When the type is array, this property is REQUIRED.
	// The empty schema for `items` indicates a media type of `application/octet-stream`.
	Items *SchemaRef `json:"items,omitzero" yaml:"items,omitempty"`

	// Object

	// For object types, defines the properties of the object
	Properties SchemaRefs `json:"properties,omitzero" yaml:"properties,omitempty"`
	// Which properties are required.
	Required             []string   `json:"required,omitempty" yaml:"required,omitempty"`
	AdditionalProperties *SchemaRef `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`

	// special encoding for binary data
	ContentMediaType string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`
	ContentEncoding  string `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`

	// Specifies the default value of the property if no value is provided.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`

	Example jsontext.Value `json:"example,omitzero" yaml:"example,omitzero"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`

	// an index to the original location of this object
	idx int

	// NOTE: consider adding:
	// Indicates whether the property can have a null value.
	// Nullable bool `json:"nullable,omitempty,omitzero" yaml:"nullable,omitempty"`
}

func getIndexSchema(s *Schema) int              { return s.idx }
func setIndexSchema(s *Schema, idx int) *Schema { s.idx = idx; return s }

func (s *Schema) Validate() error {
	s.Description = strings.TrimSpace(s.Description)

	if s.Type == "" {
		if len(s.AllOf) == 0 && len(s.OneOf) == 0 && len(s.AnyOf) == 0 && s.Not == nil {
			return &errpath.ErrField{Field: "type", Err: &errpath.ErrRequired{}}
		}
	} else if err := s.Type.Validate(); err != nil {
		return &errpath.ErrField{Field: "type", Err: err}
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
		if s.Type != TypeInteger {
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for integer type, got %s", s.Type),
			}}
		}
	case FormatFloat, FormatDouble:
		if s.Type != TypeNumber {
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
			}}
		}
	case FormatEmail, FormatPassword,
		FormatUUID, FormatURI, FormatURIRef, FormatZipCode,
		FormatIPv4, FormatIPv6:
		if s.Type != TypeString {
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for string type, got %s", s.Type),
			}}
		}
	case FormatDuration, FormatDate, FormatDateTime:
		switch s.Type {
		case TypeInteger, TypeString:
		default:
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for integer or string type, got %s", s.Type),
			}}
		}
	case FormatByte, FormatBinary:
		switch s.Type {
		case TypeString:
		default:
			return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
				Value:   s.Format,
				Message: fmt.Sprintf("only valid for string type, got %s", s.Type),
			}}
		}
	default:
		return fmt.Errorf("unimplemented format: %s", s.Format)
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
	}

	if s.Type == TypeNumber || s.Type == TypeInteger {
		if s.Min != nil && s.Max != nil && *s.Min > *s.Max {
			return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
				Value:   *s.Min,
				Message: fmt.Sprintf("minimum is greater than maximum (%v > %v)", *s.Min, *s.Max),
			}}
		}
	} else if s.Min != nil {
		return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.Min,
			Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
		}}
	} else if s.Max != nil {
		return &errpath.ErrField{Field: "maximum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.Max,
			Message: fmt.Sprintf("only valid for number type, got %s", s.Type),
		}}
	}

	// String / Enum

	// Per JSON Schema 2020-12, enum can hold any JSON type; validate each value matches the schema type.
	if s.Type != "" {
		for i, ev := range s.Enum {
			if !enumValueMatchesType(ev, s.Type) {
				return &errpath.ErrField{Field: "enum", Err: &errpath.ErrIndex{Index: i, Err: &errpath.ErrInvalid[any]{
					Value:   ev,
					Message: fmt.Sprintf("must be a %s value", s.Type),
				}}}
			}
		}
	}

	// Array

	// validate min and max items
	if s.Type == TypeArray {
		if s.MaxItems != nil && s.MinItems > *s.MaxItems {
			return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
				Value:   s.MinItems,
				Message: fmt.Sprintf("minItems is greater than maxItems (%d > %d)", s.MinItems, *s.MaxItems),
			}}
		}

		if s.Items == nil {
			return &errpath.ErrField{Field: "items", Err: &errpath.ErrRequired{}}
		}

		// empty schema for items indicates a media type of application/octet-stream.
		if !s.Items.Value.isEmpty() {
			if err := s.Items.Validate(); err != nil {
				return &errpath.ErrField{Field: "items", Err: err}
			}
		}
	} else if s.MinItems != 0 {
		return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinItems,
			Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
		}}
	} else if s.MaxItems != nil {
		return &errpath.ErrField{Field: "maxItems", Err: &errpath.ErrInvalid[uint]{
			Value:   *s.MaxItems,
			Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
		}}
	} else if s.Items != nil {
		return &errpath.ErrField{Field: "items", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for array type, got %s", s.Type),
		}}
	}

	// Object

	if s.Type == TypeObject {
		if err := s.Properties.Validate(); err != nil {
			return &errpath.ErrField{Field: "properties", Err: err}
		}

		for i, r := range s.Required {
			if _, ok := s.Properties[r]; ok {
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
	} else if s.Properties != nil {
		return &errpath.ErrField{Field: "properties", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for object type, got %s", s.Type),
		}}
	} else if s.AdditionalProperties != nil {
		return &errpath.ErrField{Field: "additionalProperties", Err: &errpath.ErrInvalid[string]{
			Message: fmt.Sprintf("only valid for object type, got %s", s.Type),
		}}
	}

	// validate default
	switch dflt := s.Default.(type) {
	case nil: // empty
	case string:
		if s.Type != TypeString {
			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[string]{
				Value:   dflt,
				Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
			}}
		}

		if s.Enum != nil {
			if !slices.Contains(s.Enum, any(dflt)) {
				return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[string]{
					Value:   dflt,
					Message: fmt.Sprintf("is not one of the enums (%q)", s.Enum),
				}}
			}
		}
	case float64:
		switch s.Type {
		case TypeNumber: // fits
		case TypeInteger:
			if asInt := int(dflt); dflt != float64(asInt) {
				return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[float64]{
					Value:   dflt,
					Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
				}}
			} else {
				s.Default = asInt // set to int version
			}
		default:
			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[float64]{
				Value:   dflt,
				Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
			}}
		}
	case int:
		switch s.Type {
		case TypeNumber, TypeInteger: // fits
		default:
			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[int]{
				Value:   dflt,
				Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
			}}
		}
	case bool:
		switch s.Type {
		case TypeBoolean: // fits
		default:
			return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[bool]{
				Value:   dflt,
				Message: fmt.Sprintf("does not match schema type, got %s", s.Type),
			}}
		}
	default:
		return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[any]{
			Value:   s.Default,
			Message: fmt.Sprintf("unknown type %T", s.Default),
		}}
	}

	return nil
}

// enumValueMatchesType reports whether a JSON-decoded value is compatible with the given DataType.
// JSON numbers unmarshal to float64, integers must additionally be whole numbers.
func enumValueMatchesType(v any, t DataType) bool {
	switch t {
	case TypeString:
		_, ok := v.(string)
		return ok
	case TypeInteger:
		f, ok := v.(float64)
		return ok && f == float64(int64(f))
	case TypeNumber:
		_, ok := v.(float64)
		return ok
	case TypeBoolean:
		_, ok := v.(bool)
		return ok
	case TypeNull:
		return v == nil
	case TypeArray:
		_, ok := v.([]any)
		return ok
	case TypeObject:
		_, ok := v.(map[string]any)
		return ok
	default:
		return true
	}
}

func (l *loader) collectSchema(s *Schema, ref ref) {
	l.schemas[ref.String()] = s // collect this schema
}

func (l *loader) resolveSchemaRef(s *SchemaRef) error {
	if s.Ref != nil && s.Value != nil {
		// $ref with sibling fields (JSON Schema 2020-12 / OAS 3.1+).
		// Sibling values take priority; the ref fills in anything not already set.
		ref, ok := l.schemas[s.Ref.Identifier]
		if !ok {
			return fmt.Errorf("couldn't resolve %q", s.Ref.Identifier)
		}
		mergeSchemaFrom(s.Value, ref)
		// Remove the now-redundant "$ref" key from Extensions so the marshaled
		// output is clean (fully inlined, no dangling $ref alongside the merged fields).
		s.Value.Extensions = removeFromExtensions(s.Value.Extensions, "$ref")
		// Clear Ref so this is treated as a resolved inline schema by Validate().
		s.Ref = nil
		return l.resolveSchema(s.Value)
	}
	return resolveRef(s, l.schemas, l.resolveSchema)
}

// removeFromExtensions returns a copy of ext with the named key removed.
// Returns nil if the result is an empty object.
func removeFromExtensions(ext Extensions, key string) Extensions {
	if ext == nil {
		return nil
	}

	dec := jsontext.NewDecoder(bytes.NewReader(ext))
	if dec.PeekKind() != '{' {
		return ext
	}
	dec.ReadToken() // consume '{'

	var buf bytes.Buffer
	enc := jsontext.NewEncoder(&buf)
	enc.WriteToken(jsontext.BeginObject)

	wrote := false
	for dec.PeekKind() != '}' {
		keyTkn, err := dec.ReadToken()
		if err != nil {
			return ext
		}
		// Save string before ReadValue voids the token.
		keyStr := keyTkn.String()

		val, err := dec.ReadValue()
		if err != nil {
			return ext
		}
		if keyStr == key {
			continue
		}
		enc.WriteToken(jsontext.String(keyStr))
		enc.WriteValue(val)
		wrote = true
	}

	enc.WriteToken(jsontext.EndObject)

	if !wrote {
		return nil
	}
	return Extensions(buf.Bytes())
}

func (l *loader) resolveSchema(s *Schema) error {
	if err := l.resolveSchemaRefList(s.AllOf); err != nil {
		return &errpath.ErrField{Field: "allOf", Err: err}
	}

	if err := l.resolveSchemaRefList(s.OneOf); err != nil {
		return &errpath.ErrField{Field: "oneOf", Err: err}
	}

	if err := l.resolveSchemaRefList(s.AnyOf); err != nil {
		return &errpath.ErrField{Field: "anyOf", Err: err}
	}

	if s.Not != nil {
		if err := l.resolveSchemaRef(s.Not); err != nil {
			return &errpath.ErrField{Field: "not", Err: err}
		}
	}

	if s.Items != nil {
		if err := l.resolveSchemaRef(s.Items); err != nil {
			return &errpath.ErrField{Field: "items", Err: err}
		}
	}

	if err := l.resolveSchemaRefs(s.Properties); err != nil {
		return &errpath.ErrField{Field: "properties", Err: err}
	}

	if s.AdditionalProperties != nil {
		if err := l.resolveSchemaRef(s.AdditionalProperties); err != nil {
			return &errpath.ErrField{Field: "additionalProperties", Err: err}
		}
	}

	return nil
}

func (s *Schema) isEmpty() bool {
	return s == nil ||
		(s.Title == "" && s.Description == "" &&
			s.Type == "" && s.Format == "" &&
			len(s.AllOf) == 0 && len(s.OneOf) == 0 && len(s.AnyOf) == 0 && s.Not == nil &&
			s.Min == nil && s.Max == nil &&
			s.Pattern == nil &&
			len(s.Enum) == 0 &&
			s.MinItems == 0 && s.MaxItems == nil && s.Items == nil &&
			s.Properties == nil && s.Required == nil &&
			s.AdditionalProperties == nil &&
			s.ContentMediaType == "" && s.ContentEncoding == "" &&
			s.Default == nil && len(s.Example) == 0)
}

// mergeSchemaFrom copies non-zero fields from src to dst, giving precedence to dst's existing values.
// Used when a $ref has sibling keywords: the sibling values take priority, the ref fills the rest.
func mergeSchemaFrom(dst, src *Schema) {
	if dst.Title == "" {
		dst.Title = src.Title
	}
	if dst.Description == "" {
		dst.Description = src.Description
	}
	if dst.Type == "" {
		dst.Type = src.Type
	}
	if dst.Format == "" {
		dst.Format = src.Format
	}
	if len(dst.AllOf) == 0 {
		dst.AllOf = src.AllOf
	}
	if len(dst.OneOf) == 0 {
		dst.OneOf = src.OneOf
	}
	if len(dst.AnyOf) == 0 {
		dst.AnyOf = src.AnyOf
	}
	if dst.Not == nil {
		dst.Not = src.Not
	}
	if dst.Min == nil {
		dst.Min = src.Min
	}
	if dst.Max == nil {
		dst.Max = src.Max
	}
	if dst.Pattern == nil {
		dst.Pattern = src.Pattern
	}
	if len(dst.Enum) == 0 {
		dst.Enum = src.Enum
	}
	if dst.MinItems == 0 {
		dst.MinItems = src.MinItems
	}
	if dst.MaxItems == nil {
		dst.MaxItems = src.MaxItems
	}
	if dst.Items == nil {
		dst.Items = src.Items
	}
	if dst.Properties == nil {
		dst.Properties = src.Properties
	}
	if len(dst.Required) == 0 {
		dst.Required = src.Required
	}
	if dst.AdditionalProperties == nil {
		dst.AdditionalProperties = src.AdditionalProperties
	}
	if dst.ContentMediaType == "" {
		dst.ContentMediaType = src.ContentMediaType
	}
	if dst.ContentEncoding == "" {
		dst.ContentEncoding = src.ContentEncoding
	}
	if dst.Default == nil {
		dst.Default = src.Default
	}
	if len(dst.Example) == 0 {
		dst.Example = src.Example
	}
}
