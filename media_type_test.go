package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
)

func TestMediaType_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		mt  openapi.MediaType
		err string
	}{
		{openapi.MediaType{
			Example:  jsontext.Value("foo"),
			Examples: openapi.Examples{},
		}, `example and examples are mutually exclusive`},
		{openapi.MediaType{
			Examples: openapi.Examples{"foo": invalidExample},
		}, `examples["foo"]: value and externalValue are mutually exclusive`},
		{openapi.MediaType{
			Encoding: openapi.Encodings{
				"foo": &openapi.Encoding{Style: "bar"},
			},
		}, `encoding["foo"].style ("bar") is invalid, must be one of: "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject"`},
		{openapi.MediaType{
			Extensions: jsontext.Value(`{"foo":"bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.mt.Validate(); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != tc.err {
				t.Fatalf("want: %v, got: %v", tc.err, err)
			}
		})
	}
}
