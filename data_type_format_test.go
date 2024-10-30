package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

const validFormats = `"int32", "int64", "float", "double", "byte", "binary", "date", "date-time", "duration", "email", "password", "uuid", "uri", "uriref", "zip-code"`

func TestFormat(t *testing.T) {
	// test a valid data type format
	if err := openapi.FormatDateTime.Validate(); err != nil {
		t.Fatal(err)
	}

	// test an invalid data type format
	err := openapi.Format("foo").Validate()
	if err == nil {
		t.Fatal("expected an error for an invalid data type")
	}

	err = &openapi.ErrField{Field: "format", Err: err}
	if want := `format ("foo") is invalid, must be one of: ` + validFormats; want != err.Error() {
		t.Fatalf("want: %s, got: %s", want, err)
	}
}
