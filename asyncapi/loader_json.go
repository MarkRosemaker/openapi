package asyncapi

import (
	"encoding/json/v2"
	"io"
	"unicode"
)

// LoadFromReaderJSON reads an AsyncAPI specification in JSON format from an io.Reader and
// parses it into a structured format.
//
// "An AsyncAPI document can be JSON or YAML format." ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#format
func (l *loader) LoadFromReaderJSON(r io.Reader) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := json.UnmarshalRead(r, doc, jsonOpts); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// LoadFromDataJSON reads an AsyncAPI specification from a byte array in JSON format and parses it into a structured format.
func LoadFromDataJSON(data []byte) (*Document, error) {
	return newLoader().LoadFromDataJSON(data)
}

// LoadFromDataJSON reads an AsyncAPI specification from a byte array in JSON format and parses it into a structured format.
func (l *loader) LoadFromDataJSON(data []byte) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := json.Unmarshal(data, doc, jsonOpts); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// isJSONRead checks if the data in the reader is JSON.
// NOTE: this is a somewhat naive check, but it should work for most cases.
func isJSONRead(r io.Reader) (bool, error) {
	for {
		var b [1]byte
		_, err := r.Read(b[:])
		if err != nil {
			return false, err
		}

		if unicode.IsSpace(rune(b[0])) {
			continue
		}

		return b[0] == '{', nil
	}
}
