package openapi_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestOperation_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		op  openapi.Operation
		err string
	}{
		{openapi.Operation{
			ExternalDocs: &openapi.ExternalDocs{},
		}, `externalDocs.url is required`},
		{openapi.Operation{
			Parameters: openapi.ParameterList{
				{Value: &openapi.Parameter{}},
			},
		}, `parameters[0].name is required`},
		{openapi.Operation{
			RequestBody: &openapi.RequestBodyRef{
				Value: &openapi.RequestBody{},
			},
		}, `requestBody.content is required`},
		{openapi.Operation{
			Responses: openapi.OperationResponses{ //nolint:exhaustive
				"foo": {},
			},
		}, `responses["foo"]: invalid status code "foo"`},
		{openapi.Operation{
			Responses: openapi.OperationResponses{ //nolint:exhaustive
				"200": {Value: &openapi.Response{}},
			},
		}, `responses["200"].description is required`},
		{openapi.Operation{
			Callbacks: openapi.CallbackRefs{
				"foo": {Value: &openapi.Callback{
					"{$request.query.callbackUrl}/data": &openapi.PathItemRef{
						Value: &openapi.PathItem{
							Extensions: jsontext.Value(`{"bar":"buz"}`),
						},
					},
				}},
			},
		}, `callbacks["foo"]["{$request.query.callbackUrl}/data"].bar: ` + openapi.ErrUnknownField.Error()},
		{openapi.Operation{
			Security: openapi.SecurityRequirements{{"": nil}},
		}, `security[0][""]: empty security scheme name`},
		{openapi.Operation{
			Servers: openapi.Servers{{}},
		}, `servers[0].url is required`},
		{openapi.Operation{
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
		{openapi.Operation{
			Responses: openapi.OperationResponses{
				openapi.StatusCodeDefault: &openapi.ResponseRef{},
			},
		}, `responses["default"]: must not be the only response`},
		{openapi.Operation{
			Responses: openapi.OperationResponses{ //nolint:exhaustive
				"500": &openapi.ResponseRef{},
			},
		}, `responses["500"]: single response must be a successful response`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.op.Validate(); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != tc.err {
				t.Fatalf("want: %v, got: %v", tc.err, err)
			}
		})
	}
}
