package openapi

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type testStruct struct {
	// A string field
	StringField string `json:"stringField,omitempty" yaml:"stringField,omitempty"`
	// An integer field
	IntegerField int `json:"integerField,omitempty" yaml:"integerField,omitempty"`
	// An Extensions field
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func TestExtensions_valid(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		ext := Extensions(nil)
		if err := validateExtensions(ext); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty", func(t *testing.T) {
		ext := Extensions(jsontext.Value{})
		if err := validateExtensions(ext); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("empty JSON", func(t *testing.T) {
		ext := Extensions(jsontext.Value{'{', '}'})
		if err := validateExtensions(ext); err != nil {
			t.Fatal(err)
		}
	})
}

func TestExtensions_invalid(t *testing.T) {
	t.Parallel()

	t.Run("does not start with x-", func(t *testing.T) {
		ext := Extensions([]byte(`{"bar":true,"x-baz":42}`))
		if err := validateExtensions(ext); err == nil {
			t.Fatal("expected error")
		} else if want := `extension key bar does not have prefix x-`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		ext := Extensions([]byte(`{"x-bar":true,"x-baz":42`))
		if err := validateExtensions(ext); err == nil {
			t.Fatal("expected error")
		} else if want := `unexpected EOF`; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}

func TestExtensions_inline(t *testing.T) {
	t.Parallel()

	ts := &testStruct{
		StringField:  "foo",
		IntegerField: 42,
		Extensions: Extensions([]byte(
			`{"x-bar":true,"x-baz":42,"x-string":"mystring","x-array":["foo","bar"]}`)),
	}

	if err := validateExtensions(ts.Extensions); err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(ts)
	if err != nil {
		t.Fatal(err)
	}

	const want = `{"stringField":"foo","integerField":42,"x-bar":true,"x-baz":42,"x-string":"mystring","x-array":["foo","bar"]}`
	if string(b) != want {
		t.Fatalf("got: %v, want: %v", string(b), want)
	}

	ts2 := &testStruct{}
	if err := json.Unmarshal(b, ts2); err != nil {
		t.Fatal(err)
	}

	b2, err := json.Marshal(ts2)
	if err != nil {
		t.Fatal(err)
	}

	if want != string(b2) {
		t.Fatalf("got: %v, want: %v", string(b2), want)
	}
}
