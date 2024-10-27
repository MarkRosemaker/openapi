package openapi

import (
	"errors"
	"fmt"
)

// ParameterList is a list of parameter references.
type ParameterList []*ParameterRef

type parameterID struct {
	Name     string
	Location ParameterLocation
}

func (p ParameterList) Validate() error {
	// The list MUST NOT include duplicated parameters. A unique parameter is defined by a combination of a name and location.
	params := make(map[parameterID]error, len(p))

	for i, param := range p {
		// check for duplicates
		id := parameterID{Name: param.Value.Name, Location: param.Value.In}

		errNotUnique := &ErrIndex{
			Index: i,
			Err: &ErrField{
				Field: "name",
				Err: &ErrInvalid[string]{
					Value:   param.Value.Name,
					Message: fmt.Sprintf("not unique in %s", param.Value.In),
				},
			},
		}

		if prevInstance := params[id]; prevInstance != nil {
			// output both instances of the parameter
			return errors.Join(prevInstance, errNotUnique)
		}
		params[id] = errNotUnique

		if err := param.Validate(); err != nil {
			return &ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
