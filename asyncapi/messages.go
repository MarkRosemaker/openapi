package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Messages describes a map of messages included in a channel.
//
// The key of each entry represents the message identifier. It is case-sensitive.
// Tools and libraries MAY use it to uniquely identify a message, therefore, it is RECOMMENDED to follow common programming naming conventions.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#messagesObject
type Messages map[string]*MessageRef

// Validate validates each message.
func (ms Messages) Validate() error {
	for name, m := range ms.ByIndex() {
		if err := m.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ms Messages) ByIndex() iter.Seq2[string, *MessageRef] {
	return ordmap.ByIndex(ms, getIndexRef[Message, *Message])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ms Messages) Sort() {
	ordmap.Sort(ms, setIndexRef[Message, *Message])
}

// Set sets a value in the map, adding it at the end of the order.
func (ms *Messages) Set(key string, m *MessageRef) {
	ordmap.Set(ms, key, m, getIndexRef[Message, *Message], setIndexRef[Message, *Message])
}

var _ json.MarshalerTo = (*Messages)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ms *Messages) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ms, enc)
}

var _ json.UnmarshalerFrom = (*Messages)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ms *Messages) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ms, dec, setIndexRef[Message, *Message])
}

func (l *loader) collectMessages(ms Messages, ref ref) {
	for name, m := range ms.ByIndex() {
		l.collectMessageRef(m, append(ref, name))
	}
}

func (l *loader) resolveMessages(ms Messages) error {
	for name, m := range ms.ByIndex() {
		if err := l.resolveMessageRef(m); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
