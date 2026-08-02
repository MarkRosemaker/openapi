package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// ReplyAddresses is a map of Operation Reply Address Objects.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsReplyAddresses
type ReplyAddresses map[string]*OperationReplyAddressRef

// Validate validates each reply address.
func (as ReplyAddresses) Validate() error {
	for name, a := range as.ByIndex() {
		if err := a.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (as ReplyAddresses) ByIndex() iter.Seq2[string, *OperationReplyAddressRef] {
	return ordmap.ByIndex(as, getIndexRef[OperationReplyAddress, *OperationReplyAddress])
}

// Sort sorts the map by key and sets the indices accordingly.
func (as ReplyAddresses) Sort() {
	ordmap.Sort(as, setIndexRef[OperationReplyAddress, *OperationReplyAddress])
}

// Set sets a value in the map, adding it at the end of the order.
func (as *ReplyAddresses) Set(key string, a *OperationReplyAddressRef) {
	ordmap.Set(as, key, a, getIndexRef[OperationReplyAddress, *OperationReplyAddress], setIndexRef[OperationReplyAddress, *OperationReplyAddress])
}

var _ json.MarshalerTo = (*ReplyAddresses)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (as *ReplyAddresses) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(as, enc)
}

var _ json.UnmarshalerFrom = (*ReplyAddresses)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (as *ReplyAddresses) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(as, dec, setIndexRef[OperationReplyAddress, *OperationReplyAddress])
}

func (l *loader) collectReplyAddresses(as ReplyAddresses, ref ref) {
	for name, a := range as.ByIndex() {
		l.collectOperationReplyAddressRef(a, append(ref, name))
	}
}

func (l *loader) resolveReplyAddresses(as ReplyAddresses) error {
	for name, a := range as.ByIndex() {
		if err := l.resolveOperationReplyAddressRef(a); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
