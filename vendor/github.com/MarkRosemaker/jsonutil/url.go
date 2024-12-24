package jsonutil

import (
	"fmt"
	"net/url"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// URLMarshal is a custom marshaler for URL values, marshaling them as strings.
func URLMarshal(enc *jsontext.Encoder, u url.URL, opts json.Options) error {
	return enc.WriteToken(jsontext.String(u.String()))
}

// URLUnmarshal is a custom unmarshaler for URL values, unmarshaling them from strings.
func URLUnmarshal(dec *jsontext.Decoder, u *url.URL, opts json.Options) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	switch tkn.Kind() {
	case '"':
		parsed, err := url.Parse(tkn.String())
		if err != nil {
			return err
		}

		*u = *parsed

		return nil
	case 'n': // null
		return nil // no url given
	default:
		return fmt.Errorf("expected string, got %s", tkn)
	}
}
