package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"time"

	"cloud.google.com/go/civil"
)

// DateMarshalIntUnix is a custom marshaler for civil.Date, marshaling them as integers representing unix time.
func DateMarshalIntUnix(enc *jsontext.Encoder, d civil.Date) error {
	if d.IsZero() {
		return enc.WriteToken(jsontext.Int(0))
	}

	return enc.WriteToken(jsontext.Int(int64(
		time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC).Unix(),
	)))
}

// DateUnmarshalIntUnix is a custom unmarshaler for civil.Date, unmarshaling them from integers and assuming they represent unix time.
func DateUnmarshalIntUnix(dec *jsontext.Decoder, d *civil.Date) error {
	var seconds int64
	if err := json.UnmarshalDecode(dec, &seconds); err != nil {
		return err
	}

	if seconds == 0 {
		*d = civil.Date{}
	} else {
		*d = civil.DateOf(time.Unix(seconds, 0))
	}

	return nil
}
