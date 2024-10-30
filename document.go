package openapi

import (
	"errors"
	"net/url"
	"regexp"
)

// ErrEmptyDocument is thrown if the OpenAPI document does not contain at least one paths field, a components field or a webhooks field.
var ErrEmptyDocument = errors.New("document must contain at least one paths field, a components field or a webhooks field")

// Document is an OpenAPI document.
// It is a self-contained or composite resource which defines or describes an API or elements of an API.
// An OpenAPI document uses and conforms to the OpenAPI Specification.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#openapi-document
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
	// The available paths and operations for the API.
	Paths Paths `json:"paths,omitempty" yaml:"paths,omitempty"`
	// The incoming webhooks that MAY be received as part of this API and that the API consumer MAY choose to implement. Closely related to the `callbacks` feature, this section describes requests initiated other than by an API call, for example by an out of band registration.
	Webhooks Webhooks `json:"webhooks,omitempty" yaml:"webhooks,omitempty"`

	// An element to hold various schemas for the document.
	Components Components `json:"components,omitempty" yaml:"components,omitempty"`
}

// reOpenAPIVersion is a regular expression that matches the OpenAPI version.
// Allowed are 3.0.x and 3.1.x.
var reOpenAPIVersion = regexp.MustCompile(`^3\.(0|1)\.\d+(-.+)?$`)

// Validate checks the OpenAPI document for correctness.
func (d *Document) Validate() error {
	if d.OpenAPI == "" {
		return &ErrField{Field: "openapi", Err: &ErrRequired{}}
	}

	if !reOpenAPIVersion.MatchString(d.OpenAPI) {
		return &ErrField{
			Field: "openapi",
			Err: &ErrInvalid[string]{
				Value:   d.OpenAPI,
				Message: "must be a valid version (3.0.x or 3.1.x)",
			},
		}
	}

	if d.Info == nil {
		return &ErrField{Field: "info", Err: &ErrRequired{}}
	}

	if err := d.Info.Validate(); err != nil {
		return &ErrField{Field: "info", Err: err}
	}

	const defaultJSONSchemaDialect = "https://spec.openapis.org/oas/3.1/dialect/base"
	if d.JSONSchemaDialect != nil &&
		d.JSONSchemaDialect.String() != defaultJSONSchemaDialect {
		return &ErrField{Field: "jsonSchemaDialect", Err: &ErrInvalid[string]{
			Value: d.JSONSchemaDialect.String(),
			Enum:  []string{defaultJSONSchemaDialect},
		}}
	}

	if err := d.Servers.Validate(); err != nil {
		return &ErrField{Field: "servers", Err: err}
	}

	// The OpenAPI document MUST contain at least one paths field, a components field or a webhooks field.
	if len(d.Paths) == 0 && len(d.Webhooks) == 0 && d.Components.isEmpty() {
		return ErrEmptyDocument
	}

	if err := d.Paths.Validate(); err != nil {
		return &ErrField{Field: "paths", Err: err}
	}

	if err := d.Webhooks.Validate(); err != nil {
		return &ErrField{Field: "webhooks", Err: err}
	}

	if err := d.Components.Validate(); err != nil {
		return &ErrField{Field: "components", Err: err}
	}

	return nil
}
