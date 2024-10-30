package openapi

type Components struct {
	Links   Links   `json:"links,omitempty" yaml:"links,omitempty"`
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

func (c *Components) Validate() error {
	if err := c.Links.Validate(); err != nil {
		return &ErrField{Field: "links", Err: err}
	}

	if err := c.Schemas.Validate(); err != nil {
		return &ErrField{Field: "schemas", Err: err}
	}

	return nil
}

func (c Components) isEmpty() bool {
	return len(c.Schemas) == 0
}
