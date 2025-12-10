package jsonutil

import (
	"encoding/json/jsontext"
	"fmt"
	"net/url"
)

// URLMarshal is a custom marshaler for URL values, marshaling them as strings.
func URLMarshal(enc *jsontext.Encoder, u url.URL) error {
	return enc.WriteToken(jsontext.String(u.String()))
}

// URLUnmarshal is a custom unmarshaler for URL values, unmarshaling them from strings.
func URLUnmarshal(dec *jsontext.Decoder, u *url.URL) error {
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
