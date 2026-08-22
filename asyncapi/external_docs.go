package asyncapi

import (
	"net/url"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// ExternalDocs allows referencing an external resource for extended documentation.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#externalDocumentationObject
type ExternalDocs struct {
	// A short description of the target documentation. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// REQUIRED. The URL for the target documentation. This MUST be in the form of an absolute URL.
	URL *url.URL `json:"url,omitempty" yaml:"url,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the external documentation for consistency.
func (ed *ExternalDocs) Validate() error {
	if ed.URL == nil {
		return &errpath.ErrField{Field: "url", Err: &errpath.ErrRequired{}}
	}

	if err := validateURL(ed.URL); err != nil {
		return &errpath.ErrField{Field: "url", Err: err}
	}

	ed.Description = strings.TrimSpace(ed.Description)

	return validateExtensions(ed.Extensions)
}

func (l *loader) collectExternalDocsRef(d *ExternalDocsRef, ref ref) {
	collectRef(l, d, l.externalDocs, ref)
}

func (l *loader) resolveExternalDocsRef(d *ExternalDocsRef) error {
	return resolveRef(d, l.externalDocs, nil)
}
