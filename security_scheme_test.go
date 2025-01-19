package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

func TestSecurityScheme_JSON(t *testing.T) {
	t.Parallel()

	// Basic Authentication Sample
	testJSON(t, []byte(`{
  "type": "http",
  "scheme": "basic"
}`), &openapi.SecurityScheme{})

	// API Key Sample
	testJSON(t, []byte(`{
  "type": "apiKey",
  "name": "api_key",
  "in": "header"
}`), &openapi.SecurityScheme{})

	// JWT Bearer Sample
	testJSON(t, []byte(`{
  "type": "http",
  "scheme": "bearer",
  "bearerFormat": "JWT"
}`), &openapi.SecurityScheme{})

	// Implicit OAuth2 Sample
	testJSON(t, []byte(`{
  "type": "oauth2",
  "flows": {
    "implicit": {
      "authorizationUrl": "https://example.com/api/oauth/dialog",
      "scopes": {
        "write:pets": "modify pets in your account",
        "read:pets": "read your pets"
      }
    }
  }
}`), &openapi.SecurityScheme{})
}

func TestSecurityScheme_Validate(t *testing.T) {
	t.Parallel()

	for _, tc := range []openapi.SecurityScheme{
		{Type: openapi.SecuritySchemeTypeMutualTLS},
	} {
		if err := tc.Validate(); err != nil {
			t.Fatalf("%#v got error: %s", tc, err)
		}
	}
}

func TestSecurityScheme_Validate_Error(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		ss  openapi.SecurityScheme
		err string
	}{
		{
			openapi.SecurityScheme{},
			`type is required`,
		},
		{
			openapi.SecurityScheme{Type: openapi.SecuritySchemeTypeAPIKey},
			`name is required`,
		},
		{
			openapi.SecurityScheme{
				Name: "api_key",
				Type: openapi.SecuritySchemeTypeAPIKey,
			},
			`in is required`,
		},
		{
			openapi.SecurityScheme{
				Name: "api_key",
				Type: openapi.SecuritySchemeTypeAPIKey,
				In:   "foo",
			},
			`in ("foo") is invalid, must be one of: "query", "header", "cookie"`,
		},
		{
			openapi.SecurityScheme{Type: openapi.SecuritySchemeTypeHTTP},
			`scheme is required`,
		},
		{
			openapi.SecurityScheme{
				Type:   openapi.SecuritySchemeTypeHTTP,
				Scheme: "BeAReR",
			},
			`bearerFormat is required`,
		},
		{
			openapi.SecurityScheme{
				Type:   openapi.SecuritySchemeTypeHTTP,
				Scheme: "BeAReR",
			},
			`bearerFormat is required`,
		},
		{
			openapi.SecurityScheme{Type: openapi.SecuritySchemeTypeOAuth2},
			`flows is required`,
		},
		{
			openapi.SecurityScheme{
				Type: openapi.SecuritySchemeTypeOAuth2,
				Flows: &openapi.OAuthFlows{
					Implicit: &openapi.OAuthFlowImplicit{},
				},
			},
			`flows.implicit.authorizationUrl is required`,
		},
		{
			openapi.SecurityScheme{
				Type: openapi.SecuritySchemeTypeOpenIDConnect,
			},
			`openIdConnectUrl is required`,
		},
	} {
		t.Run(tc.err, func(t *testing.T) {
			if err := tc.ss.Validate(); err == nil || err.Error() != tc.err {
				t.Fatalf("want: %s, got: %s", tc.err, err)
			}
		})
	}
}
