package openapi

import (
	"net/url"

	"github.com/MarkRosemaker/errpath"
)

// ExternalDocs allows referencing an external resource for extended documentation.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#external-documentation-object
type ExternalDocs struct {
	// A description of the target documentation. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// REQUIRED. The URL for the target documentation. This MUST be in the form of a URL.
	URL *url.URL `json:"url,omitempty" yaml:"url,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the external documentation for consistency.
func (ed *ExternalDocs) Validate() error {
	if ed.URL == nil {
		return &errpath.ErrField{Field: "url", Err: &errpath.ErrRequired{}}
	}

	// assume that the scheme is https and add it if it is missing
	fixScheme(ed.URL)

	return nil
}
