package openapi

import (
	"errors"
	"reflect"
	"testing"

	_json "github.com/MarkRosemaker/openapi/internal/json"
	"github.com/go-json-experiment/json"
)

var typeRefEmptyStruct = reflect.TypeFor[refEmptyStruct]()

type (
	refEmptyStruct = refOrValue[emptyStruct, *emptyStruct]
	emptyStruct    struct{}
)

func (emptyStruct) Validate() error { return nil }

func isJSONSemanticError(t *testing.T, err error, goType reflect.Type) error {
	t.Helper()

	jsonErr := &json.SemanticError{}
	if !errors.As(err, &jsonErr) {
		t.Fatalf("expected json.SemanticError, got: %v", err)
	}

	if goType != jsonErr.GoType {
		t.Fatalf("json.SemanticError is not of type %s but of %s", goType, jsonErr.GoType)
	}

	return jsonErr.Err
}

func TestRef_UnmarshalJSONV2(t *testing.T) {
	t.Parallel()

	t.Run("reference", func(t *testing.T) {
		err := json.Unmarshal([]byte(`{"$ref":"#/components/schemas/Pet"`),
			&refEmptyStruct{}, _json.Options)
		err = isJSONSemanticError(t, err, typeRefEmptyStruct)

		if want := "unexpected EOF"; err.Error() != want {
			t.Fatalf("want: %s, got: %s", want, err)
		}
	})

	t.Run("object", func(t *testing.T) {
		err := json.Unmarshal([]byte([]byte(`{"foo":"bar"}`)),
			&refEmptyStruct{}, _json.Options)
		err = isJSONSemanticError(t, err, typeRefEmptyStruct)
		err = isJSONSemanticError(t, err, reflect.TypeFor[emptyStruct]())

		if want := `unknown name "foo"`; err.Error() != want {
			t.Fatalf("want: %s, got: %s", want, err)
		}
	})
}
