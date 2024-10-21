package openapi

import (
	"net/url"
	"regexp"
)

// Document is an OpenAPI document.
// It is a self-contained or composite resource which defines or describes an API or elements of an API.
// An OpenAPI document uses and conforms to the OpenAPI Specification.
// ([Documentation])
//
// [Documentation]: https://spec.openapis.org/oas/v3.1.0#openapi-document
type Document struct {
	// REQUIRED. This string MUST be the version number of the OpenAPI Specification that the OpenAPI document uses. The openapi field SHOULD be used by tooling to interpret the OpenAPI document. This is not related to the API info.version string.
	OpenAPI string `json:"openapi" yaml:"openapi"`
	// REQUIRED. Provides metadata about the API. The metadata MAY be used by tooling as required.
	Info *Info `json:"info,omitempty" yaml:"info,omitempty"`
	// The default value for the $schema keyword within Schema Objects contained within this OAS document. This MUST be in the form of a URI.
	// Default: "https://spec.openapis.org/oas/3.1/dialect/base"
	// NOTE: Anything other than the default value is not supported.
	JSONSchemaDialect *url.URL `json:"jsonSchemaDialect,omitempty" yaml:"jsonSchemaDialect,omitempty"`
	// An array of Server Objects, which provide connectivity information to a target server. If the servers property is not provided, or is an empty array, the default value would be a Server Object with a url value of /.
	Servers Servers `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// reOpenAPIVersion is a regular expression that matches the OpenAPI version.
// Allowed are 3.0.x and 3.1.x.
var reOpenAPIVersion = regexp.MustCompile(`^3\.(0|1)\.\d+(-.+)?$`)

// Validate checks the OpenAPI document for correctness.
func (d *Document) Validate() error {
	if d.OpenAPI == "" {
		return &ErrRequired{Target: "openapi field"}
	}

	if !reOpenAPIVersion.MatchString(d.OpenAPI) {
		return &ErrInvalid{
			Target:  "openapi field",
			Value:   d.OpenAPI,
			Message: "must be a valid version (3.0.x or 3.1.x)",
		}
	}

	if d.Info == nil {
		return &ErrRequired{Target: "info"}
	}

	if err := d.Info.Validate(); err != nil {
		return &ErrField{Field: "info", Err: err}
	}

	if len(d.Servers) == 0 {
		// The default value would be a Server Object with a url value of /.
		d.Servers = Servers{{URL: "/"}}
	}

	// TODO
	// The OpenAPI document MUST contain at least one [paths](#paths-object) field, a [components](#oasComponents) field or a [webhooks](#oasWebhooks) field. An OpenAPI document uses and conforms to the OpenAPI Specification.

	return nil
}
