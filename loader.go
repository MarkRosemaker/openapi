package openapi

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-json-experiment/json/jsontext"
)

type loader struct{}

func (l *loader) reset() {}

// newLoader returns an empty Loader
func newLoader() *loader {
	return &loader{}
}

// LoadFromFile reads an OpenAPI specification from a file and parses it into a structured format
func LoadFromFile(location string) (*Document, error) {
	return newLoader().LoadFromFile(location)
}

// LoadFromFile reads an OpenAPI specification from a file and parses it into a structured format
func (l *loader) LoadFromFile(location string) (*Document, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// determine the file type and load accordingly
	switch ext := filepath.Ext(location); ext {
	case ".json":
		return l.LoadFromReaderJSON(f)
	case ".yaml", ".yml":
		return l.LoadFromReaderYAML(f)
	default:
		return nil, fmt.Errorf("unknown file extension: %s", ext)
	}
}

func LoadFromData(data []byte) (*Document, error) {
	return newLoader().LoadFromData(data)
}

// LoadFromData reads an OpenAPI specification from a byte array and parses it into a structured format.
// It will try to determine the format of the data and load it accordingly.
// If you know the format of the data, use LoadFromDataJSON or LoadFromDataYAML instead.
func (l *loader) LoadFromData(data []byte) (*Document, error) {
	if jsontext.Value(data).IsValid() {
		return l.LoadFromDataJSON(data)
	}

	return l.LoadFromDataYAML(data)
}

// LoadFromReader reads an OpenAPI specification from an io.Reader and parses it into a structured format.
// It will try to determine the format of the data and load it accordingly.
// If you know the format of the data, use LoadFromReaderJSON or LoadFromReaderYAML instead.
func LoadFromReader(r io.Reader) (*Document, error) {
	return newLoader().LoadFromReader(r)
}

// LoadFromReader reads an OpenAPI specification from an io.Reader and parses it into a structured format.
// It will try to determine the format of the data and load it accordingly.
// If you know the format of the data, use LoadFromReaderJSON or LoadFromReaderYAML instead.
func (l *loader) LoadFromReader(r io.Reader) (*Document, error) {
	l.reset()

	// by default, assume the data is JSON
	load := l.LoadFromReaderJSON

	// check if the data is JSON, save read data to buffer
	buff := &bytes.Buffer{}
	ok, err := isJSONRead(io.TeeReader(r, buff))
	if err != nil {
		return nil, err
	}

	// if the data is not JSON, use YAML
	if !ok {
		load = l.LoadFromReaderYAML
	}

	// load the document using appropriate loader
	// use multi-reader to combine what was read and the rest of the data
	doc, err := load(io.MultiReader(buff, r))
	if err != nil {
		return nil, err
	}

	// if err := l.resolveRefsIn(doc); err != nil {
	// 	return nil, err
	// }

	return doc, nil
}
