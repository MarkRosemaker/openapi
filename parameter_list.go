package openapi

import (
	"errors"
	"fmt"

	"github.com/MarkRosemaker/errpath"
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

		errNotUnique := &errpath.ErrIndex{
			Index: i,
			Err: &errpath.ErrField{
				Field: "name",
				Err: &errpath.ErrInvalid[string]{
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
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}

// In is a convenience function to filter by a specific parameter location.
func (p ParameterList) In(in ParameterLocation) ParameterList {
	var result ParameterList
	for _, param := range p {
		if param.Value.In == in {
			result = append(result, param)
		}
	}

	return result
}

// InPath returns all path parameters from the list.
func (p ParameterList) InPath() ParameterList {
	return p.In(ParameterLocationPath)
}

// InQuery returns all query parameters from the list.
func (p ParameterList) InQuery() ParameterList {
	return p.In(ParameterLocationQuery)
}

func (l *loader) resolveParameterList(p ParameterList) error {
	for i, param := range p {
		if err := l.resolveParameterRef(param); err != nil {
			return &errpath.ErrIndex{Index: i, Err: err}
		}
	}

	return nil
}
