package openapi

// The Components object holds a set of reusable objects for different aspects of the OAS.
// All objects defined within the components object will have no effect on the API unless they are explicitly referenced from properties outside the components object.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#components-object
type Components struct {
	// An object to hold reusable Schema Objects.
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`

	Links           Links           `json:"links,omitempty" yaml:"links,omitempty"`
	SecuritySchemes SecuritySchemes `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
}

func (c *Components) Validate() error {
	if err := c.Links.Validate(); err != nil {
		return &ErrField{Field: "links", Err: err}
	}

	if err := c.Schemas.Validate(); err != nil {
		return &ErrField{Field: "schemas", Err: err}
	}

	if err := c.SecuritySchemes.Validate(); err != nil {
		return &ErrField{Field: "securitySchemes", Err: err}
	}

	return nil
}

func (c Components) isEmpty() bool {
	return len(c.Schemas) == 0
}
