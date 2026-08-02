package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// MessageTraits is a map of Message Trait Objects.
type MessageTraits map[string]*MessageTraitRef

// Validate validates each message trait.
func (ts MessageTraits) Validate() error {
	for name, t := range ts.ByIndex() {
		if err := t.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ts MessageTraits) ByIndex() iter.Seq2[string, *MessageTraitRef] {
	return ordmap.ByIndex(ts, getIndexRef[MessageTrait, *MessageTrait])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ts MessageTraits) Sort() {
	ordmap.Sort(ts, setIndexRef[MessageTrait, *MessageTrait])
}

// Set sets a value in the map, adding it at the end of the order.
func (ts *MessageTraits) Set(key string, t *MessageTraitRef) {
	ordmap.Set(ts, key, t, getIndexRef[MessageTrait, *MessageTrait], setIndexRef[MessageTrait, *MessageTrait])
}

var _ json.MarshalerTo = (*MessageTraits)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ts *MessageTraits) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ts, enc)
}

var _ json.UnmarshalerFrom = (*MessageTraits)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ts *MessageTraits) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ts, dec, setIndexRef[MessageTrait, *MessageTrait])
}

func (l *loader) collectMessageTraits(ts MessageTraits, ref ref) {
	for name, t := range ts.ByIndex() {
		l.collectMessageTraitRef(t, append(ref, name))
	}
}

func (l *loader) resolveMessageTraits(ts MessageTraits) error {
	for name, t := range ts.ByIndex() {
		if err := l.resolveMessageTraitRef(t); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
