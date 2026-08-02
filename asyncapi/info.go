package asyncapi

import (
	"net/url"
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// The Info object provides metadata about the API.
// The metadata can be used by the clients if needed.
// ([Specification])
//
// [Specification]: https://www.asyncapi.com/docs/reference/specification/v3.1.0#infoObject
type Info struct {
	// REQUIRED. The title of the application.
	Title string `json:"title" yaml:"title"`
	// REQUIRED. Provides the version of the application API (not to be confused with the specification version).
	Version string `json:"version" yaml:"version"`
	// A short description of the application. CommonMark syntax can be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// A URL to the Terms of Service for the API. This MUST be in the form of an absolute URL.
	TermsOfService *url.URL `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	// The contact information for the exposed API.
	Contact *Contact `json:"contact,omitempty" yaml:"contact,omitempty"`
	// The license information for the exposed API.
	License *License `json:"license,omitempty" yaml:"license,omitempty"`
	// A list of tags for application API documentation control. Tags can be used for logical grouping of applications.
	Tags Tags `json:"tags,omitempty" yaml:"tags,omitempty"`
	// Additional external documentation of the exposed API.
	ExternalDocs *ExternalDocsRef `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:",inline"`
}

// Validate checks the info object for correctness.
func (i *Info) Validate() error {
	if i.Title == "" {
		return &errpath.ErrField{Field: "title", Err: &errpath.ErrRequired{}}
	}

	// NOTE: The version *here* is the version of the application API,
	// it can be any string, unlike the version of the specification.
	if i.Version == "" {
		return &errpath.ErrField{Field: "version", Err: &errpath.ErrRequired{}}
	}

	i.Description = strings.TrimSpace(i.Description)

	// assume that the scheme is https and add it if it is missing
	fixScheme(i.TermsOfService)

	if i.Contact != nil {
		if err := i.Contact.Validate(); err != nil {
			return &errpath.ErrField{Field: "contact", Err: err}
		}
	}

	if i.License != nil {
		if err := i.License.Validate(); err != nil {
			return &errpath.ErrField{Field: "license", Err: err}
		}
	}

	if err := i.Tags.Validate(); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if i.ExternalDocs != nil {
		if err := i.ExternalDocs.Validate(); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return validateExtensions(i.Extensions)
}

// fixScheme ensures that the URL has a scheme and that it is valid.
// If the URL is nil, it is a no-op.
func fixScheme(u *url.URL) {
	if u == nil {
		return
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}
}

func (l *loader) collectInfo(i *Info, ref ref) {
	if i == nil {
		return
	}

	l.collectTags(i.Tags, append(ref, "tags"))

	if i.ExternalDocs != nil {
		l.collectExternalDocsRef(i.ExternalDocs, append(ref, "externalDocs"))
	}
}

func (l *loader) resolveInfo(i *Info) error {
	if i == nil {
		return nil
	}

	if err := l.resolveTags(i.Tags); err != nil {
		return &errpath.ErrField{Field: "tags", Err: err}
	}

	if i.ExternalDocs != nil {
		if err := l.resolveExternalDocsRef(i.ExternalDocs); err != nil {
			return &errpath.ErrField{Field: "externalDocs", Err: err}
		}
	}

	return nil
}
