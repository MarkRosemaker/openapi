package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// MessageTrait describes a trait that MAY be applied to a [Message] object.
// This object MAY contain any property from the [Message] object, except `payload` and `traits`.
//
// If you're looking to apply traits to an operation, see the [OperationTrait] object.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#messageTraitObject
type MessageTrait struct {
	// Schema definition of the application headers. Schema MUST be a map of key-value pairs.
	// It MUST NOT define the protocol headers.
	Headers *AnySchemaRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	// Definition of the correlation ID used for message tracing or matching.
	CorrelationID *CorrelationIDRef `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	// The content type to use when encoding/decoding a message's payload.
	// When omitted, the value MUST be the one specified on the defaultContentType field of the document.
	ContentType MediaType `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	// A machine-friendly name for the message.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// A human-friendly title for the message.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short summary of what the message is about.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A verbose explanation of the message. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A list of tags for logical grouping and categorization of messages.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation for this message.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// A map where the keys describe the name of the protocol and the values describe protocol-specific definitions for the message.
	Bindings *BindingsRef `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	// List of examples.
	Examples MessageExamples `json:"examples,omitempty" yaml:"examples,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the message trait for correctness.
func (t *MessageTrait) Validate() error {
	if t.Headers != nil {
		if err := t.Headers.Validate(); err != nil {
			return &errpath.ErrField{Field: "headers", Err: err}
		}
	}

	if t.CorrelationID != nil {
		if err := t.CorrelationID.Validate(); err != nil {
			return &errpath.ErrField{Field: "correlationId", Err: err}
		}
	}

	if t.ContentType != "" {
		if err := t.ContentType.Validate(); err != nil {
			return &errpath.ErrField{Field: "contentType", Err: err}
		}
	}

	t.Description = strings.TrimSpace(t.Description)

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

	if err := t.Examples.Validate(); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	return validateExtensions(t.Extensions)
}

func (l *loader) collectMessageTraitList(ts MessageTraitList, ref ref) {
	for i, t := range ts {
		l.collectMessageTraitRef(t, append(ref, itoa(i)))
	}
}

func (l *loader) collectMessageTraitRef(t *MessageTraitRef, ref ref) {
	if !collectRef(l, t, l.messageTraits, ref) {
		return
	}

	l.collectMessageTrait(t.Value, ref)
}

func (l *loader) collectMessageTrait(t *MessageTrait, ref ref) {
	if t.Headers != nil {
		l.collectAnySchemaRef(t.Headers, append(ref, "headers"))
	}

	if t.CorrelationID != nil {
		l.collectCorrelationIDRef(t.CorrelationID, append(ref, "correlationId"))
	}

	l.collectTags(t.Tags, append(ref, "tags"))

	if t.ExternalDocs != nil {
		l.collectExternalDocsRef(t.ExternalDocs, append(ref, "externalDocs"))
	}

	if t.Bindings != nil {
		l.collectBindingsRef(t.Bindings, append(ref, "bindings"))
	}
}

func (l *loader) resolveMessageTraitRef(t *MessageTraitRef) error {
	return resolveRef(t, l.messageTraits, l.resolveMessageTrait)
}

func (l *loader) resolveMessageTrait(t *MessageTrait) error {
	if t.Headers != nil {
		if err := l.resolveAnySchemaRef(t.Headers); err != nil {
			return &errpath.ErrField{Field: "headers", Err: err}
		}
	}

	if t.CorrelationID != nil {
		if err := l.resolveCorrelationIDRef(t.CorrelationID); err != nil {
			return &errpath.ErrField{Field: "correlationId", Err: err}
		}
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
