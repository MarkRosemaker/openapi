package yaml

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

// ToJSON converts a YAML node to a JSON.
func ToJSON(n *yaml.Node) (jsontext.Value, error) {
	w := &bytes.Buffer{}
	if err := writeToJSON(w, n); err != nil {
		return nil, err
	}

	return jsontext.Value(w.Bytes()), nil
}

func writeToJSON(w *bytes.Buffer, n *yaml.Node) error {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return fmt.Errorf("expected 1 content node, got %d", len(n.Content))
		}

		return writeToJSON(w, n.Content[0])
	case yaml.SequenceNode:
		_ = w.WriteByte('[')

		lastIdx := len(n.Content) - 1
		for i, c := range n.Content {
			if err := writeToJSON(w, c); err != nil {
				return err
			}

			if i != lastIdx {
				_ = w.WriteByte(',')
			}
		}

		_ = w.WriteByte(']')

		return nil
	case yaml.MappingNode:
		l := len(n.Content)
		if l%2 != 0 {
			return fmt.Errorf("unbalanced mapping node")
		}

		_ = w.WriteByte('{')

		for i := 0; i < l; i += 2 {
			if i > 0 {
				_ = w.WriteByte(',')
			}

			if err := writeToJSON(w, n.Content[i]); err != nil {
				return err
			}

			_ = w.WriteByte(':')

			if err := writeToJSON(w, n.Content[i+1]); err != nil {
				return err
			}
		}

		_ = w.WriteByte('}')

		return nil
	case yaml.ScalarNode:
		if n.Style == 0 {
			switch n.Value {
			case "null", "true", "false":
				_, _ = w.WriteString(n.Value)
				return nil
			}

			// do not quote numbers
			if _, err := strconv.ParseFloat(n.Value, 64); err == nil {
				_, _ = w.WriteString(n.Value)
				return nil
			}
		}

		_, _ = fmt.Fprintf(w, "%q", n.Value)

		return nil
	// case yaml.AliasNode:
	default:
		return fmt.Errorf("unsupported node kind: %v", n.Kind)
	}
}
