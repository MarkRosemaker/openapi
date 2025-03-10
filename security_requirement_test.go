package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSecurityRequirement_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		s   openapi.SecurityRequirement
		err string
	}{
		{openapi.SecurityRequirement{"foo": nil}, `["foo"]: list may be empty but must not be nil`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.s.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
