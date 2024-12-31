package yaml

import (
	"errors"
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

		synErr := &jsontext.SyntacticError{}
		if err := encodeToJSON(enc, &yaml.Node{
			Kind: yaml.SequenceNode,
		}); err == nil {
			t.Fatal("expected error")
		} else if !errors.As(err, &synErr) {
			t.Fatalf("got: %T, want: %T", err, synErr)
		} else if synErr.JSONPointer != "" {
			t.Fatalf("got: %q", synErr.JSONPointer)
		} else if want := `object member name must be a string`; synErr.Err.Error() != want {
			t.Fatalf("got: %q, want: %q", want, synErr.Err)
		}
	})

	t.Run("mapping node", func(t *testing.T) {
		enc := jsontext.NewEncoder(io.Discard)
		if err := enc.WriteToken(jsontext.ObjectStart); err != nil {
			t.Fatal(err)
		}

		synErr := &jsontext.SyntacticError{}

		if err := encodeToJSON(enc, &yaml.Node{
			Kind: yaml.MappingNode,
		}); err == nil {
			t.Fatal("expected error")
		} else if !errors.As(err, &synErr) {
			t.Fatalf("got: %T, want: %T", err, synErr)
		} else if synErr.JSONPointer != "" {
			t.Fatalf("got: %q", synErr.JSONPointer)
		} else if want := `object member name must be a string`; synErr.Err.Error() != want {
			t.Fatalf("got: %q, want: %s", want, synErr.Err)
		}
	})
}
