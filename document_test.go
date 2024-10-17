package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestDocumentValidate(t *testing.T) {
	t.Parallel()

	doc := &openapi.Document{
		OpenAPI: "3.1.0",
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDocumentValidate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		doc    *openapi.Document
		errStr string
	}{
		{"empty", &openapi.Document{}, "openapi.version is required"},
		{"invalid version", &openapi.Document{
			OpenAPI: "foo",
		}, `openapi.version ("foo") is invalid: must be a valid version (3.0.x or 3.1.x)`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tc.doc.Validate(); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != tc.errStr {
				t.Fatalf("got: %v, want: %v", err, tc.errStr)
			}
		})
	}
}
