package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
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
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (cs CallbackRefs) ByIndex() iter.Seq2[string, *CallbackRef] {
	return ordmap.ByIndex(cs, getIndexRef[Callback, *Callback])
}

// Sort sorts the map by key and sets the indices accordingly.
func (cs CallbackRefs) Sort() {
	ordmap.Sort(cs, setIndexRef[Callback, *Callback])
}

// Set sets a value in the map, adding it at the end of the order.
func (cs *CallbackRefs) Set(key string, c *CallbackRef) {
	ordmap.Set(cs, key, c, getIndexRef[Callback, *Callback], setIndexRef[Callback, *Callback])
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (cs *CallbackRefs) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(cs, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (cs *CallbackRefs) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(cs, dec, opts, setIndexRef[Callback, *Callback])
}

func (l *loader) collectCallbackRefs(cs CallbackRefs, ref ref) {
	for name, c := range cs {
		l.collectCallbackRef(c, append(ref, name))
	}
}

func (l *loader) resolveCallbackRefs(cs CallbackRefs) error {
	for name, c := range cs.ByIndex() {
		if err := l.resolveCallbackRef(c); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
