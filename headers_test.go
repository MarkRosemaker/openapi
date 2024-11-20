package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestHeaders_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		header openapi.Headers
		err    string
	}{
		{
			openapi.Headers{" ": &openapi.HeaderRef{Value: &openapi.Header{}}},
			`[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
		},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.header.Validate(); err == nil {
				t.Fatal("expected error")
			} else if want := tc.err; want != err.Error() {
				t.Fatalf("want: %s, got: %s", want, err)
			}
		})
	}
}
