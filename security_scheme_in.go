package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// SecuritySchemeIn is where an API key is sent.
type SecuritySchemeIn string

// The places an API key can be sent in.
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

	return &errpath.ErrInvalid[SecuritySchemeIn]{
		Value: s,
		Enum:  allSecuritySchemeIn,
	}
}
