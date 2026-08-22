package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// Operations holds a dictionary with all the operations this application MUST implement.
//
// The key of each entry MUST be a string used to identify the operation in the document where it is defined, and its value is case-sensitive.
// Tools and libraries MAY use it to uniquely identify an operation, therefore, it is RECOMMENDED to follow common programming naming conventions.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#operationsObject
type Operations map[string]*OperationRef

// Validate validates each operation.
func (ops Operations) Validate() error {
	for name, o := range ops.ByIndex() {
		if err := o.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ops Operations) ByIndex() iter.Seq2[string, *OperationRef] {
	return ordmap.ByIndex(ops, getIndexRef[Operation, *Operation])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ops Operations) Sort() {
	ordmap.Sort(ops, setIndexRef[Operation, *Operation])
}

// Set sets a value in the map, adding it at the end of the order.
func (ops *Operations) Set(key string, o *OperationRef) {
	ordmap.Set(ops, key, o, getIndexRef[Operation, *Operation], setIndexRef[Operation, *Operation])
}

var _ json.MarshalerTo = (*Operations)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ops *Operations) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ops, enc)
}

var _ json.UnmarshalerFrom = (*Operations)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ops *Operations) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ops, dec, setIndexRef[Operation, *Operation])
}

func (l *loader) collectOperations(ops Operations, ref ref) {
	for name, o := range ops.ByIndex() {
		l.collectOperationRef(o, append(ref, name))
	}
}

func (l *loader) resolveOperations(ops Operations) error {
	for name, o := range ops.ByIndex() {
		if err := l.resolveOperationRef(o); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
