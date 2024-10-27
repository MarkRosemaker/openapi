package openapi

import "errors"

// Path represents a URL path template.
//
// Templating refers to the usage of template expressions, delimited by curly braces ({}), to mark a section of a URL path as replaceable using path parameters.
//
// Each template expression in the path MUST correspond to a path parameter that is included in the Path Item itself and/or in each of the Path Item's Operations. An exception is if the path item is empty, for example due to ACL constraints, matching path parameters are not required.
//
// The value for these path parameters MUST NOT contain any unescaped "generic syntax" characters described by RFC3986: forward slashes (/), question marks (?), or hashes (#).
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#path-templating
type Path string

// Validate checks the path for correctness.
func (p Path) Validate() error {
	if p == "" {
		return errors.New("path must not be empty")
	}

	if p[0] != '/' {
		return errors.New("path must start with a /")
	}

	return nil
}
