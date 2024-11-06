package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestExternalDocumentation_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "description": "Find more info here",
  "url": "https://example.com"
}`), &openapi.ExternalDocs{})

	testJSON(t, []byte(`{
  "description": "Find more info here",
  "url": "https://example.com",
  "x-foo": "bar",
  "x-bar": 42
}`), &openapi.ExternalDocs{})
}

func TestExternalDocumentation_Validate_Error(t *testing.T) {
	t.Parallel()

	if err := (&openapi.ExternalDocs{}).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `url is required`; want != err.Error() {
		t.Fatalf("unexpected error: %s", err)
	}
}
