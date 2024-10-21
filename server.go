package openapi

// Server is an object representing a Server.
// ([Documentation])
//
// [Documentation]: https://spec.openapis.org/oas/v3.1.0#server-object
type Server struct {
	// REQUIRED. A URL to the target host. This URL supports Server Variables and MAY be relative, to indicate that the host location is relative to the location where the OpenAPI document is being served. Variable substitutions will be made when a variable is named in {brackets}.
	URL string `json:"url,strictcase" yaml:"url"`
	// URL Path `json:"url,strictcase" yaml:"url"` TODO: Path type

	// An optional string describing the host designated by the URL. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty,strictcase" yaml:"description,omitempty"`
	// A map between a variable name and its value. The value is used for substitution in the server's URL template.
	Variables ServerVariables `json:"variables,omitempty,strictcase" yaml:"variables,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

func (s *Server) Validate() error {
	if s.URL == "" {
		return &ErrRequired{Target: "url"}
	}

	if err := s.Variables.Validate(); err != nil {
		return &ErrField{Field: "variables", Err: err}
	}

	// // validate the default URL to see if the URL is well-formed
	// if _, err := s.defaultURL(); err != nil {
	// 	return fmt.Errorf("URL: %w", err)
	// }

	return nil
}
