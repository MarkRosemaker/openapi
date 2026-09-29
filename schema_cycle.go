package openapi

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// checkSchemaCycles rejects component schemas that lead back to themselves while applying to the same value.
//
// A schema may refer to itself through a property or an item, since each step descends into the value:
// that is how a tree or a linked list is described. But $ref, allOf, anyOf, oneOf and not apply to the
// value itself, so a cycle of those alone never gets anywhere. JSON Schema 2020-12 leaves such
// "infinite recursive nesting" undefined and requires that validators not fall into an infinite loop.
func checkSchemaCycles(ss Schemas) error {
	f := &cycleFinder{names: make(map[*Schema]string, len(ss)), done: map[*Schema]bool{}}
	for name, s := range ss {
		f.names[s] = name
	}

	for _, s := range ss.ByIndex() {
		if cycle := f.visit(s); cycle != nil {
			return &errpath.ErrKey{Key: cycle[0], Err: fmt.Errorf(
				"cycle that never descends into the value: %s", strings.Join(cycle, " → "),
			)}
		}
	}

	return nil
}

type cycleFinder struct {
	names map[*Schema]string
	done  map[*Schema]bool
	stack []*Schema
}

// visit returns the names along a cycle reachable from component schema s, starting and ending with the same name.
func (f *cycleFinder) visit(s *Schema) []string {
	if f.done[s] {
		return nil
	}

	if i := slices.Index(f.stack, s); i >= 0 {
		cycle := make([]string, 0, len(f.stack)-i+1)
		for _, c := range f.stack[i:] {
			cycle = append(cycle, f.names[c])
		}

		return append(cycle, f.names[s])
	}

	f.stack = append(f.stack, s)
	cycle := f.inPlace(s)
	f.stack = f.stack[:len(f.stack)-1]
	f.done[s] = true

	return cycle
}

// inPlace follows the keywords of s, a component or a schema within one, that apply to the same value as s.
func (f *cycleFinder) inPlace(s *Schema) []string {
	if s.Ref != nil && s.Ref.Value != nil {
		if cycle := f.visit(s.Ref.Value); cycle != nil {
			return cycle
		}
	}

	for _, list := range []SchemaList{s.AllOf, s.AnyOf, s.OneOf, {s.Not}} {
		for _, sub := range list {
			if sub == nil {
				continue
			}

			if cycle := f.inPlace(sub); cycle != nil {
				return cycle
			}
		}
	}

	return nil
}
