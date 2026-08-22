package asyncapi_test

import (
	"encoding/json/v2"
	"strings"
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

// loadSchema loads a single schema by loading a document that holds it in its components.
func loadSchema(t *testing.T, schema string) *asyncapi.AnySchema {
	t.Helper()

	doc, err := asyncapi.LoadFromDataJSON([]byte(
		`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
			`"components":{"schemas":{"test":` + schema + `}}}`,
	))
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	return doc.Components.Schemas["test"].Value
}

func TestAnySchema_Schema(t *testing.T) {
	t.Parallel()

	s := loadSchema(t, `{"type":"string","format":"date-time"}`)

	if s.SchemaFormat != "" {
		t.Fatalf("got: %v, want no schema format", s.SchemaFormat)
	}

	if s.Schema == nil {
		t.Fatal("expected an AsyncAPI schema")
	}

	if got, want := s.Schema.Type.String(), "string"; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	if got, want := s.Schema.Format, asyncapi.FormatDateTime; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	if !s.Schema.Format.IsKnown() {
		t.Fatalf("%q should be a known format", s.Schema.Format)
	}
}

func TestAnySchema_MultipleTypes(t *testing.T) {
	t.Parallel()

	s := loadSchema(t, `{"type":["string","null"]}`)

	if got, want := s.Schema.Type.String(), "string, null"; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}

	if !s.Schema.Type.Contains(asyncapi.TypeNull) {
		t.Fatal("expected the schema to allow null")
	}
}

func TestAnySchema_Boolean(t *testing.T) {
	t.Parallel()

	s := loadSchema(t, `false`)

	if s.Schema == nil || s.Schema.Boolean == nil {
		t.Fatal("expected a boolean schema")
	}

	if *s.Schema.Boolean {
		t.Fatal("expected the schema to be false")
	}

	// a boolean schema is written back as a boolean
	got, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}

	if want := "false"; string(got) != want {
		t.Fatalf("got: %s, want: %s", got, want)
	}
}

func TestAnySchema_MultiFormat(t *testing.T) {
	t.Parallel()

	t.Run("avro", func(t *testing.T) {
		t.Parallel()

		s := loadSchema(t, `{"schemaFormat":"application/vnd.apache.avro;version=1.9.0",`+
			`"schema":{"type":"record","name":"User"}}`)

		if got, want := s.SchemaFormat, asyncapi.SchemaFormatAvro; got != want {
			t.Fatalf("got: %v, want: %v", got, want)
		}

		if !s.SchemaFormat.IsKnown() {
			t.Fatalf("%q should be a known schema format", s.SchemaFormat)
		}

		if s.SchemaFormat.IsAsyncAPI() {
			t.Fatalf("%q is not an AsyncAPI schema format", s.SchemaFormat)
		}

		// a schema in another format is not parsed but kept as it is
		if s.Schema != nil {
			t.Fatal("expected no AsyncAPI schema")
		}

		if got, want := string(s.Raw), `{"type":"record","name":"User"}`; got != want {
			t.Fatalf("got: %s, want: %s", got, want)
		}
	})

	t.Run("asyncapi", func(t *testing.T) {
		t.Parallel()

		s := loadSchema(t, `{"schemaFormat":"application/vnd.aai.asyncapi+json;version=3.1.0",`+
			`"schema":{"type":"string"}}`)

		if !s.SchemaFormat.IsAsyncAPI() {
			t.Fatalf("%q is an AsyncAPI schema format", s.SchemaFormat)
		}

		// an AsyncAPI schema is parsed even if the format is given explicitly
		if s.Schema == nil {
			t.Fatal("expected an AsyncAPI schema")
		}

		if got, want := s.Schema.Type.String(), "string"; got != want {
			t.Fatalf("got: %v, want: %v", got, want)
		}
	})
}

func TestAnySchema_Errors(t *testing.T) {
	t.Parallel()

	t.Run("not an object", func(t *testing.T) {
		t.Parallel()

		_, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"components":{"schemas":{"test":"string"}}}`,
		))
		if err == nil {
			t.Fatal("expected error")
		}

		want := "schema must be an object or a boolean, got string"
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("got: %v, want it to contain: %v", err, want)
		}
	})

	t.Run("multi format without a schema", func(t *testing.T) {
		t.Parallel()

		doc, err := asyncapi.LoadFromDataJSON([]byte(
			`{"asyncapi":"3.1.0","info":{"title":"foo","version":"1.0.0"},` +
				`"components":{"schemas":{"test":` +
				`{"schemaFormat":"application/vnd.apache.avro;version=1.9.0"}}}}`,
		))
		if err != nil {
			t.Fatal(err)
		}

		err = doc.Validate()
		if err == nil {
			t.Fatal("expected error")
		}

		if want := `components.schemas["test"].schema is required`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
