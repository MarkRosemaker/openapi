package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Headers map[string]*HeaderRef

func (hs Headers) Validate() error {
	for k, h := range hs.ByIndex() {
		if err := validateKey(k); err != nil {
			return err
		}

		if err := h.Validate(); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (hs Headers) ByIndex() iter.Seq2[string, *HeaderRef] {
	return ordmap.ByIndex(hs, getIndexRef[Header, *Header])
}

// Sort sorts the map by key and sets the indices accordingly.
func (hs Headers) Sort() {
	ordmap.Sort(hs, setIndexRef[Header, *Header])
}

// Set sets a value in the map, adding it at the end of the order.
func (hs *Headers) Set(key string, h *HeaderRef) {
	ordmap.Set(hs, key, h, getIndexRef[Header, *Header], setIndexRef[Header, *Header])
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (hs *Headers) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(hs, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (hs *Headers) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(hs, dec, opts, setIndexRef[Header, *Header])
}

func (l *loader) collectHeaders(hs Headers, ref ref) {
	for k, h := range hs.ByIndex() {
		l.collectHeaderRef(h, append(ref, k))
	}
}

func (l *loader) resolveHeaders(hs Headers) error {
	for k, h := range hs.ByIndex() {
		if err := l.resolveHeaderRef(h); err != nil {
			return &errpath.ErrKey{Key: k, Err: err}
		}
	}

	return nil
}
