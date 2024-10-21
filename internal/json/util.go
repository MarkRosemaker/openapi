package json

import (
	"fmt"

	"github.com/go-json-experiment/json/jsontext"
)

// skipTokenKind reads a token from the decoder and checks if it is of the expected kind
func skipTokenKind(dec *jsontext.Decoder, kind jsontext.Kind) error {
	tkn, err := dec.ReadToken()
	if err != nil {
		return err
	}

	if got := tkn.Kind(); got != kind {
		return fmt.Errorf("expected %s, got %s", kind, got)
	}

	return nil
}
