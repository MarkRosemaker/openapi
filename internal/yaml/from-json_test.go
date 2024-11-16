package yaml_test

import (
	"errors"
	"io"
	"testing"

	_yaml "github.com/MarkRosemaker/openapi/internal/yaml"
	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

func TestFromJSON(t *testing.T) {
	t.Parallel()

	n, err := _yaml.FromJSON(jsontext.Value(exampleJSON))
	if err != nil {
		t.Fatalf("from json: %v", err)
	}

	got, err := yaml.Marshal(n)
	if err != nil {
		t.Fatalf("marshalling: %v", err)
	}

	const want = `string: Hello, World!
int: 42
float: 3.14
bool: true
bool2: false
null_value: null
list:
    - item1
    - item2
    - item3
dictionary:
    key1: value1
    key2: value2
nested:
    list_of_dicts:
        - name: item1
          value: 1
        - name: item2
          value: 2
    dict_of_lists:
        key1:
            - item1
            - item2
        key2:
            - item3
            - item4
block: |
    This is a block
    style multiline string.
folded: |
    This is a folded style multiline string.
quoted: 'This string contains: special characters, colons: and dashes -'
`
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFromJSON_Error(t *testing.T) {
	t.Parallel()

	t.Run("empty JSON", func(t *testing.T) {
		if _, err := _yaml.FromJSON(jsontext.Value(``)); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, io.EOF) {
			t.Fatalf("got: %q, want: %q", err.Error(), io.EOF)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		if _, err := _yaml.FromJSON(jsontext.Value(`[{`)); err == nil {
			t.Fatal("expected error")
		} else if !errors.Is(err, io.EOF) {
			t.Fatalf("got: %q, want: %q", err.Error(), io.EOF)
		}
	})

	t.Run("invalid mapping key", func(t *testing.T) {
		if _, err := _yaml.FromJSON(jsontext.Value(`[{{}}]`)); err == nil {
			t.Fatal("expected error")
		} else if want := `unexpected kind for mapping key: {`; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("invalid JSON map", func(t *testing.T) {
		if _, err := _yaml.FromJSON(jsontext.Value(`{"foo":`)); err == nil {
			t.Fatal("expected error")
		} else if want := `unexpected EOF`; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})
}
