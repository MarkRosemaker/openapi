package openapi

import "slices"

type SecuritySchemeIn string

const (
	SecuritySchemeInQuery  SecuritySchemeIn = "query"
	SecuritySchemeInHeader SecuritySchemeIn = "header"
	SecuritySchemeInCookie SecuritySchemeIn = "cookie"
)

var allSecuritySchemeIn = []SecuritySchemeIn{
	SecuritySchemeInQuery,
	SecuritySchemeInHeader,
	SecuritySchemeInCookie,
}

// Validate validates the security location.
func (s SecuritySchemeIn) Validate() error {
	if slices.Contains(allSecuritySchemeIn, s) {
		return nil
	}

	return &ErrInvalid[SecuritySchemeIn]{
		Value: s,
		Enum:  allSecuritySchemeIn,
	}
}
