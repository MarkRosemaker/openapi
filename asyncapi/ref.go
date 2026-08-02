package asyncapi

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"

	"github.com/MarkRosemaker/errpath"
)

// Reference is "a simple object to allow referencing other components in the specification,
// internally and externally."
//
// "The Reference Object is defined by [JSON Reference] and follows the same structure,
// behavior and rules. A JSON Reference SHALL only be used to refer to a schema that is
// formatted in either JSON or YAML. In the case of a YAML-formatted Schema, the JSON Reference
// SHALL be applied to the JSON representation of that schema."
//
// "For this specification, reference resolution is done as defined by the JSON Reference
// specification and not by the JSON Schema specification."
//
// "This object cannot be extended with additional properties and any properties added SHALL be
// ignored." Additional properties are therefore dropped when a document is read, they are not
// written back.
// ([Specification])
//
// [JSON Reference]: https://tools.ietf.org/html/draft-pbryan-zyp-json-ref-03
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#referenceObject
type Reference struct {
	// REQUIRED. The reference string.
	Identifier string `json:"$ref" yaml:"$ref"`
}

// Validate checks the reference for correctness.
func (r *Reference) Validate() error {
	if r.Identifier == "" {
		return &errpath.ErrField{Field: "$ref", Err: &errpath.ErrRequired{}}
	}

	return nil
}

type referencable[T any] interface {
	Validate() error
	*T
}

// refOrValue is a reference to a component or the component itself.
type refOrValue[T any, O referencable[T]] struct {
	// The referenced object.
	Value O `json:",inline" yaml:",inline"`
	// The reference.
	Ref *Reference `json:",inline" yaml:",inline"`

	// an index to the original location of this object
	idx int
}

// Validate checks the reference or the value it points to.
func (r *refOrValue[T, O]) Validate() error {
	if r.Ref != nil {
		if r.Value == nil {
			return fmt.Errorf("%s (%T) was not resolved", r.Ref.Identifier, r.Value)
		}

		return r.Ref.Validate()
	}

	return r.Value.Validate()
}

// isRef reports whether the object was given as a reference.
func (r *refOrValue[_, _]) isRef() bool { return r.Ref != nil }

var _ json.UnmarshalerFrom = (*refOrValue[Tag, *Tag])(nil)

func (r *refOrValue[T, O]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	// we don't know if this is a reference or not, so we read the value first
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	// try to unmarshal as a reference
	// NOTE: A reference object cannot be extended with additional properties
	// and any properties added SHALL be ignored, so we don't reject them here.
	ref := &Reference{}
	if err := json.UnmarshalDecode(
		jsontext.NewDecoder(bytes.NewBuffer(val), dec.Options()), ref,
		json.RejectUnknownMembers(false),
	); err == nil && ref.Identifier != "" {
		// we successfully unmarshalled as a reference
		r.Ref = ref // set the reference
		return nil
	}

	// it is not a reference, unmarshal as object
	var v O
	if err := json.UnmarshalDecode(
		jsontext.NewDecoder(bytes.NewBuffer(val), dec.Options()), &v,
	); err != nil {
		var t T
		return fmt.Errorf("value of %T: %w", t, err)
	}

	r.Value = v // set the value

	return nil
}

var _ json.MarshalerTo = (*refOrValue[Tag, *Tag])(nil)

func (r *refOrValue[_, _]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if r.Ref == nil {
		return json.MarshalEncode(enc, r.Value)
	}

	return json.MarshalEncode(enc, r.Ref)
}
