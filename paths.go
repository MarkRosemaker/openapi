package openapi

import (
	"errors"
	"iter"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Holds the relative paths to the individual endpoints and their operations.
// The path is appended to the URL from the Server Object in order to construct the full URL. The Paths MAY be empty, due to Access Control List (ACL) constraints.
//
// Note that according to the specification, this object MAY be extended with Specification Extensions, but we do not support that in this implementation.
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#paths-object
type Paths map[Path]*PathItem

func (ps Paths) Validate() error {
	// The id of an operation MUST be unique among all operations described in the API. The operationId value is case-sensitive.
	opIDs := map[string]error{}

	for path, pathItem := range ps.ByIndex() {
		if err := path.Validate(); err != nil {
			return &ErrKey{Key: string(path), Err: err}
		}

		if err := pathItem.Validate(); err != nil {
			return &ErrKey{Key: string(path), Err: err}
		}

		for method, op := range pathItem.Operations {
			if op.OperationID == "" {
				continue
			}

			errNotUnique := &ErrKey{
				Key: string(path),
				Err: &ErrField{
					Field: method,
					Err: &ErrField{
						Field: "operationId",
						Err: &ErrInvalid[string]{
							Value: op.OperationID, Message: "must be unique",
						},
					},
				},
			}

			prevInstance := opIDs[op.OperationID]
			if prevInstance == nil {
				opIDs[op.OperationID] = errNotUnique
				continue
			}

			// output both instances of the operation ID
			return errors.Join(prevInstance, errNotUnique)
		}
	}

	return nil
}

func (l *loader) resolvePaths(ps Paths) error {
	for path, pathItem := range ps.ByIndex() {
		if err := l.resolvePathItem(pathItem); err != nil {
			return &ErrKey{Key: string(path), Err: err}
		}
	}

	return nil
}

// ByIndex returns the keys of the map in the order of the index.
func (ps Paths) ByIndex() iter.Seq2[Path, *PathItem] {
	return _json.OrderedMapByIndex(ps, func(p *PathItem) int { return p.idx })
}

// UnmarshalJSONV2 unmarshals the map from JSON and sets the index of each variable.
func (ps *Paths) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return _json.UnmarshalOrderedMap(ps, dec, opts, func(p *PathItem, i int) { p.idx = i })
}

// MarshalJSONV2 marshals the map to JSON in the order of the index.
func (ps *Paths) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return _json.MarshalOrderedMap(ps, enc, opts)
}
