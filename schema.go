package openapi

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"regexp"
	"slices"
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
	if s.Ref != nil && s.Ref.Value == nil {
		return &errpath.ErrField{Field: "$ref", Err: fmt.Errorf("%q was not resolved", s.Ref.Identifier)}
	}

	// type is optional (JSON Schema 2020-12): without one, every type's keywords apply, each to instances of its type.
	if s.Type != "" {
		if err := s.Type.Validate(); err != nil {
			return &errpath.ErrField{Field: "type", Err: err}
		}
	}

	for _, check := range []func() error{
		s.validateFormat,
		s.validateString,
		s.validateComposition,
		s.validateNumber,
		s.validateValueKinds,
		s.validateArray,
		s.validateObject,
		func() error { return s.validateDiscriminator(extended) },
		s.validateDefault,
	} {
		if err := check(); err != nil {
			return err
		}
	}

	return validateExtensions(s.Extensions)
}

// keyword is a keyword that applies to one type: its name, whether it is set, and its value, if an error should show it.
type keyword struct {
	name  string
	set   bool
	value any
}

// valueOf is what p points to, or nil.
func valueOf[T any](p *T) any {
	if p == nil {
		return nil
	}

	return *p
}

// onlyFor reports the first keyword set where the keywords of typ do not apply.
func (s *Schema) onlyFor(applies bool, typ string, kws ...keyword) error {
	if applies {
		return nil
	}

	for _, kw := range kws {
		if !kw.set {
			continue
		}

		msg := fmt.Sprintf("only valid for %s type, got %s", typ, s.typeOrNone())
		if kw.value == nil {
			return &errpath.ErrField{Field: kw.name, Err: &errpath.ErrInvalid[string]{Message: msg}}
		}

		return &errpath.ErrField{Field: kw.name, Err: &errpath.ErrInvalid[any]{Value: kw.value, Message: msg}}
	}

	return nil
}

// validateFormat checks that the format is known and suits the schema's type.
func (s *Schema) validateFormat() error {
	if s.Format == "" {
		return nil
	}

	if err := s.Format.Validate(); err != nil {
		return &errpath.ErrField{Field: "format", Err: err}
	}

	var types string

	switch s.Format {
	case FormatInt32, FormatInt64, FormatUint, FormatUint32, FormatUint64:
		if !s.allows(TypeInteger) {
			types = "integer"
		}
	case FormatFloat, FormatDouble:
		if !s.allows(TypeNumber) {
			types = "number"
		}
	case FormatEmail, FormatPassword, FormatUUID, FormatURI, FormatURIRef, FormatZipCode, FormatIPv4, FormatIPv6,
		FormatByte, FormatBinary:
		if !s.allows(TypeString) {
			types = "string"
		}
	case FormatDuration, FormatDate, FormatDateTime:
		if !s.allows(TypeInteger) && !s.allows(TypeString) {
			types = "integer or string"
		}
	default:
		return fmt.Errorf("unimplemented format: %s", s.Format)
	}

	if types == "" {
		return nil
	}

	return &errpath.ErrField{Field: "format", Err: &errpath.ErrInvalid[Format]{
		Value:   s.Format,
		Message: fmt.Sprintf("only valid for %s type, got %s", types, s.typeOrNone()),
	}}
}

// validateString checks the keywords of strings.
func (s *Schema) validateString() error {
	if err := s.onlyFor(s.allows(TypeString), "string",
		keyword{"minLength", s.MinLength != 0, nil},
		keyword{"maxLength", s.MaxLength != nil, nil},
		keyword{"pattern", s.Pattern != nil, nil},
		keyword{"contentMediaType", s.ContentMediaType != "", nil},
		keyword{"contentEncoding", s.ContentEncoding != "", nil},
	); err != nil {
		return err
	}

	if s.MaxLength != nil && s.MinLength > *s.MaxLength {
		return &errpath.ErrField{Field: "minLength", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinLength,
			Message: fmt.Sprintf("minLength is greater than maxLength (%d > %d)", s.MinLength, *s.MaxLength),
		}}
	}

	return nil
}

