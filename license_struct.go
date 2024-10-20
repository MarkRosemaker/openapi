package openapi

import (
	"errors"
	"net/url"
)

// License information for the exposed API.
// ([Documentation])
//
// [Documentation]: https://spec.openapis.org/oas/v3.1.0#license-object
type License struct {
	// REQUIRED. The license name used for the API.
	Name string `json:"name,strictcase" yaml:"name"`
	// An SPDX license expression for the API. The identifier field is mutually exclusive of the url field.
	// See: https://spdx.org/licenses/
	Identifier string `json:"identifier,omitempty,strictcase" yaml:"identifier,omitempty"`
	// A URL to the license used for the API. This MUST be in the form of a URL. The url field is mutually exclusive of the identifier field.
	URL *url.URL `json:"url,omitempty,strictcase" yaml:"url,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (l *License) Validate() error {
	if l.Name == "" {
		return &ErrRequired{Target: "name"}
	}

	if l.URL != nil && l.Identifier != "" {
		return errors.New("url and identifier are mutually exclusive")
	}

	if err := validateExtensions(l.Extensions); err != nil {
		return err
	}

	return nil
}
