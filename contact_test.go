package openapi_test

import (
	"net/url"
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-api-libs/types"
)

func TestContact_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
			"name": "API Support",
			"url": "https://www.example.com/support",
			"email": "support@example.com"
		  }`), &openapi.Contact{})

	testJSON(t, []byte(`{
			"name": "API Support",
			"url": "https://www.example.com/support",
			"email": "support@example.com",
			"x-foo": true,
			"x-bar": ["one", "two"]
		  }`), &openapi.Contact{})
}

func TestContact_Validate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		c := openapi.Contact{}
		if err := c.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("fix URL", func(t *testing.T) {
		c := openapi.Contact{URL: &url.URL{Host: "example.com"}}

		if err := c.Validate(); err != nil {
			t.Fatal(err)
		}

		if want := "https://example.com"; c.URL.String() != want {
			t.Fatalf("url not fixed, want: %q, got: %q", want, c.URL)
		}
	})

	t.Run("valid email", func(t *testing.T) {
		c := openapi.Contact{Email: types.Email("max@example.com")}

		if err := c.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		c := openapi.Contact{Email: types.Email("max[at]example.com")}

		if err := c.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `email: invalid email: "max[at]example.com"`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid extension", func(t *testing.T) {
		c := openapi.Contact{Extensions: openapi.Extensions(`{"foo":"bar"}`)}

		if err := c.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `extension key foo does not have prefix x-`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
