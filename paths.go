package openapi

import (
	"errors"
	"fmt"
	"iter"
	"slices"

	"github.com/MarkRosemaker/errpath"
	"github.com/MarkRosemaker/ordmap"
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
			return &errpath.ErrKey{Key: string(path), Err: err}
		}

		// if path has path parameter, check if path parameter is defined
		pp := path.Parse()
		for _, vn := range pp.VariableNames {
			hasPathParam := func(p *ParameterRef) bool {
				return p.Value.In == ParameterLocationPath && p.Value.Name == vn
			}

			// check if defined on the path item
			if slices.ContainsFunc(pathItem.Parameters, hasPathParam) {
				continue // automatically defined for all operations
			}

			// check if defined for all operations
			for method, op := range pathItem.Operations {
				if !slices.ContainsFunc(op.Parameters, hasPathParam) {
					return &errpath.ErrKey{
						Key: string(path),
						Err: &errpath.ErrField{
							Field: method,
							Err: &errpath.ErrField{
								Field: "parameters",
								Err:   fmt.Errorf("{%s} not defined", vn),
							},
						},
					}
				}
			}
		}

		if err := pathItem.Validate(); err != nil {
			return &errpath.ErrKey{Key: string(path), Err: err}
		}

		// check if operation ID is unique
		for method, op := range pathItem.Operations {
			if op.OperationID == "" {
				continue
			}

			errNotUnique := &errpath.ErrKey{
				Key: string(path),
				Err: &errpath.ErrField{
					Field: method,
					Err: &errpath.ErrField{
						Field: "operationId",
						Err: &errpath.ErrInvalid[string]{
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

// ByIndex returns a sequence of key-value pairs ordered by index.
func (ps Paths) ByIndex() iter.Seq2[Path, *PathItem] {
	return ordmap.ByIndex(ps, getIndexPathItem)
}

// Sort sorts the map by key and sets the indices accordingly.
func (ps Paths) Sort() {
	ordmap.Sort(ps, setIndexPathItem)
}

// Set sets a value in the map, adding it at the end of the order.
func (ps *Paths) Set(path Path, pathItem *PathItem) {
	ordmap.Set(ps, path, pathItem, getIndexPathItem, setIndexPathItem)
}

// MarshalJSONV2 marshals the key-value pairs in order.
func (ps *Paths) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	return ordmap.MarshalJSONV2(ps, enc, opts)
}

// UnmarshalJSONV2 unmarshals the key-value pairs in order and sets the indices.
func (ps *Paths) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	return ordmap.UnmarshalJSONV2(ps, dec, opts, setIndexPathItem)
}

func (l *loader) resolvePaths(ps Paths) error {
	for path, pathItem := range ps.ByIndex() {
		if err := l.resolvePathItem(pathItem); err != nil {
			return &errpath.ErrKey{Key: string(path), Err: err}
		}
	}

	return nil
}
