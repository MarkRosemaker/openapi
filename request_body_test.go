package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestRequestBody_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		p   openapi.RequestBody
		err string
	}{
		{openapi.RequestBody{}, "content is required"},
		{openapi.RequestBody{
			Content: openapi.Content{"foo; bar": &openapi.MediaType{}},
		}, `content["foo; bar"]: mime: invalid media parameter`},
		{openapi.RequestBody{
			Content:    openapi.Content{openapi.MediaRangeJSON: &openapi.MediaType{}},
			Extensions: []byte(`{"foo": "bar"}`),
		}, "foo: " + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.p.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("expected %q, got %q", tc.err, err)
			}
		})
	}
}
