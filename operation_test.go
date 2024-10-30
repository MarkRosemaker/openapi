package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
)

func TestOperation_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "tags": [
    "pet"
  ],
  "summary": "Updates a pet in the store with form data",
  "operationId": "updatePetWithForm",
  "parameters": [
    {
      "name": "petId",
      "in": "path",
      "description": "ID of pet that needs to be updated",
      "required": true,
      "schema": {
        "type": "string"
      }
    }
  ],
  "requestBody": {
    "content": {
      "application/x-www-form-urlencoded": {
        "schema": {
          "type": "object",
          "properties": {
            "name": {
              "description": "Updated name of the pet",
              "type": "string"
            },
            "status": {
              "description": "Updated status of the pet",
              "type": "string"
            }
          },
          "required": ["status"]
        }
      }
    }
  },
  "responses": {
    "200": {
      "description": "Pet updated.",
      "content": {
        "application/json": {},
        "application/xml": {}
      }
    },
    "405": {
      "description": "Method Not Allowed",
      "content": {
        "application/json": {},
        "application/xml": {}
      }
    }
  },
  "security": [
    {
      "petstore_auth": [
        "write:pets",
        "read:pets"
      ]
    }
  ]
}`), &openapi.Operation{})

	testJSON(t, []byte(`{
  "tags": [
    "pet"
  ],
  "summary": "Updates a pet in the store with form data",
  "operationId": "updatePetWithForm",
  "parameters": [
    {
      "name": "petId",
      "in": "path",
      "description": "ID of pet that needs to be updated",
      "required": true,
      "schema": {
        "type": "string"
      }
    }
  ],
  "requestBody": {
    "content": {
      "application/x-www-form-urlencoded": {
        "schema": {
          "type": "object",
          "properties": {
            "name": {
              "description": "Updated name of the pet",
              "type": "string"
            },
            "status": {
              "description": "Updated status of the pet",
              "type": "string"
            }
          },
          "required": ["status"]
        }
      }
    }
  },
  "responses": {
    "200": {
      "description": "Pet updated.",
      "content": {
        "application/json": {},
        "application/xml": {}
      }
    },
    "405": {
      "description": "Method Not Allowed",
      "content": {
        "application/json": {},
        "application/xml": {}
      }
    }
  },
  "security": [
    {
      "petstore_auth": [
        "write:pets",
        "read:pets"
      ]
    }
  ],
  "x-foo": "bar",
  "x-bar": 42
}`), &openapi.Operation{})
}

func TestOperation_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		op  openapi.Operation
		err string
	}{
		{openapi.Operation{
			ExternalDocs: &openapi.ExternalDocumentation{},
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
			Responses: openapi.OperationResponses{
				"foo": {},
			},
		}, `responses["foo"]: invalid status code "foo"`},
		{openapi.Operation{
			Responses: openapi.OperationResponses{
				"200": {Value: &openapi.Response{}},
			},
		}, `responses["200"].description is required`},
		{openapi.Operation{
			Callbacks: openapi.Callbacks{
				"foo": {
					"{$request.query.callbackUrl}/data": &openapi.PathItemRef{
						Value: &openapi.PathItem{
							Extensions: jsontext.Value(`{"bar":"buz"}`),
						},
					},
				},
			},
		}, `callbacks["foo"]["{$request.query.callbackUrl}/data"].bar: ` + openapi.ErrUnknownField.Error()},
		{openapi.Operation{
			Security: openapi.Security{{"": nil}},
		}, `security[0][""]: empty security scheme name`},
		{openapi.Operation{
			Servers: openapi.Servers{{}},
		}, `servers[0].url is required`},
		{openapi.Operation{
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: ` + openapi.ErrUnknownField.Error()},
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
