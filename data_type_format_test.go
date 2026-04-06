package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSchema_UnknownFormat(t *testing.T) {
	t.Parallel()

	// Unknown formats are accepted as annotations per spec — tools MUST NOT generate an error.
	// See: https://spec.openapis.org/oas/v3.2.0.html#data-types
	for _, f := range []openapi.Format{
		"hostname", "idn-email", "idn-hostname", "iri", "iri-reference",
		"json-pointer", "relative-json-pointer", "regex", "uri-template",
		openapi.Format("my-custom-format"),
	} {
		s := &openapi.Schema{Type: openapi.TypeString, Format: f}
		if err := s.Validate(); err != nil {
			t.Fatalf("format %q should be accepted: %s", f, err)
		}
	}
}
