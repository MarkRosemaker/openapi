package openapi_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-api-libs/types"
)

// TestDocument_Golden loads testdata/openapi.json, a document that uses each object and feature the library models once,
// validates it and writes it back: it must come out as it went in. The file is edited by hand, never regenerated, so a
// difference is a change in what the library reads, resolves or writes.
func TestDocument_Golden(t *testing.T) {
	t.Parallel()

	want, err := os.ReadFile("testdata/openapi.json")
	if err != nil {
		t.Fatal(err)
	}

	doc := testJSON(t, want)

	// references are resolved wherever they may appear, inside an operation's callbacks too
	pets := doc.Paths["/pets"]
	if pets.Post.Callbacks["onData"].Value == nil ||
		(*pets.Post.Callbacks["fixedServer"].Value)["http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}"].Value == nil ||
		doc.Webhooks["newPet"].Value.Post.RequestBody.Value == nil {
		t.Error("a reference is not resolved")
	}

	// a file holds what ToJSON returns
	path := filepath.Join(t.TempDir(), "dir", "openapi.json")
	if err := doc.WriteToFile(path); err != nil {
		t.Fatal(err)
	}

	written, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if b, err := doc.ToJSON(); err != nil || !bytes.Equal(written, b) {
		t.Errorf("the file differs from ToJSON: %v", err)
	}
}

func TestDocument_Validate(t *testing.T) {
	t.Parallel()

	// OAS 3.2.x documents must be accepted.
	// See: https://spec.openapis.org/oas/v3.2.0.html#versions-and-deprecation
	if err := (&openapi.Document{
		OpenAPI: "3.2.0",
		Info:    &openapi.Info{Title: "Test", Version: "1.0"},
		Paths:   openapi.Paths{"/": {}},
	}).Validate(); err != nil {
		t.Fatal(err)
	}

	doc := &openapi.Document{
		OpenAPI: "3.1.0",
		Info: &openapi.Info{
			Title:          "Sample Pet Store App",
			Summary:        "A pet store manager.",
			Description:    "This is a sample server for a pet store.",
			TermsOfService: mustParseURL("https://example.com/terms/"),
			Contact: &openapi.Contact{
				Name:  "API Support",
				URL:   mustParseURL("https://www.example.com/support"),
				Email: types.Email("support@example.com"),
			},
			License: &openapi.License{
				Name: "Apache 2.0",
				URL:  mustParseURL("https://www.apache.org/licenses/LICENSE-2.0.html"),
			},
			Version: "1.0.1",
		},
		Paths: openapi.Paths{"/": {}},
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDocumentValidate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		doc *openapi.Document
		err string
	}{
		{&openapi.Document{}, "openapi is required"},
		{&openapi.Document{
			OpenAPI: "foo",
		}, `openapi ("foo") is invalid: must be a valid version (3.0.x, 3.1.x or 3.2.x)`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
		}, `info is required`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{},
		}, `info.title is required`},
		{&openapi.Document{
			OpenAPI:           "3.1.0",
			Info:              &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			JSONSchemaDialect: mustParseURL("https://example.com"),
		}, `jsonSchemaDialect ("https://example.com") is invalid, must be one of: "https://spec.openapis.org/oas/3.1/dialect/base"`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Servers: openapi.Servers{{}},
		}, `servers[0].url is required`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Paths:   openapi.Paths{"": {}},
		}, `paths[""]: path must not be empty`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Webhooks: openapi.Webhooks{"myWebhook": {
				Value: &openapi.PathItem{
					Parameters: openapi.ParameterList{
						{Value: &openapi.Parameter{Name: "foo"}},
					},
				},
			}},
		}, `webhooks["myWebhook"].parameters[0].in is required`},
		{&openapi.Document{
			OpenAPI:  "3.1.0",
			Info:     &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Paths:    openapi.Paths{},
			Webhooks: openapi.Webhooks{},
			Components: openapi.Components{
				Schemas: openapi.Schemas{},
			},
		}, openapi.ErrEmptyDocument.Error()},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Components: openapi.Components{
				Schemas: openapi.Schemas{"Pet": &openapi.Schema{Type: openapi.TypeString, Required: []string{"id"}}},
			},
		}, `components.schemas["Pet"].required is invalid: only valid for object type, got string`},
		{&openapi.Document{
			OpenAPI:  "3.1.0",
			Info:     &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Paths:    openapi.Paths{"/": {}},
			Security: openapi.SecurityRequirements{{"": {}}},
		}, `security[0][""]: empty security scheme name`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Paths:   openapi.Paths{"/": {}},
			Tags:    openapi.Tags{{}},
		}, `tags[0].name is required`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Paths:   openapi.Paths{"/": {}},
			Tags:    openapi.Tags{{Name: "foo"}, {Name: "foo"}},
		}, `tags[0].name ("foo") is invalid: must be unique
tags[1].name ("foo") is invalid: must be unique`},
		{&openapi.Document{
			OpenAPI:      "3.1.0",
			Info:         &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Paths:        openapi.Paths{"/": {}},
			ExternalDocs: &openapi.ExternalDocs{},
		}, `externalDocs.url is required`},
		{&openapi.Document{
			OpenAPI: "3.1.0",
			Info:    &openapi.Info{Title: "Sample API", Version: "1.0.0"},
			Components: openapi.Components{
				Callbacks: openapi.CallbackRefs{" ": &openapi.CallbackRef{}},
			},
		}, `components.callbacks[" "] (" ") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			t.Parallel()

			if err := tc.doc.Validate(); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != tc.err {
				t.Fatalf("got: %v, want: %v", err, tc.err)
			}
		})
	}
}
