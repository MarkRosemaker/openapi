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
	enc := jsontext.NewEncoder(w)
	if err := encodeToJSON(enc, n); err != nil {
		return nil, err
	}

	return jsontext.Value(w.Bytes()), nil
}

func encodeToJSON(enc *jsontext.Encoder, n *yaml.Node) error {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return fmt.Errorf("expected 1 content node, got %d", len(n.Content))
		}

		return encodeToJSON(enc, n.Content[0])
	case yaml.SequenceNode:
		if err := enc.WriteToken(jsontext.ArrayStart); err != nil {
			return err
		}

		for _, c := range n.Content {
			if err := encodeToJSON(enc, c); err != nil {
				return err
			}
		}

		return enc.WriteToken(jsontext.ArrayEnd)
	case yaml.MappingNode:
		l := len(n.Content)
		if l%2 != 0 {
			return fmt.Errorf("unbalanced mapping node")
		}

		if err := enc.WriteToken(jsontext.ObjectStart); err != nil {
			return err
		}

		for i := 0; i < l; i += 2 {
			if err := encodeToJSON(enc, n.Content[i]); err != nil {
				return err
			}

			if err := encodeToJSON(enc, n.Content[i+1]); err != nil {
				return err
			}
		}

		return enc.WriteToken(jsontext.ObjectEnd)
	case yaml.ScalarNode:
		if n.Style == 0 {
			switch n.Value {
			case "null":
				return enc.WriteToken(jsontext.Null)
			case "true":
				return enc.WriteToken(jsontext.True)
			case "false":
				return enc.WriteToken(jsontext.False)
			}

			if n, err := strconv.ParseInt(n.Value, 10, 64); err == nil {
				return enc.WriteToken(jsontext.Int(n))
			}

			if n, err := strconv.ParseFloat(n.Value, 64); err == nil {
				return enc.WriteToken(jsontext.Float(n))
			}
		}

		return enc.WriteToken(jsontext.String(n.Value))
	// case yaml.AliasNode: // not supported yet (TODO)
	default:
		return fmt.Errorf("unsupported node kind: %v", n.Kind)
	}
}
