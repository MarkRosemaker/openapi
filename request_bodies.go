package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type RequestBodies map[string]*RequestBodyRef

func (rs RequestBodies) Validate() error {
	for k, r := range rs.ByIndex() {
		if err := validateKey(k); err != nil {
			return err
		}

		if err := r.Validate(); err != nil {
			return &ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

func (l *loader) collectRequestBodies(rs RequestBodies, ref ref) {
	for k, r := range rs.ByIndex() {
		l.collectRequestBodyRef(r, append(ref, k))
	}
}

func (l *loader) resolveRequestBodies(rs RequestBodies) error {
	for k, r := range rs.ByIndex() {
		if err := l.resolveRequestBodyRef(r); err != nil {
			return &ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (rs RequestBodies) ByIndex() iter.Seq2[string, *RequestBodyRef] {
	return _json.OrderedMapByIndex(rs, getIndexRef[RequestBody, *RequestBody])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (rs *RequestBodies) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(rs, dec, opts, setIndexRef[RequestBody, *RequestBody])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (rs *RequestBodies) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(rs, enc, opts)
}
