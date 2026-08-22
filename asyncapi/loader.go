package asyncapi

import (
	"bytes"
	"encoding/json/jsontext"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// loader helps deserialize an AsyncAPI v3 document.
//
// It remembers where every object that can be referenced is defined,
// so that the references of the document can be resolved once it was read.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#referenceObject
type loader struct {
	schemas         map[string]*AnySchema
	servers         map[string]*Server
	serverVariables map[string]*ServerVariable
	channels        map[string]*Channel
	operations      map[string]*Operation
	operationTraits map[string]*OperationTrait
	replies         map[string]*OperationReply
	replyAddresses  map[string]*OperationReplyAddress
	messages        map[string]*Message
	messageTraits   map[string]*MessageTrait
	parameters      map[string]*Parameter
	correlationIDs  map[string]*CorrelationID
	securitySchemes map[string]*SecurityScheme
	tags            map[string]*Tag
	externalDocs    map[string]*ExternalDocs
	bindings        map[string]*Bindings

	// aliases maps the path of a reference to the identifier of the object it refers to
	aliases map[string]string
}

func (l *loader) reset() {
	l.schemas = map[string]*AnySchema{}
	l.servers = map[string]*Server{}
	l.serverVariables = map[string]*ServerVariable{}
	l.channels = map[string]*Channel{}
	l.operations = map[string]*Operation{}
	l.operationTraits = map[string]*OperationTrait{}
	l.replies = map[string]*OperationReply{}
	l.replyAddresses = map[string]*OperationReplyAddress{}
	l.messages = map[string]*Message{}
	l.messageTraits = map[string]*MessageTrait{}
	l.parameters = map[string]*Parameter{}
	l.correlationIDs = map[string]*CorrelationID{}
	l.securitySchemes = map[string]*SecurityScheme{}
	l.tags = map[string]*Tag{}
	l.externalDocs = map[string]*ExternalDocs{}
	l.bindings = map[string]*Bindings{}
	l.aliases = map[string]string{}
}

// newLoader returns an empty loader.
func newLoader() *loader {
	return &loader{}
}

// LoadFromFile reads an AsyncAPI specification from a file and parses it into a structured format.
func LoadFromFile(location string) (*Document, error) {
	return newLoader().LoadFromFile(location)
}

// LoadFromFile reads an AsyncAPI specification from a file and parses it into a structured format.
func (l *loader) LoadFromFile(location string) (*Document, error) {
	f, err := os.Open(location)
	if err != nil {
		return nil, err
	}

	// determine the file type and load accordingly
	doc, err := func() (*Document, error) {
		switch ext := filepath.Ext(location); ext {
		case ".json":
			return l.LoadFromReaderJSON(f)
		case ".yaml", ".yml":
			return l.LoadFromReaderYAML(f)
		default:
			return nil, fmt.Errorf("unsupported file extension: %s", ext)
		}
	}()

	return doc, errorsJoin(err, f.Close())
}

// LoadFromData reads an AsyncAPI specification from a byte array and parses it into a structured format.
func LoadFromData(data []byte) (*Document, error) {
	return newLoader().LoadFromData(data)
}

// LoadFromData reads an AsyncAPI specification from a byte array and parses it into a structured format.
// It will try to determine the format of the data and load it accordingly.
// If you know the format of the data, use LoadFromDataJSON or LoadFromDataYAML instead.
func (l *loader) LoadFromData(data []byte) (*Document, error) {
	if jsontext.Value(data).IsValid() {
		return l.LoadFromDataJSON(data)
	}

	return l.LoadFromDataYAML(data)
}

// LoadFromReader reads an AsyncAPI specification from an io.Reader and parses it into a structured format.
// It will try to determine the format of the data and load it accordingly.
// If you know the format of the data, use LoadFromReaderJSON or LoadFromReaderYAML instead.
func LoadFromReader(r io.Reader) (*Document, error) {
	return newLoader().LoadFromReader(r)
}

// LoadFromReader reads an AsyncAPI specification from an io.Reader and parses it into a structured format.
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
