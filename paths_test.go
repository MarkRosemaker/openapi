package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestPaths_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		paths openapi.Paths
		err   string
	}{
		{openapi.Paths{"foo": {}}, `["foo"]: path must start with a /`},
		{openapi.Paths{"/": {
			Get: &openapi.Operation{
				Servers: openapi.Servers{{}},
			},
		}}, `["/"].GET.servers[0].url is required`},
		{openapi.Paths{"/": {
			Get:   &openapi.Operation{OperationID: "myOperation"},
			Patch: &openapi.Operation{OperationID: "myOperation"},
		}}, `["/"].GET.operationId ("myOperation") is invalid: must be unique` + "\n" +
			`["/"].PATCH.operationId ("myOperation") is invalid: must be unique`},
		{openapi.Paths{"/user/{id}": {
			Get: &openapi.Operation{},
		}}, `["/user/{id}"].GET.parameters: {id} not defined`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			t.Parallel()

			if err := tc.paths.Validate(); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != tc.err {
				t.Fatalf("got: %v, want: %v", err, tc.err)
			}
		})
	}
}
