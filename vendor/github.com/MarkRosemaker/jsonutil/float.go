package jsonutil

import (
	"encoding/json/jsontext"
	"fmt"
	"strconv"
)

// Float64UnmarshalStringOrNumber is a custom unmarshaler for float64, unmarshaling them from a JSON number or from a
// string holding one, such as "9.5", as an API that sends numbers as strings does, but for the odd number. Null leaves f
// as it is.
func Float64UnmarshalStringOrNumber(dec *jsontext.Decoder, f *float64) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}

	switch tok.Kind() {
	case jsontext.KindNumber:
		v, err := tok.Float()
		if err != nil {
			return err
		}

		*f = v
	case jsontext.KindString:
		v, err := strconv.ParseFloat(tok.String(), 64)
		if err != nil {
			return fmt.Errorf("parsing %q as a number: %w", tok.String(), err)
		}

		*f = v
	case jsontext.KindNull:
	default:
		return fmt.Errorf("want a number or a string of one, got %s", tok.Kind())
	}

	return nil
}