// validateComposition validates the schemas of allOf, oneOf, anyOf and not.
func (s *Schema) validateComposition() error {
	for _, l := range []struct {
		field   string
		schemas SchemaList
	}{{"allOf", s.AllOf}, {"oneOf", s.OneOf}, {"anyOf", s.AnyOf}} {
		for i, v := range l.schemas {
			if err := v.Validate(); err != nil {
				return &errpath.ErrField{Field: l.field, Err: &errpath.ErrIndex{Index: i, Err: err}}
			}
		}
	}

	if s.Not != nil {
		if err := s.Not.Validate(); err != nil {
			return &errpath.ErrField{Field: "not", Err: err}
		}
	}

	return nil
}

// validateNumber checks the keywords of numbers: whole numbers for an integer, bounds that leave some value.
func (s *Schema) validateNumber() error {
	bounds := []struct {
		field string
		value *float64
	}{
		{"minimum", s.Min},
		{"maximum", s.Max},
		{"exclusiveMinimum", s.ExclusiveMin},
		{"exclusiveMaximum", s.ExclusiveMax},
		{"multipleOf", s.MultipleOf},
	}

	if s.Type == TypeInteger {
		for _, b := range bounds {
			if b.value != nil && *b.value != float64(int(*b.value)) {
				return &errpath.ErrField{Field: b.field, Err: &errpath.ErrInvalid[float64]{
					Value:   *b.value,
					Message: "not an integer",
				}}
			}
		}
	}

	if s.allows(TypeNumber) || s.allows(TypeInteger) {
		if err := s.validateRange(); err != nil {
			return err
		}
	} else {
		kws := make([]keyword, len(bounds))
		for i, b := range bounds {
			kws[i] = keyword{b.field, b.value != nil, valueOf(b.value)}
		}

		if err := s.onlyFor(false, "number", kws...); err != nil {
			return err
		}
	}

	if s.MultipleOf != nil && *s.MultipleOf <= 0 {
		return &errpath.ErrField{Field: "multipleOf", Err: &errpath.ErrInvalid[float64]{
			Value:   *s.MultipleOf,
			Message: "must be greater than 0",
		}}
	}

	return nil
}

// validateRange checks that the lower bound of a number is below its upper bound.
func (s *Schema) validateRange() error {
	for _, r := range []struct {
		field, lowName, highName string
		low, high                *float64
		exclusive                bool
	}{
		{"minimum", "minimum", "maximum", s.Min, s.Max, false},
		{"minimum", "minimum", "exclusiveMaximum", s.Min, s.ExclusiveMax, true},
		{"exclusiveMinimum", "exclusiveMinimum", "maximum", s.ExclusiveMin, s.Max, true},
		{"exclusiveMinimum", "exclusiveMinimum", "exclusiveMaximum", s.ExclusiveMin, s.ExclusiveMax, true},
	} {
		if r.low == nil || r.high == nil {
			continue
		}

		if !r.exclusive && *r.low > *r.high {
			return &errpath.ErrField{Field: r.field, Err: &errpath.ErrInvalid[float64]{
				Value:   *r.low,
				Message: fmt.Sprintf("%s is greater than %s (%v > %v)", r.lowName, r.highName, *r.low, *r.high),
			}}
		}

		if r.exclusive && *r.low >= *r.high {
			return &errpath.ErrField{Field: r.field, Err: &errpath.ErrInvalid[float64]{
				Value:   *r.low,
				Message: fmt.Sprintf("%s is not less than %s (%v >= %v)", r.lowName, r.highName, *r.low, *r.high),
			}}
		}
	}

	return nil
}

// validateValueKinds checks that enum, const and the examples are values of the schema's type: per JSON Schema
// 2020-12 they can hold any JSON value.
func (s *Schema) validateValueKinds() error {
	if s.Type == "" {
		return nil
	}

	if err := s.valuesOfKind("enum", s.Enum); err != nil {
		return err
	}

	if err := s.valueOfKind("const", s.Const); err != nil {
		return err
	}

	if err := s.valueOfKind("example", s.Example); err != nil {
		return err
	}

	return s.valuesOfKind("examples", s.Examples)
}

// valuesOfKind reports the first of values that is not of the schema's type.
func (s *Schema) valuesOfKind(field string, values []jsontext.Value) error {
	for i, v := range values {
		if !s.allowsKindOf(v) {
			return &errpath.ErrField{Field: field, Err: &errpath.ErrIndex{Index: i, Err: s.wrongKind(v)}}
		}
	}

	return nil
}

// valueOfKind reports v if it is set and not of the schema's type.
func (s *Schema) valueOfKind(field string, v jsontext.Value) error {
	if v == nil || s.allowsKindOf(v) {
		return nil
	}

	return &errpath.ErrField{Field: field, Err: s.wrongKind(v)}
}

