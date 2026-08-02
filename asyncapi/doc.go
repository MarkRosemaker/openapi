// Package asyncapi parses, validates, formats and writes [AsyncAPI] specifications.
//
// It implements version 3.1.0 of the [AsyncAPI Specification], which describes an
// event-driven API: "The AsyncAPI Specification defines a set of files required to describe
// an application's API. These files can then be used to create utilities, such as
// documentation, code, integration, or testing tools."
//
// # Reading a document
//
// A document is read with [LoadFromFile], [LoadFromData] or [LoadFromReader], each of which
// determines whether the document is JSON or YAML, since "an AsyncAPI document can be JSON
// or YAML format" ([Format]). Use the JSON or YAML variants, e.g. [LoadFromDataJSON], if the
// format is already known.
//
// While reading, every reference is resolved, i.e. the [Reference] and the object it points
// to are both available. References that point to other references are followed, e.g. an
// operation that refers to a message of a channel which in turn refers to a message of the
// [Components] object.
//
// # Validating a document
//
// [Document.Validate] checks the document against the rules of the specification: required
// fields, enumerations, patterns of keys and addresses, absolute URLs, runtime expressions,
// as well as the rules that span several objects, e.g. that the messages of an operation
// "MUST contain a subset of the messages defined in the channel referenced in this
// operation" ([Operation Object]).
//
// An error tells where the problem is, e.g.
//
//	channels["userSignedup"].messages["userSignedUp"].contentType: mime: expected slash after first token
//
// Validation also normalizes a document where the specification allows it, e.g. by trimming
// the whitespace around a description or by adding the https scheme to a URL that is missing one.
//
// # Writing a document
//
// [Document.ToJSON], [Document.ToYAML], [Document.WriteJSON], [Document.WriteYAML] and
// [Document.WriteToFile] write the document back. Everything that was read is written back:
// the order of the keys of every map is preserved, specification extensions are kept where
// they were, and the definitions of bindings and of schemas in other formats (Avro, Protobuf,
// ...) are kept as they were given. [Document.SortMaps] sorts the maps of the document by key
// if a canonical order is preferred over the original one.
//
// [AsyncAPI]: https://www.asyncapi.com
// [AsyncAPI Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0
// [Format]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#format
// [Operation Object]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationObject
package asyncapi
