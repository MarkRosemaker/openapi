package openapi

import (
	"slices"

	"github.com/MarkRosemaker/errpath"
)

type SecurityRequirements []SecurityRequirement

func (ss SecurityRequirements) Validate() error {
	for i, s := range ss {
		if err := s.Validate(); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

func (ss SecurityRequirements) Contains(req SecurityRequirement) bool {
	return slices.ContainsFunc(ss, req.Equals)
}
