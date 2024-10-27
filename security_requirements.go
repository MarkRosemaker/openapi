package openapi

type Security []SecurityRequirement

func (ss Security) Validate() error {
	for i, s := range ss {
		if err := s.Validate(); err != nil {
			return &ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
