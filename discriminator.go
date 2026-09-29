package openapi

import "github.com/MarkRosemaker/errpath"

// Discriminator names the property that tells which of a oneOf, anyOf or allOf schema's alternatives a payload is.
// The mapping of property values to schemas is not supported yet, so a discriminator with one fails to load.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#discriminator-object
type Discriminator struct {
	// REQUIRED. The name of the property in the payload that will hold the discriminator value.
	PropertyName string `json:"propertyName" yaml:"propertyName"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",embed" yaml:",embed"`
}

// Validate checks the discriminator for correctness.
func (d *Discriminator) Validate() error {
	if d.PropertyName == "" {
		return &errpath.ErrField{Field: "propertyName", Err: &errpath.ErrRequired{}}
	}

	return validateExtensions(d.Extensions)
}
