package jsonutil

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
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

// DateMarshalBasic is a custom marshaler for civil.Date, marshaling them as strings in ISO 8601's basic form, without
// separators, such as "20260320". The zero date marshals as null.
func DateMarshalBasic(enc *jsontext.Encoder, d civil.Date) error {
	if d.IsZero() {
		return enc.WriteToken(jsontext.Null)
	}

	if d.Year < 0 || d.Year > 9999 {
		return fmt.Errorf("year %d does not fit four digits", d.Year)
	}

	b := [8]byte{
		byte('0' + d.Year/1000), byte('0' + d.Year/100%10), byte('0' + d.Year/10%10), byte('0' + d.Year%10),
		byte('0' + d.Month/10), byte('0' + d.Month%10),
		byte('0' + d.Day/10), byte('0' + d.Day%10),
	}

	return enc.WriteToken(jsontext.String(string(b[:])))
}

// DateUnmarshalBasic is a custom unmarshaler for civil.Date, unmarshaling them from strings in ISO 8601's basic form,
// such as "20260320". Null unmarshals as the zero date.
func DateUnmarshalBasic(dec *jsontext.Decoder, d *civil.Date) error {
	v, err := dec.ReadValue()
	if err != nil {
		return err
	}

	if v.Kind() == jsontext.KindNull {
		*d = civil.Date{}

		return nil
	}

	// eight digits within quotes; a digit is never escaped, so the raw value is all there is to read
	if len(v) != 10 || v[0] != '"' || v[9] != '"' {
		return fmt.Errorf("%s is no date such as \"20260320\"", v)
	}

	n := 0

	for _, c := range v[1:9] {
		if c < '0' || c > '9' {
			return fmt.Errorf("%s is no date such as \"20260320\"", v)
		}

		n = n*10 + int(c-'0')
	}

	date := civil.Date{Year: n / 10000, Month: time.Month(n / 100 % 100), Day: n % 100}
	if !date.IsValid() {
		return fmt.Errorf("%s is no valid date", v)
	}

	*d = date

	return nil
}

// DateMarshalBasicFor returns a marshaler for civil.Date that writes, as [DateMarshalBasic] does, the dates held by an
// object member of one of the names given, or by an array such a member holds; see [DateMarshalIntUnixFor].
func DateMarshalBasicFor(names ...string) func(*jsontext.Encoder, civil.Date) error {
	set := nameSet(names)

	return func(enc *jsontext.Encoder, d civil.Date) error {
		if !set[memberName(enc.StackPointer())] {
			return errors.ErrUnsupported
		}

		return DateMarshalBasic(enc, d)
	}
}

// DateUnmarshalBasicFor returns an unmarshaler for civil.Date that reads, as [DateUnmarshalBasic] does, the dates held
// by an object member of one of the names given, or by an array such a member holds; see [DateMarshalIntUnixFor].
func DateUnmarshalBasicFor(names ...string) func(*jsontext.Decoder, *civil.Date) error {
	set := nameSet(names)

	return func(dec *jsontext.Decoder, d *civil.Date) error {
		if !set[memberName(dec.StackPointer())] {
			return errors.ErrUnsupported
		}

		return DateUnmarshalBasic(dec, d)
	}
}

// DateUnmarshalEmptyAsZero is a custom unmarshaler for civil.Date that reads an empty string as the zero date, as an
// API that leaves a date out may send, and anything else as civil.Date itself does.
func DateUnmarshalEmptyAsZero(dec *jsontext.Decoder, d *civil.Date) error {
	return unmarshalEmptyAsZero(dec, d)
}
