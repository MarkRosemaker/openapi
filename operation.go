package openapi

import (
	"errors"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Operation describes a single API operation on a path.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#operation-object
type Operation struct {
	// A list of tags for API documentation control. Tags can be used for logical grouping of operations by resources or any other qualifier.
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`
	// A short summary of what the operation does.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A verbose explanation of the operation behavior. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Additional external documentation for this operation.
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// Unique string used to identify the operation. The id MUST be unique among all operations described in the API. The operationId value is **case-sensitive**. Tools and libraries MAY use the operationId to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
	OperationID string `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	// A list of parameters that are applicable for this operation. If a parameter is already defined at the Path Item, the new definition will override it but can never remove it. The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location. The list can use the Reference Object to link to parameters that are defined at the OpenAPI Object's components/parameters.
	Parameters ParameterList `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// The request body applicable for this operation. The `requestBody` is fully supported in HTTP methods where the HTTP 1.1 specification RFC7231 has explicitly defined semantics for request bodies. In other cases where the HTTP spec is vague (such as GET, HEAD and DELETE), `requestBody` is permitted but does not have well-defined semantics and SHOULD be avoided if possible.
	RequestBody *RequestBodyRef `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	// The list of possible responses as they are returned from executing this operation.
	Responses OperationResponses `json:"responses,omitempty" yaml:"responses,omitempty"`
	// A map of possible out-of band callbacks related to the parent operation. The key is a unique identifier for the Callback Object. Each value in the map is a Callback Object that describes a request that may be initiated by the API provider and the expected responses.
	Callbacks Callbacks `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	// Declares this operation to be deprecated. Consumers SHOULD refrain from usage of the declared operation. Default value is `false`.
	Deprecated bool `json:"deprecated,omitempty,omitzero" yaml:"deprecated,omitempty"`
	// A declaration of which security mechanisms can be used for this operation. The list of values includes alternative security requirement objects that can be used. Only one of the security requirement objects need to be satisfied to authorize a request. To make security optional, an empty security requirement (`{}`) can be included in the array. This definition overrides any declared top-level `security`. To remove a top-level security declaration, an empty array can be used.
	Security SecurityRequirements `json:"security,omitempty" yaml:"security,omitempty"`
	// An alternative `server` array to service this operation. If an alternative `server` object is specified at the Path Item Object or Root level, it will be overridden by this value.
	Servers Servers `json:"servers,omitempty" yaml:"servers,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate validates the operation.
func (o *Operation) Validate() error {
	o.Description = strings.TrimSpace(o.Description)

	if o.ExternalDocs != nil {
		if err := o.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if err := o.Parameters.Validate(); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	if o.RequestBody != nil {
		if err := o.RequestBody.Validate(); err != nil {
			return &errpath.ErrField{Field: "requestBody", Err: err}
		}
	}

	// validate the key: check if it is a StatusCode
	for code := range o.Responses.ByIndex() {
		if err := code.Validate(); err != nil {
			return &errpath.ErrField{
				Field: "responses",
				Err:   &errpath.ErrKey{Key: string(code), Err: err},
			}
		}
	}

	if len(o.Responses) == 1 {
		// only one response, special requirements
		for code := range o.Responses {
			if code == StatusCodeDefault {
				return &errpath.ErrField{Field: "responses", Err: &errpath.ErrKey{
					Key: string(StatusCodeDefault),
					Err: errors.New("must not be the only response"),
				}}
			}

			// must be a successful response
			if !code.IsSuccess() {
				return &errpath.ErrField{Field: "responses", Err: &errpath.ErrKey{
					Key: string(code),
					Err: errors.New("single response must be a successful response"),
				}}
			}
		}
	}

	if err := o.Responses.Validate(); err != nil {
		return &errpath.ErrField{Field: "responses", Err: err}
	}

	if err := o.Callbacks.Validate(); err != nil {
		return &errpath.ErrField{Field: "callbacks", Err: err}
	}

	if err := o.Security.Validate(); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := o.Servers.Validate(); err != nil {
		return &errpath.ErrField{Field: "servers", Err: err}
	}

	return validateExtensions(o.Extensions)
}
