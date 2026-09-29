package openapi_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/MarkRosemaker/openapi"
)

var jsonOpts = json.JoinOptions([]json.Options{
	// unevaluatedProperties is set to false in most objects according to the OpenAPI specification
	// also protect against deleting unknown fields when overwriting later
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(jsonutil.URLMarshal),
	)),
	json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(jsonutil.URLUnmarshal),
	)),
	jsontext.WithIndent("  "), // indent with two spaces
}...)

// resolveSchemaRefs points every unresolved schema reference within v at an empty schema.
func resolveSchemaRefs(v reflect.Value) {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return
		}

		if r, ok := v.Interface().(*openapi.SchemaRef); ok {
			if r.Value == nil {
				r.Value = &openapi.Schema{}
			}

			return
		}

		resolveSchemaRefs(v.Elem())
	case reflect.Struct:
		for i := range v.NumField() {
			if v.Type().Field(i).IsExported() {
				resolveSchemaRefs(v.Field(i))
			}
		}
	case reflect.Slice, reflect.Array:
		for i := range v.Len() {
			resolveSchemaRefs(v.Index(i))
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			resolveSchemaRefs(v.MapIndex(k))
		}
	default: // nothing within can hold a schema
	}
}

func resolveExamples(examples openapi.Examples) {
	for _, ex := range examples {
		if ex.Ref != nil && ex.Value == nil {
			ex.Value = &openapi.Example{}
		}
	}
}

type validator interface{ Validate() error }

func testJSON(t *testing.T, exampleJSON []byte, v validator) {
	t.Helper()

	switch v.(type) {
	case *openapi.Document:
		doc, err := openapi.LoadFromDataJSON(exampleJSON)
		if err != nil {
			t.Fatalf("load from data: %v", err)
		}

		v = doc

		if _, err = doc.ToJSON(); err != nil {
			t.Fatalf("to json: %v", err)
		}

		if err := doc.WriteToFile(filepath.Join(t.TempDir(), "foo", "openapi.json")); err != nil {
			t.Fatalf("write to file: %v", err)
		}
	default:
		if err := json.Unmarshal(exampleJSON, v, jsonOpts); err != nil {
			t.Fatalf("initial unmarshal: %v", err)
		}
	}

	// manually add unresolved references
	fixReferences(v)

	if err := v.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	// a document is written the way users write it, with the library's own options.
	var (
		b   []byte
		err error
	)
	if doc, ok := v.(*openapi.Document); ok {
		b, err = doc.ToJSON()
	} else {
		b, err = json.Marshal(v, jsonOpts)
	}

	if err != nil {
		t.Fatal(err)
	}

	got := jsontext.Value(b)
	want := jsontext.Value(exampleJSON)

	if err := want.Indent(); err != nil {
		t.Fatal(err)
	}

	if err := got.Indent(); err != nil {
		t.Fatal(err)
	}

	// NOTE: we want to avoid this dependency
	// require.Equal(t, string(want), string(got))

	if !bytes.Equal(want, got) {
		t.Fatalf("not equal, want:\n%s\ngot:\n%s", exampleJSON, got)
	}
}

func fixReferences(v validator) {
	switch v := v.(type) {
	case *openapi.Content:
		for _, mt := range *v {
			resolveExamples(mt.Examples)
		}
	case *openapi.ParameterList:
		for _, p := range *v {
			resolveExamples(p.Value.Examples)
		}
	case *openapi.OperationResponses:
		for _, r := range *v {
			resolveExamples(r.Value.Content[openapi.MediaRangeJSON].Examples)
		}
	case *openapi.Components:
		v.Responses["GeneralError"].Value.Content[openapi.MediaRangeJSON].Schema.Ref.Value = v.Schemas["GeneralError"]
	}

	resolveSchemaRefs(reflect.ValueOf(v))
}
