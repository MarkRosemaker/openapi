package openapi

import (
	"fmt"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

type ref []string

func (r ref) String() string {
	return strings.Join(r, "/")
}

// collectResolveRefs expands references in a document that was just unmarshaled
func (l *loader) collectResolveRefs(doc *Document) error {
	// collect all the references
	l.collectDocument(doc, []string{"#"})

	// resolve all the references
	if err := l.resolveDocument(doc); err != nil {
		return err
	}

	return nil
}

// resolveRef resolves a reference to a value or resolves the value itself
func resolveRef[T any, O referencable[T]](
	r *refOrValue[T, O], values map[string]*T, resolveValue func(*T) error,
) error {
	if r.Ref != nil && r.Value == nil {
		if val, ok := values[r.Ref.Identifier]; ok {
			r.Value = val
			return nil
		}

		return fmt.Errorf("couldn't resolve %q", r.Ref.Identifier)
	}

	if resolveValue == nil {
		return nil
	}

	return resolveValue(r.Value)
}

func (l *loader) collectPaths(ps Paths, ref ref) {
}

func (l *loader) resolveOperation(o *Operation) error {
	if err := l.resolveParameterList(o.Parameters); err != nil {
		return &errpath.ErrField{Field: "parameters", Err: err}
	}

	if o.RequestBody != nil {
		if err := l.resolveRequestBodyRef(o.RequestBody); err != nil {
			return &errpath.ErrField{Field: "requestBody", Err: err}
		}
	}

	if err := l.resolveOperationResponses(o.Responses); err != nil {
		return &errpath.ErrField{Field: "responses", Err: err}
	}

	if err := l.resolveCallbacks(o.Callbacks); err != nil {
		return &errpath.ErrField{Field: "callbacks", Err: err}
	}

	return nil
}

func (l *loader) resolveRequestBody(r *RequestBody) error {
	if err := l.resolveContent(r.Content); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	return nil
}

func (l *loader) resolveMediaType(mt *MediaType) error {
	if mt.Schema != nil {
		if err := l.resolveSchemaRef(mt.Schema); err != nil {
			return &errpath.ErrField{Field: "schema", Err: err}
		}
	}

	// if mt.Example != nil && mt.Examples != nil {
	// 	return errors.New("example and examples are mutually exclusive")
	// }

	// if err := mt.Examples.Validate(); err != nil {
	// 	return &errpath.ErrField{Field: "examples", Err: err}
	// }

	// if err := mt.Encoding.Validate(); err != nil {
	// 	return &errpath.ErrField{Field: "encoding", Err: err}
	// }

	return nil
}

func (l *loader) resolveOperationResponses(rs OperationResponses) error {
	for code, r := range rs.ByIndex() {
		if err := l.resolveResponseRef(r); err != nil {
			return &errpath.ErrKey{Key: string(code), Err: err}
		}
	}

	return nil
}

func (l *loader) resolveCallbacks(cs Callbacks) error {
	return nil
}

func (l *loader) collectWebhooks(ws Webhooks, ref ref) {
}
