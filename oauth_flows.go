package openapi

import "github.com/MarkRosemaker/errpath"

// The OAuthFlows object allows configuration of the supported OAuth Flows.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#oauth-flows-object
type OAuthFlows struct {
	// Configuration for the OAuth Implicit flow
	Implicit *OAuthFlowImplicit `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	// Configuration for the OAuth Resource Owner Password flow
	Password *OAuthFlowPassword `json:"password,omitempty" yaml:"password,omitempty"`
	// Configuration for the OAuth Client Credentials flow.
	// Previously called `application` in OpenAPI 2.0.
	ClientCredentials *OAuthFlowClientCredentials `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	// Configuration for the OAuth Authorization Code flow.
	// Previously called `accessCode` in OpenAPI 2.0.
	AuthorizationCode *OAuthFlowAuthorizationCode `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

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
