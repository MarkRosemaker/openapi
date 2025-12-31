package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
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
