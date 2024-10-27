package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestLinks_Validate_Error(t *testing.T) {
	t.Parallel()

	if err := (openapi.Links{"foo": {
		Value: &openapi.Link{},
	}}).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `foo: operationRef or operationId must be set`; err.Error() != want {
		t.Fatalf("want: %v, got: %v", want, err)
	}
}
