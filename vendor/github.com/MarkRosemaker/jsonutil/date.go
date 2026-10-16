package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
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
		*d = civil.DateOf(time.Unix(seconds, 0).UTC())
	}

	return nil
}

// DateMarshalIntUnixFor returns a marshaler for civil.Date that writes, as [DateMarshalIntUnix] does, the dates held by
// an object member of one of the names given, or by an array such a member holds. It leaves any other date to the
// marshaler after it, or to the default, which writes a string: for an API that writes some dates as unix time and
// others as strings.
func DateMarshalIntUnixFor(names ...string) func(*jsontext.Encoder, civil.Date) error {
	set := nameSet(names)

	return func(enc *jsontext.Encoder, d civil.Date) error {
		if !set[memberName(enc.StackPointer())] {
			return errors.ErrUnsupported
		}

		return DateMarshalIntUnix(enc, d)
	}
}

// DateUnmarshalIntUnixFor returns an unmarshaler for civil.Date that reads, as [DateUnmarshalIntUnix] does, the dates
// held by an object member of one of the names given, or by an array such a member holds; see [DateMarshalIntUnixFor].
func DateUnmarshalIntUnixFor(names ...string) func(*jsontext.Decoder, *civil.Date) error {
	set := nameSet(names)

	return func(dec *jsontext.Decoder, d *civil.Date) error {
		if !set[memberName(dec.StackPointer())] {
			return errors.ErrUnsupported
		}

		return DateUnmarshalIntUnix(dec, d)
	}
}
