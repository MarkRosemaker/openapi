package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestSchema_Validate_Errors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		schema *asyncapi.Schema
		want   string
	}{
		"unknown type": {
			&asyncapi.Schema{Type: asyncapi.DataTypes{"struct"}},
			`type ("struct") is invalid, must be one of: "integer", "number", "string", ` +
				`"array", "boolean", "object", "null"`,
		},
		"unknown type among others": {
			&asyncapi.Schema{Type: asyncapi.DataTypes{asyncapi.TypeString, "struct"}},
			`type[1] ("struct") is invalid, must be one of: "integer", "number", "string", ` +
				`"array", "boolean", "object", "null"`,
		},
		"minimum greater than maximum": {
			&asyncapi.Schema{
				Type: asyncapi.DataTypes{asyncapi.TypeInteger},
				Min:  ptr(10.0), Max: ptr(5.0),
			},
			"minimum (10) is invalid: minimum is greater than maximum (10 > 5)",
		},
		"minLength greater than maxLength": {
			&asyncapi.Schema{
				Type:      asyncapi.DataTypes{asyncapi.TypeString},
				MinLength: 10, MaxLength: ptr(uint(5)),
			},
			"minLength (10) is invalid: minLength is greater than maxLength (10 > 5)",
		},
		"minItems greater than maxItems": {
			&asyncapi.Schema{
				Type:     asyncapi.DataTypes{asyncapi.TypeArray},
				MinItems: 10, MaxItems: ptr(uint(5)),
			},
			"minItems (10) is invalid: minItems is greater than maxItems (10 > 5)",
		},
		"minProperties greater than maxProperties": {
			&asyncapi.Schema{
				Type:          asyncapi.DataTypes{asyncapi.TypeObject},
				MinProperties: 10, MaxProperties: ptr(uint(5)),
			},
			"minProperties (10) is invalid: minProperties is greater than maxProperties (10 > 5)",
		},
		"multipleOf zero": {
			&asyncapi.Schema{
				Type:       asyncapi.DataTypes{asyncapi.TypeNumber},
				MultipleOf: ptr(0.0),
			},
			"multipleOf (0) is invalid: must be greater than zero",
		},
		"discriminator that is not a property": {
			&asyncapi.Schema{
				Type:          asyncapi.DataTypes{asyncapi.TypeObject},
				Discriminator: "petType",
			},
			`discriminator ("petType") is invalid: property does not exist`,
		},
		"discriminator that is not required": {
			&asyncapi.Schema{
				Type: asyncapi.DataTypes{asyncapi.TypeObject},
				Properties: asyncapi.Schemas{
					"petType": {Value: &asyncapi.AnySchema{Schema: &asyncapi.Schema{
						Type: asyncapi.DataTypes{asyncapi.TypeString},
					}}},
				},
				Discriminator: "petType",
			},
			`discriminator ("petType") is invalid: property must be required`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.schema.Validate()
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.want {
				t.Fatalf("got: %v, want: %v", err, tc.want)
			}
		})
	}
}

func TestSchema_Validate_Discriminator(t *testing.T) {
	t.Parallel()

	s := &asyncapi.Schema{
		Type: asyncapi.DataTypes{asyncapi.TypeObject},
		Properties: asyncapi.Schemas{
			"petType": {Value: &asyncapi.AnySchema{Schema: &asyncapi.Schema{
				Type: asyncapi.DataTypes{asyncapi.TypeString},
			}}},
		},
		Required:      []string{"petType"},
		Discriminator: "petType",
	}

	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestSchema_Composition(t *testing.T) {
	t.Parallel()

	doc, err := asyncapi.LoadFromFile("examples/v3.1/anyof.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	payload := doc.Components.Messages["testMessages"].Value.Payload.Value.Schema
	if got, want := len(payload.AnyOf), 2; got != want {
		t.Fatalf("got: %d schemas, want: %d", got, want)
	}

	// the references of the composition were resolved
	if payload.AnyOf[0].Value != doc.Components.Schemas["objectWithKey"].Value {
		t.Fatal("the first schema was not resolved")
	}

	if payload.AnyOf[1].Value != doc.Components.Schemas["objectWithKey2"].Value {
		t.Fatal("the second schema was not resolved")
	}
}

func TestSchema_SortMaps(t *testing.T) {
	t.Parallel()

	s := &asyncapi.Schema{
		Type: asyncapi.DataTypes{asyncapi.TypeObject},
		Properties: asyncapi.Schemas{
			"c": {Value: &asyncapi.AnySchema{Schema: &asyncapi.Schema{}}},
			"a": {Value: &asyncapi.AnySchema{Schema: &asyncapi.Schema{}}},
			"b": {Value: &asyncapi.AnySchema{Schema: &asyncapi.Schema{}}},
		},
	}

	s.SortMaps()

	want := []string{"a", "b", "c"}

	i := 0
	for name := range s.Properties.ByIndex() {
		if name != want[i] {
			t.Fatalf("got: %v, want: %v", name, want[i])
		}

		i++
	}
}

func ptr[T any](v T) *T { return &v }
