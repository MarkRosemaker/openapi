package openapi

import (
	"net/url"

	"github.com/go-api-libs/types"
)

// Contact information for the exposed API.
// ([Source])
//
// [Source]: https://spec.openapis.org/oas/v3.1.0#contact-object
type Contact struct {
	// The identifying name of the contact person/organization.
	Name string `json:"name,omitempty,strictcase" yaml:"name,omitempty"`
	// The URL pointing to the contact information. MUST be in the format of a URL.
	URL *url.URL `json:"url,omitempty,strictcase" yaml:"url,omitempty"`
	// The email address of the contact person/organization. This MUST be in the form of an email address.
	Email types.Email `json:"email,omitempty,strictcase" yaml:"email,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the contact for consistency.
func (c *Contact) Validate() error {
	// assume that the scheme is https and add it if it is missing
	fixScheme(c.URL)

	if c.Email != "" {
		if err := c.Email.Validate(); err != nil {
			return &ErrField{Field: "email", Err: err}
		}
	}

	if err := validateExtensions(c.Extensions); err != nil {
		return err
	}

	return nil
}
