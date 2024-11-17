package openapi

import (
	"io"
	"unicode"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
)

// LoadFromReaderJSON reads an OpenAPI specification in JSON format from an io.Reader and parses it into a structured format.
func (l *loader) LoadFromReaderJSON(r io.Reader) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := json.UnmarshalRead(r, doc, _json.Options); err != nil {
		return nil, err
	}

	// if err := l.resolveRefsIn(doc); err != nil {
	// 	return nil, err
	// }

	return doc, nil
}

// LoadFromDataJSON reads an OpenAPI specification from a byte array in JSON format and parses it into a structured format.
func LoadFromDataJSON(data []byte) (*Document, error) {
	return newLoader().LoadFromDataJSON(data)
}

// LoadFromDataJSON reads an OpenAPI specification from a byte array in JSON format and parses it into a structured format.
func (l *loader) LoadFromDataJSON(data []byte) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := json.Unmarshal(data, doc, _json.Options); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, doc.Validate()
}

// isJSONRead checks if the data in the reader is JSON
// NOTE: this is a somewhat naive check, but it should work for most cases
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
