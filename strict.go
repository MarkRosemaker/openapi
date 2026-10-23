package openapi

import (
	"iter"

	"github.com/MarkRosemaker/errpath"
)

// ValidateStrict validates the document as Validate does, and then reports what the specifications allow but a
// document rarely means: a format the library does not know, which is often a typo, and a key of examples, headers,
// links or callbacks outside the components that would not do as a component's key.
func (d *Document) ValidateStrict() error {
	if err := d.Validate(); err != nil {
		return err
	}

	if err := strictMap(d.Paths.ByIndex(), strictPathItem); err != nil {
		return &errpath.ErrField{Field: "paths", Err: err}
	}

	if err := strictMap(d.Webhooks.ByIndex(), strictRef(strictPathItem)); err != nil {
		return &errpath.ErrField{Field: "webhooks", Err: err}
	}

	if err := strictComponents(&d.Components); err != nil {
		return &errpath.ErrField{Field: "components", Err: err}
	}

	return nil
}

// strictMap checks each value of a map in order, naming the key of the first that fails.
func strictMap[K ~string, V any](kvs iter.Seq2[K, V], check func(V) error) error {
	for k, v := range kvs {
		if err := check(v); err != nil {
			return &errpath.ErrKey{Key: string(k), Err: err}
		}
	}

	return nil
}

// strictKeys checks each key of a map against the components' regular expression, and then each value.
func strictKeys[V any](kvs iter.Seq2[string, V], check func(V) error) error {
	if err := validateKeys(kvs); err != nil {
		return err
	}

	return strictMap(kvs, check)
}

// strictRef checks what a value holds; a reference is checked where its target is, in the components.
func strictRef[T any, O referencable[T]](check func(O) error) func(*refOrValue[T, O]) error {
	return func(r *refOrValue[T, O]) error {
		if r == nil || r.Ref != nil || r.Value == nil {
			return nil
		}

		return check(r.Value)
	}
}

// noCheck is for a map whose values hold nothing ValidateStrict checks.
func noCheck[V any](V) error { return nil }

func strictComponents(c *Components) error {
	for _, check := range []struct {
		field string
		check func() error
	}{
		{"schemas", func() error { return strictMap(c.Schemas.ByIndex(), strictSchema) }},
		{"responses", func() error { return strictMap(c.Responses.ByIndex(), strictRef(strictResponse)) }},
		{"parameters", func() error { return strictMap(c.Parameters.ByIndex(), strictRef(strictParameter)) }},
		{"requestBodies", func() error { return strictMap(c.RequestBodies.ByIndex(), strictRef(strictRequestBody)) }},
		{"headers", func() error { return strictMap(c.Headers.ByIndex(), strictRef(strictHeader)) }},
		{"callbacks", func() error { return strictMap(c.Callbacks.ByIndex(), strictRef(strictCallback)) }},
		{"pathItems", func() error { return strictMap(c.PathItems.ByIndex(), strictRef(strictPathItem)) }},
	} {
		if err := check.check(); err != nil {
			return &errpath.ErrField{Field: check.field, Err: err}
		}
	}

	return nil
}

func strictPathItem(p *PathItem) error {
	if err := strictParameters(p.Parameters); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	for method, op := range p.Operations {
		if err := strictOperation(op); err != nil {
			return &errpath.ErrField{Field: method, Err: err}
		}
	}

	return nil
}

func strictOperation(o *Operation) error {
	if err := strictParameters(o.Parameters); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	if err := strictRef(strictRequestBody)(o.RequestBody); err != nil {
		return &errpath.ErrField{Field: "requestBody", Err: err}
	}

	if err := strictMap(o.Responses.ByIndex(), strictRef(strictResponse)); err != nil {
		return &errpath.ErrField{Field: "responses", Err: err}
	}

	if err := strictKeys(o.Callbacks.ByIndex(), strictRef(strictCallback)); err != nil {
		return &errpath.ErrField{Field: "callbacks", Err: err}
	}

	return nil
}

func strictCallback(c *Callback) error {
	return strictMap(c.ByIndex(), strictRef(strictPathItem))
}

func strictParameters(ps ParameterList) error {
	for i, p := range ps {
		if err := strictRef(strictParameter)(p); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func strictParameter(p *Parameter) error {
	if p.Schema != nil {
		if err := strictSchema(p.Schema); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	}

	if err := strictKeys(p.Examples.ByIndex(), noCheck); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := strictMap(p.Content.ByIndex(), strictMediaType); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	return nil
}

func strictRequestBody(r *RequestBody) error {
	if err := strictMap(r.Content.ByIndex(), strictMediaType); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	return nil
}

func strictResponse(r *Response) error {
	if err := strictKeys(r.Headers.ByIndex(), strictRef(strictHeader)); err != nil {
		return &errpath.ErrField{Field: "headers", Err: err}
	}

	if err := strictMap(r.Content.ByIndex(), strictMediaType); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	if err := strictKeys(r.Links.ByIndex(), noCheck); err != nil {
		return &errpath.ErrField{Field: "links", Err: err}
	}

	return nil
}

func strictHeader(h *Header) error {
	if h.Schema != nil {
		if err := strictSchema(h.Schema); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	}

	if err := strictKeys(h.Examples.ByIndex(), noCheck); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := strictMap(h.Content.ByIndex(), strictMediaType); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	return nil
}

func strictMediaType(mt *MediaType) error {
	if mt.Schema != nil {
		if err := strictSchema(mt.Schema); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	}

	if err := strictKeys(mt.Examples.ByIndex(), noCheck); err != nil {
		return &errpath.ErrField{Field: "examples", Err: err}
	}

	if err := strictMap(mt.Encoding.ByIndex(), strictEncoding); err != nil {
		return &errpath.ErrField{Field: "encoding", Err: err}
	}

	return nil
}

func strictEncoding(e *Encoding) error {
	if err := strictKeys(e.Headers.ByIndex(), strictRef(strictHeader)); err != nil {
		return &errpath.ErrField{Field: "headers", Err: err}
	}

	return nil
}

// strictSchema checks a schema and the schemas within it, but not the one its $ref points to.
func strictSchema(s *Schema) error {
	if s.Format != "" {
		if err := s.Format.Validate(); err != nil {
			return &errpath.ErrField{Field: "format", Err: err}
		}
	}

	for _, sub := range []struct {
		field  string
		schema *Schema
	}{
		{"not", s.Not},
		{"items", s.Items},
		{"additionalProperties", s.additionalPropertiesSchema()},
		{"propertyNames", s.PropertyNames},
	} {
		if sub.schema == nil {
			continue
		}

		if err := strictSchema(sub.schema); err != nil {
			return &errpath.ErrField{Field: sub.field, Err: err}
		}
	}

	for _, list := range []struct {
		field   string
		schemas SchemaList
	}{
		{"allOf", s.AllOf},
		{"oneOf", s.OneOf},
		{"anyOf", s.AnyOf},
		{"prefixItems", s.PrefixItems},
	} {
		for i, sub := range list.schemas {
			if err := strictSchema(sub); err != nil {
				return &errpath.ErrField{Field: list.field, Err: &errpath.ErrIndex{Index: i, Err: err}}
			}
		}
	}

	if err := strictMap(s.Properties.ByIndex(), strictSchema); err != nil {
		return &errpath.ErrField{Field: "properties", Err: err}
	}

	return nil
}

// additionalPropertiesSchema is the schema of the extra properties' values, or nil.
func (s *Schema) additionalPropertiesSchema() *Schema {
	if s.AdditionalProperties == nil {
		return nil
	}

	return s.AdditionalProperties.Schema
}
