package jsonutil

import (
	"encoding/json/jsontext"
	"fmt"
	"maps"
	"net/http"
	"net/textproto"
	"slices"
)

// HTTPHeaderMarshal is a custom marshaler for http.Header, marshaling values as a single strings.
// It also marshals the keys in their canonical form.
// Note that we omit keys that don't have a value.
func HTTPHeaderMarshal(enc *jsontext.Encoder, m http.Header) error {
	if m == nil {
		return enc.WriteToken(jsontext.Null)
	}

	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}

	for _, key := range slices.Sorted(maps.Keys(m)) {
		v := m[key]
		if len(v) == 0 || v[0] == "" {
			continue
		}

		if err := enc.WriteToken(jsontext.String(textproto.CanonicalMIMEHeaderKey(key))); err != nil {
			return err
		}

		if err := enc.WriteToken(jsontext.String(v[0])); err != nil {
			return err
		}
	}

	return enc.WriteToken(jsontext.EndObject)
}

// HTTPHeaderUnmarshal is a custom unmarshaler for http.Header, unmarshaling values as single strings.
func HTTPHeaderUnmarshal(dec *jsontext.Decoder, h *http.Header) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	switch tkn.Kind() {
	case jsontext.KindBeginObject: // expected, continue below
		*h = http.Header{}
	case jsontext.KindNull:
		*h = nil
		return nil // nil map
	default:
		return fmt.Errorf("expected begin object, got %s", tkn.Kind())
	}

	for dec.PeekKind() != jsontext.KindEndObject {
		keyTkn, err := dec.ReadToken()
		if err != nil {
			return err
		}

		if keyTkn.Kind() != jsontext.KindString {
			return fmt.Errorf("expected string key, got %s", keyTkn.Kind())
		}

		key := keyTkn.String()

		val, err := dec.ReadToken()
		if err != nil {
			return err
		}

		if val.Kind() != jsontext.KindString {
			return fmt.Errorf("expected string value, got %s", val.Kind())
		}

		h.Set(key, val.String())
	}

	_, err = dec.ReadToken() // consume jsontext.KindEndObject
	return err
}
