package openapi_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestTags_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		t   openapi.Tags
		err string
	}{
		{openapi.Tags{{}}, "[0].name is required"},
		{openapi.Tags{{
			Name:         "foo",
			ExternalDocs: &openapi.ExternalDocs{},
		}}, `[0].externalDocs.url is required`},
		{openapi.Tags{{
			Name:       "foo",
			Extensions: jsontext.Value(`{"foo":"bar"}`),
		}}, `[0].foo: unknown field or extension without "x-" prefix`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.t.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
