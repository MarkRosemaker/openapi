package openapi

import (
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Callback is a map of possible out-of band callbacks related to the parent operation.
// Each value in the map is a Path Item Object that describes a set of requests that may be initiated by the API provider and the expected responses.
// The key value used to identify the path item object is an expression, evaluated at runtime, that identifies a URL to use for the callback operation.
//
// To describe incoming requests from the API provider independent from another API call, use the `webhooks` field.
//
// Note that according to the specification, this object MAY be extended with Specification Extensions, but we do not support that in this implementation.
// Note that we are not validating the [runtime expression] in this implementation.
//
// [runtime expression]: https://spec.openapis.org/oas/v3.1.0#key-expression
type Callback map[RuntimeExpression]*PathItemRef

func (c Callback) Validate() error {
	for expr, p := range c.ByIndex() {
		if err := p.Validate(); err != nil {
			return &ErrKey{Key: string(expr), Err: err}
		}
	}

	return nil
}

func (l *loader) collectCallbackRef(c *CallbackRef, ref ref) {
	if c.Value != nil {
		l.collectCallback(c.Value, ref)
	}
}

func (l *loader) collectCallback(c *Callback, ref ref) {
	l.callbacks[ref.String()] = c
}

func (l *loader) resolveCallbackRef(c *CallbackRef) error {
	return resolveRef(c, l.callbacks, l.resolveCallback)
}

func (l *loader) resolveCallback(c *Callback) error {
	for expr, p := range c.ByIndex() {
		if err := l.resolvePathItemRef(p); err != nil {
			return &ErrKey{Key: string(expr), Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (c Callback) ByIndex() iter.Seq2[RuntimeExpression, *PathItemRef] {
	return _json.OrderedMapByIndex(c, getIndexRef[PathItem, *PathItem])
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (c *Callback) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(c, dec, opts, setIndexRef[PathItem, *PathItem])
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (c *Callback) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(c, enc, opts)
}
