package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestExternalDocumentation_Validate_Error(t *testing.T) {
	t.Parallel()

	if err := (&openapi.ExternalDocs{}).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `url is required`; want != err.Error() {
		t.Fatalf("unexpected error: %s", err)
	}
}
