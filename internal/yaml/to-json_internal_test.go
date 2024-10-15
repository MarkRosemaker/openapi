package yaml

import (
	"io"
	"testing"

	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

func TestEncodeToJSON_Error(t *testing.T) {
	t.Parallel()

	t.Run("sequence node", func(t *testing.T) {
		enc := jsontext.NewEncoder(io.Discard)
		if err := enc.WriteToken(jsontext.ObjectStart); err != nil {
			t.Fatal(err)
		}

		if err := encodeToJSON(enc, &yaml.Node{
			Kind: yaml.SequenceNode,
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "jsontext: missing string for object name"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("mapping node", func(t *testing.T) {
		enc := jsontext.NewEncoder(io.Discard)
		if err := enc.WriteToken(jsontext.ObjectStart); err != nil {
			t.Fatal(err)
		}

		if err := encodeToJSON(enc, &yaml.Node{
			Kind: yaml.MappingNode,
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "jsontext: missing string for object name"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})
}
