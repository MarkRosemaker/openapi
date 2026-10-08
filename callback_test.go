package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

var invalidCallback = &openapi.Callback{
	"{$request.query.callbackUrl}/data": {
		Value: &openapi.PathItem{
			Parameters: openapi.ParameterList{{
				Value: &openapi.Parameter{},
			}},
		},
	},
}

func TestCallback_Validate_Error(t *testing.T) {
	if err := invalidCallback.Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `["{$request.query.callbackUrl}/data"].parameters[0].name is required`; want != err.Error() {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}
