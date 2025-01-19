package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// The content of a request body. The key is a media type or media type range, see [RFC7231 Appendix D], and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
// [Specification]
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#fixed-fields-10
// [RFC7231 Appendix D]: https://datatracker.ietf.org/doc/html/rfc7231#appendix-D
type Content map[MediaRange]*MediaType

// Validate validates the request body content.
func (c Content) Validate() error {
	for mr, mt := range c.ByIndex() {
		if err := mr.Validate(); err != nil {
			return &errpath.ErrKey{Key: string(mr), Err: err}
		}

		if err := mt.Validate(); err != nil {
			return &errpath.ErrKey{Key: string(mr), Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (c Content) ByIndex() iter.Seq2[MediaRange, *MediaType] {
	return ordmap.ByIndex(c, getIndexMediaType)
}

// Sort sorts the map by key and sets the indices accordingly.
func (c Content) Sort() {
	ordmap.Sort(c, setIndexMediaType)
}

// Set sets a value in the map, adding it at the end of the order.
func (c *Content) Set(mr MediaRange, mt *MediaType) {
	ordmap.Set(c, mr, mt, getIndexMediaType, setIndexMediaType)
}

// MarshalJSONTo marshals the key-value pairs in order.
func (c *Content) MarshalJSONTo(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONTo(c, enc, opts)
}

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (c *Content) UnmarshalJSONFrom(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONFrom(c, dec, opts, setIndexMediaType)
}

func (l *loader) resolveContent(c Content) error {
	for mr, mt := range c.ByIndex() {
		if err := l.resolveMediaType(mt); err != nil {
			return &errpath.ErrKey{Key: string(mr), Err: err}
		}
	}

	return nil
}
