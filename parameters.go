package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Parameters map[string]*ParameterRef

func (ps Parameters) Validate() error {
	for name, p := range ps.ByIndex() {
		if err := validateKey(name); err != nil {
			return err
		}

		if err := p.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ps Parameters) ByIndex() iter.Seq2[string, *ParameterRef] {
	return ordmap.ByIndex(ps, getIndexRef[Parameter, *Parameter])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ps Parameters) Sort() {
	ordmap.Sort(ps, setIndexRef[Parameter, *Parameter])
}

// Set sets a value in the map, adding it at the end of the order.
func (ps *Parameters) Set(name string, p *ParameterRef) {
	ordmap.Set(ps, name, p, getIndexRef[Parameter, *Parameter], setIndexRef[Parameter, *Parameter])
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (ps *Parameters) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(ps, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (ps *Parameters) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(ps, dec, opts, setIndexRef[Parameter, *Parameter])
}

func (l *loader) collectParameters(ps Parameters, ref ref) {
	for name, p := range ps {
		l.collectParameterRef(p, append(ref, name))
	}
}

func (l *loader) resolveParameters(ps Parameters) error {
	for name, p := range ps.ByIndex() {
		if err := l.resolveParameterRef(p); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
