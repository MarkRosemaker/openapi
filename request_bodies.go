package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
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
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (rs RequestBodies) ByIndex() iter.Seq2[string, *RequestBodyRef] {
	return ordmap.ByIndex(rs, getIndexRef[RequestBody, *RequestBody])
}

// Sort sorts the map by key and sets the indices accordingly.
func (rs RequestBodies) Sort() {
	ordmap.Sort(rs, setIndexRef[RequestBody, *RequestBody])
}

// Set sets a value in the map, adding it at the end of the order.
func (rs *RequestBodies) Set(key string, r *RequestBodyRef) {
	ordmap.Set(rs, key, r, getIndexRef[RequestBody, *RequestBody], setIndexRef[RequestBody, *RequestBody])
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (rs *RequestBodies) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(rs, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (rs *RequestBodies) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(rs, dec, opts, setIndexRef[RequestBody, *RequestBody])
}

func (l *loader) collectRequestBodies(rs RequestBodies, ref ref) {
	for k, r := range rs.ByIndex() {
		l.collectRequestBodyRef(r, append(ref, k))
	}
}

func (l *loader) resolveRequestBodies(rs RequestBodies) error {
	for k, r := range rs.ByIndex() {
		if err := l.resolveRequestBodyRef(r); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}
