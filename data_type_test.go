package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/openapi"
)

func TestDataType(t *testing.T) {
	// test a valid data type
	if err := openapi.TypeInteger.Validate(); err != nil {
		t.Fatal(err)
	}

	// test an invalid data type
	err := openapi.DataType("foo").Validate()
	if err == nil {
		t.Fatal("expected an error for an invalid data type")
	}

	err = &errpath.ErrField{Field: "type", Err: err}
	if want := `type ("foo") is invalid, must be one of: "integer", "number", "string", "array", "boolean", "object"`; want != err.Error() {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}
