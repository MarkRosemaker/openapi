package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
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
	} {
		if err := tc.examples.Validate(); err == nil || err.Error() != tc.err {
			t.Errorf("want: %s, got: %s", tc.err, err)
		}
	}
}
