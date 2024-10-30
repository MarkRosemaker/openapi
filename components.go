package openapi

type Components struct {
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

func (c *Components) Validate() error {
	if err := c.Schemas.Validate(); err != nil {
		return &ErrField{Field: "schemas", Err: err}
	}

	return nil
}

func (c Components) isEmpty() bool {
	return len(c.Schemas) == 0
}
