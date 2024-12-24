package jsonutil

import (
	"time"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
)

// DurationMarshalIntSeconds is a custom marshaler for time.Duration, marshaling them as integers representing seconds.
func DurationMarshalIntSeconds(enc *jsontext.Encoder, d time.Duration, opts json.Options) error {
	return enc.WriteToken(jsontext.Int(int64(d / time.Second)))
}

// DurationUnmarshalIntSeconds is a custom unmarshaler for time.Duration, unmarshaling them from integers and assuming they represent seconds.
func DurationUnmarshalIntSeconds(dec *jsontext.Decoder, d *time.Duration, opts json.Options) error {
	var seconds int64
	if err := json.UnmarshalDecode(dec, &seconds, opts); err != nil {
		return err
	}

	*d = time.Duration(seconds) * time.Second

	return nil
}
