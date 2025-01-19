package openapi_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

var jsonOpts = json.JoinOptions([]json.Options{
	// unevaluatedProperties is set to false in most objects according to the OpenAPI specification
	// also protect against deleting unknown fields when overwriting later
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(jsonutil.URLMarshal))),
	json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(jsonutil.URLUnmarshal))),
	jsontext.WithIndent("  "), // indent with two spaces
}...)

func resolveSchemaRef(s *openapi.SchemaRef) {
	if s != nil && s.Ref != nil && s.Value == nil {
		s.Value = &openapi.Schema{}
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

	b, err := json.Marshal(v, jsonOpts)
	if err != nil {
		t.Fatal(err)
	}

	got := jsontext.Value(b)
	want := jsontext.Value(exampleJSON)

	if err := want.Indent("", "\t"); err != nil {
		t.Fatal(err)
	}

	if err := got.Indent("", "\t"); err != nil {
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
	case *openapi.Callback:
		for _, pi := range *v {
			resolveSchemaRef(pi.Value.Post.RequestBody.Value.Content[openapi.MediaRangeJSON].Schema)
		}
	case *openapi.Content:
		for _, mt := range *v {
			resolveSchemaRef(mt.Schema)
			resolveExamples(mt.Examples)
		}
	case *openapi.RequestBody:
		for _, c := range v.Content {
			resolveSchemaRef(c.Schema)
		}
	case *openapi.ParameterList:
		for _, p := range *v {
			resolveExamples(p.Value.Examples)
		}
	case *openapi.OperationResponses:
		for _, r := range *v {
			mt := r.Value.Content[openapi.MediaRangeJSON]
			resolveSchemaRef(mt.Schema)
			resolveExamples(mt.Examples)
		}
	case *openapi.PathItem:
		resolveSchemaRef(v.Get.Responses["default"].Value.Content[openapi.MediaRangeHTML].Schema)
	case *openapi.Components:
		v.Responses["GeneralError"].Value.Content[openapi.MediaRangeJSON].Schema.Value = v.Schemas["GeneralError"]
	}
}
