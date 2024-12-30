package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Content map[MediaRange]*MediaType

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

// MarshalJSONV2 marshals the key-value pairs in order.
func (c *Content) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(c, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (c *Content) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(c, dec, opts, setIndexMediaType)
}

func (l *loader) resolveContent(c Content) error {
	for mr, mt := range c.ByIndex() {
		if err := l.resolveMediaType(mt); err != nil {
			return &errpath.ErrKey{Key: string(mr), Err: err}
		}
	}

	return nil
}
