package types

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"regexp"
)

// regex from https://emailregex.com/
var reEmail = regexp.MustCompile(`(?m)(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])`)

// ErrInvalidEmail is returned when an email is invalid.
var ErrInvalidEmail = errors.New("invalid email")

// Email is a string that represents an email address.
type Email string

// Validate checks if the email is valid.
func (e Email) Validate() error {
	if !reEmail.MatchString(string(e)) {
		return fmt.Errorf("%w: %q", ErrInvalidEmail, e)
	}

	return nil
}

var _ json.UnmarshalerFrom = (*Email)(nil)

func (e *Email) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	s := ""
	if err := json.UnmarshalDecode(dec, &s); err != nil {
		return err
	}

	if s == "" {
		return nil // no email
	}

	*e = Email(s)

	return e.Validate()
}
