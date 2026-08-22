package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"
	"slices"
	"strings"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Bindings is a map describing protocol-specific definitions
// for a server, a channel, an operation or a message.
//
// The keys describe the name of the protocol, the values describe the protocol-specific definitions.
// The definitions themselves are described by the [bindings specification]
// and are therefore kept as raw JSON values.
// ([Specification])
//
// [bindings specification]: https://github.com/asyncapi/bindings
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#serverBindingsObject
type Bindings map[Protocol]*Binding

// Validate checks that the protocols are known and that the definitions are valid JSON.
func (bs Bindings) Validate() error {
	for protocol, b := range bs.ByIndex() {
		// the bindings object MAY be extended with specification extensions
		if strings.HasPrefix(string(protocol), "x-") {
			continue
		}

		if !slices.Contains(allBindingProtocols, protocol) {
			return &errpath.ErrKey{
				Key: string(protocol),
				Err: &errpath.ErrInvalid[Protocol]{Value: protocol, Enum: allBindingProtocols},
			}
		}

		if err := b.Validate(); err != nil {
			return &errpath.ErrKey{Key: string(protocol), Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (bs Bindings) ByIndex() iter.Seq2[Protocol, *Binding] {
	return ordmap.ByIndex(bs, getIndexBinding)
}

// Sort sorts the map by key and sets the indices accordingly.
func (bs Bindings) Sort() {
	ordmap.Sort(bs, setIndexBinding)
}

// Set sets a value in the map, adding it at the end of the order.
func (bs *Bindings) Set(key Protocol, b *Binding) {
	ordmap.Set(bs, key, b, getIndexBinding, setIndexBinding)
}

var _ json.MarshalerTo = (*Bindings)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (bs *Bindings) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(bs, enc)
}

var _ json.UnmarshalerFrom = (*Bindings)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (bs *Bindings) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(bs, dec, setIndexBinding)
}

// Binding holds the protocol-specific definitions of a single protocol.
type Binding struct {
	// The protocol-specific definition as defined by the bindings specification.
	Value jsontext.Value

	// an index to the original location of this object
	idx int
}

// Validate checks that the binding holds a value.
func (b *Binding) Validate() error {
	if len(b.Value) == 0 {
		return &errpath.ErrRequired{}
	}

	return nil
}

var _ json.UnmarshalerFrom = (*Binding)(nil)

// UnmarshalJSONFrom unmarshals the value of the binding.
func (b *Binding) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	b.Value = slices.Clone(val)

	return nil
}

var _ json.MarshalerTo = (*Binding)(nil)

// MarshalJSONTo marshals the value of the binding.
func (b *Binding) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteValue(b.Value)
}

func getIndexBinding(b *Binding) int             { return b.idx }
func setIndexBinding(b *Binding, i int) *Binding { b.idx = i; return b }

func (l *loader) collectBindingsRef(b *BindingsRef, ref ref) {
	collectRef(l, b, l.bindings, ref)
}

func (l *loader) resolveBindingsRef(b *BindingsRef) error {
	return resolveRef(b, l.bindings, nil)
}
