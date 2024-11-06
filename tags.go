package openapi

type Tags []*Tag

func (tags Tags) Validate() error {
	for i, t := range tags {
		if err := t.Validate(); err != nil {
			return &ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
