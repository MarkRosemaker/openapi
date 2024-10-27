package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestParameterList_Validate_Error(t *testing.T) {
	t.Parallel()

	err := openapi.ParameterList{{
		Value: &openapi.Parameter{
			Name: "foo", In: openapi.ParameterLocationQuery,
			Schema: &openapi.Schema{Type: openapi.TypeString},
		},
	}, {
		Value: &openapi.Parameter{
			Name: "foo", In: openapi.ParameterLocationQuery,
			Schema: &openapi.Schema{Type: openapi.TypeString},
		},
	}}.Validate()
	if err == nil {
		t.Fatal("expected error")
	} else if want := `[0].name ("foo") is invalid: not unique in query
[1].name ("foo") is invalid: not unique in query`; err.Error() != want {
		t.Fatalf("want: %v, got: %v", want, err)
	}
}
