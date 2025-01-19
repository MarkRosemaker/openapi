package openapi

import (
	"mime"
)

// MediaRange represents a media type or media type range. It is the key type in the Content map.
// See [RFC7231 Appendix D], and the value describes it. For requests that match multiple keys, only the most specific key is applicable. e.g. text/plain overrides text/*
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
