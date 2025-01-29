package openapi

import (
	"fmt"
	"strings"
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

func (l *loader) resolveCallbacks(cs Callbacks) error {
	return nil
}

func (l *loader) collectWebhooks(ws Webhooks, ref ref) {
}
