package openapi

// Servers is a list of server objects.
type Servers []Server

// Validate validates each server.
func (ss Servers) Validate() error {
	for i, s := range ss {
		if err := s.Validate(); err != nil {
			return &ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
