package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// OperationTrait describes a trait that MAY be applied to an [Operation] object.
// This object MAY contain any property from the [Operation] object, except the `action`, `channel`, `messages` and `traits` ones.
//
// If you're looking to apply traits to a message, see the [MessageTrait] object.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationTraitObject
type OperationTrait struct {
	// A human-friendly title for the operation.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short summary of what the operation is about.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A verbose explanation of the operation. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A declaration of which security schemes are associated with this operation.
	Security SecuritySchemeRefList `json:"security,omitempty" yaml:"security,omitempty"`
	// A list of tags for logical grouping and categorization of operations.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation for this operation.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// A map where the keys describe the name of the protocol and the values describe protocol-specific definitions for the operation.
	Bindings *BindingsRef `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the operation trait for correctness.
func (t *OperationTrait) Validate() error {
	t.Description = strings.TrimSpace(t.Description)

	if err := t.Security.Validate(); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := t.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if t.ExternalDocs != nil {
		if err := t.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if t.Bindings != nil {
		if err := t.Bindings.Validate(); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	return validateExtensions(t.Extensions)
}

func (l *loader) collectOperationTraitList(ts OperationTraitList, ref ref) {
	for i, t := range ts {
		l.collectOperationTraitRef(t, append(ref, itoa(i)))
	}
}

func (l *loader) collectOperationTraitRef(t *OperationTraitRef, ref ref) {
	if !collectRef(l, t, l.operationTraits, ref) {
		return
	}

	l.collectOperationTrait(t.Value, ref)
}

func (l *loader) collectOperationTrait(t *OperationTrait, ref ref) {
	l.collectSecuritySchemeRefList(t.Security, append(ref, "security"))
	l.collectTags(t.Tags, append(ref, "tags"))

	if t.ExternalDocs != nil {
		l.collectExternalDocsRef(t.ExternalDocs, append(ref, "externalDocs"))
	}

	if t.Bindings != nil {
		l.collectBindingsRef(t.Bindings, append(ref, "bindings"))
	}
}

func (l *loader) resolveOperationTraitRef(t *OperationTraitRef) error {
	return resolveRef(t, l.operationTraits, l.resolveOperationTrait)
}

func (l *loader) resolveOperationTrait(t *OperationTrait) error {
	if err := l.resolveSecuritySchemeRefList(t.Security); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := l.resolveTags(t.Tags); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if t.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(t.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if t.Bindings != nil {
		if err := l.resolveBindingsRef(t.Bindings); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	return nil
}
