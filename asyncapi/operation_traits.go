package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// OperationTraits is a map of Operation Trait Objects.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsOperationTraits
type OperationTraits map[string]*OperationTraitRef

// Validate validates each operation trait.
func (ts OperationTraits) Validate() error {
	for name, t := range ts.ByIndex() {
		if err := t.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ts OperationTraits) ByIndex() iter.Seq2[string, *OperationTraitRef] {
	return ordmap.ByIndex(ts, getIndexRef[OperationTrait, *OperationTrait])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ts OperationTraits) Sort() {
	ordmap.Sort(ts, setIndexRef[OperationTrait, *OperationTrait])
}

// Set sets a value in the map, adding it at the end of the order.
func (ts *OperationTraits) Set(key string, t *OperationTraitRef) {
	ordmap.Set(ts, key, t, getIndexRef[OperationTrait, *OperationTrait], setIndexRef[OperationTrait, *OperationTrait])
}

var _ json.MarshalerTo = (*OperationTraits)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ts *OperationTraits) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ts, enc)
}

var _ json.UnmarshalerFrom = (*OperationTraits)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ts *OperationTraits) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ts, dec, setIndexRef[OperationTrait, *OperationTrait])
}

func (l *loader) collectOperationTraits(ts OperationTraits, ref ref) {
	for name, t := range ts.ByIndex() {
		l.collectOperationTraitRef(t, append(ref, name))
	}
}

func (l *loader) resolveOperationTraits(ts OperationTraits) error {
	for name, t := range ts.ByIndex() {
		if err := l.resolveOperationTraitRef(t); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
