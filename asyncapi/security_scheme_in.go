package asyncapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// SecuritySchemeIn is "the location of the API key. Valid values are `user` and `password`
// for `apiKey` and `query`, `header` or `cookie` for `httpApiKey`."
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#securitySchemeObject
type SecuritySchemeIn string

const (
	// SecuritySchemeInUser is the location of an API key that is sent as the user of a connection.
	SecuritySchemeInUser SecuritySchemeIn = "user"
	// SecuritySchemeInPassword is the location of an API key that is sent as the password of a connection.
	SecuritySchemeInPassword SecuritySchemeIn = "password"
	// SecuritySchemeInQuery is the location of an API key that is sent as a query parameter.
	SecuritySchemeInQuery SecuritySchemeIn = "query"
	// SecuritySchemeInHeader is the location of an API key that is sent as a header.
	SecuritySchemeInHeader SecuritySchemeIn = "header"
	// SecuritySchemeInCookie is the location of an API key that is sent as a cookie.
	SecuritySchemeInCookie SecuritySchemeIn = "cookie"
)

// allSecuritySchemeInAPIKey are the valid locations for a security scheme of type apiKey.
var allSecuritySchemeInAPIKey = []SecuritySchemeIn{
	SecuritySchemeInUser,
	SecuritySchemeInPassword,
}

// allSecuritySchemeInHTTPAPIKey are the valid locations for a security scheme of type httpApiKey.
var allSecuritySchemeInHTTPAPIKey = []SecuritySchemeIn{
	SecuritySchemeInQuery,
	SecuritySchemeInHeader,
	SecuritySchemeInCookie,
}

// validate checks that the location is one of the given valid locations.
func (in SecuritySchemeIn) validate(valid []SecuritySchemeIn) error {
	if slices.Contains(valid, in) {
		return nil
	}

	return &errpath.ErrInvalid[SecuritySchemeIn]{
		Value: in,
		Enum:  valid,
	}
}
