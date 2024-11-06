package openapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-json-experiment/json/jsontext"
)

// Parameter describes a single operation parameter.
//
// A unique parameter is defined by a combination of a name and location.
//
// The rules for serialization of the parameter are specified in one of two ways:
//  1. For simpler scenarios, a `schema` and `style` can describe the structure and syntax of the parameter.
//  2. For more complex scenarios, the content property can define the media type and schema of the parameter.
//
// A parameter MUST contain either a schema property, or a content property, but not both. When example or examples are provided in conjunction with the schema object, the example MUST follow the prescribed serialization strategy for the parameter.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#parameter-object
type Parameter struct {
	//	REQUIRED. The name of the parameter. Parameter names are *case sensitive*.
	//
	// - If `in` is "path", the `name` field MUST correspond to a template expression occurring within the path field in the Paths Object. See Path Templating for further information. TODO
	// - If `in` is `"header"` and the `name` field is `"Accept"`, `"Content-Type"` or `"Authorization"`, the parameter definition SHALL be ignored.
	// - For all other cases, the `name` corresponds to the parameter name used by the `in` property.
	Name string `json:"name" yaml:"name"`
	// REQUIRED. The location of the parameter. Possible values are `"query"`, `"header"`, `"path"` or `"cookie"`.
	In ParameterLocation `json:"in" yaml:"in"`
	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Determines whether this parameter is mandatory. If the parameter location is `"path"`, this property is **REQUIRED** and its value MUST be `true`. Otherwise, the property MAY be included and its default value is `false`.
	Required bool `json:"required,omitempty,omitzero" yaml:"required,omitempty"`
	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is `false`.
	Deprecated bool `json:"deprecated,omitempty,omitzero" yaml:"deprecated,omitempty"`
	// Sets the ability to pass empty-valued parameters. This is valid only for `query` parameters and allows sending a parameter with an empty value. Default value is `false`. If `style` is used, and if behavior is `n/a` (cannot be serialized), the value of `allowEmptyValue` SHALL be ignored. Use of this property is NOT RECOMMENDED, as it is likely to be removed in a later revision.
	AllowEmptyValue bool `json:"allowEmptyValue,omitempty,omitzero" yaml:"allowEmptyValue,omitempty"`

	// Serialization via Schema

	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values (based on value of `in`): for `query` - `form`; for `path` - `simple`; for `header` - `simple`; for `cookie` - `form`.
	Style ParameterStyle `json:"style,omitempty" yaml:"style,omitempty"`
	// When this is true, parameter values of type `array` or `object` generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When `style` is `form`, the default value is `true`. For all other styles, the default value is `false`.
	Explode bool `json:"explode,omitempty,omitzero" yaml:"explode,omitempty"`
	// Determines whether the parameter value SHOULD allow reserved characters, as defined by RFC3986 `:/?#[]@!$&'()*+,;=` to be included without percent-encoding. This property only applies to parameters with an `in` value of `query`. The default value is `false`.
	AllowReserved bool `json:"allowReserved,omitempty,omitzero" yaml:"allowReserved,omitempty"`
	// The schema defining the type used for the parameter.
	Schema *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The `example` field is mutually exclusive of the `examples` field. Furthermore, if referencing a `schema` that contains an example, the `example` value SHALL _override_ the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example jsontext.Value `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The `examples` field is mutually exclusive of the `example` field. Furthermore, if referencing a `schema` that contains an example, the `examples` value SHALL _override_ the example provided by the schema.
	Examples Examples `json:"examples,omitempty" yaml:"examples,omitempty"`

	// Serialization via Content

	// A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry.
	Content Content `json:"content,omitempty" yaml:"content,omitempty"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the parameter for correctness.
func (p *Parameter) Validate() error {
	if p.Name == "" {
		return &ErrField{Field: "name", Err: &ErrRequired{}}
	}

	if err := p.In.Validate(); err != nil {
		return &ErrField{Field: "in", Err: err}
	}

	if p.In == ParameterLocationPath && !p.Required {
		return &ErrField{Field: "required", Err: &ErrInvalid[bool]{
			Value:   false,
			Message: "must be true for path parameters",
		}}
	}

	if p.In != ParameterLocationQuery {
		if p.AllowEmptyValue {
			return &ErrField{Field: "allowEmptyValue", Err: &ErrInvalid[bool]{
				Value:   true,
				Message: fmt.Sprintf("can only be true for query parameters, got %q", p.In),
			}}
		}

		if p.AllowReserved {
			return &ErrField{Field: "allowReserved", Err: &ErrInvalid[bool]{
				Value:   true,
				Message: fmt.Sprintf("only applies to query parameters, got %q", p.In),
			}}
		}
	}

	p.Description = strings.TrimSpace(p.Description)

	if p.Schema != nil {
		// A parameter MUST contain either a `schema` property, or a `content` property, but not both.
		if p.Content != nil {
			return errors.New("schema and content are mutually exclusive")
		}

		if err := p.Schema.Validate(); err != nil {
			return &ErrField{Field: "schema", Err: err}
		}
	} else {
		if p.Content == nil {
			return errors.New("schema or content is required")
		}

		if len(p.Content) != 1 {
			return &ErrField{Field: "content", Err: &ErrInvalid[string]{
				Message: fmt.Sprintf("must contain exactly one entry, got %d", len(p.Content)),
			}}
		}

		if err := p.Content.Validate(); err != nil {
			return &ErrField{Field: "content", Err: err}
		}
	}

	if p.Style != "" {
		if err := p.Style.Validate(); err != nil {
			return &ErrField{Field: "style", Err: err}
		}
	}

	if p.Explode {
		if p.Schema == nil {
			return &ErrField{Field: "explode", Err: &ErrInvalid[bool]{
				Value:   true,
				Message: "property has no effect when schema is not present",
			}}
		}

		if p.Schema.Type != TypeArray && p.Schema.Type != TypeObject {
			return &ErrField{Field: "explode", Err: &ErrInvalid[bool]{
				Value:   true,
				Message: fmt.Sprintf("property has no effect when schema type is not array or object, got %q", p.Schema.Type),
			}}
		}
	}

	if p.Example != nil && p.Examples != nil {
		return errors.New("example and examples are mutually exclusive")
	}

	if err := p.Examples.Validate(); err != nil {
		return &ErrField{Field: "examples", Err: err}
	}

	return validateExtensions(p.Extensions)
}
