package openapi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-api-libs/types"
)

func TestDocument_JSON(t *testing.T) {
	t.Parallel()

	// security optional via an empty security requirement
	testJSON(t, []byte(`{
	  "openapi": "3.1.0",
	  "info": {
	    "title": "Sample Pet Store App",
		"version": "1.0.0"
	  },
	  "paths": {"/": {}},
	  "security": [{}]
}`), &openapi.Document{})
}

func TestDocument_Validate(t *testing.T) {
	t.Parallel()

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

func TestDocument_Examples(t *testing.T) {
	t.Parallel()

	if err := filepath.Walk("examples", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) == ".txt" {
			return err
		}

		t.Run(path, func(t *testing.T) {
			original, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}

			switch filepath.Ext(path) {
			case ".json":
				testJSON(t, original, &openapi.Document{})
			case ".yaml":
				doc, err := openapi.LoadFromData(original)
				if err != nil {
					t.Fatal(err)
				}

				if err := doc.Validate(); err != nil {
					t.Fatal(err)
				}
			}
		})

		return nil
	}); err != nil {
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
		}, `openapi ("foo") is invalid: must be a valid version (3.0.x or 3.1.x)`},
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
				Schemas: openapi.Schemas{"Pet": &openapi.Schema{}},
			},
		}, `components.schemas["Pet"].type is required`},
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
