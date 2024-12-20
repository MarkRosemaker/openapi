package openapi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json/jsontext"
)

// Header represents a single header parameter to be included in the operation.
type Header struct {
	// A brief description of the parameter. This could contain examples of use. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Determines whether this parameter is mandatory. If the parameter location is `"path"`, this property is **REQUIRED** and its value MUST be `true`. Otherwise, the property MAY be included and its default value is `false`.
	Required bool `json:"required,omitempty,omitzero" yaml:"required,omitempty"`
	// Specifies that a parameter is deprecated and SHOULD be transitioned out of usage. Default value is `false`.
	Deprecated bool `json:"deprecated,omitempty,omitzero" yaml:"deprecated,omitempty"`
	// The schema defining the type used for the parameter.
	Schema *Schema `json:"schema,omitempty" yaml:"schema,omitempty"`
	// Describes how the parameter value will be serialized depending on the type of the parameter value. Default values is `simple`.
	Style ParameterStyle `json:"style,omitempty" yaml:"style,omitempty"`
	// When this is true, parameter values of type `array` or `object` generate separate parameters for each value of the array or key-value pair of the map. For other types of parameters this property has no effect. When `style` is `form`, the default value is `true`. For all other styles, the default value is `false`.
	Explode bool `json:"explode,omitempty,omitzero" yaml:"explode,omitempty"`
	// Example of the parameter's potential value. The example SHOULD match the specified schema and encoding properties if present. The `example` field is mutually exclusive of the `examples` field. Furthermore, if referencing a `schema` that contains an example, the `example` value SHALL _override_ the example provided by the schema. To represent examples of media types that cannot naturally be represented in JSON or YAML, a string value can contain the example with escaping where necessary.
	Example jsontext.Value `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples of the parameter's potential value. Each example SHOULD contain a value in the correct format as specified in the parameter encoding. The `examples` field is mutually exclusive of the `example` field. Furthermore, if referencing a `schema` that contains an example, the `examples` value SHALL _override_ the example provided by the schema.
	Examples Examples `json:"examples,omitempty" yaml:"examples,omitempty"`
	// A map containing the representations for the parameter. The key is the media type and the value describes it. The map MUST only contain one entry.
	Content Content `json:"content,omitempty" yaml:"content,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`
}

func (h *Header) Validate() error {
	h.Description = strings.TrimSpace(h.Description)

	if h.Schema != nil {
		// A parameter MUST contain either a `schema` property, or a `content` property, but not both.
		if h.Content != nil {
			return errors.New("schema and content are mutually exclusive")
		}

		if err := h.Schema.Validate(); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	} else {
		if h.Content == nil {
			return errors.New("schema or content is required")
		}

		if len(h.Content) != 1 {
			return &errpath.ErrField{Field: "content", Err: &errpath.ErrInvalid[string]{
				Message: fmt.Sprintf("must contain exactly one entry, got %d", len(h.Content)),
			}}
		}

		if err := h.Content.Validate(); err != nil {
			return &errpath.ErrField{Field: "content", Err: err}
		}
	}

	if h.Style != "" {
		if err := h.Style.Validate(); err != nil {
			return &errpath.ErrField{Field: "style", Err: err}
		}
	}

	if h.Explode {
		if h.Schema == nil {
			return &errpath.ErrField{Field: "explode", Err: &errpath.ErrInvalid[bool]{
				Value:   true,
				Message: "property has no effect when schema is not present",
			}}
		}

		if h.Schema.Type != TypeArray && h.Schema.Type != TypeObject {
			return &errpath.ErrField{Field: "explode", Err: &errpath.ErrInvalid[bool]{
				Value:   true,
				Message: fmt.Sprintf("property has no effect when schema type is not array or object, got %q", h.Schema.Type),
			}}
		}
	}

	if h.Example != nil && h.Examples != nil {
		return errors.New("example and examples are mutually exclusive")
	}

	if err := h.Examples.Validate(); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	// When `example` or `examples` are provided in conjunction with the `schema` object, the example MUST follow the prescribed serialization strategy for the parameter.
	// TODO

	return validateExtensions(h.Extensions)
}

func (l *loader) collectHeaderRef(h *HeaderRef, ref ref) {
	if h.Value != nil {
		l.collectHeader(h.Value, ref)
	}
}

func (l *loader) collectHeader(h *Header, ref ref) {
	l.headers[ref.String()] = h
}

func (l *loader) resolveHeaderRef(h *HeaderRef) error {
	return resolveRef(h, l.headers, l.resolveHeader)
}

func (l *loader) resolveHeader(h *Header) error {
	return nil // TODO
}
