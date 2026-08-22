package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Message describes a message received on a given channel and operation.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#messageObject
type Message struct {
	// Schema definition of the application headers. Schema MUST be a map of key-value pairs.
	// It MUST NOT define the protocol headers.
	Headers *AnySchemaRef `json:"headers,omitempty" yaml:"headers,omitempty"`
	// Definition of the message payload.
	Payload *AnySchemaRef `json:"payload,omitempty" yaml:"payload,omitempty"`
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
	// A list of traits to apply to the message object.
	// Traits MUST be merged using the traits merge mechanism.
	// The resulting object MUST be a valid Message Object.
	Traits MessageTraitList `json:"traits,omitempty" yaml:"traits,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the message for correctness.
func (m *Message) Validate() error {
	if m.Headers != nil {
		if err := m.Headers.Validate(); err != nil {
			return &errpath.ErrField{Field: "headers", Err: err}
		}
	}

	if m.Payload != nil {
		if err := m.Payload.Validate(); err != nil {
			return &errpath.ErrField{Field: "payload", Err: err}
		}
	}

	if m.CorrelationID != nil {
		if err := m.CorrelationID.Validate(); err != nil {
			return &errpath.ErrField{Field: "correlationId", Err: err}
		}
	}

	if m.ContentType != "" {
		if err := m.ContentType.Validate(); err != nil {
			return &errpath.ErrField{Field: "contentType", Err: err}
		}
	}

	m.Description = strings.TrimSpace(m.Description)

	if err := m.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if m.ExternalDocs != nil {
		if err := m.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if m.Bindings != nil {
		if err := m.Bindings.Validate(); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	if err := m.Examples.Validate(); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := m.Traits.Validate(); err != nil {
		return &errpath.ErrField{Field: "traits", Err: err}
	}

	return validateExtensions(m.Extensions)
}

func (l *loader) collectMessageRef(m *MessageRef, ref ref) {
	if !collectRef(l, m, l.messages, ref) {
		return
	}

	l.collectMessage(m.Value, ref)
}

func (l *loader) collectMessage(m *Message, ref ref) {
	if m.Headers != nil {
		l.collectAnySchemaRef(m.Headers, append(ref, "headers"))
	}

	if m.Payload != nil {
		l.collectAnySchemaRef(m.Payload, append(ref, "payload"))
	}

	if m.CorrelationID != nil {
		l.collectCorrelationIDRef(m.CorrelationID, append(ref, "correlationId"))
	}

	l.collectTags(m.Tags, append(ref, "tags"))

	if m.ExternalDocs != nil {
		l.collectExternalDocsRef(m.ExternalDocs, append(ref, "externalDocs"))
	}

	if m.Bindings != nil {
		l.collectBindingsRef(m.Bindings, append(ref, "bindings"))
	}

	l.collectMessageTraitList(m.Traits, append(ref, "traits"))
}

func (l *loader) resolveMessageRef(m *MessageRef) error {
	return resolveRef(m, l.messages, l.resolveMessage)
}

func (l *loader) resolveMessage(m *Message) error {
	if m.Headers != nil {
		if err := l.resolveAnySchemaRef(m.Headers); err != nil {
			return &errpath.ErrField{Field: "headers", Err: err}
		}
	}

	if m.Payload != nil {
		if err := l.resolveAnySchemaRef(m.Payload); err != nil {
			return &errpath.ErrField{Field: "payload", Err: err}
		}
	}

	if m.CorrelationID != nil {
		if err := l.resolveCorrelationIDRef(m.CorrelationID); err != nil {
			return &errpath.ErrField{Field: "correlationId", Err: err}
		}
	}

	if err := l.resolveTags(m.Tags); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if m.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(m.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if m.Bindings != nil {
		if err := l.resolveBindingsRef(m.Bindings); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	if err := l.resolveMessageTraitList(m.Traits); err != nil {
		return &errpath.ErrField{Field: "traits", Err: err}
	}

	return nil
}
