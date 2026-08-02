package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Tag allows adding meta data to a single tag.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#tagObject
type Tag struct {
	// REQUIRED. The name of the tag.
	Name string `json:"name" yaml:"name"`
	// A short description for the tag. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Additional external documentation for this tag.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the tag for correctness.
func (t *Tag) Validate() error {
	if t.Name == "" {
		return &errpath.ErrField{Field: "name", Err: &errpath.ErrRequired{}}
	}

	t.Description = strings.TrimSpace(t.Description)

	if t.ExternalDocs != nil {
		if err := t.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return validateExtensions(t.Extensions)
}

func (l *loader) collectTagRef(t *TagRef, ref ref) {
	if !collectRef(l, t, l.tags, ref) {
		return
	}

	if t.Value.ExternalDocs != nil {
		l.collectExternalDocsRef(t.Value.ExternalDocs, append(ref, "externalDocs"))
	}
}

func (l *loader) resolveTagRef(t *TagRef) error {
	return resolveRef(t, l.tags, func(t *Tag) error {
		if t.ExternalDocs != nil {
			if err := l.resolveExternalDocsRef(t.ExternalDocs); err != nil {
				return &errpath.ErrField{Field: "externalDocs", Err: err}
			}
		}

		return nil
	})
}
