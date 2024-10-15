package yaml_test

import (
	"bufio"
	"bytes"
	_ "embed"
	"testing"

	_yaml "github.com/MarkRosemaker/openapi/internal/yaml"
	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed example.yaml
	example []byte
	//go:embed example.json
	wantJSON jsontext.Value
)

func TestToJSON(t *testing.T) {
	t.Parallel()

	n := &yaml.Node{}
	if err := yaml.Unmarshal(example, n); err != nil {
		t.Fatal(err)
	}

	got, err := _yaml.ToJSON(n)
	if err != nil {
		t.Fatal(err)
	}

	equalJSON(t, got, wantJSON)
}

func TestToJSON_Error(t *testing.T) {
	t.Parallel()

	t.Run("invalid node kind", func(t *testing.T) {
		if _, err := _yaml.ToJSON(&yaml.Node{
			Kind: yaml.SequenceNode,
			Content: []*yaml.Node{
				{}, // invalid node kind
			},
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "unsupported node kind: 0"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("doc with invalid num of content nodes", func(t *testing.T) {
		if _, err := _yaml.ToJSON(&yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{
				{Kind: yaml.MappingNode},
				{Kind: yaml.ScalarNode},
			},
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "expected 1 content node, got 2"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("unbalanced mapping node", func(t *testing.T) {
		if _, err := _yaml.ToJSON(&yaml.Node{
			Kind: yaml.DocumentNode,
			Content: []*yaml.Node{{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode},
				},
			}},
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "unbalanced mapping node"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("even mapping child node fails", func(t *testing.T) {
		if _, err := _yaml.ToJSON(&yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{Kind: yaml.ScalarNode},
				{}, // invalid node kind
			},
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "unsupported node kind: 0"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})

	t.Run("odd mapping child node fails", func(t *testing.T) {
		if _, err := _yaml.ToJSON(&yaml.Node{
			Kind: yaml.MappingNode,
			Content: []*yaml.Node{
				{}, // invalid node kind
				{Kind: yaml.ScalarNode},
			},
		}); err == nil {
			t.Fatal("expected error")
		} else if want := "unsupported node kind: 0"; err.Error() != want {
			t.Fatalf("got: %q, want: %q", err.Error(), want)
		}
	})
}


func equalJSON(t *testing.T, got, want jsontext.Value) {
	t.Helper()

	if err := got.Indent("", "\t"); err != nil {
		t.Fatalf("formatting got: %v", err)
	}

	if err := want.Indent("", "\t"); err != nil {
		t.Fatalf("formatting want: %v", err)
	}

	gotSc := bufio.NewScanner(bytes.NewReader(got))
	wantSc := bufio.NewScanner(bytes.NewReader(want))

	line := 0
	for gotSc.Scan() {
		line++
		if !wantSc.Scan() {
			t.Fatalf("line %d:\ngot:  %q\nwant: %s", line, gotSc.Text(), "<EOF>")
		}

		if gotSc.Text() != wantSc.Text() {
			t.Fatalf("line %d:\ngot:  %q\nwant: %s", line, gotSc.Text(), wantSc.Text())
		}
	}

	if wantSc.Scan() {
		t.Fatalf("line %d:\ngot:  %s\nwant: %q", line, "<EOF>", wantSc.Text())
	}
}
