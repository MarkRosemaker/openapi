package openapi

import (
	"errors"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// The `Link object` represents a possible design-time link for a response.
// The presence of a link does not guarantee the caller's ability to successfully invoke it, rather it provides a known relationship and traversal mechanism between responses and other operations.
//
// Unlike _dynamic_ links (i.e. links provided **in** the response payload), the OAS linking mechanism does not require link information in the runtime response.
//
// For computing links, and providing instructions to execute them, a [runtime expression](#runtime-expressions) is used for accessing values in an operation and using them as parameters while invoking the linked operation.
//
// Clients follow all links at their discretion.
// Neither permissions, nor the capability to make a successful call to that link, is guaranteed solely by the existence of a relationship.
type Link struct {
	// A relative or absolute URI reference to an OAS operation. This field is mutually exclusive of the `operationId` field, and MUST point to an Operation Object. Relative `operationRef` values MAY be used to locate an existing Operation Object in the OpenAPI definition. See the rules for resolving Relative References.
	// Note that in the use of `operationRef`, the _escaped forward-slash_ is necessary when using JSON references.
	// Because of the potential for name clashes, the `operationRef` syntax is preferred for OpenAPI documents with external references.
	OperationRef string `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	// The name of an _existing_, resolvable OAS operation, as defined with a unique `operationId`.
	// This field is mutually exclusive of the `operationRef` field.
	OperationID string `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	// A map representing parameters to pass to an operation as specified with `operationId` or identified via `operationRef`. The key is the parameter name to be used, whereas the value can be a constant or an expression to be evaluated and passed to the linked operation.
	// The parameter name can be qualified using the parameter location `[{in}.]{name}` for operations that use the same parameter name in different locations (e.g. path.id).
	Parameters LinkParameters `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	// A literal value or {expression} to use as a request body when calling the target operation.
	RequestBody RuntimeExpression `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	// A description of the link. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`
}

func (l *Link) Validate() error {
	if l.OperationRef != "" && l.OperationID != "" {
		return errors.New("operationRef and operationId are mutually exclusive")
	}

	// A linked operation MUST be identified using either an `operationRef` or `operationId`.
	if l.OperationRef == "" && l.OperationID == "" {
		return errors.New("operationRef or operationId must be set")
	}

	if err := l.Parameters.Validate(); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	// NOTE: We don't check RequestBody yet.

	l.Description = strings.TrimSpace(l.Description)

	return validateExtensions(l.Extensions)
}

func (l *loader) collectLinkRef(lr *LinkRef, ref ref) {
	if lr.Value != nil {
		l.collectLink(lr.Value, ref)
	}
}

func (l *loader) collectLink(link *Link, ref ref) { l.links[ref.String()] = link }

func (l *loader) resolveLinkRef(lr *LinkRef) error {
	return resolveRef(lr, l.links, nil)
}
