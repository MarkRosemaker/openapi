package openapi

import (
	"net/url"

	"github.com/MarkRosemaker/errpath"
)

// XML says how a schema is written as XML: the name of its element or attribute, its namespace, and whether it is an
// attribute or, for a list, wrapped in an element of its own. ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#xml-object
type XML struct {
	// The name of the element or attribute, in place of the property's or, for a list's items, the list's.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// The absolute URI of the namespace.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// The prefix of the name.
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	// Whether the property is an attribute of its element rather than an element of its own.
	Attribute bool `json:"attribute,omitzero" yaml:"attribute,omitempty"`
	// Whether a list is wrapped in an element of its own, as <books><book/><book/></books>, rather than its items
	// standing side by side.
	Wrapped bool `json:"wrapped,omitzero" yaml:"wrapped,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",embed" yaml:",embed"`
}

// Validate checks the XML object for correctness.
func (x *XML) Validate() error {
	if x.Namespace != "" {
		if u, err := url.Parse(x.Namespace); err != nil || !u.IsAbs() {
			return &errpath.ErrField{Field: "namespace", Err: &errpath.ErrInvalid[string]{
				Value:   x.Namespace,
				Message: "must be an absolute URI",
			}}
		}
	}

	return validateExtensions(x.Extensions)
}
