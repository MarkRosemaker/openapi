package openapi

import (
	"errors"

	"github.com/MarkRosemaker/errpath"
	"github.com/go-json-experiment/json/jsontext"
)

// Each Media Type Object provides schema and examples for the media type identified by its key.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#media-type-object
type MediaType struct {
	// The schema defining the content of the request, response, or parameter.
	Schema *SchemaRef `json:"schema,omitempty" yaml:"schema,omitempty"`
	// Example of the media type.
	// The example object SHOULD be in the correct format as specified by the media type.
	// The `example` field is mutually exclusive of the `examples` field.  Furthermore, if referencing a `schema` which contains an example, the `example` value SHALL _override_ the example provided by the schema.
	Example jsontext.Value `json:"example,omitempty" yaml:"example,omitempty"`
	// Examples of the media type.
	// Each example object SHOULD match the media type and specified schema if present.
	// The `examples` field is mutually exclusive of the `example` field.  Furthermore, if referencing a `schema` which contains an example, the `examples` value SHALL _override_ the example provided by the schema.
	Examples Examples `json:"examples,omitempty" yaml:"examples,omitempty"`
	// A map between a property name and its encoding information. The key, being the property name, MUST exist in the schema as a property. The encoding object SHALL only apply to `requestBody` objects when the media type is `multipart` or `application/x-www-form-urlencoded`.
	Encoding Encodings `json:"encoding,omitempty" yaml:"encoding,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`

	// an index to the original location of this object
	idx int
}

func getIndexMediaType(mt *MediaType) int      { return mt.idx }
func setIndexMediaType(mt *MediaType, idx int) { mt.idx = idx }

func (mt *MediaType) Validate() error {
	if mt.Schema != nil {
		if err := mt.Schema.Validate(); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	}

	if mt.Example != nil && mt.Examples != nil {
		return errors.New("example and examples are mutually exclusive")
	}

	if err := mt.Examples.Validate(); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := mt.Encoding.Validate(); err != nil {
		return &errpath.ErrField{Field: "encoding", Err: err}
	}

	return validateExtensions(mt.Extensions)
}
