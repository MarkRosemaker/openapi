package asyncapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"iter"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
)

// ExternalDocsByName is a map of External Documentation Objects.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#componentsExternalDocs
type ExternalDocsByName map[string]*ExternalDocsRef

// Validate validates each external documentation object.
func (ds ExternalDocsByName) Validate() error {
	for name, d := range ds.ByIndex() {
		if err := d.Validate(); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ds ExternalDocsByName) ByIndex() iter.Seq2[string, *ExternalDocsRef] {
	return ordmap.ByIndex(ds, getIndexRef[ExternalDocs, *ExternalDocs])
}

// Sort sorts the map by key and sets the indices accordingly.
func (ds ExternalDocsByName) Sort() {
	ordmap.Sort(ds, setIndexRef[ExternalDocs, *ExternalDocs])
}

// Set sets a value in the map, adding it at the end of the order.
func (ds *ExternalDocsByName) Set(key string, d *ExternalDocsRef) {
	ordmap.Set(ds, key, d, getIndexRef[ExternalDocs, *ExternalDocs], setIndexRef[ExternalDocs, *ExternalDocs])
}

var _ json.MarshalerTo = (*ExternalDocsByName)(nil)

// MarshalJSONTo marshals the key-value pairs in order.
func (ds *ExternalDocsByName) MarshalJSONTo(enc *jsontext.Encoder) error {
	return ordmap.MarshalJSONTo(ds, enc)
}

var _ json.UnmarshalerFrom = (*ExternalDocsByName)(nil)

// UnmarshalJSONFrom unmarshals the key-value pairs in order and sets the indices.
func (ds *ExternalDocsByName) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	return ordmap.UnmarshalJSONFrom(ds, dec, setIndexRef[ExternalDocs, *ExternalDocs])
}

func (l *loader) collectExternalDocsByName(ds ExternalDocsByName, ref ref) {
	for name, d := range ds.ByIndex() {
		l.collectExternalDocsRef(d, append(ref, name))
	}
}

func (l *loader) resolveExternalDocsByName(ds ExternalDocsByName) error {
	for name, d := range ds.ByIndex() {
		if err := l.resolveExternalDocsRef(d); err != nil {
			return &errpath.ErrKey{Key: name, Err: err}
		}
	}

	return nil
}
