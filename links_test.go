package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestLinks_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		ls  openapi.Links
		err string
	}{
		{
			openapi.Links{"foo": {Value: &openapi.Link{}}},
			`foo: operationRef or operationId must be set`,
		},
		{
			openapi.Links{" ": {Value: &openapi.Link{}}},
			`[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
		},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.ls.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
