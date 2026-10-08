package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestContent_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.Content
		err string
	}{
		{openapi.Content{
			"not a real media type": &openapi.MediaType{},
		}, `["not a real media type"]: mime: expected slash after first token`},
		{openapi.Content{
			openapi.MediaRangeJSON: &openapi.MediaType{
				Schema: &openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}},
			},
		}, `["application/json"].schema.required is invalid: only valid for object type, got string`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("expected %q, got %q", tc.err, err)
			}
		})
	}
}
