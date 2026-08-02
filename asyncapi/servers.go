package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Servers is a map of Server Objects.
type Servers map[string]*ServerRef

// Validate validates each server.
func (ss Servers) Validate() error {
	for name, s := range ss.ByIndex() {
		if err := s.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ss Servers) ByIndex() iter.Seq2[string, *ServerRef] {
	return ordmap.ByIndex(ss, getIndexRef[Server, *Server])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ss Servers) Sort() {
	ordmap.Sort(ss, setIndexRef[Server, *Server])
}

// Set sets a value in the map, adding it at the end of the order.
func (ss *Servers) Set(key string, s *ServerRef) {
	ordmap.Set(ss, key, s, getIndexRef[Server, *Server], setIndexRef[Server, *Server])
}

var _ json.MarshalerTo = (*Servers)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ss *Servers) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ss, enc)
}

var _ json.UnmarshalerFrom = (*Servers)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ss *Servers) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ss, dec, setIndexRef[Server, *Server])
}

func (l *loader) collectServers(ss Servers, ref ref) {
	for name, s := range ss.ByIndex() {
		l.collectServerRef(s, append(ref, name))
	}
}

func (l *loader) resolveServers(ss Servers) error {
	for name, s := range ss.ByIndex() {
		if err := l.resolveServerRef(s); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
