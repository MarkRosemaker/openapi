package asyncapi

import (
	"net/url"

	"github.com/MarkRosemaker/errpath"
)

// License information for the exposed API.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#licenseObject
type License struct {
	// REQUIRED. The license name used for the API.
	Name string `json:"name" yaml:"name"`
	// A URL to the license used for the API. This MUST be in the form of an absolute URL.
	URL *url.URL `json:"url,omitempty" yaml:"url,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the license for correctness.
func (l *License) Validate() error {
	if l.Name == "" {
		return &errpath.ErrField{Field: "name", Err: &errpath.ErrRequired{}}
	}

	// assume that the scheme is https and add it if it is missing
	fixScheme(l.URL)

	return validateExtensions(l.Extensions)
}
