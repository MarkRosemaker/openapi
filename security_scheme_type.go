package openapi

import "slices"

type SecuritySchemeType string

const (
	SecuritySchemeTypeAPIKey        SecuritySchemeType = "apiKey"
	SecuritySchemeTypeHTTP          SecuritySchemeType = "http"
	SecuritySchemeTypeMutualTLS     SecuritySchemeType = "mutualTLS"
	SecuritySchemeTypeOAuth2        SecuritySchemeType = "oauth2"
	SecuritySchemeTypeOpenIDConnect SecuritySchemeType = "openIdConnect"
)

var allSecurityTypes = []SecuritySchemeType{
	SecuritySchemeTypeAPIKey,
	SecuritySchemeTypeHTTP,
	SecuritySchemeTypeMutualTLS,
	SecuritySchemeTypeOAuth2,
	SecuritySchemeTypeOpenIDConnect,
}

// Validate validates the security scheme type.
func (tp SecuritySchemeType) Validate() error {
	if slices.Contains(allSecurityTypes, tp) {
		return nil
	}

	return &ErrInvalid[SecuritySchemeType]{
		Value: tp,
		Enum:  allSecurityTypes,
	}
}
