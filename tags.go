package openapi

import "errors"

// A list of tags used by the document with additional metadata. The order of the tags can be used to reflect on their order by the parsing tools. Not all tags that are used by the Operation Object must be declared. The tags that are not declared MAY be organized randomly or based on the tools' logic.
type Tags []*Tag

func (tags Tags) Validate() error {
	// Each tag name in the list MUST be unique.
	names := map[string]error{}

	for i, t := range tags {
		errNotUnique := &ErrIndex{
			Index: i,
			Err: &ErrField{
				Field: "name",
				Err:   &ErrInvalid[string]{Value: t.Name, Message: "must be unique"},
			},
		}

		prevInstance := names[t.Name]
		if prevInstance == nil {
			names[t.Name] = errNotUnique
		} else { // output both instances of the name
			return errors.Join(prevInstance, errNotUnique)
		}

		if err := t.Validate(); err != nil {
			return &ErrIndex{Index: i, Err: err}
		}

	}

	return nil
}
