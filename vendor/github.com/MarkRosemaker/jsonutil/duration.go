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

// DurationMarshalString encodes a time.Duration as a JSON string
// using the canonical units format (e.g. "1h30m", "500ms").
func DurationMarshalString(enc *jsontext.Encoder, d time.Duration) error {
	return enc.WriteToken(jsontext.String(d.String()))
}

// DurationUnmarshalString decodes a JSON string into a time.Duration
// by parsing it with time.ParseDuration.
func DurationUnmarshalString(dec *jsontext.Decoder, d *time.Duration) error {
	var s string
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return err
	}

	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}
