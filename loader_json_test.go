package openapi

import (
	"encoding/json/jsontext"
	"errors"
	"testing"
)

func TestLoadFromDataJSON_Error(t *testing.T) {
	if _, err := LoadFromDataJSON([]byte(`{`)); err == nil {
		t.Fatal("expected error, got nil")
	} else if syntErr, ok := errors.AsType[*jsontext.SyntacticError](err); !ok {
		t.Fatalf("want: jsontext.SyntacticError, got: %T", err)
	} else if want := `unexpected EOF`; syntErr.Err == nil || syntErr.Err.Error() != want {
		t.Fatalf("want: %s, got: %v", want, syntErr.Err)
	}
}
