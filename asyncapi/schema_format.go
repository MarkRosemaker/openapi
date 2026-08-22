package asyncapi

import (
	"slices"
	"strings"
)

// SchemaFormat is the name of the schema format that is used to define the information
// of a [MultiFormatSchema].
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#multiFormatSchemaObject
type SchemaFormat string

const (
	// SchemaFormatAsyncAPI is the AsyncAPI 3.1.0 Schema Object format.
	// It is the default when a schema format is not provided.
	SchemaFormatAsyncAPI SchemaFormat = "application/vnd.aai.asyncapi;version=3.1.0"
	// SchemaFormatAsyncAPIJSON is the AsyncAPI 3.1.0 Schema Object format, given as JSON.
	SchemaFormatAsyncAPIJSON SchemaFormat = "application/vnd.aai.asyncapi+json;version=3.1.0"
	// SchemaFormatAsyncAPIYAML is the AsyncAPI 3.1.0 Schema Object format, given as YAML.
	SchemaFormatAsyncAPIYAML SchemaFormat = "application/vnd.aai.asyncapi+yaml;version=3.1.0"
	// SchemaFormatJSONSchema is the JSON Schema Draft 07 format.
	SchemaFormatJSONSchema SchemaFormat = "application/schema+json;version=draft-07"
	// SchemaFormatJSONSchemaYAML is the JSON Schema Draft 07 format, given as YAML.
	SchemaFormatJSONSchemaYAML SchemaFormat = "application/schema+yaml;version=draft-07"
	// SchemaFormatAvro is the Avro 1.9.0 schema format.
	SchemaFormatAvro SchemaFormat = "application/vnd.apache.avro;version=1.9.0"
	// SchemaFormatAvroJSON is the Avro 1.9.0 schema format, given as JSON.
	SchemaFormatAvroJSON SchemaFormat = "application/vnd.apache.avro+json;version=1.9.0"
	// SchemaFormatAvroYAML is the Avro 1.9.0 schema format, given as YAML.
	SchemaFormatAvroYAML SchemaFormat = "application/vnd.apache.avro+yaml;version=1.9.0"
	// SchemaFormatOpenAPI is the OpenAPI 3.0.0 Schema Object format.
	SchemaFormatOpenAPI SchemaFormat = "application/vnd.oai.openapi;version=3.0.0"
	// SchemaFormatOpenAPIJSON is the OpenAPI 3.0.0 Schema Object format, given as JSON.
	SchemaFormatOpenAPIJSON SchemaFormat = "application/vnd.oai.openapi+json;version=3.0.0"
	// SchemaFormatOpenAPIYAML is the OpenAPI 3.0.0 Schema Object format, given as YAML.
	SchemaFormatOpenAPIYAML SchemaFormat = "application/vnd.oai.openapi+yaml;version=3.0.0"
	// SchemaFormatRAML is the RAML 1.0 data type format.
	SchemaFormatRAML SchemaFormat = "application/raml+yaml;version=1.0"
	// SchemaFormatProtobuf2 is the Protocol Buffers version 2 format.
	SchemaFormatProtobuf2 SchemaFormat = "application/vnd.google.protobuf;version=2"
	// SchemaFormatProtobuf3 is the Protocol Buffers version 3 format.
	SchemaFormatProtobuf3 SchemaFormat = "application/vnd.google.protobuf;version=3"
)

// allSchemaFormats are the schema formats listed in the specification.
// Custom values are allowed but their implementation is optional.
var allSchemaFormats = []SchemaFormat{
	SchemaFormatAsyncAPI, SchemaFormatAsyncAPIJSON, SchemaFormatAsyncAPIYAML,
	SchemaFormatJSONSchema, SchemaFormatJSONSchemaYAML,
	SchemaFormatAvro, SchemaFormatAvroJSON, SchemaFormatAvroYAML,
	SchemaFormatOpenAPI, SchemaFormatOpenAPIJSON, SchemaFormatOpenAPIYAML,
	SchemaFormatRAML,
	SchemaFormatProtobuf2, SchemaFormatProtobuf3,
}

// IsKnown reports whether the schema format is one of the formats listed in the specification.
//
// Custom values are allowed, so a format that is not known is not necessarily invalid.
func (f SchemaFormat) IsKnown() bool { return slices.Contains(allSchemaFormats, f) }

// IsAsyncAPI reports whether the schema format denotes an AsyncAPI Schema Object,
// i.e. whether the schema can be parsed into a [Schema].
func (f SchemaFormat) IsAsyncAPI() bool {
	return strings.HasPrefix(string(f), "application/vnd.aai.asyncapi")
}
