package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Webhooks describes requests initiated other than by an API call, for example by an out of band registration.
// The key name is a unique string to refer to each webhook, while the (optionally referenced) Path Item Object describes a request that may be initiated by the API provider and the expected responses.
type Webhooks map[string]*PathItemRef

// Validate checks the Webhooks for correctness.
func (ws Webhooks) Validate() error {
	for name, w := range ws.ByIndex() {
		if err := w.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ws Webhooks) ByIndex() iter.Seq2[string, *PathItemRef] {
	return ordmap.ByIndex(ws, getIndexRef[PathItem, *PathItem])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ws Webhooks) Sort() {
	ordmap.Sort(ws, setIndexRef[PathItem, *PathItem])
}

// Set sets a value in the map, adding it at the end of the order.
func (ws *Webhooks) Set(key string, v *PathItemRef) {
	ordmap.Set(ws, key, v, getIndexRef[PathItem, *PathItem], setIndexRef[PathItem, *PathItem])
}

var _ json.MarshalerTo = (*Webhooks)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ws *Webhooks) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ws, enc)
}

var _ json.UnmarshalerFrom = (*Webhooks)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ws *Webhooks) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ws, dec, setIndexRef[PathItem, *PathItem])
}

func (l *loader) resolveWebhooks(ws Webhooks) error {
	for name, w := range ws.ByIndex() {
		if err := l.resolvePathItemRef(w); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
