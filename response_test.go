package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestResponse_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		p   openapi.Response
		err string
	}{
		{openapi.Response{}, "description is required"},
		{openapi.Response{
			Description: "some description",
			Headers:     openapi.Headers{"foo": {Value: &openapi.Header{}}},
		}, `headers["foo"]: schema or content is required`},
		{openapi.Response{
			Description: "some description",
			Content: openapi.Content{openapi.MediaRangeJSON: {
				Schema: &openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			}},
		}, `content["application/json"].schema.required is invalid: only valid for object type, got string`},
		{openapi.Response{
			Description: "some description",
			Links:       openapi.Links{"address": {Value: &openapi.Link{}}},
		}, `links.address: operationRef or operationId must be set`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.p.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("expected: %s, got: %s", tc.err, err)
			}
		})
	}
}
