package openapi_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

var invalidExample = &openapi.ExampleRef{
	Value: &openapi.Example{
		Value:         jsontext.Value(`"foo"`),
		ExternalValue: mustParseURL("https://example.com/examples/foo"),
	},
}

func TestExamples_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		examples openapi.Examples
		err      string
	}{
		{
			openapi.Examples{"foo": invalidExample},
			`["foo"]: value and externalValue are mutually exclusive`,
		},
		{openapi.Examples{
			"foo": &openapi.ExampleRef{
				Value: &openapi.Example{
					Extensions: jsontext.Value(`{"bar": "buz"}`),
				},
			},
		}, `["foo"].bar: ` + openapi.ErrUnknownField.Error()},
		{
			openapi.Examples{" ": &openapi.ExampleRef{Value: &openapi.Example{}}},
			`[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
		},
	} {
		if err := tc.examples.Validate(); err == nil || err.Error() != tc.err {
			t.Fatalf("want: %s, got: %s", tc.err, err)
		}
	}
}
