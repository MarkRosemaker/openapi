package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

type PathItems map[string]*PathItemRef

// Validate checks that all keys and values are valid.
func (ps PathItems) Validate() error {
	for name, p := range ps.ByIndex() {
		if err := validateKey(name); err != nil {
			return err
		}

		if err := p.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ps PathItems) ByIndex() iter.Seq2[string, *PathItemRef] {
	return ordmap.ByIndex(ps, getIndexRef[PathItem, *PathItem])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ps PathItems) Sort() {
	ordmap.Sort(ps, setIndexRef[PathItem, *PathItem])
}

// Set sets a value in the map, adding it at the end of the order.
func (ps *PathItems) Set(key string, v *PathItemRef) {
	ordmap.Set(ps, key, v, getIndexRef[PathItem, *PathItem], setIndexRef[PathItem, *PathItem])
}

var _ json.MarshalerTo = (*PathItems)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ps *PathItems) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ps, enc)
}

var _ json.UnmarshalerFrom = (*PathItems)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ps *PathItems) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ps, dec, setIndexRef[PathItem, *PathItem])
}

func (l *loader) collectPathItems(ps PathItems, ref ref) {
	for name, p := range ps.ByIndex() {
		l.collectPathItemRef(p, append(ref, name))
	}
}

func (l *loader) resolvePathItems(ps PathItems) error {
	for name, p := range ps.ByIndex() {
		if err := l.resolvePathItemRef(p); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
