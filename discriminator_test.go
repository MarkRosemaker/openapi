package openapi_test

import (
	"fmt"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func discriminatorDoc(mapping string) []byte {
	return fmt.Appendf(nil, `{
  "openapi": "3.1.0",
  "info": {
    "title": "t",
    "version": "1"
  },
  "components": {
    "schemas": {
      "Pet": {
        "oneOf": [
          {
            "$ref": "#/components/schemas/Dog"
          },
          {
            "$ref": "#/components/schemas/Cat"
          }
        ],
        "discriminator": {
          "propertyName": "kind",
          "mapping": %s
        }
      },
      "Dog": {
        "type": "object"
      },
      "Cat": {
        "type": "object"
      }
    }
  }
}`, mapping)
}

func TestDiscriminator_Mapping(t *testing.T) {
	t.Parallel()

	// a bare name and a reference both stand for a component schema, and the order is kept
	testJSON(t, discriminatorDoc(`{
            "woof": "Dog",
            "meow": "#/components/schemas/Cat"
          }`))
}

func TestDiscriminator_MappingUnresolved(t *testing.T) {
	t.Parallel()

	_, err := openapi.LoadFromDataJSON(discriminatorDoc(`{"moo": "Cow"}`))

	want := `components.schemas["Pet"].discriminator.mapping["moo"]: couldn't resolve "Cow"`
	if err == nil || err.Error() != want {
		t.Fatalf("want: %s\ngot:  %v", want, err)
	}
}

func TestMappingRef(t *testing.T) {
	t.Parallel()

	for value, want := range map[string]string{
		"Dog":                      "#/components/schemas/Dog",
		"#/components/schemas/Cat": "#/components/schemas/Cat",
	} {
		if got := openapi.MappingRef(value); got != want {
			t.Errorf("MappingRef(%q) = %q, want %q", value, got, want)
		}
	}
}

func TestDiscriminator_OnAllOfParent(t *testing.T) {
	t.Parallel()

	// the example the specification gives for allOf: Pet carries the discriminator, Cat and Dog extend it
	const spec = `{
  "openapi": "3.1.0",
  "info": {"title": "t", "version": "1"},
  "components": {
    "schemas": {
      "Pet": {
        "type": "object",
        "required": ["petType"],
        "properties": {"petType": {"type": "string"}},
        "discriminator": {"propertyName": "petType", "mapping": {"dog": "Dog"}}
      },
      "Cat": {"allOf": [{"$ref": "#/components/schemas/Pet"}, {"type": "object", "properties": {"name": {"type": "string"}}}]},
      "Dog": {"allOf": [{"$ref": "#/components/schemas/Pet"}, {"type": "object", "properties": {"bark": {"type": "string"}}}]}
    }
  }
}`

	doc, err := openapi.LoadFromDataJSON([]byte(spec))
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	// with nothing extending it, the discriminator has no alternatives to tell apart
	doc.Components.Schemas["Cat"].AllOf = doc.Components.Schemas["Cat"].AllOf[1:]
	doc.Components.Schemas["Dog"].AllOf = doc.Components.Schemas["Dog"].AllOf[1:]

	want := `components.schemas["Pet"].discriminator is invalid: only valid with oneOf, anyOf or allOf, or on a component schema another extends through allOf`
	if err := doc.Validate(); err == nil || err.Error() != want {
		t.Fatalf("got %v, want %s", err, want)
	}
}
