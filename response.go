package openapi

import (
	"strings"

	"github.com/MarkRosemaker/errpath"
)

// Response describes a single response from an API Operation, including design-time, static `links` to operations based on the response.
//
// ([Specification])
//
// [Specification]: https://spec.openapis.org/oas/v3.1.0#response-object
type Response struct {
	// REQUIRED. A description of the response. CommonMark syntax MAY be used for rich text representation.
	Description string `json:"description" yaml:"description"`
	// Maps a header name to its definition. RFC7230 states header names are case insensitive. If a response header is defined with the name `"Content-Type"`, it SHALL be ignored.
	Headers Headers `json:"headers,omitempty" yaml:"headers,omitempty"`
	// A map containing descriptions of potential response payloads. The key is a media type or media type range and the value describes it. For responses that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
	Content Content `json:"content,omitempty" yaml:"content,omitempty"`
	// A map of operations links that can be followed from the response. The key of the map is a short name for the link, following the naming constraints of the names for Component Objects.
	Links Links `json:"links,omitempty" yaml:"links,omitempty"`
	// This object MAY be extended with Specification Extensions.
	Extensions Extensions `json:",inline" yaml:"-"`
}

func (r *Response) Validate() error {
	if r.Description == "" {
		return &errpath.ErrField{Field: "description", Err: &errpath.ErrRequired{}}
	}

	r.Description = strings.TrimSpace(r.Description)

	if err := r.Headers.Validate(); err != nil {
		return &errpath.ErrField{Field: "headers", Err: err}
	}

	if err := r.Content.Validate(); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	if err := r.Links.Validate(); err != nil {
		return &errpath.ErrField{Field: "links", Err: err}
	}

	return validateExtensions(r.Extensions)
}

func (l *loader) collectResponseRef(r *ResponseRef, ref ref) {
	if r.Value != nil {
		l.collectResponse(r.Value, ref)
	}
}

func (l *loader) collectResponse(r *Response, ref ref) {
	l.responses[ref.String()] = r
}

func (l *loader) resolveResponseRef(r *ResponseRef) error {
	return resolveRef(r, l.responses, l.resolveResponse)
}

func (l *loader) resolveResponse(r *Response) error {
	if err := l.resolveHeaders(r.Headers); err != nil {
		return &errpath.ErrField{Field: "headers", Err: err}
	}

	if err := l.resolveContent(r.Content); err != nil {
		return &errpath.ErrField{Field: "content", Err: err}
	}

	if err := l.resolveLinks(r.Links); err != nil {
		return &errpath.ErrField{Field: "links", Err: err}
	}

	return nil
}
