package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestDocument_ValidateStrict(t *testing.T) {
	t.Parallel()

	const paths = `{"openapi": "3.1.0", "info": {"title": "t", "version": "1"}, "paths": {"/pets": {"get": {"responses": {"200": `

	for _, tc := range []struct {
		doc, err string
	}{
		{
			paths + `{"description": "d", "content": {"application/json": {"schema": {"type": "array", "items": {"type": "string", "format": "json"}}}}}}}}}}`,
			`paths["/pets"].GET.responses["200"].content["application/json"].schema.items.format ("json") is invalid, must be one of: ` + validFormats,
		},
		{
			paths + `{"description": "d", "content": {"application/json": {"examples": {"Request Example": {"value": 1}}}}}}}}}}`,
			`paths["/pets"].GET.responses["200"].content["application/json"].examples["Request Example"] ("Request Example") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
		},
		{
			paths + `{"description": "d", "headers": {"X Rate": {"schema": {"type": "integer"}}}}}}}}}`,
			`paths["/pets"].GET.responses["200"].headers["X Rate"] ("X Rate") is invalid: must match the regular expression "^[a-zA-Z0-9\\.\\-_]+$"`,
		},
	} {
		doc, err := openapi.LoadFromData([]byte(tc.doc))
		if err != nil {
			t.Fatal(err)
		}

		if err := doc.Validate(); err != nil {
			t.Fatalf("Validate: %v", err)
		}

		if err := doc.ValidateStrict(); err == nil || err.Error() != tc.err {
			t.Errorf("want: %s, got: %v", tc.err, err)
		}
	}
}
