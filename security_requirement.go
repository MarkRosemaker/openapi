package openapi

import (
	"errors"
	"maps"
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// SecurityRequirement lists the required security schemes to execute this operation.
// The name used for each property MUST correspond to a security scheme declared in the Security Schemes under the Components Object.
//
// Security Requirement Objects that contain multiple schemes require that all schemes MUST be satisfied for a request to be authorized.
// This enables support for scenarios where multiple query parameters or HTTP headers are required to convey security information.
//
// When a list of Security Requirement Objects is defined on the OpenAPI Object or Operation Object, only one of the Security Requirement Objects in the list needs to be satisfied to authorize the request.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#security-requirement-object
type SecurityRequirement map[SecuritySchemeName][]string

// SecuritySchemeName is the name of a security scheme defined in the Security Schemes under the Components Object.
type SecuritySchemeName string

func (sr SecurityRequirement) Validate() error {
	for name, scopes := range sr {
		if name == "" {
			return &errpath.ErrKey{
				Key: string(name),
				Err: errors.New("empty security scheme name"),
			}

			// NOTE: Each name MUST correspond to a security scheme which is declared in the Security Schemes under the Components Object.
		}

		if scopes == nil {
			return &errpath.ErrKey{
				Key: string(name),
				Err: errors.New("list may be empty but must not be nil"),
			}

			// NOTE: If the security scheme is of type "oauth2" or "openIdConnect", then the value is a list of scope names required for the execution, and the list MAY be empty if authorization does not require a specified scope. For other security scheme types, the array MAY contain a list of role names which are required for the execution, but are not otherwise defined or exchanged in-band.
		}
	}

	return nil
}

func (sr SecurityRequirement) Equals(other SecurityRequirement) bool {
	return maps.EqualFunc(sr, other, slices.Equal)
}
