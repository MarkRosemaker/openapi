package openapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func docWithSchemas(schemas string) []byte {
	return fmt.Appendf(nil, `{
  "openapi": "3.1.0",
  "info": {
    "title": "t",
    "version": "1"
  },
  "components": {
    "schemas": %s
  }
}`, schemas)
}

// Recursion that descends into the value describes a tree or a list, and is fine.
func TestSchema_Recursion(t *testing.T) {
	t.Parallel()

	for name, schemas := range map[string]string{
		"a list": `{
      "Node": {
        "type": "object",
        "properties": {
          "next": {
            "$ref": "#/components/schemas/Node"
          }
        }
      }
    }`,
		"A to B to C to A": `{
      "A": {
        "type": "object",
        "properties": {
          "b": {
            "$ref": "#/components/schemas/B"
          }
        }
      },
      "B": {
        "type": "array",
        "items": {
          "$ref": "#/components/schemas/C"
        }
      },
      "C": {
        "type": "object",
        "additionalProperties": {
          "$ref": "#/components/schemas/A"
        }
      }
    }`,
		"through allOf": `{
      "A": {
        "allOf": [
          {
            "type": "object",
            "properties": {
              "b": {
                "$ref": "#/components/schemas/B"
              }
            }
          }
        ]
      },
      "B": {
        "$ref": "#/components/schemas/A"
      }
    }`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testJSON(t, docWithSchemas(schemas), &openapi.Document{})
		})
	}
}

// A cycle of $ref, allOf, anyOf, oneOf or not never gets anywhere: it is rejected when loading.
func TestSchema_Cycle(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name, schemas, err string
	}{
		{"A to A", `{
      "A": {"$ref": "#/components/schemas/A"}
    }`, `A → A`},
		{"A to B to A", `{
      "A": {"$ref": "#/components/schemas/B"},
      "B": {"$ref": "#/components/schemas/A"}
    }`, `A → B → A`},
		{"A to B to C to A", `{
      "A": {"$ref": "#/components/schemas/B"},
      "B": {"$ref": "#/components/schemas/C"},
      "C": {"$ref": "#/components/schemas/A"}
    }`, `A → B → C → A`},
		{"the specification's allOf", `{
      "alice": {"allOf": [{"$ref": "#/components/schemas/bob"}]},
      "bob": {"allOf": [{"$ref": "#/components/schemas/alice"}]}
    }`, `alice → bob → alice`},
		{"through not", `{
      "A": {"not": {"anyOf": [{"$ref": "#/components/schemas/B"}]}},
      "B": {"oneOf": [{"$ref": "#/components/schemas/A"}]}
    }`, `A → B → A`},
		{"with sibling keywords", `{
      "A": {"description": "a", "deprecated": true, "$ref": "#/components/schemas/B"},
      "B": {"$ref": "#/components/schemas/A"}
    }`, `A → B → A`},
		{"reached from outside it", `{
      "X": {"$ref": "#/components/schemas/A"},
      "A": {"$ref": "#/components/schemas/B"},
      "B": {"$ref": "#/components/schemas/A"}
    }`, `A → B → A`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := openapi.LoadFromDataJSON(docWithSchemas(tc.schemas))

			first, _, _ := strings.Cut(tc.err, " ")
			want := fmt.Sprintf(`components.schemas[%q]: cycle that never descends into the value: %s`, first, tc.err)
			if err == nil || err.Error() != want {
				t.Fatalf("want: %s\ngot:  %v", want, err)
			}
		})
	}
}

// A document built in code is not checked for cycles, but validating it still ends.
func TestSchema_CycleInCode(t *testing.T) {
	t.Parallel()

	a := &openapi.Schema{}
	b := &openapi.Schema{Ref: &openapi.SchemaRef{Identifier: "#/components/schemas/A", Value: a}}
	a.Ref = &openapi.SchemaRef{Identifier: "#/components/schemas/B", Value: b}

	p := &openapi.Parameter{
		Name:    "p",
		In:      openapi.ParameterLocationQuery,
		Schema:  &openapi.Schema{Ref: &openapi.SchemaRef{Identifier: "#/components/schemas/A", Value: a}},
		Explode: new(true),
	}

	if err := p.Validate(); err == nil {
		t.Fatal("want an error: explode needs an array or object, and the cycle has no type")
	}
}
