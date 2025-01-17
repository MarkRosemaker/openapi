package openapi

import (
	"net/url"

	"github.com/MarkRosemaker/errpath"
)

// OAuthFlowImplicit allows configuration details for the OAuth Implicit flow.
type OAuthFlowImplicit struct {
	// REQUIRED. The authorization URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	AuthorizationURL *url.URL `json:"authorizationUrl" yaml:"authorizationUrl"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	RefreshURL *url.URL `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	// REQUIRED. The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it. The map MAY be empty.
	Scopes MapOfStrings `json:"scopes" yaml:"scopes"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (f *OAuthFlowImplicit) Validate() error {
	if f.AuthorizationURL == nil {
		return &errpath.ErrField{Field: "authorizationUrl", Err: &errpath.ErrRequired{}}
	}

	if f.Scopes == nil {
		return &errpath.ErrField{Field: "scopes", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(f.Extensions)
}

// OAuthFlowPassword allows configuration details for the OAuth Resource Owner Password flow.
type OAuthFlowPassword struct {
	// REQUIRED. The token URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	TokenURL *url.URL `json:"tokenUrl" yaml:"tokenUrl"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	RefreshURL *url.URL `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	// REQUIRED. The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it. The map MAY be empty.
	Scopes MapOfStrings `json:"scopes" yaml:"scopes"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (f *OAuthFlowPassword) Validate() error {
	if f.TokenURL == nil {
		return &errpath.ErrField{Field: "tokenUrl", Err: &errpath.ErrRequired{}}
	}

	if f.Scopes == nil {
		return &errpath.ErrField{Field: "scopes", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(f.Extensions)
}

// OAuthFlowClientCredentials allows configuration details for the OAuth Client Credentials flow.
type OAuthFlowClientCredentials = OAuthFlowPassword

// OAuthFlowAuthorizationCode allows configuration details for the OAuth Authorization Code flow.
type OAuthFlowAuthorizationCode struct {
	// REQUIRED. The authorization URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	AuthorizationURL *url.URL `json:"authorizationUrl" yaml:"authorizationUrl"`
	// REQUIRED. The token URL to be used for this flow. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	TokenURL *url.URL `json:"tokenUrl" yaml:"tokenUrl"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of a URL. The OAuth2 standard requires the use of TLS.
	RefreshURL *url.URL `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	// REQUIRED. The available scopes for the OAuth2 security scheme. A map between the scope name and a short description for it. The map MAY be empty.
	Scopes MapOfStrings `json:"scopes" yaml:"scopes"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (f *OAuthFlowAuthorizationCode) Validate() error {
	if f.AuthorizationURL == nil {
		return &errpath.ErrField{Field: "authorizationUrl", Err: &errpath.ErrRequired{}}
	}

	if f.TokenURL == nil {
		return &errpath.ErrField{Field: "tokenUrl", Err: &errpath.ErrRequired{}}
	}

	if f.Scopes == nil {
		return &errpath.ErrField{Field: "scopes", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(f.Extensions)
}
