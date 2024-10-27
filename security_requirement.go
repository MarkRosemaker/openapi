package openapi

import (
	"errors"
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
			return &ErrKey{
				Key: string(name),
				Err: errors.New("empty security scheme name"),
			}
		}

		if len(scopes) == 0 {
			return &ErrKey{
				Key: string(name),
				Err: errors.New("security requirement must have at least one scope"),
			}
		}
	}

	return nil
}
