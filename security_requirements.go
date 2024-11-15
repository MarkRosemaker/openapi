package openapi

type SecurityRequirements []SecurityRequirement

func (ss SecurityRequirements) Validate() error {
	for i, s := range ss {
		if err := s.Validate(); err != nil {
			return &ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
