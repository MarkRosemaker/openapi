package yaml

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

// FromJSON converts a JSON value to a YAML node.
func FromJSON(b jsontext.Value) (*yaml.Node, error) {
	n := &yaml.Node{Kind: yaml.DocumentNode}
	content := &yaml.Node{}
	n.Content = append(n.Content, content)

	dec := jsontext.NewDecoder(bytes.NewReader(b))
	if err := decodeFromJSON(dec, content); err != nil {
		return nil, err
	}

	return n, nil
}

func decodeFromJSON(dec *jsontext.Decoder, n *yaml.Node) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	switch tkn.Kind() {
	case '"':
		n.Kind = yaml.ScalarNode
		n.Value = tkn.String()
	case '0':
		n.Kind = yaml.ScalarNode
		n.Value = fmt.Sprintf("%v", tkn.Float()) // or tkn.Uint()?
	case 't':
		n.Kind = yaml.ScalarNode
		n.Value = "true"
	case 'f':
		n.Kind = yaml.ScalarNode
		n.Value = "false"
	case 'n':
		n.Kind = yaml.ScalarNode
		n.Value = "null"
	case '{':
		n.Kind = yaml.MappingNode
		return decodeMapFromJSON(dec, n)
	case '[':
		n.Kind = yaml.SequenceNode

		for {
			if dec.PeekKind() == ']' {
				_, err := dec.ReadToken() // read the ']' we peeked at
				return err
			}

			el := &yaml.Node{}
			if err := decodeFromJSON(dec, el); err != nil {
				return err
			}

			n.Content = append(n.Content, el)
		}
	default:
		return fmt.Errorf("unsupported token kind: %v", tkn.Kind())
	}

	return nil
}

func decodeMapFromJSON(dec *jsontext.Decoder, n *yaml.Node) error {
	for {
		switch k := dec.PeekKind(); k {
		case '"': // string
			key := &yaml.Node{}
			// ignore error, we know it's a string
			_ = decodeFromJSON(dec, key)
			n.Content = append(n.Content, key)
		case '}':
			_, err := dec.ReadToken() // read the '}' we peeked at
			return err                // done
		case 0:
			return io.EOF
		default:
			return fmt.Errorf("unexpected kind for mapping key: %s", k)
		}

		// write the value
		val := &yaml.Node{}
		if err := decodeFromJSON(dec, val); err != nil {
			return err
		}

		n.Content = append(n.Content, val)
	}
}
