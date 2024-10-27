package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestReference_JSON(t *testing.T) {
	t.Parallel()

	// Reference Object Example
	testJSON(t, []byte(`{
	"$ref": "#/components/schemas/Pet"
}`), &openapi.Reference{})

	// Relative Schema Document Example
	testJSON(t, []byte(`{
  "$ref": "Pet.json"
}`), &openapi.Reference{})

	// Relative Documents With Embedded Schema Example
	testJSON(t, []byte(`{
  "$ref": "definitions.json#/Pet"
}`), &openapi.Reference{})
}

func TestReference(t *testing.T) {
	if err := (&openapi.Reference{}).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `$ref is required`; err.Error() != want {
		t.Fatalf("want: %s, got: %s", want, err)
	}
}
