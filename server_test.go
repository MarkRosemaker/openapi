package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestServer_Validate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		s := openapi.Server{}
		if err := s.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "url is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid variables", func(t *testing.T) {
		s := openapi.Server{URL: "example.com/{foo}/bar", Variables: openapi.ServerVariables{
			"baz": &openapi.ServerVariable{Default: "mydefault", Enum: []string{"enum1", "enum2"}},
		}}
		if err := s.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `variables["baz"]: default value "mydefault" must exist in the enum's values`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
