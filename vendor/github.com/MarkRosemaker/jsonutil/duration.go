package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"time"
)

// DurationMarshalIntSeconds is a custom marshaler for time.Duration, marshaling them as integers representing seconds.
func DurationMarshalIntSeconds(enc *jsontext.Encoder, d time.Duration) error {
	return enc.WriteToken(jsontext.Int(int64(d / time.Second)))
}

// DurationUnmarshalIntSeconds is a custom unmarshaler for time.Duration, unmarshaling them from integers and assuming they represent seconds.
func DurationUnmarshalIntSeconds(dec *jsontext.Decoder, d *time.Duration) error {
	var seconds int64
	if err := json.UnmarshalDecode(dec, &seconds); err != nil {
		return err
	}

	*d = time.Duration(seconds) * time.Second

	return nil
}
