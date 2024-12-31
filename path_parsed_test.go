package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestParsedPath(t *testing.T) {
	t.Parallel()
	for _, path := range []openapi.Path{
		"foo/bar/baz",
		"/users/{userid}/address",
		"/2.0/repositories/{username}/{slug}",
	} {
		t.Run(string(path), func(t *testing.T) {
			pp := path.Parse()

			if pp.String() != string(path) {
				t.Fatalf("got: %q, want: %q", pp, path)
			}
		})
	}
}
