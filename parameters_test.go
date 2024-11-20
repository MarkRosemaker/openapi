package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestParameters_Validate_Error(t *testing.T) {
	want := `[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`
	if err := (openapi.Parameters{
		" ": &openapi.ParameterRef{},
	}).Validate(); err == nil || err.Error() != want {
		t.Fatalf("want: %s, got: %s", want, err)
	}
}
