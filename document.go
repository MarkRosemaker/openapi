package openapi

import "regexp"

// Document is an OpenAPI document.
// It is a self-contained or composite resource which defines or describes an API or elements of an API.
// An OpenAPI document uses and conforms to the OpenAPI Specification.
// ([Source])
//
// [Source]: https://spec.openapis.org/oas/v3.1.0#openapi-document
type Document struct {
	// REQUIRED. This string MUST be the version number of the OpenAPI Specification that the OpenAPI document uses. The openapi field SHOULD be used by tooling to interpret the OpenAPI document. This is not related to the API info.version string.
	OpenAPI string `json:"openapi,strictcase" yaml:"openapi"`
}

// reOpenAPIVersion is a regular expression that matches the OpenAPI version.
// Allowed are 3.0.x and 3.1.x.
var reOpenAPIVersion = regexp.MustCompile(`^3\.(0|1)\.\d+(-.+)?$`)

// Validate checks the OpenAPI document for correctness.
func (d *Document) Validate() error {
	// validate and canonicalize the openapi version
	if d.OpenAPI == "" {
		return &ErrRequired{Target: "openapi.version"}
	}

	if !reOpenAPIVersion.MatchString(d.OpenAPI) {
		return &ErrInvalid{
			Target:  "openapi.version",
			Value:   d.OpenAPI,
			Message: "must be a valid version (3.0.x or 3.1.x)",
		}
	}

	// TODO
	// The OpenAPI document MUST contain at least one [paths](#paths-object) field, a [components](#oasComponents) field or a [webhooks](#oasWebhooks) field. An OpenAPI document uses and conforms to the OpenAPI Specification.

	return nil
}
