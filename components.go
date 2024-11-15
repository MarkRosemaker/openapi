package openapi

// The Components object holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#components-object
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	// An object to hold reusable Response Objects.
	Responses ResponsesByName `json:"responses,omitempty" yaml:"responses,omitempty"`
	// An object to hold reusable Parameter Objects.
	Parameters Parameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// An object to hold reusable Example Objects.
	Examples Examples `json:"examples,omitempty" yaml:"examples,omitempty"`
	// An object to hold reusable Request Body Objects.
	RequestBodies RequestBodies `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	// An object to hold reusable Header Objects.
	Headers Headers `json:"headers,omitempty" yaml:"headers,omitempty"`
	// An object to hold reusable Security Scheme Objects.
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	// An object to hold reusable Link Objects.
	Links Links `json:"links,omitempty" yaml:"links,omitempty"`
	// An object to hold reusable Callback Objects.
	Callbacks CallbackRefs `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
	// An object to hold reusable Path Item Object.
	PathItems PathItems `json:"pathItems,omitempty" yaml:"pathItems,omitempty"`
}

func (c *Components) Validate() error {
	if err := c.Schemas.Validate(); err != nil {
		return &ErrField{Field: "schemas", Err: err}
	}

	if err := c.Responses.Validate(); err != nil {
		return &ErrField{Field: "responses", Err: err}
	}

	if err := c.Parameters.Validate(); err != nil {
		return &ErrField{Field: "parameters", Err: err}
	}

	if err := c.Examples.Validate(); err != nil {
		return &ErrField{Field: "examples", Err: err}
	}

	if err := c.RequestBodies.Validate(); err != nil {
		return &ErrField{Field: "requestBodies", Err: err}
	}

	if err := c.Headers.Validate(); err != nil {
		return &ErrField{Field: "headers", Err: err}
	}

	if err := c.SecuritySchemes.Validate(); err != nil {
		return &ErrField{Field: "securitySchemes", Err: err}
	}

	if err := c.Links.Validate(); err != nil {
		return &ErrField{Field: "links", Err: err}
	}

	if err := c.Callbacks.Validate(); err != nil {
		return &ErrField{Field: "callbacks", Err: err}
	}

	if err := c.PathItems.Validate(); err != nil {
		return &ErrField{Field: "pathItems", Err: err}
	}

	return nil
}

func (c Components) isEmpty() bool {
	return len(c.Schemas) == 0 && len(c.Responses) == 0 && len(c.Parameters) == 0 &&
		len(c.Examples) == 0 && len(c.RequestBodies) == 0 && len(c.Headers) == 0 &&
		len(c.SecuritySchemes) == 0 && len(c.Links) == 0 && len(c.Callbacks) == 0 &&
		len(c.PathItems) == 0
}
