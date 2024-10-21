package openapi_test

import (
	"testing"

	"github.com/MarkRosemaker/openapi"
)

type objectWithServers struct {
	Servers openapi.Servers `json:"servers"`
}

func (o objectWithServers) Validate() error {
	return o.Servers.Validate()
}

func TestServers_JSON(t *testing.T) {
	t.Parallel()

	// from the documentation:
	// The following shows how multiple servers can be described, for example, at the OpenAPI Object's servers:
	testJSON(t, []byte(`{
  "servers": [
    {
      "url": "https://development.gigantic-server.com/v1",
      "description": "Development server"
    },
    {
      "url": "https://staging.gigantic-server.com/v1",
      "description": "Staging server"
    },
    {
      "url": "https://api.gigantic-server.com/v1",
      "description": "Production server"
    }
  ]
}`), &objectWithServers{})

	// from the documentation:
	// The following shows how variables can be used for a server configuration:
	testJSON(t, []byte(`{
  "servers": [
    {
      "url": "https://{username}.gigantic-server.com:{port}/{basePath}",
      "description": "The production API server",
      "variables": {
        "username": {
          "default": "demo",
          "description": "this value is assigned by the service provider, in this example `+"`"+`gigantic-server.com`+"`"+`"
        },
        "port": {
          "enum": [
            "8443",
            "443"
          ],
          "default": "8443"
        },
        "basePath": {
          "default": "v2"
        }
      }
    }
  ]
}`), &objectWithServers{})
}

func TestValidate(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		s := openapi.Servers{}
		if err := s.Validate(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid server", func(t *testing.T) {
		s := openapi.Servers{{}}
		if err := s.Validate(); err == nil {
			t.Fatal("expected error")
		} else if want := "0: url is required"; err.Error() != want {
			t.Fatalf("got: %v, want: %v", err, want)
		}
	})
}
