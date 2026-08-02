package asyncapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// SecuritySchemeType is the type of a security scheme.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#securitySchemeObject
type SecuritySchemeType string

const (
	// SecuritySchemeTypeUserPassword is the user/password authentication.
	SecuritySchemeTypeUserPassword SecuritySchemeType = "userPassword"
	// SecuritySchemeTypeAPIKey is an API key, either as user or as password.
	SecuritySchemeTypeAPIKey SecuritySchemeType = "apiKey"
	// SecuritySchemeTypeX509 is an X.509 certificate.
	SecuritySchemeTypeX509 SecuritySchemeType = "X509"
	// SecuritySchemeTypeSymmetricEncryption is a symmetric end-to-end encryption.
	SecuritySchemeTypeSymmetricEncryption SecuritySchemeType = "symmetricEncryption"
	// SecuritySchemeTypeAsymmetricEncryption is an asymmetric end-to-end encryption.
	SecuritySchemeTypeAsymmetricEncryption SecuritySchemeType = "asymmetricEncryption"
	// SecuritySchemeTypeHTTPAPIKey is an API key that is sent as an HTTP header, query or cookie parameter.
	SecuritySchemeTypeHTTPAPIKey SecuritySchemeType = "httpApiKey"
	// SecuritySchemeTypeHTTP is an HTTP authentication.
	SecuritySchemeTypeHTTP SecuritySchemeType = "http"
	// SecuritySchemeTypeOAuth2 is one of OAuth2's common flows.
	SecuritySchemeTypeOAuth2 SecuritySchemeType = "oauth2"
	// SecuritySchemeTypeOpenIDConnect is OpenID Connect Discovery.
	SecuritySchemeTypeOpenIDConnect SecuritySchemeType = "openIdConnect"
	// SecuritySchemeTypePlain is the SASL PLAIN mechanism.
	SecuritySchemeTypePlain SecuritySchemeType = "plain"
	// SecuritySchemeTypeScramSha256 is the SASL SCRAM-SHA-256 mechanism.
	SecuritySchemeTypeScramSha256 SecuritySchemeType = "scramSha256"
	// SecuritySchemeTypeScramSha512 is the SASL SCRAM-SHA-512 mechanism.
	SecuritySchemeTypeScramSha512 SecuritySchemeType = "scramSha512"
	// SecuritySchemeTypeGSSAPI is the SASL GSSAPI mechanism.
	SecuritySchemeTypeGSSAPI SecuritySchemeType = "gssapi"
)

var allSecuritySchemeTypes = []SecuritySchemeType{
	SecuritySchemeTypeUserPassword,
	SecuritySchemeTypeAPIKey,
	SecuritySchemeTypeX509,
	SecuritySchemeTypeSymmetricEncryption,
	SecuritySchemeTypeAsymmetricEncryption,
	SecuritySchemeTypeHTTPAPIKey,
	SecuritySchemeTypeHTTP,
	SecuritySchemeTypeOAuth2,
	SecuritySchemeTypeOpenIDConnect,
	SecuritySchemeTypePlain,
	SecuritySchemeTypeScramSha256,
	SecuritySchemeTypeScramSha512,
	SecuritySchemeTypeGSSAPI,
}

// Validate validates the security scheme type.
func (tp SecuritySchemeType) Validate() error {
	if slices.Contains(allSecuritySchemeTypes, tp) {
		return nil
	}

	return &errpath.ErrInvalid[SecuritySchemeType]{
		Value: tp,
		Enum:  allSecuritySchemeTypes,
	}
}
