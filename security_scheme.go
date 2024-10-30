package openapi

import (
	"fmt"
	"net/url"
	"strings"
)

//
// []: https://spec.openapis.org/oas/v3.1.0#security-scheme-object
type SecurityScheme struct {
	Type        SecuritySchemeType `json:"type" yaml:"type"`
	Description string             `json:"description,omitempty" yaml:"description,omitempty"`
	// The name of the API key.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The location of the API key.
	In SecuritySchemeIn `json:"in,omitempty" yaml:"in,omitempty"`
	// The HTTP scheme to use, e.g. "bearer".
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	// The format of the bearer token, e.g. "jwt".
	BearerFormat string `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	// The flows supported by the OAuth2 security scheme.
	Flows *OAuthFlows `json:"flows,omitempty" yaml:"flows,omitempty"`
	// The OpenID connect URL to use
	OpenIdConnectURL *url.URL `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

const SecuritySchemeBearer = "bearer"

func (s *SecurityScheme) Validate() error {
	if err := s.Type.Validate(); err != nil {
		return &ErrField{Field: "type", Err: err}
	}

	s.Description = strings.TrimSpace(s.Description)

	switch s.Type {
	case SecuritySchemeTypeAPIKey:
		if s.Name == "" {
			return &ErrField{Field: "name", Err: &ErrRequired{}}
		}

		if s.In == "" {
			return &ErrField{Field: "in", Err: &ErrRequired{}}
		}

		if err := s.In.Validate(); err != nil {
			return &ErrField{Field: "in", Err: err}
		}
	case SecuritySchemeTypeHTTP:
		if s.Scheme == "" {
			return &ErrField{Field: "scheme", Err: &ErrRequired{}}
		}

		if SecuritySchemeBearer == strings.ToLower(s.Scheme) {
			s.Scheme = SecuritySchemeBearer // unify

			if s.BearerFormat == "" {
				return &ErrField{Field: "bearerFormat", Err: &ErrRequired{}}
			}
		}
	case SecuritySchemeTypeMutualTLS: // nothing to do
	case SecuritySchemeTypeOAuth2:
		if s.Flows == nil {
			return &ErrField{Field: "flows", Err: &ErrRequired{}}
		}

		if err := s.Flows.Validate(); err != nil {
			return &ErrField{Field: "flows", Err: err}
		}
	case SecuritySchemeTypeOpenIDConnect:
		if s.OpenIdConnectURL == nil {
			return &ErrField{Field: "openIdConnectUrl", Err: &ErrRequired{}}
		}
	default:
		return fmt.Errorf("unimplemented type %q", s.Type)
	}

	if err := validateExtensions(s.Extensions); err != nil {
		return err
	}

	return nil
}
