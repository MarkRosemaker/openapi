package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestReference(t *testing.T) {
	if err := (&openapi.Reference{}).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `$ref is required`; err.Error() != want {
		t.Fatalf("want: %s, got: %s", want, err)
	}
}
