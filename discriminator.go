package openapi

import (
	"fmt"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Discriminator names the property that tells which of a oneOf, anyOf or allOf schema's alternatives a payload is.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#discriminator-object
type Discriminator struct {
	// REQUIRED. The name of the property in the payload that will hold the discriminator value.
	PropertyName string `json:"propertyName" yaml:"propertyName"`
	// Maps values of the property to the schemas they stand for, each a component schema's name or a reference to one.
	Mapping MapOfStrings `json:"mapping,omitempty" yaml:"mapping,omitempty"`
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

// MappingRef is the reference a mapping value stands for: the value itself, or for a bare name, the component schema of that name.
func MappingRef(value string) string {
	if strings.Contains(value, "#") || strings.Contains(value, "/") {
		return value
	}

	return "#/components/schemas/" + value
}

func (l *loader) resolveDiscriminator(d *Discriminator) error {
	for key, value := range d.Mapping.ByIndex() {
		if _, ok := l.schemas[MappingRef(value.Value)]; !ok {
			return &errpath.ErrField{Field: "mapping", Err: &errpath.ErrKey{
				Key: key, Err: fmt.Errorf("couldn't resolve %q", value.Value),
			}}
		}
	}

	return nil
}
