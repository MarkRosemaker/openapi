package openapi

import (
	"net/url"
)

// The Info object provides metadata about the API. The metadata MAY be used by the clients if needed, and MAY be presented in editing or documentation generation tools for convenience.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#info-object
type Info struct {
	// REQUIRED. The title of the API.
	Title string `json:"title" yaml:"title"`
	// A short summary of the API.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A description of the API. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A URL to the Terms of Service for the API.
	TermsOfService *url.URL `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	// The license information for the exposed API.
	License *License `json:"license,omitempty" yaml:"license,omitempty"`
	// REQUIRED. The version of the OpenAPI document (which is distinct from the OpenAPI Specification version or the API implementation version).
	Version string `json:"version" yaml:"version"`
	// The object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (i *Info) Validate() error {
	if i.Title == "" {
		return &ErrField{Field: "title", Err: &ErrRequired{}}
	}

	// NOTE: The version *here* can be any string, but the version in the OpenAPI document must be a valid semantic version.
	if i.Version == "" {
		return &ErrField{Field: "version", Err: &ErrRequired{}}
	}

	// assume that the scheme is https and add it if it is missing
	fixScheme(i.TermsOfService)

	if i.Contact != nil {
		if err := i.Contact.Validate(); err != nil {
			return &ErrField{Field: "contact", Err: err}
		}
	}

	if i.License != nil {
		if err := i.License.Validate(); err != nil {
			return &ErrField{Field: "license", Err: err}
		}
	}

	return validateExtensions(i.Extensions)
}

// fixScheme ensures that the URL has a scheme and that it is valid.
// If the URL is nil, it is a no-op
func fixScheme(u *url.URL) {
	if u == nil {
		return
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}
}
