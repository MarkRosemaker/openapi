package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

// SecurityRequirements lists alternative security requirements, any one of which authorizes a request.
type SecurityRequirements []SecurityRequirement

// Validate returns an error if the SecurityRequirements breaks the specification.
func (ss SecurityRequirements) Validate() error {
	for i, s := range ss {
		if err := s.Validate(); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

// Contains reports whether the list holds a requirement equal to req.
func (ss SecurityRequirements) Contains(req SecurityRequirement) bool {
	return slices.ContainsFunc(ss, req.Equals)
}
