package openapi

import "slices"

// ParameterLocation defines the location of the parameter.
// There are four possible parameter locations specified by the `in` field of a Parameter.
type ParameterLocation string

const (
	// Used together with Path Templating, where the parameter value is actually part of the operation's URL. This does not include the host or base path of the API. For example, in `/items/{itemId}`, the path parameter is `itemId`.
	ParameterLocationPath ParameterLocation = "path"
	// Parameters that are appended to the URL. For example, in `/items?id=###`, the query parameter is `id`.
	ParameterLocationQuery ParameterLocation = "query"
	// Custom headers that are expected as part of the request. Note that RFC7230 states header names are case insensitive.
	ParameterLocationHeader ParameterLocation = "header"
	// Used to pass a specific cookie value to the API.
	ParameterLocationCookie ParameterLocation = "cookie"
)

var allParameterLocations = []ParameterLocation{
	ParameterLocationPath,
	ParameterLocationQuery,
	ParameterLocationHeader,
	ParameterLocationCookie,
}

func (p ParameterLocation) Validate() error {
	if p == "" {
		return &ErrRequired{}
	}

	if slices.Contains(allParameterLocations, p) {
		return nil
	}

	return &ErrInvalid[ParameterLocation]{
		Value: p,
		Enum:  allParameterLocations,
	}
}
