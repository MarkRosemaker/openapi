package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

type Links map[string]*LinkRef

func (ls Links) Validate() error {
	for expr, l := range ls.ByIndex() {
		if err := validateKey(expr); err != nil {
			return err
		}

		if err := l.Validate(); err != nil {
			return &errpath.ErrField{Field: expr, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ls Links) ByIndex() iter.Seq2[string, *LinkRef] {
	return ordmap.ByIndex(ls, getIndexRef[Link, *Link])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ls Links) Sort() {
	ordmap.Sort(ls, setIndexRef[Link, *Link])
}

// Set sets a value in the map, adding it at the end of the order.
func (ls *Links) Set(key string, l *LinkRef) {
	ordmap.Set(ls, key, l, getIndexRef[Link, *Link], setIndexRef[Link, *Link])
}

var _ json.MarshalerTo = (*Links)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ls *Links) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ls, enc)
}

var _ json.UnmarshalerFrom = (*Links)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ls *Links) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ls, dec, setIndexRef[Link, *Link])
}

func (l *loader) collectLinks(ls Links, ref ref) {
	for expr, lr := range ls.ByIndex() {
		l.collectLinkRef(lr, append(ref, expr))
	}
}

func (l *loader) resolveLinks(ls Links) error {
	for expr, lr := range ls.ByIndex() {
		if err := l.resolveLinkRef(lr); err != nil {
			return &errpath.ErrField{Field: expr, Err: err}
		}
	}

	return nil
}
