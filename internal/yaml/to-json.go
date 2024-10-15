package yaml

import (
	"fmt"
	"strconv"

	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

// ToJSON converts a YAML node to a JSON.
func ToJSON(n *yaml.Node) (jsontext.Value, error) {
	val := jsontext.Value{}
	if err := writeToJSON(&val, n); err != nil {
		return nil, err
	}

	return val, nil
}

func writeToJSON(val *jsontext.Value, n *yaml.Node) error {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return fmt.Errorf("expected 1 content node, got %d", len(n.Content))
		}

		return writeToJSON(val, n.Content[0])
	case yaml.SequenceNode:
		*val = append(*val, '[')

		lastIdx := len(n.Content) - 1
		for i, c := range n.Content {
			if err := writeToJSON(val, c); err != nil {
				return err
			}

			if i != lastIdx {
				*val = append(*val, ',')
			}
		}

		*val = append(*val, ']')

		return nil
	case yaml.MappingNode:
		l := len(n.Content)
		if l%2 != 0 {
			return fmt.Errorf("unbalanced mapping node")
		}

		*val = append(*val, '{')

		for i := 0; i < l; i += 2 {
			if i > 0 {
				*val = append(*val, ',')
			}

			if err := writeToJSON(val, n.Content[i]); err != nil {
				return err
			}

			*val = append(*val, ':')

			if err := writeToJSON(val, n.Content[i+1]); err != nil {
				return err
			}
		}

		*val = append(*val, '}')

		return nil
	case yaml.ScalarNode:
		if n.Style == 0 {
			switch n.Value {
			case "null", "true", "false":
				*val = append(*val, n.Value...)
				return nil
			}

			// do not quote numbers
			if _, err := strconv.ParseFloat(n.Value, 64); err == nil {
				*val = append(*val, n.Value...)
				return nil
			}
		}

		*val = append(*val, fmt.Sprintf("%q", n.Value)...)

		return nil
	// case yaml.AliasNode:
	default:
		return fmt.Errorf("unsupported node kind: %v", n.Kind)
	}
}
