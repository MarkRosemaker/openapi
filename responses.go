package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// OperationsResponses is a container for the expected responses of an operation.
// The container maps a HTTP response code to the expected response.
//
// The documentation is not necessarily expected to cover all possible HTTP response codes because they may not be known in advance.
// However, documentation is expected to cover a successful operation response and any known errors.
//
// The `default` MAY be used as a default response object for all HTTP codes
// that are not covered individually by the `Responses Object`.
//
// The `Responses Object` MUST contain at least one response code, and if only one
// response code is provided it SHOULD be the response for a successful operation
// call.
//
// Note that according to the specification, this object MAY be extended with Specification Extensions, but we do not support that in this implementation.
type OperationResponses = Responses[StatusCode]

// ResponsesByName is a map of response names to response objects.
type ResponsesByName = Responses[string]

// Responses is a map of either response name or status code to a response object.
type Responses[K ~string] map[K]*ResponseRef

// Validate checks that each response is valid.
// It does not check the validity of the keys as they could be either status codes or response names.
func (rs Responses[K]) Validate() error {
	for keyOrCode, r := range rs.ByIndex() {
		if err := r.Validate(); err != nil {
			return &errpath.ErrKey{Key: string(keyOrCode), Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (rs Responses[K]) ByIndex() iter.Seq2[K, *ResponseRef] {
	return ordmap.ByIndex(rs, getIndexRef[Response, *Response])
}

// Sort sorts the map by key and sets the indices accordingly.
func (rs Responses[_]) Sort() {
	ordmap.Sort(rs, setIndexRef[Response, *Response])
}

// Set sets a value in the map, adding it at the end of the order.
func (rs *Responses[K]) Set(key K, v *ResponseRef) {
	ordmap.Set(rs, key, v, getIndexRef[Response, *Response], setIndexRef[Response, *Response])
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (rs *Responses[_]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(rs, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (rs *Responses[_]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(rs, dec, opts, setIndexRef[Response, *Response])
}

func (l *loader) collectResponses(rs ResponsesByName, ref ref) {
	for name, r := range rs.ByIndex() {
		l.collectResponseRef(r, append(ref, name))
	}
}

func (l *loader) resolveResponses(rs ResponsesByName) error {
	for name, r := range rs.ByIndex() {
		if err := l.resolveResponseRef(r); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
