package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type CallbackRefs map[string]*CallbackRef

func (cs CallbackRefs) Validate() error {
	for name, value := range cs.ByIndex() {
		if err := value.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (cs CallbackRefs) ByIndex() iter.Seq2[string, *CallbackRef] {
	return _json.OrderedMapByIndex(cs, getIndexRef[Callback, *Callback])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (cs *CallbackRefs) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(cs, dec, opts, setIndexRef[Callback, *Callback])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (cs *CallbackRefs) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(cs, enc, opts)
}
