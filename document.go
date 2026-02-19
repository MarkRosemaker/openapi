package openapi

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/MarkRosemaker/errpath"
)

// ErrEmptyDocument is thrown if the OpenAPI document does not contain at least one paths field, a components field or a webhooks field.
var ErrEmptyDocument = errors.New("document must contain at least one paths field, a components field or a webhooks field")

// Document is an OpenAPI document, the root object.
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
	Components Components `json:"components,omitzero" yaml:"components,omitempty"`
	// A declaration of which security mechanisms can be used across the API. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. Individual operations can override this definition. To make security optional, an empty security requirement (`{}`) can be included in the array.
	Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	// A list of tags used by the document with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic. Each tag name in the list MUST be unique.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// reOpenAPIVersion is a regular expression that matches the OpenAPI version.
// Allowed are 3.0.x and 3.1.x.
var reOpenAPIVersion = regexp.MustCompile(`^3\.(0|1)\.\d+(-.+)?$`)

// Validate checks the OpenAPI document for correctness.
func (d *Document) Validate() error {
	if d.OpenAPI == "" {
		return &errpath.ErrField{Field: "openapi", Err: &errpath.ErrRequired{}}
	}

	if !reOpenAPIVersion.MatchString(d.OpenAPI) {
		return &errpath.ErrField{
			Field: "openapi",
			Err: &errpath.ErrInvalid[string]{
				Value:   d.OpenAPI,
				Message: "must be a valid version (3.0.x or 3.1.x)",
			},
		}
	}

	if d.Info == nil {
		return &errpath.ErrField{Field: "info", Err: &errpath.ErrRequired{}}
	}

	if err := d.Info.Validate(); err != nil {
		return &errpath.ErrField{Field: "info", Err: err}
	}

	const defaultJSONSchemaDialect = "https://spec.openapis.org/oas/3.1/dialect/base"
	if d.JSONSchemaDialect != nil &&
		d.JSONSchemaDialect.String() != defaultJSONSchemaDialect {
		return &errpath.ErrField{Field: "jsonSchemaDialect", Err: &errpath.ErrInvalid[string]{
			Value: d.JSONSchemaDialect.String(),
			Enum:  []string{defaultJSONSchemaDialect},
		}}
	}

	if err := d.Servers.Validate(); err != nil {
		return &errpath.ErrField{Field: "servers", Err: err}
	}

	// The OpenAPI document MUST contain at least one paths field, a components field or a webhooks field.
	if len(d.Paths) == 0 && len(d.Webhooks) == 0 && d.Components.isEmpty() {
		return ErrEmptyDocument
	}

	if err := d.Paths.Validate(); err != nil {
		return &errpath.ErrField{Field: "paths", Err: err}
	}

	if err := d.Webhooks.Validate(); err != nil {
		return &errpath.ErrField{Field: "webhooks", Err: err}
	}

	if err := d.Components.Validate(); err != nil {
		return &errpath.ErrField{Field: "components", Err: err}
	}

	if err := d.Security.Validate(); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := d.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if d.ExternalDocs != nil {
		if err := d.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return validateExtensions(d.Extensions)
}

// Sorts the paths and fields of components that are maps by key.
func (d *Document) SortMaps() {
	d.Paths.Sort()
	d.Components.SortMaps()
}

func (l *loader) collectDocument(doc *Document, ref ref) {
	l.collectPaths(doc.Paths, append(ref, "paths"))
	l.collectWebhooks(doc.Webhooks, append(ref, "webhooks"))
	l.collectComponents(doc.Components, append(ref, "components"))
}

func (l *loader) resolveDocument(doc *Document) error {
	// fields that don't need to be resolved:
	// - Info
	// - JSONSchemaDialect
	// - Servers
	// - Security
	// - Tags
	// - ExternalDocs

	if err := l.resolvePaths(doc.Paths); err != nil {
		return &errpath.ErrField{Field: "paths", Err: err}
	}

	if err := l.resolveWebhooks(doc.Webhooks); err != nil {
		return &errpath.ErrField{Field: "webhooks", Err: err}
	}

	if err := l.resolveComponents(doc.Components); err != nil {
		return &errpath.ErrField{Field: "components", Err: err}
	}

	return nil
}
