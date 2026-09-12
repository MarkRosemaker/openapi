package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
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

// timeAltLayout is a fallback layout for time strings that don't conform to RFC3339.
const timeAltLayout = "Mon Jan 2 2006 15:04:05 MST-0700"

// TimeUnmarshalStringOrIntUnix unmarshals a time.Time from either an RFC3339 or
// timeAltLayout string, or an integer representing unix seconds. Nulls decode as the zero time.
func TimeUnmarshalStringOrIntUnix(dec *jsontext.Decoder, d *time.Time) error {
	return timeUnmarshalStringOrIntUnix([]string{time.RFC3339, timeAltLayout})(dec, d)
}

// timeUnmarshalStringOrIntUnix returns a custom unmarshaler for time.Time that unmarshals
// from either an integer representing unix seconds or a string, tried against each of the
// given layouts (see time.Parse) in order until one succeeds. Nulls decode as the zero time.
func timeUnmarshalStringOrIntUnix(layouts []string) func(dec *jsontext.Decoder, d *time.Time) error {
	return func(dec *jsontext.Decoder, d *time.Time) error {
		tkn, err := dec.ReadToken()
		if err != nil {
			return err
		}

		switch tkn.Kind() {
		case jsontext.KindNumber:
			seconds, err := tkn.Int()
			if err != nil {
				return err
			}

			if seconds == 0 {
				*d = time.Time{}
			} else {
				*d = time.Unix(seconds, 0)
			}
		case jsontext.KindString:
			s := tkn.String()

			var errs error
			for _, layout := range layouts {
				t, err := time.Parse(layout, s)
				if err == nil {
					*d = t
					return nil
				}

				errs = errors.Join(errs, err)
			}

			return fmt.Errorf("could not parse %q using any of the given layouts: %w", s, errs)
		case jsontext.KindNull: // ok, nothing to do
		default:
			return fmt.Errorf("unknown token kind %s", tkn.Kind())
		}

		return nil
	}
}
