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

func TestParameterList_In(t *testing.T) {
	t.Parallel()

	var list openapi.ParameterList

	if inPath := list.InPath(); inPath != nil {
		t.Fatalf("want: nil, got: %v", inPath)
	}

	if inQuery := list.InQuery(); inQuery != nil {
		t.Fatalf("want: nil, got: %v", inQuery)
	}

	list = append(list, &openapi.ParameterRef{
		Value: &openapi.Parameter{
			Name: "foo", In: openapi.ParameterLocationQuery,
			Schema: &openapi.Schema{Type: openapi.TypeString},
		},
	}, &openapi.ParameterRef{
		Value: &openapi.Parameter{
			Name: "bar", In: openapi.ParameterLocationPath,
			Schema: &openapi.Schema{Type: openapi.TypeString},
		},
	})

	if inPath := list.InPath(); len(inPath) != 1 {
		t.Fatalf("want: 1, got: %v", len(inPath))
	} else if want := "bar"; inPath[0].Value.Name != want {
		t.Fatalf("want: %q, got: %q", want, inPath[0].Value.Name)
	}

	if inQuery := list.InQuery(); len(inQuery) != 1 {
		t.Fatalf("want: 1, got: %v", len(inQuery))
	} else if want := "foo"; inQuery[0].Value.Name != want {
		t.Fatalf("want: %q, got: %q", want, inQuery[0].Value.Name)
	}
}
