package asyncapi_test

import (
	"net/url"
	"testing"

	"github.com/MarkRosemaker/asyncapi"
)

func TestSecurityScheme_Validate(t *testing.T) {
	t.Parallel()

	for name, s := range map[string]*asyncapi.SecurityScheme{
		"user/password": {Type: asyncapi.SecuritySchemeTypeUserPassword},
		"api key": {
			Type: asyncapi.SecuritySchemeTypeAPIKey,
			In:   asyncapi.SecuritySchemeInUser,
		},
		"X.509": {Type: asyncapi.SecuritySchemeTypeX509},
		"http api key": {
			Type: asyncapi.SecuritySchemeTypeHTTPAPIKey,
			Name: "api_key",
			In:   asyncapi.SecuritySchemeInHeader,
		},
		"http": {Type: asyncapi.SecuritySchemeTypeHTTP, Scheme: "basic"},
		"oauth2": {
			Type: asyncapi.SecuritySchemeTypeOAuth2,
			Flows: &asyncapi.OAuthFlows{
				ClientCredentials: &asyncapi.OAuthFlowClientCredentials{
					TokenURL:        mustParseURL("https://example.com/api/oauth/token"),
					AvailableScopes: asyncapi.MapOfStrings{},
				},
			},
			Scopes: []string{},
		},
		"openid connect": {
			Type:             asyncapi.SecuritySchemeTypeOpenIDConnect,
			OpenIDConnectURL: mustParseURL("https://example.com/.well-known/openid-configuration"),
		},
		"sasl": {Type: asyncapi.SecuritySchemeTypeScramSha256},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if err := s.Validate(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestSecurityScheme_Validate_Errors(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		scheme *asyncapi.SecurityScheme
		want   string
	}{
		"no type": {
			&asyncapi.SecurityScheme{},
			"type is required",
		},
		"unknown type": {
			&asyncapi.SecurityScheme{Type: "carrierPigeon"},
			`type ("carrierPigeon") is invalid, must be one of: "userPassword", "apiKey", ` +
				`"X509", "symmetricEncryption", "asymmetricEncryption", "httpApiKey", "http", ` +
				`"oauth2", "openIdConnect", "plain", "scramSha256", "scramSha512", "gssapi"`,
		},
		"api key without a location": {
			&asyncapi.SecurityScheme{Type: asyncapi.SecuritySchemeTypeAPIKey},
			"in is required",
		},
		"api key in the wrong location": {
			&asyncapi.SecurityScheme{
				Type: asyncapi.SecuritySchemeTypeAPIKey,
				In:   asyncapi.SecuritySchemeInHeader,
			},
			`in ("header") is invalid, must be one of: "user", "password"`,
		},
		"http api key without a name": {
			&asyncapi.SecurityScheme{
				Type: asyncapi.SecuritySchemeTypeHTTPAPIKey,
				In:   asyncapi.SecuritySchemeInHeader,
			},
			"name is required",
		},
		"http api key in the wrong location": {
			&asyncapi.SecurityScheme{
				Type: asyncapi.SecuritySchemeTypeHTTPAPIKey,
				Name: "api_key",
				In:   asyncapi.SecuritySchemeInUser,
			},
			`in ("user") is invalid, must be one of: "query", "header", "cookie"`,
		},
		"http without a scheme": {
			&asyncapi.SecurityScheme{Type: asyncapi.SecuritySchemeTypeHTTP},
			"scheme is required",
		},
		"oauth2 without flows": {
			&asyncapi.SecurityScheme{Type: asyncapi.SecuritySchemeTypeOAuth2},
			"flows is required",
		},
		"oauth2 flow without scopes": {
			&asyncapi.SecurityScheme{
				Type: asyncapi.SecuritySchemeTypeOAuth2,
				Flows: &asyncapi.OAuthFlows{
					Implicit: &asyncapi.OAuthFlowImplicit{
						AuthorizationURL: mustParseURL("https://example.com/api/oauth/dialog"),
					},
				},
			},
			"flows.implicit.availableScopes is required",
		},
		"openid connect without a URL": {
			&asyncapi.SecurityScheme{Type: asyncapi.SecuritySchemeTypeOpenIDConnect},
			"openIdConnectUrl is required",
		},
		"scopes of a scheme that has none": {
			&asyncapi.SecurityScheme{
				Type:   asyncapi.SecuritySchemeTypeX509,
				Scopes: []string{"read"},
			},
			`scopes is invalid: only valid for the types "oauth2" and "openIdConnect"`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := tc.scheme.Validate()
			if err == nil {
				t.Fatal("expected error")
			}

			if err.Error() != tc.want {
				t.Fatalf("got: %v, want: %v", err, tc.want)
			}
		})
	}
}

func TestSecurityScheme_UnifyBearer(t *testing.T) {
	t.Parallel()

	s := &asyncapi.SecurityScheme{Type: asyncapi.SecuritySchemeTypeHTTP, Scheme: "Bearer"}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}

	if got, want := s.Scheme, asyncapi.SecuritySchemeBearer; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}

func TestSecurityScheme_FixScheme(t *testing.T) {
	t.Parallel()

	s := &asyncapi.SecurityScheme{
		Type:             asyncapi.SecuritySchemeTypeOpenIDConnect,
		OpenIDConnectURL: &url.URL{Host: "example.com", Path: "/.well-known/openid-configuration"},
	}
	if err := s.Validate(); err != nil {
		t.Fatal(err)
	}

	if got, want := s.OpenIDConnectURL.Scheme, "https"; got != want {
		t.Fatalf("got: %v, want: %v", got, want)
	}
}
