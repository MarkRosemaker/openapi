package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSecuritySchemes_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		ss  openapi.SecuritySchemes
		err string
	}{
		{openapi.SecuritySchemes{" ": &openapi.SecuritySchemeRef{}},
			`[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.ss.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
