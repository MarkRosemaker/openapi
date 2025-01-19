package openapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Adds metadata to a single tag that is used by the [Operation] object.
// It is not mandatory to have a Tag object per tag defined in the Operation object instances.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#tag-object
type Tag struct {
	// The name of the tag.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// A description for the tag. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate validates the tag.
func (t *Tag) Validate() error {
	if t.Name == "" {
		return &errpath.ErrField{Field: "name", Err: &errpath.ErrRequired{}}
	}

	t.Description = strings.TrimSpace(t.Description)

	if t.ExternalDocs != nil {
		if err := t.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return validateExtensions(t.Extensions)
}
