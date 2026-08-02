package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Operation describes a specific operation.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationObject
type Operation struct {
	// REQUIRED. Use `send` when it's expected that the application will send a message to the given channel,
	// and `receive` when the application should expect receiving messages from the given channel.
	Action OperationAction `json:"action" yaml:"action"`
	// REQUIRED. A $ref pointer to the definition of the channel in which this operation is performed.
	// If the operation is located in the root Operations Object, it MUST point to a channel definition located in the root Channels Object.
	Channel *ChannelRef `json:"channel" yaml:"channel"`
	// A human-friendly title for the operation.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short summary of what the operation is about.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A verbose explanation of the operation. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A declaration of which security schemes are associated with this operation.
	// Only one of the security scheme objects MUST be satisfied to authorize an operation.
	// In cases where server security also applies, it MUST also be satisfied.
	Security SecuritySchemeRefList `json:"security,omitempty" yaml:"security,omitempty"`
	// A list of tags for logical grouping and categorization of operations.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation for this operation.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// A map where the keys describe the name of the protocol and the values describe protocol-specific definitions for the operation.
	Bindings *BindingsRef `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	// A list of traits to apply to the operation object.
	// Traits MUST be merged using the traits merge mechanism.
	// The resulting object MUST be a valid Operation Object.
	Traits OperationTraitList `json:"traits,omitempty" yaml:"traits,omitempty"`
	// A list of $ref pointers pointing to the supported Message Objects that can be processed by this operation.
	// It MUST contain a subset of the messages defined in the channel referenced in this operation.
	//
	// Note: excluding this property from the operation implies that all messages from the channel will be included.
	// Explicitly set it to an empty, non-nil list if this operation should contain no messages.
	Messages MessageRefList `json:"messages,omitempty" yaml:"messages,omitempty"`
	// The definition of the reply in a request-reply operation.
	Reply *OperationReplyRef `json:"reply,omitempty" yaml:"reply,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the operation for correctness.
func (o *Operation) Validate() error {
	if o.Action == "" {
		return &errpath.ErrField{Field: "action", Err: &errpath.ErrRequired{}}
	}

	if err := o.Action.Validate(); err != nil {
		return &errpath.ErrField{Field: "action", Err: err}
	}

	if o.Channel == nil {
		return &errpath.ErrField{Field: "channel", Err: &errpath.ErrRequired{}}
	}

	if !o.Channel.isRef() {
		return &errpath.ErrField{Field: "channel", Err: ErrMustBeReference}
	}

	if err := o.Channel.Validate(); err != nil {
		return &errpath.ErrField{Field: "channel", Err: err}
	}

	o.Description = strings.TrimSpace(o.Description)

	if err := o.Security.Validate(); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := o.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if o.ExternalDocs != nil {
		if err := o.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if o.Bindings != nil {
		if err := o.Bindings.Validate(); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	if err := o.Traits.Validate(); err != nil {
		return &errpath.ErrField{Field: "traits", Err: err}
	}

	if err := o.Messages.Validate(); err != nil {
		return &errpath.ErrField{Field: "messages", Err: err}
	}

	if o.Reply != nil {
		if err := o.Reply.Validate(); err != nil {
			return &errpath.ErrField{Field: "reply", Err: err}
		}
	}

	return validateExtensions(o.Extensions)
}

func (l *loader) collectOperationRef(o *OperationRef, ref ref) {
	if !collectRef(l, o, l.operations, ref) {
		return
	}

	l.collectOperation(o.Value, ref)
}

func (l *loader) collectOperation(o *Operation, ref ref) {
	l.collectSecuritySchemeRefList(o.Security, append(ref, "security"))
	l.collectTags(o.Tags, append(ref, "tags"))

	if o.ExternalDocs != nil {
		l.collectExternalDocsRef(o.ExternalDocs, append(ref, "externalDocs"))
	}

	if o.Bindings != nil {
		l.collectBindingsRef(o.Bindings, append(ref, "bindings"))
	}

	l.collectOperationTraitList(o.Traits, append(ref, "traits"))

	if o.Reply != nil {
		l.collectOperationReplyRef(o.Reply, append(ref, "reply"))
	}
}

func (l *loader) resolveOperationRef(o *OperationRef) error {
	return resolveRef(o, l.operations, l.resolveOperation)
}

func (l *loader) resolveOperation(o *Operation) error {
	if o.Channel != nil {
		if err := l.resolveChannelRef(o.Channel); err != nil {
			return &errpath.ErrField{Field: "channel", Err: err}
		}
	}

	if err := l.resolveSecuritySchemeRefList(o.Security); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := l.resolveTags(o.Tags); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if o.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(o.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if o.Bindings != nil {
		if err := l.resolveBindingsRef(o.Bindings); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	if err := l.resolveOperationTraitList(o.Traits); err != nil {
		return &errpath.ErrField{Field: "traits", Err: err}
	}

	if err := l.resolveMessageRefList(o.Messages); err != nil {
		return &errpath.ErrField{Field: "messages", Err: err}
	}

	if o.Reply != nil {
		if err := l.resolveOperationReplyRef(o.Reply); err != nil {
			return &errpath.ErrField{Field: "reply", Err: err}
		}
	}

	return nil
}
