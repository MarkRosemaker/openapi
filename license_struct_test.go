package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestLicense_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "name": "Apache 2.0",
  "identifier": "Apache-2.0"
}`), &openapi.License{})

	testJSON(t, []byte(`{
  "name": "Apache 2.0",
  "identifier": "Apache-2.0",
  "x-foo": true,
  "x-bar": ["one", "two"]
}`), &openapi.License{})
}

func TestLicense_Validate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		l := openapi.License{}
		if err := l.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `name is required`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("valid", func(t *testing.T) {
		l := openapi.License{Name: "Apache 2.0"}
		if err := l.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("both URL and identifier", func(t *testing.T) {
		l := openapi.License{
			Name:       "Apache 2.0",
			Identifier: "Apache-2.0",
			URL:        mustParseURL("https://www.apache.org/licenses/LICENSE-2.0"),
		}
		if err := l.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `url and identifier are mutually exclusive`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid extension", func(t *testing.T) {
		l := openapi.License{
			Name:       "Apache 2.0",
			Extensions: openapi.Extensions(`{"foo":"bar"}`),
		}

		if err := l.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `foo: ` + openapi.ErrUnknownField.Error(); err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
