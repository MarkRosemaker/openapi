package openapi_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestLink_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		link openapi.Link
		err  string
	}{
		{openapi.Link{}, `operationRef or operationId must be set`},
		{openapi.Link{
			OperationRef: "foo",
			OperationID:  "bar",
		}, `operationRef and operationId are mutually exclusive`},
		{openapi.Link{
			OperationRef: "myRef",
			Extensions:   jsontext.Value(`{"foo":"bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.link.Validate(); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != tc.err {
				t.Fatalf("want: %v, got: %v", tc.err, err)
			}
		})
	}
}