func (s *Schema) wrongKind(v jsontext.Value) error {
	return &errpath.ErrInvalid[any]{Value: jsonDisplayValue(v), Message: fmt.Sprintf("must be a %s value", s.Type)}
}

// validateArray checks the keywords of arrays.
func (s *Schema) validateArray() error {
	if !s.allows(TypeArray) {
		return s.onlyFor(false, "array",
			keyword{"minItems", s.MinItems != 0, s.MinItems},
			keyword{"maxItems", s.MaxItems != nil, valueOf(s.MaxItems)},
			keyword{"uniqueItems", s.UniqueItems, true},
			keyword{"prefixItems", len(s.PrefixItems) != 0, nil},
			keyword{"items", s.Items != nil, nil},
		)
	}

	if s.MaxItems != nil && s.MinItems > *s.MaxItems {
		return &errpath.ErrField{Field: "minItems", Err: &errpath.ErrInvalid[uint]{
			Value:   s.MinItems,
			Message: fmt.Sprintf("minItems is greater than maxItems (%d > %d)", s.MinItems, *s.MaxItems),
		}}
	}

	for i, v := range s.PrefixItems {
		if err := v.Validate(); err != nil {
			return &errpath.ErrField{Field: "prefixItems", Err: &errpath.ErrIndex{Index: i, Err: err}}
		}
	}

	if s.Items != nil {
		if err := s.Items.Validate(); err != nil {
			return &errpath.ErrField{Field: "items", Err: err}
		}
	}

	return nil
}

// validateObject checks the keywords of objects.
func (s *Schema) validateObject() error {
	if !s.allows(TypeObject) {
		return s.onlyFor(false, "object",
			keyword{"properties", s.Properties != nil, nil},
			keyword{"required", s.Required != nil, nil},
			keyword{"additionalProperties", s.AdditionalProperties != nil, nil},
			keyword{"maxProperties", s.MaxProperties != nil, valueOf(s.MaxProperties)},
			keyword{"propertyNames", s.PropertyNames != nil, nil},
		)
	}

	if err := s.Properties.Validate(); err != nil {
		return &errpath.ErrField{Field: "properties", Err: err}
	}

	// without a type, required is a constraint on whatever object properties hold
	if s.Type != "" {
		for i, r := range s.Required {
			if _, ok := s.Properties[r]; !ok {
				return &errpath.ErrField{Field: "required", Err: &errpath.ErrIndex{Index: i, Err: &errpath.ErrInvalid[string]{
					Value:   r,
					Message: "property does not exist",
				}}}
			}
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

	return nil
}

// validateDiscriminator checks where the discriminator may be: beside a composition, or on a parent another schema
// extends through allOf.
func (s *Schema) validateDiscriminator(extended bool) error {
	if s.Discriminator == nil {
		return nil
	}

	if len(s.OneOf) == 0 && len(s.AnyOf) == 0 && len(s.AllOf) == 0 && !extended {
		return &errpath.ErrField{Field: "discriminator", Err: &errpath.ErrInvalid[string]{
			Message: "only valid with oneOf, anyOf or allOf, or on a component schema another extends through allOf",
		}}
	}

	if err := s.Discriminator.Validate(); err != nil {
		return &errpath.ErrField{Field: "discriminator", Err: err}
	}

	return nil
}

// validateDefault checks that the default is a value of the schema's type and, with an enum, one of its values.
func (s *Schema) validateDefault() error {
	if len(s.Default) == 0 {
		return nil
	}

	if !enumKindMatchesType(s.Default, s.Type) {
		return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[any]{
			Value:   jsonDisplayValue(s.Default),
			Message: fmt.Sprintf("does not match schema type, got %s", s.typeOrNone()),
		}}
	}

	if len(s.Enum) == 0 || slices.ContainsFunc(s.Enum, func(v jsontext.Value) bool { return bytes.Equal(v, s.Default) }) {
		return nil
	}

	parts := make([]string, len(s.Enum))
	for i, v := range s.Enum {
		parts[i] = v.String()
	}

	return &errpath.ErrField{Field: "default", Err: &errpath.ErrInvalid[any]{
		Value:   jsonDisplayValue(s.Default),
		Message: fmt.Sprintf("is not one of the enums ([%s])", strings.Join(parts, " ")),
	}}
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
