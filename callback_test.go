package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestCallback_JSON(t *testing.T) {
	t.Parallel()

	// The following example uses the user provided `queryUrl` query string parameter to define the callback URL.
	// This is an example of how to use a callback object to describe a WebHook callback that goes with the subscription operation to enable registering for the WebHook.
	testJSON(t, []byte(`{
    "{$request.query.queryUrl}": {
      "post": {
        "requestBody": {
          "description": "Callback payload",
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SomePayload"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "callback successfully processed"
          }
        }
      }
    }
}
`), &openapi.Callback{})

	// The following example shows a callback where the server is hard-coded, but the query string parameters are populated from the `id` and `email` property in the request body.
	testJSON(t, []byte(`{
    "http://notificationServer.com?transactionId={$request.body#/id}&email={$request.body#/email}": {
      "post": {
        "requestBody": {
          "description": "Callback payload",
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/SomePayload"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "callback successfully processed"
          }
        }
      }
    }
}`), &openapi.Callback{})
}

func TestCallback_Validate_Error(t *testing.T) {
	t.Parallel()

	c := openapi.Callback{
		"{$request.query.callbackUrl}/data": {
			Value: &openapi.PathItem{
				Parameters: openapi.ParameterList{{
					Value: &openapi.Parameter{},
				}},
			},
		},
	}

	if err := c.Validate(); err == nil {
		t.Fatal("expected error")
	} else if want := `["{$request.query.callbackUrl}/data"].parameters[0].name is required`; want != err.Error() {
		t.Fatalf("expected %q, got %q", want, err.Error())
	}
}
