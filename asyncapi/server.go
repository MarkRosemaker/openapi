package asyncapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Server is an object representing a message broker, a server or any other kind of computer program capable of sending and/or receiving data.
// This object is used to capture details such as URIs, protocols and security configuration.
// Variable substitution can be used so that some details, for example usernames and passwords, can be injected by code generation tools.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#serverObject
type Server struct {
	// REQUIRED. The server host name. It MAY include the port.
	// This field supports Server Variables. Variable substitutions will be made when a variable is named in {braces}.
	Host string `json:"host" yaml:"host"`
	// REQUIRED. The protocol this server supports for connection.
	Protocol Protocol `json:"protocol" yaml:"protocol"`
	// The version of the protocol used for connection. For instance: AMQP 0.9.1, HTTP 2.0, Kafka 1.0.0, etc.
	ProtocolVersion string `json:"protocolVersion,omitempty" yaml:"protocolVersion,omitempty"`
	// The path to a resource in the host.
	// This field supports Server Variables. Variable substitutions will be made when a variable is named in {braces}.
	Pathname string `json:"pathname,omitempty" yaml:"pathname,omitempty"`
	// An optional string describing the server. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A human-friendly title for the server.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	// A short summary of the server.
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`
	// A map between a variable name and its value. The value is used for substitution in the server's host and pathname template.
	Variables ServerVariables `json:"variables,omitempty" yaml:"variables,omitempty"`
	// A declaration of which security schemes can be used with this server.
	// The list of values includes alternative security scheme objects that can be used.
	// Only one of the security scheme objects need to be satisfied to authorize a connection or operation.
	Security SecuritySchemeRefList `json:"security,omitempty" yaml:"security,omitempty"`
	// A list of tags for logical grouping and categorization of servers.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation for this server.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// A map where the keys describe the name of the protocol and the values describe protocol-specific definitions for the server.
	Bindings *BindingsRef `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the server for correctness.
func (s *Server) Validate() error {
	if s.Host == "" {
		return &errpath.ErrField{Field: "host", Err: &errpath.ErrRequired{}}
	}

	if s.Protocol == "" {
		return &errpath.ErrField{Field: "protocol", Err: &errpath.ErrRequired{}}
	}

	s.Description = strings.TrimSpace(s.Description)

	if err := s.Variables.Validate(); err != nil {
		return &errpath.ErrField{Field: "variables", Err: err}
	}

	if err := s.Security.Validate(); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := s.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if s.ExternalDocs != nil {
		if err := s.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if s.Bindings != nil {
		if err := s.Bindings.Validate(); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	return validateExtensions(s.Extensions)
}

func (l *loader) collectServerRef(s *ServerRef, ref ref) {
	if !collectRef(l, s, l.servers, ref) {
		return
	}

	l.collectServer(s.Value, ref)
}

func (l *loader) collectServer(s *Server, ref ref) {
	l.collectServerVariables(s.Variables, append(ref, "variables"))
	l.collectSecuritySchemeRefList(s.Security, append(ref, "security"))
	l.collectTags(s.Tags, append(ref, "tags"))

	if s.ExternalDocs != nil {
		l.collectExternalDocsRef(s.ExternalDocs, append(ref, "externalDocs"))
	}

	if s.Bindings != nil {
		l.collectBindingsRef(s.Bindings, append(ref, "bindings"))
	}
}

func (l *loader) resolveServerRef(s *ServerRef) error {
	return resolveRef(s, l.servers, l.resolveServer)
}

func (l *loader) resolveServer(s *Server) error {
	if err := l.resolveServerVariables(s.Variables); err != nil {
		return &errpath.ErrField{Field: "variables", Err: err}
	}

	if err := l.resolveSecuritySchemeRefList(s.Security); err != nil {
		return &errpath.ErrField{Field: "security", Err: err}
	}

	if err := l.resolveTags(s.Tags); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if s.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(s.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	if s.Bindings != nil {
		if err := l.resolveBindingsRef(s.Bindings); err != nil {
			return &errpath.ErrField{Field: "bindings", Err: err}
		}
	}

	return nil
}
