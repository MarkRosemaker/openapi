package openapi

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// Reference is a simple object to allow referencing other components in the OpenAPI document, internally and externally.
//
// The `$ref` string value contains a URI RFC3986, which identifies the location of the value being referenced.
//
// See the rules for resolving Relative References.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#reference-object
type Reference struct {
	// REQUIRED. The reference identifier. This MUST be in the form of a URI.
	Identifier string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	// A short summary which by default SHOULD override that of the referenced component. If the referenced object-type does not allow a `summary` field, then this field has no effect.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A description which by default SHOULD override that of the referenced component. CommonMark syntax MAY be used for rich text representation. If the referenced object-type does not allow a `description` field, then this field has no effect.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

func (r *Reference) Validate() error {
	if r.Identifier == "" {
		return &ErrField{Field: "$ref", Err: &ErrRequired{}}
	}

	r.Description = strings.TrimSpace(r.Description)

	return nil
}

type referencable[T any] interface {
	Validate() error
	*T
}

// refOrValue is a reference to a component.
type refOrValue[T any, O referencable[T]] struct {
	// The referenced object.
	Value O `json:",inline" yaml:",inline"`
	// The reference.
	Ref *Reference `json:",inline" yaml:",inline"`

	// an index to the original location of this object
	idx int
}

func (r *refOrValue[T, O]) Validate() error {
	if r.Ref != nil {
		return r.Ref.Validate()
	}

	return r.Value.Validate()
}

func (r *refOrValue[T, O]) UnmarshalJSONV2(dec *jsontext.Decoder, opts json.Options) error {
	// we don't know if this is a reference or not, so we read the value first
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	// try to unmarshal as a reference
	ref := &Reference{}
	if err := json.UnmarshalDecode(
		jsontext.NewDecoder(bytes.NewBuffer(val), opts), ref, opts,
	); err == nil && ref.Identifier != "" {
		// we successfully unmarshalled as a reference
		r.Ref = ref // set the reference
		return nil
	}

	// it is not a reference, unmarshal as object
	var v O
	if err := json.UnmarshalDecode(
		jsontext.NewDecoder(bytes.NewBuffer(val), opts), &v, opts,
	); err != nil {
		var t T
		return fmt.Errorf("value of %T: %w", t, err)
	}

	r.Value = v // set the value

	return nil
}

func (r *refOrValue[_, _]) MarshalJSONV2(enc *jsontext.Encoder, opts json.Options) error {
	if r.Ref == nil {
		return json.MarshalEncode(enc, r.Value, opts)
	}

	return json.MarshalEncode(enc, r.Ref, opts)
}
