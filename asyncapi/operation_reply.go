package asyncapi

import "github.com/MarkRosemaker/errpath"

// OperationReply describes the reply part that MAY be applied to an [Operation] object.
// If an operation implements the request/reply pattern, the reply object represents the response message.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationReplyObject
type OperationReply struct {
	// Definition of the address that implementations MUST use for the reply.
	Address *OperationReplyAddressRef `json:"address,omitempty" yaml:"address,omitempty"`
	// A $ref pointer to the definition of the channel in which this operation is performed.
	// When address is specified, the address property of the channel referenced by this property MUST be either null or not defined.
	Channel *ChannelRef `json:"channel,omitempty" yaml:"channel,omitempty"`
	// A list of $ref pointers pointing to the supported Message Objects that can be processed by this operation as reply.
	// It MUST contain a subset of the messages defined in the channel referenced in this operation reply.
	Messages MessageRefList `json:"messages,omitempty" yaml:"messages,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the operation reply for correctness.
func (r *OperationReply) Validate() error {
	if r.Address != nil {
		if err := r.Address.Validate(); err != nil {
			return &errpath.ErrField{Field: "address", Err: err}
		}
	}

	if r.Channel != nil {
		if !r.Channel.isRef() {
			return &errpath.ErrField{Field: "channel", Err: ErrMustBeReference}
		}

		if err := r.Channel.Validate(); err != nil {
			return &errpath.ErrField{Field: "channel", Err: err}
		}

		// when an address is specified, the address of the referenced channel MUST be empty
		if r.Address != nil && r.Channel.Value != nil && r.Channel.Value.Address != "" {
			return &errpath.ErrField{Field: "channel", Err: &errpath.ErrField{
				Field: "address",
				Err: &errpath.ErrInvalid[string]{
					Value:   r.Channel.Value.Address,
					Message: "must be empty when the reply defines an address",
				},
			}}
		}
	}

	if err := r.Messages.Validate(); err != nil {
		return &errpath.ErrField{Field: "messages", Err: err}
	}

	// "It MUST contain a subset of the messages defined in the channel
	// referenced in this operation reply."
	if r.Channel != nil {
		if err := r.Messages.mustBeOfChannel(r.Channel.Value); err != nil {
			return &errpath.ErrField{Field: "messages", Err: err}
		}
	}

	return validateExtensions(r.Extensions)
}

func (l *loader) collectOperationReplyRef(r *OperationReplyRef, ref ref) {
	if !collectRef(l, r, l.replies, ref) {
		return
	}

	if r.Value.Address != nil {
		l.collectOperationReplyAddressRef(r.Value.Address, append(ref, "address"))
	}
}

func (l *loader) resolveOperationReplyRef(r *OperationReplyRef) error {
	return resolveRef(r, l.replies, l.resolveOperationReply)
}

func (l *loader) resolveOperationReply(r *OperationReply) error {
	if r.Address != nil {
		if err := l.resolveOperationReplyAddressRef(r.Address); err != nil {
			return &errpath.ErrField{Field: "address", Err: err}
		}
	}

	if r.Channel != nil {
		if err := l.resolveChannelRef(r.Channel); err != nil {
			return &errpath.ErrField{Field: "channel", Err: err}
		}
	}

	if err := l.resolveMessageRefList(r.Messages); err != nil {
		return &errpath.ErrField{Field: "messages", Err: err}
	}

	return nil
}
