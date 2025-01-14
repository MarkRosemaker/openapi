package openapi

import (
	"io"

	"github.com/MarkRosemaker/yaml"
)

// LoadFromReaderYAML reads an OpenAPI specification in YAML format from an io.Reader and parses it into a structured format.
func (l *loader) LoadFromReaderYAML(r io.Reader) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := yaml.UnmarshalRead(r, doc, jsonOpts); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, doc.Validate()
}

// LoadFromDataYAML reads an OpenAPI specification from a byte array in YAML format and parses it into a structured format.
func LoadFromDataYAML(data []byte) (*Document, error) {
	return newLoader().LoadFromDataYAML(data)
}

// LoadFromDataYAML reads an OpenAPI specification from a byte array in YAML format and parses it into a structured format.
func (l *loader) LoadFromDataYAML(data []byte) (*Document, error) {
	l.reset()

	doc := &Document{}
	if err := yaml.Unmarshal(data, doc, jsonOpts); err != nil {
		return nil, err
	}

	if err := l.collectResolveRefs(doc); err != nil {
		return nil, err
	}

	return doc, doc.Validate()
}
