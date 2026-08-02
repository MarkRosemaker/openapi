package asyncapi

import (
	"errors"
	"net/url"
	"regexp"

	"github.com/MarkRosemaker/errpath"
)

// ErrEmptyDocument is thrown if the AsyncAPI document neither defines channels nor operations nor components.
var ErrEmptyDocument = errors.New("document must contain at least a channels field, an operations field or a components field")

// reServerKey is the regular expression the keys of the servers object must match.
var reServerKey = regexp.MustCompile(`^[A-Za-z0-9_\-]+$`)

// Document is an AsyncAPI document, the root object.
// It combines resource listing and API declaration together into one document.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#A2SObject
type Document struct {
	// REQUIRED. Specifies the AsyncAPI Specification version being used.
	// It can be used by tooling Specifications and clients to interpret the version.
	// The structure shall be `major`.`minor`.`patch`, where `patch` versions MUST be compatible
	// with the existing `major`.`minor` tooling.
	AsyncAPI string `json:"asyncapi" yaml:"asyncapi"`
	// Identifier of the application the AsyncAPI document is defining.
	// It must conform to the URI format.
	// It is RECOMMENDED to use a URN to globally and uniquely identify the application
	// during long periods of time, even after it becomes unavailable or ceases to exist.
	ID *url.URL `json:"id,omitempty" yaml:"id,omitempty"`
	// REQUIRED. Provides metadata about the API. The metadata can be used by the clients if needed.
	Info *Info `json:"info,omitempty" yaml:"info,omitempty"`
	// Provides connection details of servers.
	Servers Servers `json:"servers,omitempty" yaml:"servers,omitempty"`
	// Default content type to use when encoding/decoding a message's payload.
	// The value MUST be a specific media type (e.g. `application/json`).
	// This value MUST be used by schema parsers when the contentType property is omitted.
	DefaultContentType MediaType `json:"defaultContentType,omitempty" yaml:"defaultContentType,omitempty"`
	// The channels used by this application.
	Channels Channels `json:"channels,omitempty" yaml:"channels,omitempty"`
	// The operations this application MUST implement.
	Operations Operations `json:"operations,omitempty" yaml:"operations,omitempty"`
	// An element to hold various reusable objects for the specification.
	// Everything that is defined inside this object represents a resource that MAY or MAY NOT be
	// used in the rest of the document and MAY or MAY NOT be used by the implemented application.
	Components Components `json:"components,omitzero" yaml:"components,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// reAsyncAPIVersion is a regular expression that matches the AsyncAPI version.
// Only version 3.x.y is supported.
// See: https://www.asyncapi.com/docs/reference/specification/v3.1.0#A2SVersionString
var reAsyncAPIVersion = regexp.MustCompile(`^3\.\d+\.\d+(-.+)?$`)

// Validate checks the AsyncAPI document for correctness.
func (d *Document) Validate() error {
	if d.AsyncAPI == "" {
		return &errpath.ErrField{Field: "asyncapi", Err: &errpath.ErrRequired{}}
	}

	if !reAsyncAPIVersion.MatchString(d.AsyncAPI) {
		return &errpath.ErrField{
			Field: "asyncapi",
			Err: &errpath.ErrInvalid[string]{
				Value:   d.AsyncAPI,
				Message: "must be a valid version (3.x.y)",
			},
		}
	}

	if d.Info == nil {
		return &errpath.ErrField{Field: "info", Err: &errpath.ErrRequired{}}
	}

	if err := d.Info.Validate(); err != nil {
		return &errpath.ErrField{Field: "info", Err: err}
	}

	// the keys of the servers object must match a certain pattern
	for name := range d.Servers {
		if !reServerKey.MatchString(name) {
			return &errpath.ErrField{Field: "servers", Err: &errpath.ErrKey{
				Key: name,
				Err: &errpath.ErrInvalid[string]{
					Value:   name,
					Message: "must match the regular expression " + `"` + reServerKey.String() + `"`,
				},
			}}
		}
	}

	if err := d.Servers.Validate(); err != nil {
		return &errpath.ErrField{Field: "servers", Err: err}
	}

	if d.DefaultContentType != "" {
		if err := d.DefaultContentType.Validate(); err != nil {
			return &errpath.ErrField{Field: "defaultContentType", Err: err}
		}
	}

	// an AsyncAPI document that neither describes channels nor operations nor
	// reusable components doesn't describe anything at all
	if len(d.Channels) == 0 && len(d.Operations) == 0 && d.Components.isEmpty() {
		return ErrEmptyDocument
	}

	if err := d.Channels.Validate(); err != nil {
		return &errpath.ErrField{Field: "channels", Err: err}
	}

	if err := d.Operations.Validate(); err != nil {
		return &errpath.ErrField{Field: "operations", Err: err}
	}

	if err := d.Components.Validate(); err != nil {
		return &errpath.ErrField{Field: "components", Err: err}
	}

	return validateExtensions(d.Extensions)
}

// SortMaps sorts the servers, channels, operations and the fields of the components that are maps by key.
func (d *Document) SortMaps() {
	d.Servers.Sort()
	d.Channels.Sort()
	d.Operations.Sort()
	d.Components.SortMaps()
}

func (l *loader) collectDocument(doc *Document, ref ref) {
	l.collectInfo(doc.Info, append(ref, "info"))
	l.collectServers(doc.Servers, append(ref, "servers"))
	l.collectChannels(doc.Channels, append(ref, "channels"))
	l.collectOperations(doc.Operations, append(ref, "operations"))
	l.collectComponents(doc.Components, append(ref, "components"))
}

func (l *loader) resolveDocument(doc *Document) error {
	// fields that don't need to be resolved:
	// - AsyncAPI
	// - ID
	// - DefaultContentType

	if err := l.resolveInfo(doc.Info); err != nil {
		return &errpath.ErrField{Field: "info", Err: err}
	}

	if err := l.resolveServers(doc.Servers); err != nil {
		return &errpath.ErrField{Field: "servers", Err: err}
	}

	if err := l.resolveChannels(doc.Channels); err != nil {
		return &errpath.ErrField{Field: "channels", Err: err}
	}

	if err := l.resolveOperations(doc.Operations); err != nil {
		return &errpath.ErrField{Field: "operations", Err: err}
	}

	if err := l.resolveComponents(doc.Components); err != nil {
		return &errpath.ErrField{Field: "components", Err: err}
	}

	return nil
}
