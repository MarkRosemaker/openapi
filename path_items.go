package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type PathItems map[string]*PathItemRef

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

// MarshalJSONV2 marshals the key-value pairs in order.
func (ps *PathItems) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(ps, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (ps *PathItems) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(ps, dec, opts, setIndexRef[PathItem, *PathItem])
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
