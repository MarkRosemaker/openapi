package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Content map[MediaRange]*MediaType

func (c Content) Validate() error {
	for mr, mt := range c.ByIndex() {
		if err := mr.Validate(); err != nil {
			return &ErrKey{Key: string(mr), Err: err}
		}

		if err := mt.Validate(); err != nil {
			return &ErrKey{Key: string(mr), Err: err}
		}
	}

	return nil
}

func (l *loader) resolveContent(c Content) error {
	for mr, mt := range c.ByIndex() {
		if err := l.resolveMediaType(mt); err != nil {
			return &ErrKey{Key: string(mr), Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (c Content) ByIndex() iter.Seq2[MediaRange, *MediaType] {
	return _json.OrderedMapByIndex(c, getIndexMediaType)
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (c *Content) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(c, dec, opts, setIndexMediaType)
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (c *Content) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(c, enc, opts)
}
