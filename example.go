package openapi

import (
	"fmt"
	"net/url"

	"github.com/go-json-experiment/json/jsontext"
)

// Example represents an example of a schema.
//
// In all cases, the example value is expected to be compatible with the type schema of its associated value.
// Tooling implementations MAY choose to validate compatibility automatically, and reject the example value(s) if incompatible.
type Example struct {
	// Short description for the example.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// Long description for the example. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Embedded literal example. The `value` field and `externalValue` field are mutually exclusive. To represent examples of media types that cannot naturally represented in JSON or YAML, use a string value to contain the example, escaping where necessary.
	Value jsontext.Value `json:"value,omitempty" yaml:"value,omitempty"`
	// A URI that points to the literal example. This provides the capability to reference examples that cannot easily be included in JSON or YAML documents.
	// The `value` field and `externalValue` field are mutually exclusive. See the rules for resolving [Relative References](#relative-references-in-uris).
	ExternalValue *url.URL `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`
}

func (ex *Example) Validate() error {
	if ex.Value != nil && ex.ExternalValue != nil {
		return fmt.Errorf("value and externalValue are mutually exclusive")
	}

	return validateExtensions(ex.Extensions)
}

func (l *loader) collectExampleRef(ex *ExampleRef, ref ref) {
	if ex.Value != nil {
		l.collectExample(ex.Value, ref)
	}
}

func (l *loader) collectExample(ex *Example, ref ref) {
	l.examples[ref.String()] = ex
}

func (l *loader) resolveExampleRef(ex *ExampleRef) error {
	return resolveRef(ex, l.examples, nil)
}
