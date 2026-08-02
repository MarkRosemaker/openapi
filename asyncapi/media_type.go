package asyncapi

import "mime"

// MediaType is the content type to use when encoding/decoding a message's payload,
// e.g. `application/json`.
type MediaType string

const (
	// MediaTypeJSON is the media type for JSON payloads.
	MediaTypeJSON MediaType = "application/json"
	// MediaTypeYAML is the media type for YAML payloads.
	MediaTypeYAML MediaType = "application/yaml"
	// MediaTypeAvro is the media type for Avro payloads.
	MediaTypeAvro MediaType = "avro/binary"
	// MediaTypeProtobuf is the media type for Protocol Buffers payloads.
	MediaTypeProtobuf MediaType = "application/protobuf"
	// MediaTypeText is the media type for plain text payloads.
	MediaTypeText MediaType = "text/plain"
)

// Validate checks that the media type is well-formed.
func (mt MediaType) Validate() error {
	_, _, err := mime.ParseMediaType(string(mt))
	return err
}
