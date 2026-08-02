package asyncapi

import (
	"net/url"

	"github.com/MarkRosemaker/errpath"
)

// OAuthFlowImplicit holds the configuration details for the OAuth Implicit flow.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#oauthFlowObject
type OAuthFlowImplicit struct {
	// REQUIRED. The authorization URL to be used for this flow. This MUST be in the form of an absolute URL.
	AuthorizationURL *url.URL `json:"authorizationUrl" yaml:"authorizationUrl"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of an absolute URL.
	RefreshURL *url.URL `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	// REQUIRED. The available scopes for the OAuth2 security scheme.
	// A map between the scope name and a short description for it.
	AvailableScopes MapOfStrings `json:"availableScopes" yaml:"availableScopes"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the OAuth flow for correctness.
func (f *OAuthFlowImplicit) Validate() error {
	if f.AuthorizationURL == nil {
		return &errpath.ErrField{Field: "authorizationUrl", Err: &errpath.ErrRequired{}}
	}

	if f.AvailableScopes == nil {
		return &errpath.ErrField{Field: "availableScopes", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(f.Extensions)
}

// OAuthFlowPassword holds the configuration details for the OAuth Resource Owner Protected Credentials flow.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#oauthFlowObject
type OAuthFlowPassword struct {
	// REQUIRED. The token URL to be used for this flow. This MUST be in the form of an absolute URL.
	TokenURL *url.URL `json:"tokenUrl" yaml:"tokenUrl"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of an absolute URL.
	RefreshURL *url.URL `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	// REQUIRED. The available scopes for the OAuth2 security scheme.
	// A map between the scope name and a short description for it.
	AvailableScopes MapOfStrings `json:"availableScopes" yaml:"availableScopes"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the OAuth flow for correctness.
func (f *OAuthFlowPassword) Validate() error {
	if f.TokenURL == nil {
		return &errpath.ErrField{Field: "tokenUrl", Err: &errpath.ErrRequired{}}
	}

	if f.AvailableScopes == nil {
		return &errpath.ErrField{Field: "availableScopes", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(f.Extensions)
}

// OAuthFlowClientCredentials holds the configuration details for the OAuth Client Credentials flow.
type OAuthFlowClientCredentials = OAuthFlowPassword

// OAuthFlowAuthorizationCode holds the configuration details for the OAuth Authorization Code flow.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#oauthFlowObject
type OAuthFlowAuthorizationCode struct {
	// REQUIRED. The authorization URL to be used for this flow. This MUST be in the form of an absolute URL.
	AuthorizationURL *url.URL `json:"authorizationUrl" yaml:"authorizationUrl"`
	// REQUIRED. The token URL to be used for this flow. This MUST be in the form of an absolute URL.
	TokenURL *url.URL `json:"tokenUrl" yaml:"tokenUrl"`
	// The URL to be used for obtaining refresh tokens. This MUST be in the form of an absolute URL.
	RefreshURL *url.URL `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	// REQUIRED. The available scopes for the OAuth2 security scheme.
	// A map between the scope name and a short description for it.
	AvailableScopes MapOfStrings `json:"availableScopes" yaml:"availableScopes"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the OAuth flow for correctness.
func (f *OAuthFlowAuthorizationCode) Validate() error {
	if f.AuthorizationURL == nil {
		return &errpath.ErrField{Field: "authorizationUrl", Err: &errpath.ErrRequired{}}
	}

	if f.TokenURL == nil {
		return &errpath.ErrField{Field: "tokenUrl", Err: &errpath.ErrRequired{}}
	}

	if f.AvailableScopes == nil {
		return &errpath.ErrField{Field: "availableScopes", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(f.Extensions)
}
