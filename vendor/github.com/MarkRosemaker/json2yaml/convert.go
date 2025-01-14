package json2yaml

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-json-experiment/json/jsontext"
	"gopkg.in/yaml.v3"
)

// Convert converts a JSON value to a YAML node.
func Convert(b jsontext.Value) (*yaml.Node, error) {
	dec := jsontext.NewDecoder(bytes.NewReader(b))

	n := &yaml.Node{}
	if err := decodeFromJSON(dec, n); err != nil {
		return nil, err
	}

	// check if we reached the end
	if dec.PeekKind() != 0 {
		return nil, fmt.Errorf("expected EOF, got %v", dec.PeekKind())
	}

	return n, nil
}

func decodeFromJSON(dec *jsontext.Decoder, n *yaml.Node) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	switch tkn.Kind() {
	case '"', 't', 'f', 'n', '0':
		n.Kind = yaml.ScalarNode
		n.Value = tkn.String()
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
