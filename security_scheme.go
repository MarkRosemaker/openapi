package openapi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/MarkRosemaker/errpath"
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
	Extensions Extensions `json:",inline" yaml:"-"`
}

const (
	SecuritySchemeBearer = "bearer"
	SecuritySchemeBasic  = "basic"
)

func (s *SecurityScheme) Validate() error {
	if err := s.Type.Validate(); err != nil {
		return &errpath.ErrField{Field: "type", Err: err}
	}

	s.Description = strings.TrimSpace(s.Description)

	switch s.Type {
	case SecuritySchemeTypeAPIKey:
		if s.Name == "" {
			return &errpath.ErrField{Field: "name", Err: &errpath.ErrRequired{}}
		}

		if s.In == "" {
			return &errpath.ErrField{Field: "in", Err: &errpath.ErrRequired{}}
		}

		if err := s.In.Validate(); err != nil {
			return &errpath.ErrField{Field: "in", Err: err}
		}
	case SecuritySchemeTypeHTTP:
		if s.Scheme == "" {
			return &errpath.ErrField{Field: "scheme", Err: &errpath.ErrRequired{}}
		}

		if SecuritySchemeBearer == strings.ToLower(s.Scheme) {
			s.Scheme = SecuritySchemeBearer // unify

			if s.BearerFormat == "" {
				return &errpath.ErrField{Field: "bearerFormat", Err: &errpath.ErrRequired{}}
			}
		}
	case SecuritySchemeTypeMutualTLS: // nothing to do
	case SecuritySchemeTypeOAuth2:
		if s.Flows == nil {
			return &errpath.ErrField{Field: "flows", Err: &errpath.ErrRequired{}}
		}

		if err := s.Flows.Validate(); err != nil {
			return &errpath.ErrField{Field: "flows", Err: err}
		}
	case SecuritySchemeTypeOpenIDConnect:
		if s.OpenIdConnectURL == nil {
			return &errpath.ErrField{Field: "openIdConnectUrl", Err: &errpath.ErrRequired{}}
		}
	default:
		return fmt.Errorf("unimplemented type %q", s.Type)
	}

	return validateExtensions(s.Extensions)
}

func (l *loader) collectSecuritySchemeRef(r *SecuritySchemeRef, ref ref) {
	if r.Value != nil {
		l.collectSecurityScheme(r.Value, ref)
	}
}

func (l *loader) collectSecurityScheme(s *SecurityScheme, ref ref) {
	l.securitySchemes[ref.String()] = s
}

func (l *loader) resolveSecuritySchemeRef(r *SecuritySchemeRef) error {
	return resolveRef(r, l.securitySchemes, nil)
}
