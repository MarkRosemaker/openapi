package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestPath_Validate_Error(t *testing.T) {
	t.Parallel()

	if err := (openapi.Path("")).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := "path must not be empty"; err.Error() != want {
		t.Fatalf("want: %s, got: %s", want, err)
	}

	if err := (openapi.Path("foo")).Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := "path must start with a /"; err.Error() != want {
		t.Fatalf("want: %s, got: %s", want, err)
	}
}
