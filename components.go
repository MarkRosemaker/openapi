package openapi

type Components struct {
	Schemas Schemas `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

func (c Components) isEmpty() bool {
	return len(c.Schemas) == 0
}
