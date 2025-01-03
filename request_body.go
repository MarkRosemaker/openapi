package openapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// RequestBody describes a single request body.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#request-body-object
type RequestBody struct {
	// A brief description of the request body. This could contain examples of use.
	// CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Determines if the request body is required in the request. Defaults to `false`.
	Required bool `json:"required,omitempty,omitzero" yaml:"required,omitempty"`
	// REQUIRED. The content of the request body. The key is a media type or media type range and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content Content `json:"content" yaml:"content"`

	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`
}

func (r *RequestBody) Validate() error {
	r.Description = strings.TrimSpace(r.Description)

	if len(r.Content) == 0 {
		return &errpath.ErrField{Field: "content", Err: &errpath.ErrRequired{}}
	}

	if err := r.Content.Validate(); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	return validateExtensions(r.Extensions)
}

func (l *loader) collectRequestBodyRef(r *RequestBodyRef, ref ref) {
	if r.Value != nil {
		l.collectRequestBody(r.Value, ref)
	}
}

func (l *loader) collectRequestBody(r *RequestBody, ref ref) {
	l.requestBodies[ref.String()] = r
}

func (l *loader) resolveRequestBodyRef(r *RequestBodyRef) error {
	return resolveRef(r, l.requestBodies, l.resolveRequestBody)
}
