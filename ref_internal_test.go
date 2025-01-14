package openapi

import (
	"errors"
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
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

func errAs[T any, E interface {
	*T
	error
}](t *testing.T, err error,
) E {
	t.Helper()

	var zero T
	target := E(&zero)
	if !errors.As(err, &target) {
		t.Fatalf("want: %T, got: %T", target, err)
	}

	return target
}

func TestRef_UnmarshalJSONV2(t *testing.T) {
	t.Parallel()

	t.Run("reference", func(t *testing.T) {
		err := json.Unmarshal([]byte(`{"$ref":"#/components/schemas/Pet"`),
			&refEmptyStruct{}, jsonOpts)
		synErr := errAs[jsontext.SyntacticError](t, err)
		if synErr.JSONPointer != "" || synErr.ByteOffset != 34 ||
			synErr.Err.Error() != "unexpected EOF" {
			t.Fatalf("got: %#v", synErr.Err)
		}
	})

	t.Run("object", func(t *testing.T) {
		err := json.Unmarshal([]byte([]byte(`{"foo":"bar"}`)),
			&refEmptyStruct{}, jsonOpts)
		semErr := errAs[json.SemanticError](t, err)
		if semErr.GoType != typeRefEmptyStruct {
			t.Fatalf("want: %s, got: %s", typeRefEmptyStruct, semErr.GoType)
		} else if semErr = errAs[json.SemanticError](t, semErr.Err); semErr.GoType != typeEmptyStruct {
			t.Fatalf("want: %s, got: %s", typeEmptyStruct, semErr.GoType)
		} else if want := "unknown object member name"; semErr.Err.Error() != want {
			t.Fatalf("want: %s, got: %q", want, semErr.Err)
		}
	})
}
