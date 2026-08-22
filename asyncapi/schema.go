package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Schema allows the definition of input and output data types.
// These types can be objects, but also primitives and arrays.
// This object is a superset of the [JSON Schema Specification Draft 07].
//
// The empty schema (which allows any instance to validate) MAY be represented by the boolean
// value `true` and a schema which allows no instance to validate MAY be represented by the
// boolean value `false`. Both are represented by the [Schema.Boolean] field.
//
// For other formats (e.g. Avro, RAML, etc.) see the [MultiFormatSchema] object.
// ([Specification])
//
// [JSON Schema Specification Draft 07]: https://json-schema.org/specification-links.html#draft-7
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#schemaObject
type Schema struct {
	// A boolean schema. `true` allows any instance to validate, `false` allows none.
	// When it is set, all other fields are ignored.
	Boolean *bool `json:"-" yaml:"-"`

	// The URI of the schema, used to identify it and to resolve relative references against.
	ID *url.URL `json:"$id,omitempty" yaml:"$id,omitempty"`
	// The dialect of the schema.
	Dialect *url.URL `json:"$schema,omitempty" yaml:"$schema,omitempty"`
	// A comment for the schema that is not meant to be displayed to end users.
	Comment string `json:"$comment,omitempty" yaml:"$comment,omitempty"`

	// The name of the schema.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short description of the schema. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Specifies the data type of the schema. It is either a single type or a list of types.
	Type DataTypes `json:"type,omitempty" yaml:"type,omitempty"`
	// Further refines the data type. See [Format] for the formats defined by the specification.
	Format Format `json:"format,omitempty" yaml:"format,omitempty"`

	// Composition

	// AllOf validates the value against ALL of the given schemas.
	AllOf AnySchemaRefList `json:"allOf,omitempty" yaml:"allOf,omitempty"`
	// OneOf validates the value against EXACTLY ONE of the given schemas.
	OneOf AnySchemaRefList `json:"oneOf,omitempty" yaml:"oneOf,omitempty"`
	// AnyOf validates the value against AT LEAST ONE of the given schemas.
	AnyOf AnySchemaRefList `json:"anyOf,omitempty" yaml:"anyOf,omitempty"`
	// Not validates the value against the negation of the given schema.
	Not *AnySchemaRef `json:"not,omitempty" yaml:"not,omitempty"`
	// If is the condition of a conditional schema.
	If *AnySchemaRef `json:"if,omitempty" yaml:"if,omitempty"`
	// Then is applied when the value validates against the schema given in the `if` keyword.
	Then *AnySchemaRef `json:"then,omitempty" yaml:"then,omitempty"`
	// Else is applied when the value does not validate against the schema given in the `if` keyword.
	Else *AnySchemaRef `json:"else,omitempty" yaml:"else,omitempty"`

	// Integer / Number

	// The value must be a multiple of this number.
	MultipleOf *float64 `json:"multipleOf,omitempty" yaml:"multipleOf,omitempty"`
	// The minimum value of the number.
	Min *float64 `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	// The exclusive minimum value of the number.
	ExclusiveMin *float64 `json:"exclusiveMinimum,omitempty" yaml:"exclusiveMinimum,omitempty"`
	// The maximum value of the number.
	Max *float64 `json:"maximum,omitempty" yaml:"maximum,omitempty"`
	// The exclusive maximum value of the number.
	ExclusiveMax *float64 `json:"exclusiveMaximum,omitempty" yaml:"exclusiveMaximum,omitempty"`

	// String

	// The minimum length of the string.
	MinLength uint `json:"minLength,omitzero" yaml:"minLength,omitempty"`
	// The maximum length of the string.
	MaxLength *uint `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`
	// The pattern is used to validate the string.
	// This string SHOULD be a valid regular expression, according to the ECMA 262 regular expression dialect.
	// NOTE: We simply use text unmarshalling for this field. This guarantees that the regular expression is valid or we can't unmarshal.
	Pattern *regexp.Regexp `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Array

	// The minimum number of items in the array.
	MinItems uint `json:"minItems,omitzero" yaml:"minItems,omitempty"`
	// The maximum number of items in the array.
	MaxItems *uint `json:"maxItems,omitempty" yaml:"maxItems,omitempty"`
	// Whether the items of the array must be unique.
	UniqueItems bool `json:"uniqueItems,omitzero" yaml:"uniqueItems,omitempty"`
	// The schema the items of the array must validate against.
	Items *AnySchemaRef `json:"items,omitempty" yaml:"items,omitempty"`
	// The schema the additional items of the array must validate against.
	AdditionalItems *AnySchemaRef `json:"additionalItems,omitempty" yaml:"additionalItems,omitempty"`
	// The schema at least one item of the array must validate against.
	Contains *AnySchemaRef `json:"contains,omitempty" yaml:"contains,omitempty"`

	// Object

	// The minimum number of properties of the object.
	MinProperties uint `json:"minProperties,omitzero" yaml:"minProperties,omitempty"`
	// The maximum number of properties of the object.
	MaxProperties *uint `json:"maxProperties,omitempty" yaml:"maxProperties,omitempty"`
	// Which properties are required.
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`
	// The properties of the object.
	Properties Schemas `json:"properties,omitempty" yaml:"properties,omitempty"`
	// The properties of the object whose names match a regular expression.
	PatternProperties Schemas `json:"patternProperties,omitempty" yaml:"patternProperties,omitempty"`
	// The schema the additional properties of the object must validate against.
	AdditionalProperties *AnySchemaRef `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	// The schema the property names of the object must validate against.
	PropertyNames *AnySchemaRef `json:"propertyNames,omitempty" yaml:"propertyNames,omitempty"`
	// Reusable schemas that are referenced from within this schema.
	Definitions Schemas `json:"definitions,omitempty" yaml:"definitions,omitempty"`

	// Values

	// A list of possible values.
	Enum []jsontext.Value `json:"enum,omitempty" yaml:"enum,omitempty"`
	// The only possible value.
	Const jsontext.Value `json:"const,omitempty" yaml:"const,omitempty"`
	// The value that is used if no other value is present.
	// Unlike JSON Schema, the value MUST conform to the defined type for the schema defined at the same level.
	Default jsontext.Value `json:"default,omitempty" yaml:"default,omitempty"`
	// A list of examples of the value.
	Examples []jsontext.Value `json:"examples,omitempty" yaml:"examples,omitempty"`

	// special encoding for binary data
	ContentEncoding  string `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`
	ContentMediaType string `json:"contentMediaType,omitempty" yaml:"contentMediaType,omitempty"`

	// Whether the value is only sent by the server and must not be sent by the client.
	ReadOnly bool `json:"readOnly,omitzero" yaml:"readOnly,omitempty"`
	// Whether the value is only sent by the client and must not be sent by the server.
	WriteOnly bool `json:"writeOnly,omitzero" yaml:"writeOnly,omitempty"`

	// AsyncAPI vocabulary

	// Adds support for polymorphism.
	// The discriminator is the schema property name that is used to differentiate between other schemas that inherit this schema.
	// The property name used MUST be defined at this schema and it MUST be in the required property list.
	Discriminator string `json:"discriminator,omitempty" yaml:"discriminator,omitempty"`
	// Additional external documentation for this schema.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// Specifies that a schema is deprecated and SHOULD be transitioned out of usage.
	Deprecated bool `json:"deprecated,omitzero" yaml:"deprecated,omitempty"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// schemaValue is a [Schema] without the custom marshalling, used to avoid infinite recursion.
type schemaValue Schema

var _ json.UnmarshalerFrom = (*Schema)(nil)

// UnmarshalJSONFrom unmarshals the schema, which may be a boolean schema.
func (s *Schema) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case 't', 'f':
		var b bool
		if err := json.UnmarshalDecode(dec, &b); err != nil {
			return err
		}

		s.Boolean = &b

		return nil
	}

	return json.UnmarshalDecode(dec, (*schemaValue)(s))
}

var _ json.MarshalerTo = (*Schema)(nil)

// MarshalJSONTo marshals the schema, which may be a boolean schema.
func (s *Schema) MarshalJSONTo(enc *jsontext.Encoder) error {
	if s.Boolean != nil {
		return json.MarshalEncode(enc, *s.Boolean)
	}

	return json.MarshalEncode(enc, (*schemaValue)(s))
}

// Validate checks the schema for correctness.
func (s *Schema) Validate() error {
	if s.Boolean != nil {
		return nil // a boolean schema has no other fields
	}

	s.Description = strings.TrimSpace(s.Description)

	if err := s.Type.Validate(); err != nil {
		return &errpath.ErrField{Field: "type", Err: err}
	}

	// NOTE: The format is an open string-valued property,
	// so a format that is not defined by the specification is not an error.

	if err := s.validateComposition(); err != nil {
		return err
	}

	if err := s.validateRanges(); err != nil {
		return err
	}

	if err := s.validateProperties(); err != nil {
		return err
	}

	for _, sub := range s.subSchemas() {
		if sub.schema == nil {
			continue
		}

		if err := sub.schema.Validate(); err != nil {
			return &errpath.ErrField{Field: sub.field, Err: err}
		}
	}

	if s.ExternalDocs != nil {
		if err := s.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return validateExtensions(s.Extensions)
}

func (s *Schema) validateComposition() error {
	if err := s.AllOf.Validate(); err != nil {
		return &errpath.ErrField{Field: "allOf", Err: err}
	}

	if err := s.OneOf.Validate(); err != nil {
		return &errpath.ErrField{Field: "oneOf", Err: err}
	}

	if err := s.AnyOf.Validate(); err != nil {
		return &errpath.ErrField{Field: "anyOf", Err: err}
	}

	return nil
}

// subSchema pairs a field name with the subschema it holds.
type subSchema struct {
	field  string
	schema *AnySchemaRef
}

// subSchemas returns the subschemas that are held by a single field, in the order of the fields.
func (s *Schema) subSchemas() []subSchema {
	return []subSchema{
		{"not", s.Not}, {"if", s.If}, {"then", s.Then}, {"else", s.Else},
		{"items", s.Items}, {"additionalItems", s.AdditionalItems}, {"contains", s.Contains},
		{"additionalProperties", s.AdditionalProperties}, {"propertyNames", s.PropertyNames},
	}
}

// validateRanges checks that the lower bounds are not greater than the upper bounds.
func (s *Schema) validateRanges() error {
	if s.Min != nil && s.Max != nil && *s.Min > *s.Max {
		return &errpath.ErrField{Field: "minimum", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.Min,
			Message: fmt.Sprintf("minimum is greater than maximum (%v > %v)", *s.Min, *s.Max),
		}}
	}

	if s.MaxLength != nil && s.MinLength > *s.MaxLength {
		return &errpath.ErrField{Field: "minLength", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinLength,
			Message: fmt.Sprintf("minLength is greater than maxLength (%d > %d)", s.MinLength, *s.MaxLength),
		}}
	}

	if s.MaxItems != nil && s.MinItems > *s.MaxItems {
		return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinItems,
			Message: fmt.Sprintf("minItems is greater than maxItems (%d > %d)", s.MinItems, *s.MaxItems),
		}}
	}

	if s.MaxProperties != nil && s.MinProperties > *s.MaxProperties {
		return &errpath.ErrField{Field: "minProperties", Err: &errpath.ErrInvalid[uint]{
			Value: s.MinProperties,
			Message: fmt.Sprintf("minProperties is greater than maxProperties (%d > %d)",
				s.MinProperties, *s.MaxProperties),
		}}
	}

	if s.MultipleOf != nil && *s.MultipleOf <= 0 {
		return &errpath.ErrField{Field: "multipleOf", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.MultipleOf,
			Message: "must be greater than zero",
		}}
	}

	return nil
}

func (s *Schema) validateProperties() error {
	if err := s.Properties.Validate(); err != nil {
		return &errpath.ErrField{Field: "properties", Err: err}
	}

	if err := s.PatternProperties.Validate(); err != nil {
		return &errpath.ErrField{Field: "patternProperties", Err: err}
	}

	if err := s.Definitions.Validate(); err != nil {
		return &errpath.ErrField{Field: "definitions", Err: err}
	}

	// the discriminator MUST be defined at this schema and MUST be in the required property list
	if s.Discriminator != "" {
		if _, ok := s.Properties[s.Discriminator]; !ok {
			return &errpath.ErrField{Field: "discriminator", Err: &errpath.ErrInvalid[string]{
				Value:   s.Discriminator,
				Message: "property does not exist",
			}}
		}

		if !slices.Contains(s.Required, s.Discriminator) {
			return &errpath.ErrField{Field: "discriminator", Err: &errpath.ErrInvalid[string]{
				Value:   s.Discriminator,
				Message: "property must be required",
			}}
		}
	}

	return nil
}

// SortMaps sorts the properties of the schema and of all of its subschemas by key.
func (s *Schema) SortMaps() {
	if s == nil {
		return
	}

	s.Properties.Sort()
	s.PatternProperties.Sort()
	s.Definitions.Sort()

	for _, ss := range []Schemas{s.Properties, s.PatternProperties, s.Definitions} {
		for _, sub := range ss {
			sub.Value.SortMaps()
		}
	}

	for _, ss := range []AnySchemaRefList{s.AllOf, s.OneOf, s.AnyOf} {
		for _, sub := range ss {
			sub.Value.SortMaps()
		}
	}

	for _, sub := range s.subSchemas() {
		if sub.schema != nil {
			sub.schema.Value.SortMaps()
		}
	}
}

func (l *loader) collectSchema(s *Schema, ref ref) {
	l.collectAnySchemaRefList(s.AllOf, append(ref, "allOf"))
	l.collectAnySchemaRefList(s.OneOf, append(ref, "oneOf"))
	l.collectAnySchemaRefList(s.AnyOf, append(ref, "anyOf"))

	for _, sub := range s.subSchemas() {
		if sub.schema != nil {
			l.collectAnySchemaRef(sub.schema, append(ref, sub.field))
		}
	}

	l.collectSchemas(s.Properties, append(ref, "properties"))
	l.collectSchemas(s.PatternProperties, append(ref, "patternProperties"))
	l.collectSchemas(s.Definitions, append(ref, "definitions"))

	if s.ExternalDocs != nil {
		l.collectExternalDocsRef(s.ExternalDocs, append(ref, "externalDocs"))
	}
}

func (l *loader) resolveSchema(s *Schema) error {
	if err := l.resolveAnySchemaRefList(s.AllOf); err != nil {
		return &errpath.ErrField{Field: "allOf", Err: err}
	}

	if err := l.resolveAnySchemaRefList(s.OneOf); err != nil {
		return &errpath.ErrField{Field: "oneOf", Err: err}
	}

	if err := l.resolveAnySchemaRefList(s.AnyOf); err != nil {
		return &errpath.ErrField{Field: "anyOf", Err: err}
	}

	for _, sub := range s.subSchemas() {
		if sub.schema == nil {
			continue
		}

		if err := l.resolveAnySchemaRef(sub.schema); err != nil {
			return &errpath.ErrField{Field: sub.field, Err: err}
		}
	}

	if err := l.resolveSchemas(s.Properties); err != nil {
		return &errpath.ErrField{Field: "properties", Err: err}
	}

	if err := l.resolveSchemas(s.PatternProperties); err != nil {
		return &errpath.ErrField{Field: "patternProperties", Err: err}
	}

	if err := l.resolveSchemas(s.Definitions); err != nil {
		return &errpath.ErrField{Field: "definitions", Err: err}
	}

	if s.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(s.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return nil
}
