package openapi

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-json-experiment/json/jsontext"
)

// loader helps deserialize an OpenAPI v3 document
type loader struct {
	schemas         map[string]*Schema
	headers         map[string]*Header
	responses       map[string]*Response
	parameters      map[string]*Parameter
	requestBodies   map[string]*RequestBody
	links           map[string]*Link
	pathItems       map[string]*PathItem
	examples        map[string]*Example
	securitySchemes map[string]*SecurityScheme
	callbacks       map[string]*Callback
}

func (l *loader) reset() {
	l.schemas = map[string]*Schema{}
	l.headers = map[string]*Header{}
	l.responses = map[string]*Response{}
	l.parameters = map[string]*Parameter{}
	l.requestBodies = map[string]*RequestBody{}
	l.links = map[string]*Link{}
	l.pathItems = map[string]*PathItem{}
	l.examples = map[string]*Example{}
	l.securitySchemes = map[string]*SecurityScheme{}
	l.callbacks = map[string]*Callback{}
}

// newLoader returns an empty Loader
func newLoader() *loader {
	return &loader{}
}

// LoadFromFile reads an OpenAPI specification from a file and parses it into a structured format
func LoadFromFile(location string) (*Document, error) {
	return newLoader().LoadFromFile(location)
}

// LoadFromFile reads an OpenAPI specification from a file and parses it into a structured format.
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
	case ".yaml", ".yml": // NOTE: openapi.yaml is recommended by the spec
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
	return load(io.MultiReader(buff, r)) // already includes resolving of references
}
