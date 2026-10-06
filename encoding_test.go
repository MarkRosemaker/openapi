package openapi_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestEncoding_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.Encoding
		err string
	}{
		{openapi.Encoding{
			Headers: openapi.Headers{
				"foo": {Value: &openapi.Header{}},
			},
		}, `headers["foo"]: schema or content is required`},
		{openapi.Encoding{
			Style: "not a valid style",
		}, `style ("not a valid style") is invalid, must be one of: "matrix", "label", "form", "simple", "spaceDelimited", "pipeDelimited", "deepObject"`},
		{openapi.Encoding{
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
