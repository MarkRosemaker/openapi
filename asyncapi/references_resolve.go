package asyncapi

import (
	"fmt"
	"strconv"
	"strings"
)

// ref is the path to an object within the document, i.e. a JSON pointer split into its tokens.
//
// "The Reference Object is defined by JSON Reference and follows the same structure, behavior
// and rules. [...] For this specification, reference resolution is done as defined by the JSON
// Reference specification and not by the JSON Schema specification." ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#referenceObject
type ref []string

func (r ref) String() string {
	return strings.Join(r, "/")
}

// itoa converts an index of a list to a token of a JSON pointer.
func itoa(i int) string { return strconv.Itoa(i) }

// collectResolveRefs expands references in a document that was just unmarshaled.
func (l *loader) collectResolveRefs(doc *Document) error {
	// collect all the objects that can be referenced
	l.collectDocument(doc, []string{"#"})

	// let the paths that hold a reference point to the object that is ultimately referenced
	l.dereference()
	l.linkAliases()

	// resolve all the references
	return l.resolveDocument(doc)
}

// linkAliases makes the objects that are referenced available under the paths
// of the references that point to them.
func (l *loader) linkAliases() {
	linkAliases(l, l.schemas)
	linkAliases(l, l.servers)
	linkAliases(l, l.serverVariables)
	linkAliases(l, l.channels)
	linkAliases(l, l.operations)
	linkAliases(l, l.operationTraits)
	linkAliases(l, l.replies)
	linkAliases(l, l.replyAddresses)
	linkAliases(l, l.messages)
	linkAliases(l, l.messageTraits)
	linkAliases(l, l.parameters)
	linkAliases(l, l.correlationIDs)
	linkAliases(l, l.securitySchemes)
	linkAliases(l, l.tags)
	linkAliases(l, l.externalDocs)
	linkAliases(l, l.bindings)
}

func linkAliases[T any](l *loader, values map[string]*T) {
	for path, target := range l.aliases {
		if v, ok := values[target]; ok {
			values[path] = v
		}
	}
}

// collectRef remembers where an object is defined so that references to it can be resolved.
//
// It reports whether the object itself was given. If it was given as a reference, the reference
// is remembered as well, because it may be referenced in turn, e.g. an operation refers to a
// message of a channel which in turn refers to a message of the components object.
func collectRef[T any, O referencable[T]](
	l *loader, r *refOrValue[T, O], values map[string]*T, ref ref,
) bool {
	if r.Ref != nil {
		l.aliases[ref.String()] = r.Ref.Identifier
		return false
	}

	if r.Value == nil {
		return false
	}

	values[ref.String()] = (*T)(r.Value)

	return true
}

// resolveRef resolves a reference to a value or resolves the value itself.
func resolveRef[T any, O referencable[T]](
	r *refOrValue[T, O], values map[string]*T, resolveValue func(*T) error,
) error {
	if r.Ref != nil && r.Value == nil {
		val, ok := values[r.Ref.Identifier]
		if !ok {
			return fmt.Errorf("couldn't resolve %q", r.Ref.Identifier)
		}

		r.Value = val

		return nil
	}

	if resolveValue == nil || r.Value == nil {
		return nil
	}

	return resolveValue(r.Value)
}

// dereference follows the references that point to other references
// so that every reference points to the object it ultimately refers to.
func (l *loader) dereference() {
	for from, to := range l.aliases {
		seen := map[string]bool{from: true}

		for {
			next, ok := l.aliases[to]
			if !ok || seen[to] {
				break
			}

			seen[to] = true
			to = next
		}

		l.aliases[from] = to
	}
}
