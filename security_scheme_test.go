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
