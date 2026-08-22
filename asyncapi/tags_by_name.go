package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// TagsByName is a map of Tag Objects.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsTags
type TagsByName map[string]*TagRef

// Validate validates each tag.
func (ts TagsByName) Validate() error {
	for name, t := range ts.ByIndex() {
		if err := t.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ts TagsByName) ByIndex() iter.Seq2[string, *TagRef] {
	return ordmap.ByIndex(ts, getIndexRef[Tag, *Tag])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ts TagsByName) Sort() {
	ordmap.Sort(ts, setIndexRef[Tag, *Tag])
}

// Set sets a value in the map, adding it at the end of the order.
func (ts *TagsByName) Set(key string, t *TagRef) {
	ordmap.Set(ts, key, t, getIndexRef[Tag, *Tag], setIndexRef[Tag, *Tag])
}

var _ json.MarshalerTo = (*TagsByName)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ts *TagsByName) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ts, enc)
}

var _ json.UnmarshalerFrom = (*TagsByName)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ts *TagsByName) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ts, dec, setIndexRef[Tag, *Tag])
}

func (l *loader) collectTagsByName(ts TagsByName, ref ref) {
	for name, t := range ts.ByIndex() {
		l.collectTagRef(t, append(ref, name))
	}
}

func (l *loader) resolveTagsByName(ts TagsByName) error {
	for name, t := range ts.ByIndex() {
		if err := l.resolveTagRef(t); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
