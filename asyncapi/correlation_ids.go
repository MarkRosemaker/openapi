package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// CorrelationIDs is a map of Correlation ID Objects.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsCorrelationIDs
type CorrelationIDs map[string]*CorrelationIDRef

// Validate validates each correlation ID.
func (cs CorrelationIDs) Validate() error {
	for name, c := range cs.ByIndex() {
		if err := c.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (cs CorrelationIDs) ByIndex() iter.Seq2[string, *CorrelationIDRef] {
	return ordmap.ByIndex(cs, getIndexRef[CorrelationID, *CorrelationID])
}

// Sort sorts the map by key and sets the indices accordingly.
func (cs CorrelationIDs) Sort() {
	ordmap.Sort(cs, setIndexRef[CorrelationID, *CorrelationID])
}

// Set sets a value in the map, adding it at the end of the order.
func (cs *CorrelationIDs) Set(key string, c *CorrelationIDRef) {
	ordmap.Set(cs, key, c, getIndexRef[CorrelationID, *CorrelationID], setIndexRef[CorrelationID, *CorrelationID])
}

var _ json.MarshalerTo = (*CorrelationIDs)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (cs *CorrelationIDs) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(cs, enc)
}

var _ json.UnmarshalerFrom = (*CorrelationIDs)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (cs *CorrelationIDs) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(cs, dec, setIndexRef[CorrelationID, *CorrelationID])
}

func (l *loader) collectCorrelationIDs(cs CorrelationIDs, ref ref) {
	for name, c := range cs.ByIndex() {
		l.collectCorrelationIDRef(c, append(ref, name))
	}
}

func (l *loader) resolveCorrelationIDs(cs CorrelationIDs) error {
	for name, c := range cs.ByIndex() {
		if err := l.resolveCorrelationIDRef(c); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
