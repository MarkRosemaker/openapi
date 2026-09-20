package openapi

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"reflect"
	"testing"
)

var (
	typeRefEmptyStruct = reflect.TypeFor[refEmptyStruct]()
	typeEmptyStruct    = reflect.TypeFor[emptyStruct]()
)

type (
	refEmptyStruct = refOrValue[emptyStruct, *emptyStruct]
	emptyStruct    struct{}
)

func (emptyStruct) Validate() error { return nil }

func TestRef_UnmarshalJSONV2(t *testing.T) {
	t.Parallel()

	t.Run("reference", func(t *testing.T) {
		if err := json.Unmarshal([]byte(`{"$ref":"#/components/schemas/Pet"`), &refEmptyStruct{}, jsonOpts); err == nil {
			t.Fatal("expected error, got nil")
		} else if synErr, ok := errors.AsType[*jsontext.SyntacticError](err); !ok {
			t.Fatalf("want: *jsontext.SyntacticError, got: %T", err)
		} else if synErr.JSONPointer != "" || synErr.ByteOffset != 34 ||
			synErr.Err.Error() != "unexpected EOF" {
			t.Fatalf("got: %#v", synErr.Err)
		}
	})

	t.Run("object", func(t *testing.T) {
		if err := json.Unmarshal([]byte([]byte(`{"foo":"bar"}`)), &refEmptyStruct{}, jsonOpts); err == nil {
			t.Fatal("expected error, got nil")
		} else if semErr, ok := errors.AsType[*json.SemanticError](err); !ok {
			t.Fatalf("want: *json.SemanticError, got: %T", err)
		} else if semErr.GoType != typeRefEmptyStruct {
			t.Fatalf("want: %s, got: %s", typeRefEmptyStruct, semErr.GoType)
		} else if semErr, ok = errors.AsType[*json.SemanticError](semErr.Err); !ok {
			t.Fatalf("want: *json.SemanticError, got: %T", err)
		} else if semErr.GoType != typeEmptyStruct {
			t.Fatalf("want: %s, got: %s", typeEmptyStruct, semErr.GoType)
		} else if want := "unknown object member name"; semErr.Err.Error() != want {
			t.Fatalf("want: %s, got: %q", want, semErr.Err)
		}
	})
}
