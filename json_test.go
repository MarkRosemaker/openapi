package openapi_test

import (
	"bytes"
	"testing"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

func testJSON(t *testing.T, exampleJSON []byte, v interface {
	Validate() error
},
) {
	t.Helper()

	if err := json.Unmarshal(exampleJSON, v, _json.Options); err != nil {
		t.Fatal(err)
	}

	if err := v.Validate(); err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(v, _json.Options)
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

	if !bytes.Equal(want, got) {
		t.Fatalf("not equal, want: %s, got: %s", exampleJSON, got)
	}
}
