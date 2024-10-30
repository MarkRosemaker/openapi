package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
	"github.com/go-json-experiment/json/jsontext"
)

func TestLink_JSON(t *testing.T) {
	t.Parallel()

	// Computing a link from a request operation where the `$request.path.id` is used to pass a request parameter to the linked operation.
	testJSON(t, []byte(`{
    "/users/{id}": {
      "parameters": [
        {
          "name": "id",
          "in": "path",
          "description": "the user identifier, as userId",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "get": {
        "responses": {
          "200": {
            "description": "the user being returned",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "uuid": {
                      "type": "string",
                      "format": "uuid"
                    }
                  }
                }
              }
            },
            "links": {
              "address": {
                "operationId": "getUserAddress",
                "parameters": {
                  "userId": "$request.path.id"
                }
              }
            }
          }
        }
      }
    },
    "/users/{userid}/address": {
      "parameters": [
        {
          "name": "userid",
          "in": "path",
          "description": "the user identifier, as userId",
          "required": true,
          "schema": {
            "type": "string"
          }
        }
      ],
      "get": {
        "operationId": "getUserAddress",
        "responses": {
          "200": {
            "description": "the user's address"
          }
        }
      }
    }
  }`), &openapi.Paths{})

	// When a runtime expression fails to evaluate, no parameter value is passed to the target operation.
	// Values from the response body can be used to drive a linked operation.
	testJSON(t, []byte(`{
    "address": {
      "operationId": "getUserAddressByUUID",
      "parameters": {
        "userUuid": "$response.body#/uuid"
      }
    }
  }`), &openapi.Links{})

	// As references to `operationId` MAY NOT be possible (the `operationId` is an optional field in an [Operation Object](#operation-object)), references MAY also be made through a relative `operationRef`:
	testJSON(t, []byte(`{
    "UserRepositories": {
      "operationRef": "#/paths/~12.0~1repositories~1{username}/get",
      "parameters": {
        "username": "$response.body#/username"
      }
    }
  }`), &openapi.Links{})

	// or an absolute `operationRef`:
	testJSON(t, []byte(`{
    "UserRepositories": {
      "operationRef": "https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get",
      "parameters": {
        "username": "$response.body#/username"
      }
    }
  }`), &openapi.Links{})
}

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
		{openapi.Link{
			OperationRef: "foo",
			Parameters: openapi.LinkParameters{
				"": &openapi.LinkParameter{},
			},
		}, `parameters[""] is required`},
		{openapi.Link{
			OperationRef: "foo",
			Parameters: openapi.LinkParameters{
				"username": &openapi.LinkParameter{},
			},
		}, `parameters["username"] is required`},
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
