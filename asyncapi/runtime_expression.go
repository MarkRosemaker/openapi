package asyncapi

import (
	"regexp"

	"github.com/MarkRosemaker/errpath"
)

// A runtime expression allows values to be defined based on information that will be available within the message.
// This mechanism is used by the [CorrelationID] object and the [OperationReplyAddress] object.
//
// The runtime expression is defined by the following [ABNF] syntax:
//
//	expression = ( "$message" "." source )
//	source = ( header-reference | payload-reference )
//	header-reference = "header" ["#" fragment]
//	payload-reference = "payload" ["#" fragment]
//	fragment = a JSON Pointer [RFC6901]
//
// Examples:
//
//	| Source Location         | Example expression               |
//	|-------------------------|----------------------------------|
//	| Message Header Property | `$message.header#/MQMD/CorrelId` |
//	| Message Payload Property | `$message.payload#/messageId`    |
//
// Runtime expressions preserve the type of the referenced value.
// ([Specification])
//
// [ABNF]: https://tools.ietf.org/html/rfc5234
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#runtimeExpression
//
// [RFC6901]: https://tools.ietf.org/html/rfc6901
type RuntimeExpression string

// reRuntimeExpression matches a runtime expression.
var reRuntimeExpression = regexp.MustCompile(`^\$message\.(header|payload)(#(/[^/]*)*)?$`)

// Validate checks that the runtime expression is well-formed.
func (expr RuntimeExpression) Validate() error {
	if expr == "" {
		return &errpath.ErrRequired{}
	}

	if !reRuntimeExpression.MatchString(string(expr)) {
		return &errpath.ErrInvalid[RuntimeExpression]{
			Value:   expr,
			Message: `must be a runtime expression, e.g. "$message.header#/correlationId"`,
		}
	}

	return nil
}
