package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// SecuritySchemeType is the kind of a security scheme.
type SecuritySchemeType string

// The kinds of security scheme.
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

	return &errpath.ErrInvalid[SecuritySchemeType]{
		Value: tp,
		Enum:  allSecurityTypes,
	}
}
