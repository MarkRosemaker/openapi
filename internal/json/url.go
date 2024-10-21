package json

import (
	"fmt"
	"net/url"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// URLMarshal is a custom marshaler for URL values, marshaling them as strings.
var URLMarshal = json.MarshalFuncV2(func(
	enc *jsontext.Encoder, u *url.URL, _ json.Options,
) error {
	return enc.WriteToken(jsontext.String(u.String()))
})

// URLUnmarshal is a custom unmarshaler for URL values, unmarshaling them from strings.
var URLUnmarshal = json.UnmarshalFuncV2(func(
	dec *jsontext.Decoder, u *url.URL, _ json.Options,
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
	default:
		return fmt.Errorf("expected string, got %s", tkn)
	}
})
