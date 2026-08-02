package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// OperationReplyAddress is an object that specifies where an operation has to send the reply.
//
// For specifying and computing the location of a reply address, a [RuntimeExpression] is used.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationReplyAddressObject
type OperationReplyAddress struct {
	// An optional description of the address. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// REQUIRED. A runtime expression that specifies the location of the reply address.
	Location RuntimeExpression `json:"location" yaml:"location"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the reply address for correctness.
func (a *OperationReplyAddress) Validate() error {
	a.Description = strings.TrimSpace(a.Description)

	if err := a.Location.Validate(); err != nil {
		return &errpath.ErrField{Field: "location", Err: err}
	}

	return validateExtensions(a.Extensions)
}

func (l *loader) collectOperationReplyAddressRef(a *OperationReplyAddressRef, ref ref) {
	collectRef(l, a, l.replyAddresses, ref)
}

func (l *loader) resolveOperationReplyAddressRef(a *OperationReplyAddressRef) error {
	return resolveRef(a, l.replyAddresses, nil)
}
