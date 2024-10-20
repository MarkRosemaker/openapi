package jsonutil

import (
	"fmt"
	"net/url"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// URLMarshal is a custom marshaler for URL values, marshaling them as strings.
var URLMarshal = json.MarshalFuncV2(func(
	enc *jsontext.Encoder, u *url.URL, opts json.Options,
) error {
	if u == nil {
		return enc.WriteToken(jsontext.Null)
	}

	return enc.WriteToken(jsontext.String(u.String()))
})

// URLUnmarshal is a custom unmarshaler for URL values, unmarshaling them from strings.
var URLUnmarshal = json.UnmarshalFuncV2(func(
	dec *jsontext.Decoder, u *url.URL, opts json.Options,
) error {
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
	case 'n':
		return nil // no URL given
	default:
		return fmt.Errorf("expected string, got %s", tkn)
	}
})
