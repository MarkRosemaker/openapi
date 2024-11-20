package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestRequestBodies_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		in  openapi.RequestBodies
		err string
	}{
		{
			openapi.RequestBodies{" ": &openapi.RequestBodyRef{}},
			`[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
		},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.in.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %v", tc.err, err)
			}
		})
	}
}
