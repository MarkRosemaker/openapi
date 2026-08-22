package asyncapi

import (
	"errors"

	"github.com/MarkRosemaker/errpath"
)

// Tags is a list of Tag Objects. A Tag Object in a list can be referenced by a Reference Object.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#tagsObject
type Tags []*TagRef

// Validate validates each tag and makes sure that the tag names are unique.
func (tags Tags) Validate() error {
	names := map[string]error{}

	for i, t := range tags {
		if err := t.Validate(); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}

		// a tag that was given as a reference is validated where it is defined
		if t.isRef() {
			continue
		}

		errNotUnique := &errpath.ErrIndex{
			Index: i,
			Err: &errpath.ErrField{
				Field: "name",
				Err:   &errpath.ErrInvalid[string]{Value: t.Value.Name, Message: "must be unique"},
			},
		}

		prevInstance := names[t.Value.Name]
		if prevInstance == nil {
			names[t.Value.Name] = errNotUnique
		} else { // output both instances of the name
			return errors.Join(prevInstance, errNotUnique)
		}
	}

	return nil
}

func (l *loader) collectTags(tags Tags, ref ref) {
	for i, t := range tags {
		l.collectTagRef(t, append(ref, itoa(i)))
	}
}

func (l *loader) resolveTags(tags Tags) error {
	for i, t := range tags {
		if err := l.resolveTagRef(t); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
