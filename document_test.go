package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-api-libs/types"
)

func TestDocument_Validate(t *testing.T) {
	t.Parallel()

	doc := &openapi.Document{
		OpenAPI: "3.1.0",
		Info: &openapi.Info{
			Title:          "Sample Pet Store App",
			Summary:        "A pet store manager.",
			Description:    "This is a sample server for a pet store.",
			TermsOfService: mustParseURL("https://example.com/terms/"),
			Contact: &openapi.Contact{
				Name:  "API Support",
				URL:   mustParseURL("https://www.example.com/support"),
				Email: types.Email("support@example.com"),
			},
			License: &openapi.License{
				Name: "Apache 2.0",
				URL:  mustParseURL("https://www.apache.org/licenses/LICENSE-2.0.html"),
			},
			Version: "1.0.1",
		},
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
		{"empty", &openapi.Document{}, "openapi field is required"},
		{"invalid openapi version", &openapi.Document{
			OpenAPI: "foo",
		}, `openapi field ("foo") is invalid: must be a valid version (3.0.x or 3.1.x)`},
		{"no info field", &openapi.Document{
			OpenAPI: "3.1.0",
		}, `info is required`},
		{"invalid info", &openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{},
		}, `info.title is required`},
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
