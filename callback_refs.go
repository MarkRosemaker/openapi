package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type CallbackRefs map[string]*CallbackRef

func (cs CallbackRefs) Validate() error {
	for name, c := range cs.ByIndex() {
		if err := validateKey(name); err != nil {
			return err
		}

		if err := c.Validate(); err != nil {
			return &ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

func (l *loader) collectCallbackRefs(cs CallbackRefs, ref ref) {
	for name, c := range cs {
		l.collectCallbackRef(c, append(ref, name))
	}
}

func (l *loader) resolveCallbackRefs(cs CallbackRefs) error {
	for name, c := range cs.ByIndex() {
		if err := l.resolveCallbackRef(c); err != nil {
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
