package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSchema_Ref(t *testing.T) {
	t.Parallel()

	doc, err := openapi.LoadFromDataJSON([]byte(`{
  "openapi": "3.1.0",
  "info": {
    "title": "t",
    "version": "1"
  },
  "paths": {
    "/things": {
      "get": {
        "parameters": [
          {
            "name": "ids",
            "in": "query",
            "schema": {
              "$ref": "#/components/schemas/Ids"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/Alias"
                }
              }
            }
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Alias": {
        "$ref": "#/components/schemas/Thing"
      },
      "Ids": {
        "type": "array",
        "items": {
          "type": "string"
        }
      },
      "Thing": {
        "type": "object",
        "properties": {
          "old": {
            "$ref": "#/components/schemas/Ids",
            "description": "Use ids instead.",
            "deprecated": true
          }
        }
      }
    }
  }
}`))
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	schemas := doc.Components.Schemas

	// a component can be just a reference to another one
	if got := schemas["Alias"].Ref.Value; got != schemas["Thing"] {
		t.Errorf("Alias resolves to %v, want Thing", got)
	}

	// sibling keywords stay on the reference, which still resolves
	old := schemas["Thing"].Properties["old"]
	if old.Ref.Value != schemas["Ids"] || !old.Deprecated || old.Description != "Use ids instead." {
		t.Errorf("got %+v, want a deprecated, described reference to Ids", old)
	}

	// a reference to an array gets the defaults of an array parameter
	p := doc.Paths["/things"].Get.Parameters[0].Value
	if p.Style != openapi.ParameterStyleForm || p.Explode == nil || !*p.Explode {
		t.Errorf("got style %q and explode %v, want form and true", p.Style, p.Explode)
	}

	// a reference to a reference resolves one step at a time
	resp := doc.Paths["/things"].Get.Responses["200"].Value
	if got := resp.Content[openapi.MediaRangeJSON].Schema.Ref.Value; got != schemas["Alias"] {
		t.Errorf("response schema resolves to %v, want Alias", got)
	}
}
