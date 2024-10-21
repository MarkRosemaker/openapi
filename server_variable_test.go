package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestServerVariable_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
          "enum": [
            "8443",
            "443"
          ],
          "default": "8443"
        }`), &openapi.ServerVariable{})

	testJSON(t, []byte(`{
          "enum": [
            "8443",
            "443"
          ],
          "default": "8443",
  "x-foo": true,
  "x-bar": ["one", "two"]
}`), &openapi.ServerVariable{})
}

func TestServerVariable_Validate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		sv := openapi.ServerVariable{}
		if err := sv.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "default is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("empty enum array", func(t *testing.T) {
		sv := openapi.ServerVariable{Enum: []string{}}
		if err := sv.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "enum array must not be empty"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("empty enum array", func(t *testing.T) {
		sv := openapi.ServerVariable{
			Default: "foo",
			Enum:    []string{"bar", "baz"},
		}
		if err := sv.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `default value "foo" must exist in the enum's values`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
