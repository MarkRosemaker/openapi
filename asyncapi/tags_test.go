package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestTags_Validate(t *testing.T) {
	t.Parallel()

	t.Run("no name", func(t *testing.T) {
		t.Parallel()

		tags := asyncapi.Tags{{Value: &asyncapi.Tag{}}}

		err := tags.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		if want := "[0].name is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		t.Parallel()

		tags := asyncapi.Tags{
			{Value: &asyncapi.Tag{Name: "user"}},
			{Value: &asyncapi.Tag{Name: "user"}},
		}

		err := tags.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		want := `[0].name ("user") is invalid: must be unique` + "\n" +
			`[1].name ("user") is invalid: must be unique`
		if err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("external docs without a URL", func(t *testing.T) {
		t.Parallel()

		tags := asyncapi.Tags{{Value: &asyncapi.Tag{
			Name:         "user",
			ExternalDocs: &asyncapi.ExternalDocsRef{Value: &asyncapi.ExternalDocs{}},
		}}}

		err := tags.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		if want := "[0].externalDocs.url is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("valid", func(t *testing.T) {
		t.Parallel()

		tags := asyncapi.Tags{
			{Value: &asyncapi.Tag{Name: "user", Description: " Messages about users. "}},
			{Value: &asyncapi.Tag{
				Name: "signup",
				ExternalDocs: &asyncapi.ExternalDocsRef{Value: &asyncapi.ExternalDocs{
					URL: mustParseURL("https://example.com/docs"),
				}},
			}},
		}

		if err := tags.Validate(); err != nil {
			t.Fatal(err)
		}

		// the description is trimmed
		if got, want := tags[0].Value.Description, "Messages about users."; got != want {
			t.Fatalf("got: %q, want: %q", got, want)
		}
	})
}

func TestExternalDocs_FixScheme(t *testing.T) {
	t.Parallel()

	// the scheme is added if it is missing
	ed := &asyncapi.ExternalDocs{URL: mustParseURL("//example.com/docs")}
	if err := ed.Validate(); err != nil {
		t.Fatal(err)
	}

	if got, want := ed.URL.String(), "https://example.com/docs"; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}
