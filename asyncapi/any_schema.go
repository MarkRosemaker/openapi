package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/MarkRosemaker/errpath"
)

// AnySchema is a schema definition, wherever the specification allows
// a [Schema] object as well as a [MultiFormatSchema] object.
//
// A Schema Object is equivalent to a Multi Format Schema Object with the default schema format,
// which is why both are represented by this one type: if no schema format is given, the schema
// itself is the AsyncAPI schema, otherwise the schema is wrapped in a Multi Format Schema Object.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#multiFormatSchemaObject
type AnySchema struct {
	// The name of the schema format that is used to define the information.
	// If it is missing, it defaults to an AsyncAPI schema format, i.e. the schema is a [Schema] object.
	SchemaFormat SchemaFormat `json:"schemaFormat,omitempty" yaml:"schemaFormat,omitempty"`
	// The schema as an AsyncAPI Schema Object.
	// It is set if the schema format is missing or denotes an AsyncAPI schema format.
	Schema *Schema `json:"-" yaml:"-"`
	// The schema in a format other than the AsyncAPI Schema Object, e.g. Avro or Protobuf.
	// Non-JSON-based schemas (e.g. Protobuf or XSD) are inlined as a string.
	Raw jsontext.Value `json:"-" yaml:"-"`
	// A Multi Format Schema Object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// MultiFormatSchema represents a schema definition in a specific schema format.
// It is the wire format of an [AnySchema] that has a schema format.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#multiFormatSchemaObject
type MultiFormatSchema struct {
	// REQUIRED. A string containing the name of the schema format that is used to define the information.
	SchemaFormat SchemaFormat `json:"schemaFormat" yaml:"schemaFormat"`
	// REQUIRED. Definition of the message payload.
	// It can be of any type but defaults to a Schema Object.
	Schema jsontext.Value `json:"schema" yaml:"schema"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the schema for correctness.
func (s *AnySchema) Validate() error {
	if s.SchemaFormat == "" {
		if s.Schema == nil {
			return &errpath.ErrRequired{}
		}

		return s.Schema.Validate()
	}

	if s.Schema == nil && len(s.Raw) == 0 {
		return &errpath.ErrField{Field: "schema", Err: &errpath.ErrRequired{}}
	}

	if s.Schema != nil {
		if err := s.Schema.Validate(); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	}

	return validateExtensions(s.Extensions)
}

// SortMaps sorts the maps of the underlying AsyncAPI schema by key.
func (s *AnySchema) SortMaps() {
	if s == nil {
		return
	}

	s.Schema.SortMaps()
}

var _ json.UnmarshalerFrom = (*AnySchema)(nil)

// UnmarshalJSONFrom unmarshals either a Schema Object or a Multi Format Schema Object.
func (s *AnySchema) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	// a schema may be a boolean, in which case it is an AsyncAPI schema
	switch val.Kind() {
	case 't', 'f':
		return s.unmarshalSchema(val, dec.Options())
	case '{': // check below whether a schema format is given
	default:
		return fmt.Errorf("schema must be an object or a boolean, got %s", val.Kind())
	}

	// we don't know whether this is a multi format schema, so we peek at the members
	members := map[string]jsontext.Value{}
	if err := json.Unmarshal(val, &members); err != nil {
		return err
	}

	if _, ok := members["schemaFormat"]; !ok {
		// without a schema format, the object is an AsyncAPI schema
		return s.unmarshalSchema(val, dec.Options())
	}

	mf := &MultiFormatSchema{}
	if err := json.Unmarshal(val, mf, dec.Options()); err != nil {
		return err
	}

	s.SchemaFormat = mf.SchemaFormat
	s.Extensions = mf.Extensions

	// only an AsyncAPI schema can be parsed into a Schema object
	if !mf.SchemaFormat.IsAsyncAPI() {
		s.Raw = mf.Schema
		return nil
	}

	return s.unmarshalSchema(mf.Schema, dec.Options())
}

// unmarshalSchema unmarshals the value as an AsyncAPI schema.
func (s *AnySchema) unmarshalSchema(val jsontext.Value, opts json.Options) error {
	if len(val) == 0 {
		return &errpath.ErrField{Field: "schema", Err: &errpath.ErrRequired{}}
	}

	schema := &Schema{}
	if err := json.Unmarshal(val, schema, opts); err != nil {
		return err
	}

	s.Schema = schema

	return nil
}

var _ json.MarshalerTo = (*AnySchema)(nil)

// MarshalJSONTo marshals either a Schema Object or a Multi Format Schema Object.
func (s *AnySchema) MarshalJSONTo(enc *jsontext.Encoder) error {
	// the schema of a multi format schema object, either the AsyncAPI schema or the raw schema
	schema := s.Raw

	if s.Schema != nil {
		var err error
		if schema, err = json.Marshal(s.Schema, enc.Options()); err != nil {
			return err
		}

		// without a schema format, the schema itself is the object
		if s.SchemaFormat == "" {
			return enc.WriteValue(schema)
		}
	}

	return json.MarshalEncode(enc, &MultiFormatSchema{
		SchemaFormat: s.SchemaFormat,
		Schema:       schema,
		Extensions:   s.Extensions,
	})
}

func (l *loader) collectAnySchemaRefList(ss AnySchemaRefList, ref ref) {
	for i, s := range ss {
		l.collectAnySchemaRef(s, append(ref, itoa(i)))
	}
}

func (l *loader) collectAnySchemaRef(s *AnySchemaRef, ref ref) {
	if !collectRef(l, s, l.schemas, ref) {
		return
	}

	if s.Value.Schema != nil {
		l.collectSchema(s.Value.Schema, ref)
	}
}

func (l *loader) resolveAnySchemaRef(s *AnySchemaRef) error {
	return resolveRef(s, l.schemas, func(s *AnySchema) error {
		if s.Schema == nil {
			return nil
		}

		return l.resolveSchema(s.Schema)
	})
}
