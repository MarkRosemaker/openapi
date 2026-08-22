package asyncapi

import (
	"errors"

	"github.com/MarkRosemaker/errpath"
)

// ErrMustBeReference is returned when the specification demands a reference object
// but the object itself was given, e.g. for the channel of an operation:
// "Please note the `channel` property value MUST be a Reference Object and, therefore,
// MUST NOT contain a Channel Object." ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationObject
var ErrMustBeReference = errors.New("must be a reference object")

// ErrMessageNotOfChannel is returned when an operation or an operation reply refers to a message
// that is not one of the messages of the channel it operates on.
var ErrMessageNotOfChannel = errors.New("must be a message of the channel of this operation")

type (
	// AnySchemaRef is a reference to a schema or an actual schema.
	AnySchemaRef = refOrValue[AnySchema, *AnySchema]
	// ServerRef is a reference to a Server or an actual Server.
	ServerRef = refOrValue[Server, *Server]
	// ServerVariableRef is a reference to a ServerVariable or an actual ServerVariable.
	ServerVariableRef = refOrValue[ServerVariable, *ServerVariable]
	// ChannelRef is a reference to a Channel or an actual Channel.
	ChannelRef = refOrValue[Channel, *Channel]
	// OperationRef is a reference to an Operation or an actual Operation.
	OperationRef = refOrValue[Operation, *Operation]
	// OperationTraitRef is a reference to an OperationTrait or an actual OperationTrait.
	OperationTraitRef = refOrValue[OperationTrait, *OperationTrait]
	// OperationReplyRef is a reference to an OperationReply or an actual OperationReply.
	OperationReplyRef = refOrValue[OperationReply, *OperationReply]
	// OperationReplyAddressRef is a reference to an OperationReplyAddress or an actual OperationReplyAddress.
	OperationReplyAddressRef = refOrValue[OperationReplyAddress, *OperationReplyAddress]
	// MessageRef is a reference to a Message or an actual Message.
	MessageRef = refOrValue[Message, *Message]
	// MessageTraitRef is a reference to a MessageTrait or an actual MessageTrait.
	MessageTraitRef = refOrValue[MessageTrait, *MessageTrait]
	// ParameterRef is a reference to a Parameter or an actual Parameter.
	ParameterRef = refOrValue[Parameter, *Parameter]
	// CorrelationIDRef is a reference to a CorrelationID or an actual CorrelationID.
	CorrelationIDRef = refOrValue[CorrelationID, *CorrelationID]
	// SecuritySchemeRef is a reference to a SecurityScheme or an actual SecurityScheme.
	SecuritySchemeRef = refOrValue[SecurityScheme, *SecurityScheme]
	// TagRef is a reference to a Tag or an actual Tag.
	TagRef = refOrValue[Tag, *Tag]
	// ExternalDocsRef is a reference to an ExternalDocs or an actual ExternalDocs.
	ExternalDocsRef = refOrValue[ExternalDocs, *ExternalDocs]
	// BindingsRef is a reference to a Bindings object or an actual Bindings object.
	BindingsRef = refOrValue[Bindings, *Bindings]

	// AnySchemaRefList is a slice of AnySchemaRef.
	AnySchemaRefList []*AnySchemaRef
	// ServerRefList is a slice of ServerRef.
	ServerRefList []*ServerRef
	// MessageRefList is a slice of MessageRef.
	MessageRefList []*MessageRef
	// SecuritySchemeRefList is a slice of SecuritySchemeRef.
	SecuritySchemeRefList []*SecuritySchemeRef
	// OperationTraitList is a slice of OperationTraitRef.
	OperationTraitList []*OperationTraitRef
	// MessageTraitList is a slice of MessageTraitRef.
	MessageTraitList []*MessageTraitRef
)

func getIndexRef[T any, O referencable[T]](ref *refOrValue[T, O]) int { return ref.idx }
func setIndexRef[T any, O referencable[T]](
	ref *refOrValue[T, O], i int,
) *refOrValue[T, O] {
	ref.idx = i
	return ref
}

// contains reports whether the map holds the given object, no matter whether it was
// defined there or whether the entry is a reference that was resolved to it.
func contains[T any, O referencable[T]](m map[string]*refOrValue[T, O], v O) bool {
	if v == nil {
		return false
	}

	for _, r := range m {
		if r.Value == v {
			return true
		}
	}

	return false
}

// validateRefList validates every entry of a list of references or values.
func validateRefList[T any, O referencable[T]](rs []*refOrValue[T, O]) error {
	for i, r := range rs {
		if err := r.Validate(); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

// mustBeRefs makes sure that every entry of a list is given as a reference.
//
// The specification demands that some properties, e.g. the servers of a channel,
// hold reference objects and not the objects themselves.
func mustBeRefs[T any, O referencable[T]](rs []*refOrValue[T, O]) error {
	for i, r := range rs {
		if !r.isRef() {
			return &errpath.ErrIndex{Index: i, Err: ErrMustBeReference}
		}
	}

	return nil
}

// Validate validates each schema of the list.
func (ss AnySchemaRefList) Validate() error { return validateRefList(ss) }

// Validate validates each server of the list and makes sure they are references.
func (ss ServerRefList) Validate() error {
	if err := mustBeRefs(ss); err != nil {
		return err
	}

	return validateRefList(ss)
}

// Validate validates each message of the list and makes sure they are references.
func (ms MessageRefList) Validate() error {
	if err := mustBeRefs(ms); err != nil {
		return err
	}

	return validateRefList(ms)
}

// mustBeOfChannel makes sure that every message of the list is a message of the given channel,
// as the specification demands of the messages of an operation and of an operation reply.
//
// If the channel is not known, e.g. because it wasn't given, there is nothing to check.
func (ms MessageRefList) mustBeOfChannel(c *Channel) error {
	if c == nil {
		return nil
	}

	for i, m := range ms {
		if contains(c.Messages, m.Value) {
			continue
		}

		return &errpath.ErrIndex{Index: i, Err: ErrMessageNotOfChannel}
	}

	return nil
}

// Validate validates each security scheme of the list.
func (ss SecuritySchemeRefList) Validate() error { return validateRefList(ss) }

// Validate validates each operation trait of the list.
func (ts OperationTraitList) Validate() error { return validateRefList(ts) }

// Validate validates each message trait of the list.
func (ts MessageTraitList) Validate() error { return validateRefList(ts) }

func (l *loader) resolveAnySchemaRefList(ss AnySchemaRefList) error {
	for i, s := range ss {
		if err := l.resolveAnySchemaRef(s); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func (l *loader) resolveServerRefList(ss ServerRefList) error {
	for i, s := range ss {
		if err := l.resolveServerRef(s); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func (l *loader) resolveMessageRefList(ms MessageRefList) error {
	for i, m := range ms {
		if err := l.resolveMessageRef(m); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func (l *loader) resolveSecuritySchemeRefList(ss SecuritySchemeRefList) error {
	for i, s := range ss {
		if err := l.resolveSecuritySchemeRef(s); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func (l *loader) resolveOperationTraitList(ts OperationTraitList) error {
	for i, t := range ts {
		if err := l.resolveOperationTraitRef(t); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func (l *loader) resolveMessageTraitList(ts MessageTraitList) error {
	for i, t := range ts {
		if err := l.resolveMessageTraitRef(t); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
