package openapi_test

import (
	"strings"
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

func TestSchema_Replace(t *testing.T) {
	t.Parallel()

	var props openapi.Schemas
	props.Set("a", &openapi.Schema{Type: openapi.TypeString})
	props.Set("b", &openapi.Schema{Type: openapi.TypeInteger})
	props.Set("c", &openapi.Schema{Type: openapi.TypeBoolean})

	// replacing a with the last of four schemas elsewhere keeps it first
	var other openapi.Schemas
	for _, k := range []string{"w", "x", "y", "z"} {
		other.Set(k, &openapi.Schema{Type: openapi.TypeNumber})
	}

	props["a"].Replace(other["z"])

	var order []string
	for k, s := range props.ByIndex() {
		order = append(order, k+":"+string(s.Type))
	}

	if got, want := strings.Join(order, " "), "a:number b:integer c:boolean"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestSchema_PropertyNamesRef(t *testing.T) {
	t.Parallel()

	doc, err := openapi.LoadFromDataJSON([]byte(`{
  "openapi": "3.1.0",
  "info": {"title": "t", "version": "1"},
  "components": {
    "schemas": {
      "Direction": {"type": "string", "enum": ["north", "south"]},
      "Directions": {
        "type": "object",
        "propertyNames": {"$ref": "#/components/schemas/Direction"}
      }
    }
  }
}`))
	if err != nil {
		t.Fatal(err)
	}

	s := doc.Components.Schemas
	if got := s["Directions"].PropertyNames.Ref.Value; got != s["Direction"] {
		t.Errorf("propertyNames resolves to %v, want Direction", got)
	}
}
