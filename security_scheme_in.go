package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

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

	return &errpath.ErrInvalid[SecuritySchemeIn]{
		Value: s,
		Enum:  allSecuritySchemeIn,
	}
}
