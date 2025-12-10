package openapi_test

import (
	"encoding/json/jsontext"
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestOAuthFlows_JSON(t *testing.T) {
	t.Parallel()

	testJSON(t, []byte(`{
  "type": "oauth2",
  "flows": {
    "implicit": {
      "authorizationUrl": "https://example.com/api/oauth/dialog",
      "scopes": {
        "write:pets": "modify pets in your account",
        "read:pets": "read your pets"
      }
    },
    "authorizationCode": {
      "authorizationUrl": "https://example.com/api/oauth/dialog",
      "tokenUrl": "https://example.com/api/oauth/token",
      "scopes": {
        "write:pets": "modify pets in your account",
        "read:pets": "read your pets"
      }
    }
  }
}`), &openapi.SecurityScheme{})
}

func TestOAuthFlows_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		c   openapi.OAuthFlows
		err string
	}{
		{openapi.OAuthFlows{
			Implicit: &openapi.OAuthFlowImplicit{},
		}, `implicit.authorizationUrl is required`},
		{openapi.OAuthFlows{
			Password: &openapi.OAuthFlowPassword{},
		}, `password.tokenUrl is required`},
		{openapi.OAuthFlows{
			ClientCredentials: &openapi.OAuthFlowClientCredentials{},
		}, `clientCredentials.tokenUrl is required`},
		{openapi.OAuthFlows{
			AuthorizationCode: &openapi.OAuthFlowAuthorizationCode{},
		}, `authorizationCode.authorizationUrl is required`},
		{openapi.OAuthFlows{
			Extensions: jsontext.Value(`{"foo": "bar"}`),
		}, `foo: unknown field or extension without "x-" prefix`},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.c.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
