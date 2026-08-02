package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// BindingsByName is a map of bindings objects.
type BindingsByName map[string]*BindingsRef

// Validate validates each bindings object.
func (bs BindingsByName) Validate() error {
	for name, b := range bs.ByIndex() {
		if err := b.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (bs BindingsByName) ByIndex() iter.Seq2[string, *BindingsRef] {
	return ordmap.ByIndex(bs, getIndexRef[Bindings, *Bindings])
}

// Sort sorts the map by key and sets the indices accordingly.
func (bs BindingsByName) Sort() {
	ordmap.Sort(bs, setIndexRef[Bindings, *Bindings])
}

// Set sets a value in the map, adding it at the end of the order.
func (bs *BindingsByName) Set(key string, b *BindingsRef) {
	ordmap.Set(bs, key, b, getIndexRef[Bindings, *Bindings], setIndexRef[Bindings, *Bindings])
}

var _ json.MarshalerTo = (*BindingsByName)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (bs *BindingsByName) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(bs, enc)
}

var _ json.UnmarshalerFrom = (*BindingsByName)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (bs *BindingsByName) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(bs, dec, setIndexRef[Bindings, *Bindings])
}

func (l *loader) collectBindingsByName(bs BindingsByName, ref ref) {
	for name, b := range bs.ByIndex() {
		l.collectBindingsRef(b, append(ref, name))
	}
}

func (l *loader) resolveBindingsByName(bs BindingsByName) error {
	for name, b := range bs.ByIndex() {
		if err := l.resolveBindingsRef(b); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
