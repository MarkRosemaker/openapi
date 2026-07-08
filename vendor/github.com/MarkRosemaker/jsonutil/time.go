package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"time"
)

// TimeMarshalIntUnix is a custom marshaler for time.Time, marshaling them as integers representing unix time.
func TimeMarshalIntUnix(enc *jsontext.Encoder, t time.Time) error {
	if t.IsZero() {
		return enc.WriteToken(jsontext.Int(0))
	}

	return enc.WriteToken(jsontext.Int(int64(t.Unix())))
}

// TimeUnmarshalIntUnix is a custom unmarshaler for time.Time, unmarshaling them from integers and assuming they represent unix time.
func TimeUnmarshalIntUnix(dec *jsontext.Decoder, d *time.Time) error {
	var seconds int64
	if err := json.UnmarshalDecode(dec, &seconds); err != nil {
		return err
	}

	if seconds == 0 {
		*d = time.Time{}
	} else {
		*d = time.Unix(seconds, 0)
	}

	return nil
}

// TimeUnmarshalStringOrIntUnix unmarshals a time.Time from either an RFC3339
// string or an integer representing unix seconds. Nulls decode as the zero time.
func TimeUnmarshalStringOrIntUnix(dec *jsontext.Decoder, d *time.Time) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	switch tkn.Kind() {
	case jsontext.KindNumber:
		if seconds := tkn.Int(); seconds == 0 {
			*d = time.Time{}
		} else {
			*d = time.Unix(seconds, 0)
		}
	case jsontext.KindString:
		return d.UnmarshalText([]byte(tkn.String()))
	case jsontext.KindNull: // ok, nothing to do
	default:
		return fmt.Errorf("unknown token kind %s", tkn.Kind())
	}

	return nil
}
