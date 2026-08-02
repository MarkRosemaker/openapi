package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Replies is a map of Operation Reply Objects.
type Replies map[string]*OperationReplyRef

// Validate validates each reply.
func (rs Replies) Validate() error {
	for name, r := range rs.ByIndex() {
		if err := r.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (rs Replies) ByIndex() iter.Seq2[string, *OperationReplyRef] {
	return ordmap.ByIndex(rs, getIndexRef[OperationReply, *OperationReply])
}

// Sort sorts the map by key and sets the indices accordingly.
func (rs Replies) Sort() {
	ordmap.Sort(rs, setIndexRef[OperationReply, *OperationReply])
}

// Set sets a value in the map, adding it at the end of the order.
func (rs *Replies) Set(key string, r *OperationReplyRef) {
	ordmap.Set(rs, key, r, getIndexRef[OperationReply, *OperationReply], setIndexRef[OperationReply, *OperationReply])
}

var _ json.MarshalerTo = (*Replies)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (rs *Replies) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(rs, enc)
}

var _ json.UnmarshalerFrom = (*Replies)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (rs *Replies) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(rs, dec, setIndexRef[OperationReply, *OperationReply])
}

func (l *loader) collectReplies(rs Replies, ref ref) {
	for name, r := range rs.ByIndex() {
		l.collectOperationReplyRef(r, append(ref, name))
	}
}

func (l *loader) resolveReplies(rs Replies) error {
	for name, r := range rs.ByIndex() {
		if err := l.resolveOperationReplyRef(r); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
