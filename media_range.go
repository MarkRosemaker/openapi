package openapi

import (
	"mime"
)

// MediaRange represents a media type or media type range. It is the key type in the Content map.
// See [RFC7231 Appendix D], and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
//
// Some examples of possible media type definitions:
//
//	text/plain; charset=utf-8
//	application/json
//	application/vnd.github+json
//	application/vnd.github.v3+json
//	application/vnd.github.v3.raw+json
//	application/vnd.github.v3.text+json
//	application/vnd.github.v3.html+json
//	application/vnd.github.v3.full+json
//	application/vnd.github.v3.diff
//	application/vnd.github.v3.patch
//
// [RFC7231]: https://datatracker.ietf.org/doc/html/rfc7231#appendix-D
type MediaRange string

const (
	MediaRangeJSON = "application/json"
	MediaRangeHTML = "text/html"
)

func (mr MediaRange) Validate() error {
	_, _, err := mime.ParseMediaType(string(mr))
	return err
}
