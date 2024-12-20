package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

type Links map[string]*LinkRef

func (ls Links) Validate() error {
	for expr, l := range ls.ByIndex() {
		if err := validateKey(expr); err != nil {
			return err
		}

		if err := l.Validate(); err != nil {
			return &errpath.ErrField{Field: expr, Err: err}
		}
	}

	return nil
}

func (l *loader) collectLinks(ls Links, ref ref) {
	for expr, lr := range ls.ByIndex() {
		l.collectLinkRef(lr, append(ref, expr))
	}
}

func (l *loader) resolveLinks(ls Links) error {
	for expr, lr := range ls.ByIndex() {
		if err := l.resolveLinkRef(lr); err != nil {
			return &errpath.ErrField{Field: expr, Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (l Links) ByIndex() iter.Seq2[string, *LinkRef] {
	return _json.OrderedMapByIndex(l, getIndexRef[Link, *Link])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (l *Links) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(l, dec, opts, setIndexRef[Link, *Link])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (l *Links) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(l, enc, opts)
}
