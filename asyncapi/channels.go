package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Channels is an object containing all the Channel Object definitions the application MUST use during runtime.
//
// The key of each entry is an identifier for the described channel. The channel ID is case-sensitive.
// Tools and libraries MAY use it to uniquely identify a channel, therefore, it is RECOMMENDED to follow common programming naming conventions.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#channelsObject
type Channels map[string]*ChannelRef

// Validate validates each channel.
func (cs Channels) Validate() error {
	for name, c := range cs.ByIndex() {
		if err := c.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (cs Channels) ByIndex() iter.Seq2[string, *ChannelRef] {
	return ordmap.ByIndex(cs, getIndexRef[Channel, *Channel])
}

// Sort sorts the map by key and sets the indices accordingly.
func (cs Channels) Sort() {
	ordmap.Sort(cs, setIndexRef[Channel, *Channel])
}

// Set sets a value in the map, adding it at the end of the order.
func (cs *Channels) Set(key string, c *ChannelRef) {
	ordmap.Set(cs, key, c, getIndexRef[Channel, *Channel], setIndexRef[Channel, *Channel])
}

var _ json.MarshalerTo = (*Channels)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (cs *Channels) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(cs, enc)
}

var _ json.UnmarshalerFrom = (*Channels)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (cs *Channels) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(cs, dec, setIndexRef[Channel, *Channel])
}

func (l *loader) collectChannels(cs Channels, ref ref) {
	for name, c := range cs.ByIndex() {
		l.collectChannelRef(c, append(ref, name))
	}
}

func (l *loader) resolveChannels(cs Channels) error {
	for name, c := range cs.ByIndex() {
		if err := l.resolveChannelRef(c); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
