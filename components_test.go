package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestComponents_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.Components
		err string
	}{
		{openapi.Components{
			Schemas: openapi.Schemas{"Pet": &openapi.Schema{}},
		}, `schemas["Pet"].type is required`},
		{openapi.Components{
			Responses: openapi.ResponsesByName{"PetResponse": &openapi.ResponseRef{
				Value: &openapi.Response{},
			}},
		}, `responses["PetResponse"].description is required`},
		{openapi.Components{
			Parameters: openapi.Parameters{"MyParameter": &openapi.ParameterRef{
				Value: &openapi.Parameter{},
			}},
		}, `parameters["MyParameter"].name is required`},
		{openapi.Components{
			Examples: openapi.Examples{"MyExample": invalidExample},
		}, `examples["MyExample"]: value and externalValue are mutually exclusive`},
		{openapi.Components{
			RequestBodies: openapi.RequestBodies{"MyRequestBody": &openapi.RequestBodyRef{
				Value: &openapi.RequestBody{},
			}},
		}, `requestBodies["MyRequestBody"].content is required`},
		{openapi.Components{
			Headers: openapi.Headers{"MyRequestBody": &openapi.HeaderRef{
				Value: &openapi.Header{},
			}},
		}, `headers["MyRequestBody"]: schema or content is required`},
		{openapi.Components{
			SecuritySchemes: openapi.SecuritySchemes{"MyRequestBody": &openapi.SecuritySchemeRef{
				Value: &openapi.SecurityScheme{},
			}},
		}, `securitySchemes["MyRequestBody"].type is invalid, must be one of: "apiKey", "http", "mutualTLS", "oauth2", "openIdConnect"`},
		{openapi.Components{
			Links: openapi.Links{"MyLink": &openapi.LinkRef{
				Value: &openapi.Link{},
			}},
		}, `links.MyLink: operationRef or operationId must be set`},
		{openapi.Components{
			Callbacks: openapi.CallbackRefs{"MyCallback": &openapi.CallbackRef{
				Value: invalidCallback,
			}},
		}, `callbacks["MyCallback"]["{$request.query.callbackUrl}/data"].parameters[0].name is required`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
