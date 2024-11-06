package openapi

import "strings"

//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#tag-object
type Tag struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string        `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (t *Tag) Validate() error {
	if t.Name == "" {
		return &ErrField{Field: "name", Err: &ErrRequired{}}
	}

	t.Description = strings.TrimSpace(t.Description)

	if t.ExternalDocs != nil {
		if err := t.ExternalDocs.Validate(); err != nil {
			return &ErrField{Field: "externalDocs", Err: err}
		}
	}

	return validateExtensions(t.Extensions)
}
