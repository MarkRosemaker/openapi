package asyncapi

import "github.com/MarkRosemaker/errpath"

// OAuthFlows allows configuration of the supported OAuth Flows.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#oauthFlowsObject
type OAuthFlows struct {
	// Configuration for the OAuth Implicit flow.
	Implicit *OAuthFlowImplicit `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Protected Credentials flow.
	Password *OAuthFlowPassword `json:"password,omitempty" yaml:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow.
	ClientCredentials *OAuthFlowClientCredentials `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow.
	AuthorizationCode *OAuthFlowAuthorizationCode `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the OAuth flows for correctness.
func (f *OAuthFlows) Validate() error {
	if f.Implicit != nil {
		if err := f.Implicit.Validate(); err != nil {
			return &errpath.ErrField{Field: "implicit", Err: err}
		}
	}

	if f.Password != nil {
		if err := f.Password.Validate(); err != nil {
			return &errpath.ErrField{Field: "password", Err: err}
		}
	}

	if f.ClientCredentials != nil {
		if err := f.ClientCredentials.Validate(); err != nil {
			return &errpath.ErrField{Field: "clientCredentials", Err: err}
		}
	}

	if f.AuthorizationCode != nil {
		if err := f.AuthorizationCode.Validate(); err != nil {
			return &errpath.ErrField{Field: "authorizationCode", Err: err}
		}
	}

	return validateExtensions(f.Extensions)
}
