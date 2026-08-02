package asyncapi

import (
	"io"

	"github.com/MarkRosemaker/yaml"
)

// LoadFromReaderYAML reads an AsyncAPI specification in YAML format from an io.Reader and
// parses it into a structured format.
//
// "An AsyncAPI document can be JSON or YAML format." ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#format
func (l *loader) LoadFromReaderYAML(r io.Reader) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := yaml.UnmarshalRead(r, doc, jsonOpts); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// LoadFromDataYAML reads an AsyncAPI specification from a byte array in YAML format and parses it into a structured format.
func LoadFromDataYAML(data []byte) (*Document, error) {
	return newLoader().LoadFromDataYAML(data)
}

// LoadFromDataYAML reads an AsyncAPI specification from a byte array in YAML format and parses it into a structured format.
func (l *loader) LoadFromDataYAML(data []byte) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := yaml.Unmarshal(data, doc, jsonOpts); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, nil
}
