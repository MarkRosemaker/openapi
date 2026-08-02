package asyncapi_test

import (
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestExtensions(t *testing.T) {
	t.Parallel()

	data := []byte(`{
  "asyncapi": "3.1.0",
  "info": {
    "title": "Account Service",
    "version": "1.0.0",
    "x-twitter": "@asyncapispec"
  },
  "channels": {
    "userSignedup": {
      "address": "user/signedup"
    }
  },
  "x-linkedin": "async-api"
}`)

	doc, err := asyncapi.LoadFromDataJSON(data)
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.Validate(); err != nil {
		t.Fatal(err)
	}

	if got, want := string(doc.Info.Extensions), `{"x-twitter":"@asyncapispec"}`; got != want {
		t.Fatalf("got: %s, want: %s", got, want)
	}

	// the extensions are written back where they were
	got, err := doc.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(data) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, data)
	}
}

func TestExtensions_Errors(t *testing.T) {
	t.Parallel()

	doc := minimalDocument()
	doc.Info.Extensions = []byte(`{"twitter":"@asyncapispec"}`)

	err := doc.Validate()
	if err == nil {
		t.Fatal("expected error")
	}

	if want := "info.twitter: " + asyncapi.ErrUnknownField.Error(); err.Error() != want {
		t.Fatalf("got: %v, want: %v", err, want)
	}
}
