package openapi_test

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"testing"

	"github.com/MarkRosemaker/jsonutil"
	"github.com/MarkRosemaker/openapi"
)

// jsonOpts are the library's own options, which tests of a single object encode it with.
var jsonOpts = json.JoinOptions(
	json.RejectUnknownMembers(true),
	json.WithMarshalers(json.MarshalToFunc(jsonutil.URLMarshal)),
	json.WithUnmarshalers(json.UnmarshalFromFunc(jsonutil.URLUnmarshal)),
	jsontext.WithIndent("  "),
)

// testJSON loads and validates the document data and writes it back, which must give data again, up to indentation.
func testJSON(t *testing.T, data []byte) *openapi.Document {
	t.Helper()

	doc, err := openapi.LoadFromDataJSON(data)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}

	b, err := doc.ToJSON()
	if err != nil {
		t.Fatalf("to json: %v", err)
	}

	got, want := jsontext.Value(b), jsontext.Value(bytes.Clone(data))
	if err := got.Indent(); err != nil {
		t.Fatal(err)
	}

	if err := want.Indent(); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(got, want) {
		return doc
	}

	gotLines, wantLines := bytes.Split(got, []byte("\n")), bytes.Split(want, []byte("\n"))
	for i := range min(len(gotLines), len(wantLines)) {
		if !bytes.Equal(gotLines[i], wantLines[i]) {
			t.Fatalf("line %d: got %s, want %s", i+1, bytes.TrimSpace(gotLines[i]), bytes.TrimSpace(wantLines[i]))
		}
	}

	t.Fatalf("got %d lines, want %d", len(gotLines), len(wantLines))

	return nil
}
