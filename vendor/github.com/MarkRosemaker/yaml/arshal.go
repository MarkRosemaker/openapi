package yaml

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"

	"github.com/MarkRosemaker/json2yaml"
	"github.com/MarkRosemaker/yaml2json"
	"gopkg.in/yaml.v3"
)

// Marshal serializes a Go value as a []byte according to the provided
// marshal and encode options (while ignoring unmarshal or decode options).
func Marshal(in any, opts ...json.Options) ([]byte, error) {
	// marshal as json with the given options
	out, err := json.Marshal(in, opts...)
	if err != nil {
		return nil, err
	}

	// convert to yaml
	n, err := json2yaml.Convert(out)
	if err != nil {
		return nil, err
	}

	tagStrings(jsontext.NewDecoder(bytes.NewReader(out)), n)

	// marshal the yaml to bytes
	return yaml.Marshal(n)
}

// tagStrings tags the JSON strings in n as !!str, so yaml quotes those that would read back as another type.
func tagStrings(dec *jsontext.Decoder, n *yaml.Node) {
	tkn, _ := dec.ReadToken() // json2yaml has already decoded the same bytes without error

	switch tkn.Kind() {
	case '"':
		n.Tag = "!!str"
	case '{', '[':
		for _, c := range n.Content {
			tagStrings(dec, c)
		}

		_, _ = dec.ReadToken() // the closing '}' or ']'
	default: // other scalars keep the tag yaml resolves from their value
	}
}

// Unmarshal decodes a []byte input into a Go value according to the provided
// unmarshal and decode options (while ignoring marshal or encode options).
// The output must be a non-nil pointer.
func Unmarshal(in []byte, out any, opts ...json.Options) error {
	// parse input into a yaml document node
	n := &yaml.Node{}
	if err := yaml.Unmarshal(in, n); err != nil {
		return err
	}

	// convert the yaml to json
	val, err := yaml2json.Convert(n)
	if err != nil {
		return err
	}

	return json.Unmarshal(val, out, opts...)
}

// UnmarshalRead deserializes a Go value from an [io.Reader] according to the
// provided unmarshal and decode options (while ignoring marshal or encode options).
// It consumes the entirety of [io.Reader] until [io.EOF] is encountered,
// without reporting an error for EOF. The output must be a non-nil pointer.
// See [Unmarshal] for details about the conversion of JSON into a Go value.
func UnmarshalRead(in io.Reader, out any, opts ...json.Options) error {
	// read the input into a byte slice
	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}

	return Unmarshal(b, out, opts...)
}
