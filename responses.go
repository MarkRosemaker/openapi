package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
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

// Responses is a map of either response name or status code to a response object.
type Responses[K ~string] map[K]*ResponseRef

// Validate checks that each response is valid.
// It does not check the validity of the keys as they could be either status codes or response names.
func (rs Responses[K]) Validate() error {
	for keyOrCode, r := range rs.ByIndex() {
		if err := r.Validate(); err != nil {
			return &ErrKey{Key: string(keyOrCode), Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (rs Responses[K]) ByIndex() iter.Seq2[K, *ResponseRef] {
	return _json.OrderedMapByIndex(rs, getIndexRef[Response, *Response])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (r *Responses[_]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(r, dec, opts, setIndexRef[Response, *Response])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (r *Responses[_]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(r, enc, opts)
}
