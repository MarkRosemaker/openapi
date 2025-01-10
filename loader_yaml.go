package openapi

import (
	"bytes"
	"fmt"
	"io"

	_yaml "github.com/MarkRosemaker/openapi/internal/yaml"
	"gopkg.in/yaml.v3"
)

// LoadFromReaderYAML reads an OpenAPI specification in YAML format from an io.Reader and parses it into a structured format.
func (l *loader) LoadFromReaderYAML(r io.Reader) (*Document, error) {
	// decode YAML into a Node
	n := &yaml.Node{}
	if err := yaml.NewDecoder(r).Decode(n); err != nil {
		return nil, err
	}

	// convert YAML to JSON
	// we have to use this hack until we have a proper YAML decoder similar to the JSON v2 decoder
	js, err := _yaml.ToJSON(n)
	if err != nil {
		return nil, fmt.Errorf("convert YAML to JSON: %w", err)
	}

	return l.LoadFromDataJSON(js)
}

// LoadFromDataYAML reads an OpenAPI specification from a byte array in YAML format and parses it into a structured format.
func LoadFromDataYAML(data []byte) (*Document, error) {
	return newLoader().LoadFromDataYAML(data)
}

// LoadFromDataYAML reads an OpenAPI specification from a byte array in YAML format and parses it into a structured format.
func (l *loader) LoadFromDataYAML(data []byte) (*Document, error) {
	return l.LoadFromReaderYAML(bytes.NewReader(data))
}
