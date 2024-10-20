package openapi_test

import (
	"net/url"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestInfo_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "title": "Sample Pet Store App",
  "summary": "A pet store manager.",
  "description": "This is a sample server for a pet store.",
  "termsOfService": "https://example.com/terms/",
  "contact": {
    "name": "API Support",
    "url": "https://www.example.com/support",
    "email": "support@example.com"
  },
  "license": {
    "name": "Apache 2.0",
    "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
  },
  "version": "1.0.1"
}`), &openapi.Info{})

	testJSON(t, []byte(`{
  "title": "Sample Pet Store App",
  "summary": "A pet store manager.",
  "description": "This is a sample server for a pet store.",
  "termsOfService": "https://example.com/terms/",
  "contact": {
    "name": "API Support",
    "url": "https://www.example.com/support",
    "email": "support@example.com"
  },
  "license": {
    "name": "Apache 2.0",
    "url": "https://www.apache.org/licenses/LICENSE-2.0.html"
  },
  "version": "1.0.1",
  "x-foo": true,
  "x-bar": ["one", "two"]
}`), &openapi.Info{})
}

func TestInfo_Validate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		i := openapi.Info{}
		if err := i.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "title is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("no version", func(t *testing.T) {
		i := openapi.Info{Title: "Sample Pet Store App"}
		if err := i.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "version is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("valid", func(t *testing.T) {
		i := openapi.Info{Title: "Sample Pet Store App", Version: "1.0.1"}
		if err := i.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("fix TOS URL", func(t *testing.T) {
		i := openapi.Info{
			Title:          "Sample Pet Store App",
			TermsOfService: &url.URL{Host: "example.com"},
			Version:        "1.0.1",
		}

		if err := i.Validate(); err != nil {
			t.Fatal(err)
		}

		if want := "https://example.com"; i.TermsOfService.String() != want {
			t.Fatalf("url not fixed, want: %q, got: %q", want, i.TermsOfService)
		}
	})

	t.Run("invalid contact", func(t *testing.T) {
		i := openapi.Info{
			Title:   "Sample Pet Store App",
			Contact: &openapi.Contact{Email: "foo"},
			Version: "1.0.1",
		}

		if err := i.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `contact.email: invalid email: "foo"`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid license", func(t *testing.T) {
		i := openapi.Info{
			Title:   "Sample Pet Store App",
			License: &openapi.License{},
			Version: "1.0.1",
		}

		if err := i.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `license.name is required`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid extension", func(t *testing.T) {
		i := openapi.Info{
			Title:      "Sample Pet Store App",
			Version:    "1.0.1",
			Extensions: openapi.Extensions(`{"foo":"bar"}`),
		}

		if err := i.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := `extension key foo does not have prefix x-`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
