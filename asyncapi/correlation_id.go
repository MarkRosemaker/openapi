package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// CorrelationID is an object that specifies an identifier at design time that can used for message tracing and correlation.
//
// For specifying and computing the location of a Correlation ID, a [RuntimeExpression] is used.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#correlationIdObject
type CorrelationID struct {
	// An optional description of the identifier. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// REQUIRED. A runtime expression that specifies the location of the correlation ID.
	Location RuntimeExpression `json:"location" yaml:"location"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the correlation ID for correctness.
func (c *CorrelationID) Validate() error {
	c.Description = strings.TrimSpace(c.Description)

	if err := c.Location.Validate(); err != nil {
		return &errpath.ErrField{Field: "location", Err: err}
	}

	return validateExtensions(c.Extensions)
}

func (l *loader) collectCorrelationIDRef(c *CorrelationIDRef, ref ref) {
	collectRef(l, c, l.correlationIDs, ref)
}

func (l *loader) resolveCorrelationIDRef(c *CorrelationIDRef) error {
	return resolveRef(c, l.correlationIDs, nil)
}
